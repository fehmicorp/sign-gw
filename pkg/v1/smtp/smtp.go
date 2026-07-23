package smtp

import (
	"crypto/tls"
	"fmt"

	esmtp "github.com/emersion/go-smtp"
	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
)

// ----------------------------------------------------------------------
// SMTP Server
// ----------------------------------------------------------------------

func Start() error {
	cfg := config.SmtpC
	logger.Info("server initiated")

	server := esmtp.NewServer(NewBackend())

	// ------------------------------------------------------------------
	// Listener
	// ------------------------------------------------------------------

	server.Addr = fmt.Sprintf(
		"%s:%d",
		cfg.ListenHost,
		cfg.ListenPort,
	)

	server.Domain = cfg.Hostname

	server.AllowInsecureAuth = cfg.AllowInsecure

	server.MaxRecipients = cfg.MaxRecipients

	server.MaxMessageBytes = cfg.MaxMessageSize

	// ------------------------------------------------------------------
	// STARTTLS
	// ------------------------------------------------------------------

	if cfg.UseTLS {

		cert, err := tls.LoadX509KeyPair(
			"certs/server.crt",
			"certs/server.key",
		)

		if err != nil {
			return fmt.Errorf(
				"load TLS certificate: %w",
				err,
			)
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			ServerName:   cfg.Hostname,
		}
	}

	// ------------------------------------------------------------------
	// Start Server
	// ------------------------------------------------------------------

	return server.ListenAndServe()
}
