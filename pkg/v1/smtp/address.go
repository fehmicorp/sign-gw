package smtp

import (
	"fmt"
	"strings"
)

type Address struct {
	Name    string
	Address string
}

func (a Address) String() string {

	if a.Address == "" {
		return ""
	}

	if a.Name == "" {
		return a.Address
	}

	return fmt.Sprintf("%s <%s>", a.Name, a.Address)
}

func JoinAddresses(list []Address) string {

	out := make([]string, 0, len(list))

	for _, a := range list {

		if s := a.String(); s != "" {
			out = append(out, s)
		}
	}

	return strings.Join(out, ", ")
}
