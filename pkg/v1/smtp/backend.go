package smtp

import (
	"net"

	esmtp "github.com/emersion/go-smtp"
	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
	"go.uber.org/zap"
)

// Backend implements the SMTP backend.
type Backend struct{}

// NewBackend creates a new SMTP backend.
func NewBackend() *Backend {
	return &Backend{}
}

// NewSession is called for every incoming SMTP connection.
func (b *Backend) NewSession(conn *esmtp.Conn) (esmtp.Session, error) {

	remote := ""
	local := ""

	if c, ok := conn.Conn().(net.Conn); ok {
		if c.RemoteAddr() != nil {
			remote = c.RemoteAddr().String()
		}
		if c.LocalAddr() != nil {
			local = c.LocalAddr().String()
		}
	}

	logger.Info(
		"New SMTP Session",
		zap.String("client", remote),
		zap.String("server", local),
	)

	return &Session{
		mail: &config.Email{
			Headers: make(map[string][]string),
		},
	}, nil
}
