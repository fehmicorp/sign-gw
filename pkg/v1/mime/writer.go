package mime

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/textproto"
	"sort"
	"strings"
)

// Write serializes an Entity tree into RFC822/MIME bytes.
//
// It automatically:
//
//   - encodes modified bodies
//   - preserves multipart boundaries
//   - writes all headers
//   - recursively writes child entities
func Write(root *Entity) ([]byte, error) {

	if root == nil {
		return nil, fmt.Errorf("nil entity")
	}

	if err := EncodeTree(root); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	w := bufio.NewWriter(&buf)

	if err := writeEntity(
		w,
		root,
	); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func writeEntity(
	w *bufio.Writer,
	e *Entity,
) error {

	if e == nil {
		return nil
	}

	if err := writeHeaders(
		w,
		e.Header,
	); err != nil {
		return err
	}

	if _, err := w.WriteString("\r\n"); err != nil {
		return err
	}

	if e.IsMultipart() {

		return writeMultipart(
			w,
			e,
		)
	}

	return writeBody(
		w,
		e,
	)
}

func writeHeaders(
	w io.Writer,
	h textproto.MIMEHeader,
) error {

	if h == nil {
		return nil
	}

	keys := make([]string, 0, len(h))

	for k := range h {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {

		values := h[key]

		for _, value := range values {

			if _, err := fmt.Fprintf(
				w,
				"%s: %s\r\n",
				key,
				value,
			); err != nil {

				return err
			}
		}
	}

	return nil
}

func UpdateHeaders(e *Entity) {

	if e.Header == nil {

		e.Header = make(textproto.MIMEHeader)
	}

	if e.ContentType != "" {

		e.Header.Set(
			"Content-Type",
			e.ContentType,
		)
	}

	if e.Encoding != "" {

		e.Header.Set(
			"Content-Transfer-Encoding",
			e.Encoding,
		)
	}

	if e.Disposition != "" {

		e.Header.Set(
			"Content-Disposition",
			e.Disposition,
		)
	}

	if e.ContentID != "" {

		e.Header.Set(
			"Content-ID",
			"<"+e.ContentID+">",
		)
	}
}

func writeCRLF(
	w io.Writer,
) error {

	_, err := io.WriteString(
		w,
		"\r\n",
	)

	return err
}

func writeBoundary(
	w io.Writer,
	boundary string,
) error {

	_, err := io.WriteString(
		w,
		"--"+boundary+"\r\n",
	)

	return err
}

func writeClosingBoundary(
	w io.Writer,
	boundary string,
) error {

	_, err := io.WriteString(
		w,
		"--"+boundary+"--\r\n",
	)

	return err
}

func normalizeCRLF(
	s string,
) string {

	s = strings.ReplaceAll(
		s,
		"\r\n",
		"\n",
	)

	s = strings.ReplaceAll(
		s,
		"\r",
		"\n",
	)

	s = strings.ReplaceAll(
		s,
		"\n",
		"\r\n",
	)

	return s
}

// ------------------------------------------------------------
// Multipart Writer
// ------------------------------------------------------------

func writeMultipart(
	w *bufio.Writer,
	e *Entity,
) error {

	if e == nil {
		return nil
	}

	if e.Boundary == "" {
		return fmt.Errorf("multipart entity missing boundary")
	}

	for _, child := range e.Children {

		//----------------------------------------------------
		// Start Boundary
		//----------------------------------------------------

		if err := writeBoundary(
			w,
			e.Boundary,
		); err != nil {
			return err
		}

		//----------------------------------------------------
		// Child Entity
		//----------------------------------------------------

		if err := writeEntity(
			w,
			child,
		); err != nil {
			return err
		}

		//----------------------------------------------------
		// Blank line between MIME parts
		//----------------------------------------------------

		if err := writeCRLF(w); err != nil {
			return err
		}
	}

	//--------------------------------------------------------
	// Closing Boundary
	//--------------------------------------------------------

	return writeClosingBoundary(
		w,
		e.Boundary,
	)
}

func GenerateBoundary() string {

	var b [16]byte

	_, _ = rand.Read(b[:])

	return "----=_FEHMI_" + hex.EncodeToString(b[:])
}

func NewMultipart(
	contentType string,
) *Entity {

	boundary := GenerateBoundary()

	h := make(textproto.MIMEHeader)

	h.Set(
		"Content-Type",
		fmt.Sprintf(
			"%s; boundary=\"%s\"",
			contentType,
			boundary,
		),
	)

	return &Entity{
		Header:      h,
		ContentType: h.Get("Content-Type"),
		MediaType:   contentType,
		Boundary:    boundary,
	}
}

func AddPart(
	parent *Entity,
	child *Entity,
) {

	if parent == nil || child == nil {
		return
	}

	child.Parent = parent

	parent.Children = append(
		parent.Children,
		child,
	)
}

func RemovePart(
	parent *Entity,
	index int,
) {

	if parent == nil {
		return
	}

	if index < 0 ||
		index >= len(parent.Children) {
		return
	}

	parent.Children = append(
		parent.Children[:index],
		parent.Children[index+1:]...,
	)
}

func CloneMultipart(
	root *Entity,
) *Entity {

	if root == nil {
		return nil
	}

	return root.Clone()
}
func CountParts(
	root *Entity,
) int {

	if root == nil {
		return 0
	}

	count := 0

	root.Walk(func(e *Entity) {

		count++
	})

	return count
}
func CountAttachments(
	root *Entity,
) int {

	total := 0

	root.Walk(func(e *Entity) {

		if e.IsAttachment() {
			total++
		}
	})

	return total
}
func CountInlineImages(
	root *Entity,
) int {

	total := 0

	root.Walk(func(e *Entity) {

		if e.IsInline() &&
			e.ContentID != "" {

			total++
		}
	})

	return total
}

func DumpTree(
	root *Entity,
) string {

	var buf bytes.Buffer

	var walk func(*Entity, int)

	walk = func(
		e *Entity,
		level int,
	) {

		if e == nil {
			return
		}

		buf.WriteString(
			strings.Repeat(
				"  ",
				level,
			),
		)

		buf.WriteString(e.MediaType)

		if e.FileName != "" {

			buf.WriteString(
				" (" +
					e.FileName +
					")",
			)
		}

		buf.WriteString("\n")

		for _, c := range e.Children {
			walk(c, level+1)
		}
	}

	walk(root, 0)

	return buf.String()
}

// ------------------------------------------------------------
// Leaf Entity Writer
// ------------------------------------------------------------

func writeBody(
	w *bufio.Writer,
	e *Entity,
) error {

	if e == nil {
		return nil
	}

	//----------------------------------------------------------
	// Multipart entities never have a body here.
	//----------------------------------------------------------

	if e.IsMultipart() {
		return nil
	}

	//----------------------------------------------------------
	// Body
	//----------------------------------------------------------

	if len(e.RawBody) == 0 {
		return nil
	}

	if _, err := w.Write(e.RawBody); err != nil {
		return err
	}

	return nil
}

func WriteHeaders(
	e *Entity,
) ([]byte, error) {

	if e == nil {
		return nil, fmt.Errorf("nil entity")
	}

	var buf bytes.Buffer

	if err := writeHeaders(
		&buf,
		e.Header,
	); err != nil {
		return nil, err
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func SyncHeaders(
	e *Entity,
) {

	if e == nil {
		return
	}

	if e.Header == nil {
		e.Header = make(textproto.MIMEHeader)
	}

	if e.ContentType != "" {
		e.Header.Set(
			"Content-Type",
			e.ContentType,
		)
	}

	if e.Encoding != "" {
		e.Header.Set(
			"Content-Transfer-Encoding",
			e.Encoding,
		)
	}

	if e.Disposition != "" {
		e.Header.Set(
			"Content-Disposition",
			e.Disposition,
		)
	}

	if e.ContentID != "" {
		e.Header.Set(
			"Content-ID",
			"<"+e.ContentID+">",
		)
	}
}

func SyncTree(
	root *Entity,
) {

	if root == nil {
		return
	}

	root.Walk(func(e *Entity) {

		SyncHeaders(e)
	})
}

func Validate(
	root *Entity,
) error {

	if root == nil {
		return fmt.Errorf("nil MIME tree")
	}

	root.Walk(func(e *Entity) {

		if e.IsMultipart() &&
			e.Boundary == "" {

			e.Boundary = GenerateBoundary()
		}
	})

	return nil
}
