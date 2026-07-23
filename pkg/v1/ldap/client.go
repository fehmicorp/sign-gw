package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/go-ldap/ldap/v3"
)

var Conn *ldap.Conn

func Connect() (*ldap.Conn, error) {
	cfg := config.LdapC
	address := fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)

	var err error

	if cfg.UseTLS {

		Conn, err = ldap.DialTLS(
			"tcp",
			address,
			&tls.Config{
				InsecureSkipVerify: true,
			},
		)

	} else {

		Conn, err = ldap.Dial("tcp", address)

	}

	if err != nil {
		return nil, err
	}

	err = Conn.Bind(cfg.BindDN, cfg.Password)

	if err != nil {
		Conn.Close()
		return nil, err
	}

	return Conn, nil
}
