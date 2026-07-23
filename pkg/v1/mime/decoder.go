package mime

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"strings"
)

// Decode decodes the entity body according to the
// Content-Transfer-Encoding header.
//
// RawBody -> Body
//
// Supported:
//
//	base64
//	quoted-printable
//	7bit
//	8bit
//	binary
func Decode(e *Entity) error {

	if e == nil {
		return fmt.Errorf("nil entity")
	}

	if len(e.RawBody) == 0 {
		e.Body = nil
		return nil
	}

	encoding := strings.ToLower(strings.TrimSpace(e.Encoding))

	switch encoding {

	case "", "7bit", "8bit", "binary":

		e.Body = append([]byte{}, e.RawBody...)
		return nil

	case "base64":

		return decodeBase64(e)

	case "quoted-printable":

		return decodeQuotedPrintable(e)

	default:

		// Unknown encoding.
		// Preserve original data.
		e.Body = append([]byte{}, e.RawBody...)
		return nil
	}
}

// DecodeTree recursively decodes an entire MIME tree.
func DecodeTree(root *Entity) error {

	if root == nil {
		return nil
	}

	if !root.IsMultipart() {

		if err := Decode(root); err != nil {
			return err
		}
	}

	for _, child := range root.Children {

		if err := DecodeTree(child); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------
// Base64
// ----------------------------------------------------------------------

func decodeBase64(e *Entity) error {

	src := bytes.TrimSpace(e.RawBody)

	if len(src) == 0 {
		e.Body = nil
		return nil
	}

	decoder := base64.NewDecoder(
		base64.StdEncoding,
		bytes.NewReader(src),
	)

	data, err := io.ReadAll(decoder)
	if err != nil {
		return err
	}

	e.Body = data

	return nil
}

// ----------------------------------------------------------------------
// Quoted Printable
// ----------------------------------------------------------------------

func decodeQuotedPrintable(e *Entity) error {

	reader := quotedprintable.NewReader(
		bytes.NewReader(e.RawBody),
	)

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	e.Body = data

	return nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

// DecodedString returns the decoded body as UTF-8 text.
func DecodedString(e *Entity) string {

	if e == nil {
		return ""
	}

	return string(e.Body)
}

// SetDecodedString replaces the decoded body.
// Encoding is preserved until Encode() is called.
func SetDecodedString(
	e *Entity,
	text string,
) {

	if e == nil {
		return
	}

	e.Body = []byte(text)
}

// DecodeChildren decodes only direct children.
func DecodeChildren(e *Entity) error {

	if e == nil {
		return nil
	}

	for _, child := range e.Children {

		if child.IsMultipart() {

			if err := DecodeChildren(child); err != nil {
				return err
			}

			continue
		}

		if err := Decode(child); err != nil {
			return err
		}
	}

	return nil
}

// IsDecoded reports whether the entity has decoded data.
func IsDecoded(e *Entity) bool {

	if e == nil {
		return false
	}

	return len(e.Body) > 0
}
