package smtp

import (
	"fmt"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// HTMLSignature fetches user details via LDAP and renders the HTML signature template
func HTMLSignature(email *config.Email) (string, error) {
	if email == nil {
		return "", fmt.Errorf("email is nil")
	}

	sender := email.EnvelopeFrom
	if sender == "" {
		return "", fmt.Errorf("envelope sender is empty")
	}

	// 1. Fetch LDAP
	// user details
	user, err := GetUser(sender)
	if err != nil {
		return "", fmt.Errorf("failed to get user for %s: %w", sender, err)
	}
	if user == nil {
		return "", fmt.Errorf("user object returned nil for %s", sender)
	}

	// 2. Load Office or Default Template
	tmpl := config.Get(user.Office)
	if tmpl == nil {
		tmpl = config.Get("default")
	}

	if tmpl == nil {
		return "", fmt.Errorf("no signature template found for office %s or default", user.Office)
	}

	// 3. Render HTML Signature using template and user details
	signature := Render(tmpl.HTML, user)

	return signature, nil
}
