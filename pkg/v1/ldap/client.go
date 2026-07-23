package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/go-ldap/ldap/v3"
)

func Connect() (*ldap.Conn, error) {
	cfg := config.LdapC
	address := fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)

	var conn *ldap.Conn
	var err error

	if cfg.UseTLS {

		conn, err = ldap.DialTLS(
			"tcp",
			address,
			&tls.Config{
				InsecureSkipVerify: true,
			},
		)

	} else {

		conn, err = ldap.Dial("tcp", address)

	}

	if err != nil {
		return nil, err
	}

	err = conn.Bind(cfg.BindDN, cfg.Password)

	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
