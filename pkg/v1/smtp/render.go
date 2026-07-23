package smtp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// Render replaces all %%Field%% placeholders with LDAP values.
func Render(html string, user *config.User) string {

	if user == nil {
		return html
	}

	v := reflect.ValueOf(*user)
	t := reflect.TypeOf(*user)

	for i := 0; i < v.NumField(); i++ {

		field := t.Field(i)

		value := ""

		if !v.Field(i).IsZero() {
			value = fmt.Sprint(v.Field(i).Interface())
		}

		// %%DisplayName%%
		html = strings.ReplaceAll(
			html,
			fmt.Sprintf("%%%%%s%%%%", field.Name),
			value,
		)

		// %%displayname%%
		html = strings.ReplaceAll(
			html,
			fmt.Sprintf("%%%%%s%%%%", strings.ToLower(field.Name)),
			value,
		)

		// %%DISPLAYNAME%%
		html = strings.ReplaceAll(
			html,
			fmt.Sprintf("%%%%%s%%%%", strings.ToUpper(field.Name)),
			value,
		)
	}

	return html
}
