package smtp

import (
	"io"
	"strings"

	esmtp "github.com/emersion/go-smtp"
	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
	"go.uber.org/zap"
)

type Session struct {
	mail *config.Email
}

func (s *Session) Mail(
	from string,
	opts *esmtp.MailOptions,
) error {

	if s.mail == nil {
		s.Reset()
	}

	s.mail.EnvelopeFrom = from

	logger.Info(
		"MAIL FROM",
		zap.String("from", from),
	)

	return nil
}

func (s *Session) Rcpt(
	to string,
	opts *esmtp.RcptOptions,
) error {

	if s.mail == nil {
		s.Reset()
	}

	s.mail.EnvelopeTo = append(
		s.mail.EnvelopeTo,
		to,
	)

	logger.Info(
		"RCPT TO",
		zap.String("to", to),
	)

	return nil
}

func (s *Session) Data(r io.Reader) error {

	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	s.mail.Raw = raw

	//----------------------------------------------------------
	// Loop Detection
	//----------------------------------------------------------

	if strings.Contains(
		strings.ToLower(string(raw)),
		"x-fehmi-gateway: processed",
	) {

		logger.Info(
			"Message already processed. Skipping.",
		)

		return nil
	}

	//----------------------------------------------------------
	// Save Original
	//----------------------------------------------------------

	// if Conf.SaveRawEML {

	// 	if err := SaveEML(s.mail); err != nil {

	// 		logger.Error(
	// 			"save original eml",
	// 			zap.Error(err),
	// 		)
	// 	}
	// }

	//----------------------------------------------------------
	// Parse Message
	//----------------------------------------------------------

	email, err := ParseMessage(raw)
	if err != nil {

		logger.Error(
			"parse message",
			zap.Error(err),
		)

		return err
	}

	email.Raw = raw
	email.EnvelopeFrom = s.mail.EnvelopeFrom
	email.EnvelopeTo = s.mail.EnvelopeTo

	//----------------------------------------------------------
	// Inject HTML Signature
	//----------------------------------------------------------

	if strings.TrimSpace(email.HTML) != "" {

		email.HTML = ParseBody(
			email.HTML,
			email.EnvelopeFrom,
		)
	}

	//----------------------------------------------------------
	// Inject Plain Signature
	//----------------------------------------------------------

	if strings.TrimSpace(email.Text) != "" {

		email.Text = ParseBody(
			email.Text,
			email.EnvelopeFrom,
		)
	}

	//----------------------------------------------------------
	// Build RFC822
	//----------------------------------------------------------

	newRaw, err := Build(email)
	if err != nil {

		logger.Error(
			"build message",
			zap.Error(err),
		)

		return err
	}

	email.Raw = newRaw

	s.mail = email

	//----------------------------------------------------------
	// Save Edited
	//----------------------------------------------------------

	if config.SmtpC.SaveRawEML {

		if err := SaveEditedEML(email); err != nil {

			logger.Error(
				"save edited eml",
				zap.Error(err),
			)
		}
	}

	logger.Info(
		"Message Processed",
		zap.String("from", email.EnvelopeFrom),
		zap.Strings("to", email.EnvelopeTo),
		zap.String("subject", email.Subject),
	)

	//----------------------------------------------------------
	// Relay
	//----------------------------------------------------------

	if err := Relay(email); err != nil {

		logger.Error(
			"relay failed",
			zap.Error(err),
		)

		return err
	}

	logger.Info("Relay Successful")

	return nil
}

func (s *Session) Reset() {

	logger.Info("SMTP Session Reset")

	s.mail = &config.Email{
		Headers: make(map[string][]string),
	}
}

func (s *Session) Logout() error {

	logger.Info("SMTP Session Closed")

	return nil
}
