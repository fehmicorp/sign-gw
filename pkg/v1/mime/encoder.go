package mime

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/quotedprintable"
	"strings"
)

// Encode encodes Body back into RawBody using the original
// Content-Transfer-Encoding.
//
// Body -> RawBody
func Encode(e *Entity) error {

	if e == nil {
		return fmt.Errorf("nil entity")
	}

	if e.IsMultipart() {
		return nil
	}

	encoding := strings.ToLower(strings.TrimSpace(e.Encoding))

	switch encoding {

	case "", "7bit", "8bit", "binary":

		e.RawBody = append([]byte{}, e.Body...)
		return nil

	case "base64":

		return encodeBase64(e)

	case "quoted-printable":

		return encodeQuotedPrintable(e)

	default:

		// Preserve readable body
		e.RawBody = append([]byte{}, e.Body...)
		return nil
	}
}

// EncodeTree recursively encodes every leaf MIME entity.
func EncodeTree(root *Entity) error {

	if root == nil {
		return nil
	}

	if !root.IsMultipart() {

		if err := Encode(root); err != nil {
			return err
		}
	}

	for _, child := range root.Children {

		if err := EncodeTree(child); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------
// Base64
// ----------------------------------------------------------------------

func encodeBase64(e *Entity) error {

	if len(e.Body) == 0 {

		e.RawBody = nil
		return nil
	}

	var buf bytes.Buffer

	enc := base64.NewEncoder(
		base64.StdEncoding,
		&buf,
	)

	if _, err := enc.Write(e.Body); err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return err
	}

	e.RawBody = wrapBase64(buf.Bytes())

	return nil
}

// ----------------------------------------------------------------------
// Quoted Printable
// ----------------------------------------------------------------------

func encodeQuotedPrintable(e *Entity) error {

	var buf bytes.Buffer

	w := quotedprintable.NewWriter(&buf)

	if _, err := w.Write(e.Body); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	e.RawBody = buf.Bytes()

	return nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

// EncodedString returns encoded body.
func EncodedString(e *Entity) string {

	if e == nil {
		return ""
	}

	return string(e.RawBody)
}

// SetEncodedString replaces encoded body directly.
func SetEncodedString(
	e *Entity,
	s string,
) {

	if e == nil {
		return
	}

	e.RawBody = []byte(s)
}

// EncodeChildren encodes only immediate children.
func EncodeChildren(e *Entity) error {

	if e == nil {
		return nil
	}

	for _, child := range e.Children {

		if child.IsMultipart() {

			if err := EncodeChildren(child); err != nil {
				return err
			}

			continue
		}

		if err := Encode(child); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------
// RFC2045 Base64 Line Wrapping
// ----------------------------------------------------------------------

// Base64 MIME bodies should be wrapped at 76 chars.
func wrapBase64(src []byte) []byte {

	if len(src) == 0 {
		return src
	}

	var out bytes.Buffer

	for len(src) > 76 {

		out.Write(src[:76])
		out.WriteString("\r\n")

		src = src[76:]
	}

	out.Write(src)

	return out.Bytes()
}
