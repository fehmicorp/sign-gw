package smtp

import (
	"bytes"
	"fmt"
	"net/mail"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// ParseMessage extracts RFC822 metadata from raw EML bytes
func ParseMessage(raw []byte) (*config.Email, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	from := msg.Header.Get("From")
	toHeader := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")

	if addr, err := mail.ParseAddress(from); err == nil {
		from = addr.Address
	}

	var toList []string
	if toHeader != "" {
		if addrs, err := mail.ParseAddressList(toHeader); err == nil {
			for _, a := range addrs {
				toList = append(toList, a.Address)
			}
		} else {
			toList = append(toList, toHeader)
		}
	}

	return &config.Email{
		Raw:          raw,
		EnvelopeFrom: from,
		EnvelopeTo:   toList,
		Subject:      subject,
	}, nil
}
