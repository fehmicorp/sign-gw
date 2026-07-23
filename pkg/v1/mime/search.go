package mime

import "strings"

// FindFirst returns the first entity matching fn.
func FindFirst(
	root *Entity,
	fn func(*Entity) bool,
) *Entity {

	if root == nil {
		return nil
	}

	if fn(root) {
		return root
	}

	for _, child := range root.Children {

		if e := FindFirst(child, fn); e != nil {
			return e
		}
	}

	return nil
}

// FindAll returns all matching entities.
func FindAll(
	root *Entity,
	fn func(*Entity) bool,
) []*Entity {

	var list []*Entity

	if root == nil {
		return list
	}

	root.Walk(func(e *Entity) {

		if fn(e) {
			list = append(list, e)
		}
	})

	return list
}

// ------------------------------------------------------------
// HTML
// ------------------------------------------------------------

func FindHTML(root *Entity) *Entity {

	return FindFirst(root, func(e *Entity) bool {

		return strings.EqualFold(
			e.MediaType,
			"text/html",
		)
	})
}

// ------------------------------------------------------------
// Plain Text
// ------------------------------------------------------------

func FindText(root *Entity) *Entity {

	return FindFirst(root, func(e *Entity) bool {

		return strings.EqualFold(
			e.MediaType,
			"text/plain",
		)
	})
}

// ------------------------------------------------------------
// Multipart
// ------------------------------------------------------------

func FindMultipart(root *Entity) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		return strings.HasPrefix(
			e.MediaType,
			"multipart/",
		)
	})
}

// ------------------------------------------------------------
// Attachments
// ------------------------------------------------------------

func FindAttachments(root *Entity) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		return e.IsAttachment()
	})
}

// ------------------------------------------------------------
// Inline Images
// ------------------------------------------------------------

func FindInlineImages(root *Entity) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		return e.IsInline() &&
			e.ContentID != ""
	})
}

// ------------------------------------------------------------
// By Media Type
// ------------------------------------------------------------

func FindMediaType(
	root *Entity,
	mediaType string,
) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		return strings.EqualFold(
			e.MediaType,
			mediaType,
		)
	})
}

// ------------------------------------------------------------
// By Content-ID
// ------------------------------------------------------------

func FindContentID(
	root *Entity,
	cid string,
) *Entity {

	return FindFirst(root, func(e *Entity) bool {

		return strings.EqualFold(
			e.ContentID,
			cid,
		)
	})
}

// ------------------------------------------------------------
// By File Name
// ------------------------------------------------------------

func FindFilename(
	root *Entity,
	name string,
) *Entity {

	return FindFirst(root, func(e *Entity) bool {

		return strings.EqualFold(
			e.FileName,
			name,
		)
	})
}

// ------------------------------------------------------------
// By Header
// ------------------------------------------------------------

func FindHeader(
	root *Entity,
	header string,
	value string,
) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		v := e.Header.Get(header)

		if value == "" {
			return v != ""
		}

		return strings.EqualFold(v, value)
	})
}

// ------------------------------------------------------------
// Leaf Nodes
// ------------------------------------------------------------

func FindLeafNodes(root *Entity) []*Entity {

	return FindAll(root, func(e *Entity) bool {

		return len(e.Children) == 0
	})
}
