package smtp

import (
	"bytes"
	"io"
	"net/mail"
	"strings"

	_ "github.com/emersion/go-message/charset"
	mmail "github.com/emersion/go-message/mail"
	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// ParseMessage parses a raw RFC822/EML byte slice into an *Email struct.
func ParseMessage(raw []byte) (*config.Email, error) {
	email := &config.Email{
		Raw:     raw,
		Headers: make(map[string][]string),
	}

	// 1. Read top-level RFC822 headers using standard net/mail
	parsedMsg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err == nil {
		for k, v := range parsedMsg.Header {
			email.Headers[k] = v
		}
		email.Subject = parsedMsg.Header.Get("Subject")
	}

	// 2. Traverse MIME parts with go-message
	mr, err := mmail.CreateReader(bytes.NewReader(raw))
	if err != nil {
		return email, err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return email, err
		}

		switch h := part.Header.(type) {
		case *mmail.InlineHeader:
			contentType, _, _ := h.ContentType()
			bodyBytes, err := io.ReadAll(part.Body)
			if err != nil {
				continue
			}

			if strings.HasPrefix(contentType, "text/plain") {
				email.Text = string(bodyBytes)
			} else if strings.HasPrefix(contentType, "text/html") {
				email.HTML = string(bodyBytes)
			} else if cid := h.Get("Content-ID"); cid != "" {
				// Inline image with CID
				cidClean := strings.Trim(cid, "<>")
				email.InlineImages = append(email.InlineImages, config.InlineImage{
					ContentID:   cidClean,
					ContentType: contentType,
					Data:        bodyBytes,
				})
			}

		case *mmail.AttachmentHeader:
			filename, _ := h.Filename()
			contentType, _, _ := h.ContentType()
			bodyBytes, err := io.ReadAll(part.Body)
			if err != nil {
				continue
			}

			email.Attachments = append(email.Attachments, config.Attachment{
				FileName:    filename,
				ContentType: contentType,
				Data:        bodyBytes,
			})
		}
	}

	return email, nil
}
