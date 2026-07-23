package smtp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
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
	// Header & Body Separation
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

	if !strings.Contains(strings.ToLower(headers), "x-fehmi-gateway:") {
		headers += "\r\nX-FEHMI-Gateway: Processed"
	}

	body := raw[headerEnd+len(sep):]

	//----------------------------------------------------------------------
	// Parse Body via SMTP Package (LDAP + Template Rendering)
	//----------------------------------------------------------------------
	parsedBodyStr := ParseBody(string(body), email.EnvelopeFrom)
	body = []byte(parsedBodyStr)

	//----------------------------------------------------------------------
	// Replace HTML & Text Body (Fallback if ParseBody didn't handle tokens)
	//----------------------------------------------------------------------
	if email.HTML != "" {
		body = bodyHtml(body, email)
	}

	if email.Text != "" {
		body = bodyText(body, email)
	}

	//----------------------------------------------------------------------
	// Final RFC822 Assembly
	//----------------------------------------------------------------------
	var out bytes.Buffer
	out.WriteString(headers)
	out.Write(sep)
	out.Write(body)

	return out.Bytes(), nil
}

func bodyText(body []byte, email *config.Email) []byte {
	logger.Info("Text body processing")
	old := string(body)

	if strings.Contains(old, "%%SIGN%%") || strings.Contains(old, "%%SIGNATURE%%") {
		old = strings.ReplaceAll(old, "%%SIGN%%", email.Text)
		old = strings.ReplaceAll(old, "%%SIGNATURE%%", email.Text)
		return []byte(old)
	}

	return body
}

func bodyHtml(body []byte, email *config.Email) []byte {
	logger.Info("Html body processing")

	bodyStr := string(body)
	if strings.Contains(bodyStr, "%%SIGN%%") || strings.Contains(bodyStr, "%%SIGNATURE%%") {
		bodyStr = strings.ReplaceAll(bodyStr, "%%SIGN%%", email.HTML)
		bodyStr = strings.ReplaceAll(bodyStr, "%%SIGNATURE%%", email.HTML)
		return []byte(bodyStr)
	}

	htmlStart := bytes.Index(body, []byte("<html"))
	if htmlStart >= 0 {
		htmlEnd := bytes.Index(body[htmlStart:], []byte("</html>"))
		if htmlEnd >= 0 {
			htmlEnd += htmlStart + len("</html>")

			newBody := make([]byte, 0, len(body)+len(email.HTML))
			newBody = append(newBody, body[:htmlStart]...)
			newBody = append(newBody, []byte(email.HTML)...)
			newBody = append(newBody, body[htmlEnd:]...)

			return newBody
		}
	}

	return body
}
