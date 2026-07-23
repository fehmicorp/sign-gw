package mime

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/textproto"
	"strings"
)

// Parse parses an RFC822/MIME message into an Entity tree.
func Parse(raw []byte) (*Entity, error) {

	if len(raw) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	root, err := parseEntity(raw)
	if err != nil {
		return nil, err
	}

	if err := DecodeTree(root); err != nil {
		return nil, err
	}

	return root, nil
}

// ------------------------------------------------------------
// Parse Single MIME Entity
// ------------------------------------------------------------

func parseEntity(data []byte) (*Entity, error) {

	reader := textproto.NewReader(
		bufio.NewReader(bytes.NewReader(data)),
	)

	header, err := reader.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	entity := &Entity{
		Header: header,
	}

	entity.loadHeaders()

	body, err := io.ReadAll(reader.R)
	if err != nil {
		return nil, err
	}

	//--------------------------------------------------------
	// Multipart
	//--------------------------------------------------------

	if strings.HasPrefix(
		strings.ToLower(entity.ContentType),
		"multipart/",
	) {

		children, err := parseMultipart(
			body,
			entity.Boundary,
			entity,
		)

		if err != nil {
			return nil, err
		}

		entity.Children = children

		return entity, nil
	}

	//--------------------------------------------------------
	// Leaf Entity
	//--------------------------------------------------------

	entity.RawBody = body

	return entity, nil
}

// loadHeaders extracts useful MIME metadata.
func (e *Entity) loadHeaders() {

	ct := e.Header.Get("Content-Type")

	if ct == "" {
		ct = "text/plain"
	}

	e.ContentType = ct

	mediaType, params, _ := mime.ParseMediaType(ct)

	e.MediaType = strings.ToLower(mediaType)

	if v, ok := params["charset"]; ok {
		e.Charset = v
	}

	if v, ok := params["boundary"]; ok {
		e.Boundary = v
	}

	if v, ok := params["name"]; ok {
		e.Name = v
	}

	e.Encoding = strings.TrimSpace(
		e.Header.Get("Content-Transfer-Encoding"),
	)

	cd := e.Header.Get("Content-Disposition")

	if cd != "" {

		e.Disposition = cd

		_, params, _ := mime.ParseMediaType(cd)

		if v, ok := params["filename"]; ok {
			e.FileName = v
		}
	}

	e.ContentID = strings.Trim(
		e.Header.Get("Content-ID"),
		"<>",
	)
}

func ParseBytes(raw []byte) (*Entity, error) {
	return Parse(raw)
}

// ------------------------------------------------------------
// Multipart Parser
// ------------------------------------------------------------

func parseMultipart(
	body []byte,
	boundary string,
	parent *Entity,
) ([]*Entity, error) {

	if boundary == "" {
		return nil, fmt.Errorf("multipart without boundary")
	}

	var children []*Entity

	startBoundary := "--" + boundary
	endBoundary := "--" + boundary + "--"

	lines := splitLines(body)

	var part bytes.Buffer
	inPart := false

	flush := func() error {

		if part.Len() == 0 {
			return nil
		}

		child, err := parseEntity(part.Bytes())
		if err != nil {
			return err
		}

		child.Parent = parent
		children = append(children, child)

		part.Reset()

		return nil
	}

	for _, line := range lines {

		switch {

		//------------------------------------------
		// New Part
		//------------------------------------------

		case line == startBoundary:

			if inPart {

				if err := flush(); err != nil {
					return nil, err
				}
			}

			inPart = true

		//------------------------------------------
		// Last Boundary
		//------------------------------------------

		case line == endBoundary:

			if err := flush(); err != nil {
				return nil, err
			}

			inPart = false

			break

		//------------------------------------------
		// Body
		//------------------------------------------

		default:

			if inPart {

				part.WriteString(line)
				part.WriteString("\r\n")
			}
		}
	}

	if part.Len() > 0 {

		if err := flush(); err != nil {
			return nil, err
		}
	}

	return children, nil
}

func splitLines(data []byte) []string {

	s := string(data)

	s = strings.ReplaceAll(
		s,
		"\r\n",
		"\n",
	)

	return strings.Split(
		s,
		"\n",
	)
}
func isMultipart(contentType string) bool {

	return strings.HasPrefix(
		strings.ToLower(contentType),
		"multipart/",
	)
}

func hasBoundary(e *Entity) bool {

	return e != nil &&
		e.Boundary != ""
}
