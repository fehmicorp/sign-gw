package smtp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// SaveEML stores the original RFC822 message.
func SaveEML(email *config.Email) error {

	if !config.SmtpC.SaveRawEML {
		return nil
	}

	dir := filepath.Join(
		"data",
		"eml",
		"orignal",
		time.Now().Format("2006-01-02"),
	)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	subject := sanitizeFilename(email.Subject)

	if subject == "" {
		subject = "No Subject"
	}

	file := fmt.Sprintf(
		"%s_%s.eml",
		time.Now().Format("150405.000"),
		subject,
	)

	path := filepath.Join(dir, file)

	return os.WriteFile(path, email.Raw, 0644)
}

func SaveEditedEML(email *config.Email) error {

	if !config.SmtpC.SaveRawEML {
		return nil
	}

	dir := filepath.Join(
		"data",
		"eml",
		time.Now().Format("2006-01-02"),
	)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	subject := sanitizeFilename(email.Subject)

	if subject == "" {
		subject = "No Subject"
	}

	file := fmt.Sprintf(
		"%s_%s.eml",
		time.Now().Format("150405.000"),
		subject,
	)

	path := filepath.Join(dir, file)

	return os.WriteFile(path, email.Raw, 0644)
}

func sanitizeFilename(name string) string {

	name = strings.TrimSpace(name)

	r := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	name = r.Replace(name)

	if len(name) > 100 {
		name = name[:100]
	}

	return name
}
