package smtp

import (
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/ldap"
)

func GetUser(username string) (*config.User, error) {

	if i := strings.Index(username, "@"); i > 0 {
		username = username[:i]
	}

	conn, err := ldap.Connect()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	config.Conn = conn

	return ldap.GetUser(username)
}
