package mime

import (
	"net/textproto"
)

// Entity represents a single MIME entity.
//
// A MIME message is a tree:
//
//	multipart/mixed
//	├── multipart/alternative
//	│   ├── text/plain
//	│   └── text/html
//	└── application/pdf
//
// Every node is an Entity.
type Entity struct {

	// Original MIME headers.
	Header textproto.MIMEHeader

	// MIME information.
	ContentType string
	MediaType   string
	Charset     string
	Boundary    string
	Name        string

	// Content disposition.
	Disposition string
	FileName    string

	// Transfer encoding.
	Encoding string

	// Inline Content-ID.
	ContentID string

	// Raw encoded body.
	RawBody []byte

	// Decoded body.
	Body []byte

	// Child MIME entities.
	Children []*Entity

	// Parent entity.
	Parent *Entity
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

// IsMultipart reports whether this entity is multipart.
func (e *Entity) IsMultipart() bool {
	return len(e.Children) > 0
}

// IsHTML reports text/html.
func (e *Entity) IsHTML() bool {
	return e.MediaType == "text/html"
}

// IsText reports text/plain.
func (e *Entity) IsText() bool {
	return e.MediaType == "text/plain"
}

// IsAttachment reports attachment.
func (e *Entity) IsAttachment() bool {
	return e.Disposition == "attachment"
}

// IsInline reports inline content.
func (e *Entity) IsInline() bool {
	return e.Disposition == "inline"
}

// Walk performs depth-first traversal.
func (e *Entity) Walk(fn func(*Entity)) {

	if e == nil {
		return
	}

	fn(e)

	for _, child := range e.Children {
		child.Walk(fn)
	}
}

// AddChild appends a child entity.
func (e *Entity) AddChild(child *Entity) {

	if child == nil {
		return
	}

	child.Parent = e
	e.Children = append(e.Children, child)
}

// Clone performs a deep copy.
func (e *Entity) Clone() *Entity {

	if e == nil {
		return nil
	}

	n := &Entity{
		ContentType: e.ContentType,
		MediaType:   e.MediaType,
		Charset:     e.Charset,
		Boundary:    e.Boundary,
		Name:        e.Name,
		Disposition: e.Disposition,
		FileName:    e.FileName,
		Encoding:    e.Encoding,
		ContentID:   e.ContentID,
	}

	if e.RawBody != nil {
		n.RawBody = append([]byte{}, e.RawBody...)
	}

	if e.Body != nil {
		n.Body = append([]byte{}, e.Body...)
	}

	if e.Header != nil {

		n.Header = make(textproto.MIMEHeader)

		for k, v := range e.Header {

			values := make([]string, len(v))
			copy(values, v)

			n.Header[k] = values
		}
	}

	for _, c := range e.Children {
		n.AddChild(c.Clone())
	}

	return n
}
