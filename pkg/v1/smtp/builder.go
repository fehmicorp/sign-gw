package smtp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

func Build(email *config.Email) ([]byte, error) {

	if email == nil {
		return nil, fmt.Errorf("email is nil")
	}

	raw := email.Raw

	if len(raw) == 0 {
		return nil, fmt.Errorf("raw message is empty")
	}

	//----------------------------------------------------------------------
	// Add FEHMI Header
	//----------------------------------------------------------------------

	headerEnd := bytes.Index(raw, []byte("\r\n\r\n"))

	sep := []byte("\r\n\r\n")

	if headerEnd < 0 {

		headerEnd = bytes.Index(raw, []byte("\n\n"))
		sep = []byte("\n\n")

		if headerEnd < 0 {
			return nil, fmt.Errorf("invalid RFC822 message")
		}
	}

	headers := string(raw[:headerEnd])

	if !strings.Contains(
		strings.ToLower(headers),
		"x-fehmi-gateway:",
	) {

		headers += "\r\nX-FEHMI-Gateway: Processed"
	}

	body := raw[headerEnd+len(sep):]

	//----------------------------------------------------------------------
	// Replace HTML Body
	//----------------------------------------------------------------------

	if email.HTML != "" {

		htmlStart := bytes.Index(
			body,
			[]byte("<html"),
		)

		if htmlStart >= 0 {

			htmlEnd := bytes.Index(
				body[htmlStart:],
				[]byte("</html>"),
			)

			if htmlEnd >= 0 {

				htmlEnd += htmlStart + len("</html>")

				newBody := append(
					[]byte{},
					body[:htmlStart]...,
				)

				newBody = append(
					newBody,
					[]byte(email.HTML)...,
				)

				newBody = append(
					newBody,
					body[htmlEnd:]...,
				)

				body = newBody
			}
		}
	}

	//----------------------------------------------------------------------
	// Replace Plain Text (%%SIGN%%)
	//----------------------------------------------------------------------

	if email.Text != "" {

		old := string(body)

		if strings.Contains(old, "%%SIGN%%") {

			old = strings.ReplaceAll(
				old,
				"%%SIGN%%",
				email.Text,
			)

			body = []byte(old)
		}
	}

	//----------------------------------------------------------------------
	// Final RFC822
	//----------------------------------------------------------------------

	var out bytes.Buffer

	out.WriteString(headers)

	out.Write(sep)

	out.Write(body)

	return out.Bytes(), nil
}
