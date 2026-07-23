package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
	"go.uber.org/zap"
)

// loginAuth implements the RFC-compliant 2-step LOGIN authentication scheme required by Exchange Online
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Send mechanism name only; no initial client response payload
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// Handle base64 encoded challenges ("Username:" -> "VXNZXJuYW1lOiA=", "Password:" -> "UGFzc3dvcmQ6")
		switch string(fromServer) {
		case "VXNZXJuYW1lOiA=", "Username:":
			return []byte(a.username), nil
		case "UGFzc3dvcmQ6", "Password:":
			return []byte(a.password), nil
		default:
			return []byte(a.password), nil
		}
	}
	return nil, nil
}

// Relay sends the email message via the configured upstream SMTP server
func Relay(email *config.Email) error {
	addr := fmt.Sprintf(
		"%s:%d",
		config.SmtpC.RelayHost,
		config.SmtpC.RelayPort,
	)

	logger.Info(
		"Relaying Email",
		zap.String("server", addr),
		zap.String("from", email.EnvelopeFrom),
		zap.Strings("to", email.EnvelopeTo),
		zap.String("subject", email.Subject),
	)

	// 1. Dial TCP Connection
	c, err := smtp.Dial(addr)
	if err != nil {
		logger.Error("Relay Dial Failed", zap.Error(err))
		return fmt.Errorf("failed to dial relay host %s: %w", addr, err)
	}
	defer c.Close()

	// 2. Issue EHLO
	if err := c.Hello("localhost"); err != nil {
		logger.Error("Relay EHLO Failed", zap.Error(err))
		return fmt.Errorf("EHLO failed: %w", err)
	}

	// 3. Upgrade to TLS via STARTTLS (if enabled or on port 587)
	if config.SmtpC.RelayTLS || config.SmtpC.RelayPort == 587 {
		tlsConfig := &tls.Config{
			ServerName: config.SmtpC.RelayHost,
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			logger.Error("Relay STARTTLS Failed", zap.Error(err))
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	// 4. Authenticate using LOGIN auth (with PLAIN auth fallback)
	if config.SmtpC.RelayUsername != "" && config.SmtpC.RelayPassword != "" {
		auth := LoginAuth(config.SmtpC.RelayUsername, config.SmtpC.RelayPassword)
		if err := c.Auth(auth); err != nil {
			// Fallback to PLAIN auth if LOGIN fails
			plainAuth := smtp.PlainAuth("", config.SmtpC.RelayUsername, config.SmtpC.RelayPassword, config.SmtpC.RelayHost)
			if err := c.Auth(plainAuth); err != nil {
				logger.Error("Relay Auth Failed", zap.Error(err))
				return fmt.Errorf("smtp authentication failed: %w", err)
			}
		}
	}

	// 5. Envelope Sender
	if err := c.Mail(email.EnvelopeFrom); err != nil {
		logger.Error("Relay MAIL FROM Failed", zap.Error(err))
		return fmt.Errorf("MAIL FROM command failed: %w", err)
	}

	// 6. Envelope Recipients
	for _, recipient := range email.EnvelopeTo {
		if err := c.Rcpt(recipient); err != nil {
			logger.Error("Relay RCPT TO Failed", zap.String("recipient", recipient), zap.Error(err))
			return fmt.Errorf("RCPT TO command failed for %s: %w", recipient, err)
		}
	}

	// 7. Stream Raw EML Data
	w, err := c.Data()
	if err != nil {
		logger.Error("Relay DATA Command Failed", zap.Error(err))
		return fmt.Errorf("DATA command failed: %w", err)
	}

	if _, err := w.Write(email.Raw); err != nil {
		logger.Error("Relay Write Body Failed", zap.Error(err))
		return fmt.Errorf("failed writing raw message body: %w", err)
	}

	if err := w.Close(); err != nil {
		logger.Error("Relay Close Data Writer Failed", zap.Error(err))
		return fmt.Errorf("failed closing data writer: %w", err)
	}

	// 8. Gracefully terminate connection
	if err := c.Quit(); err != nil {
		// Log warning but don't fail transaction if QUIT fails after DATA stream succeeds
		logger.Info("Relay QUIT Warning", zap.Error(err))
	}

	logger.Info("Relay Successful")
	return nil
}
