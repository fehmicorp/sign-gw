package smtp

import (
	"fmt"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

func FormatAddress(a config.Address) string {

	if a.Address == "" {
		return ""
	}

	if a.Name == "" {
		return a.Address
	}

	return fmt.Sprintf(`"%s" <%s>`, a.Name, a.Address)
}
