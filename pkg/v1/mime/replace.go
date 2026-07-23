package mime

import (
	"bytes"
	"strings"
)

// ReplaceOptions controls how signatures are injected.
type ReplaceOptions struct {

	// Placeholder to replace.
	Placeholder string

	// HTML signature.
	HTMLSignature string

	// Plain text signature.
	TextSignature string

	// Insert signature if placeholder doesn't exist.
	InsertIfMissing bool

	// HTML insertion position.
	InsertHTMLBeforeBodyEnd bool

	// Plain text insertion position.
	AppendText bool
}

// ----------------------------------------------------------------------
// Replace entire MIME tree
// ----------------------------------------------------------------------

func Replace(root *Entity, opt ReplaceOptions) {

	if root == nil {
		return
	}

	if opt.Placeholder == "" {
		opt.Placeholder = "%%SIGN%%"
	}

	root.Walk(func(e *Entity) {

		switch {

		case e.IsHTML():

			replaceHTML(e, opt)

		case e.IsText():

			replaceText(e, opt)
		}
	})
}

// ----------------------------------------------------------------------
// HTML
// ----------------------------------------------------------------------

func replaceHTML(
	e *Entity,
	opt ReplaceOptions,
) {

	if e == nil {
		return
	}

	html := string(e.Body)

	//--------------------------------------------------------
	// Placeholder
	//--------------------------------------------------------

	if strings.Contains(
		html,
		opt.Placeholder,
	) {

		html = strings.ReplaceAll(
			html,
			opt.Placeholder,
			opt.HTMLSignature,
		)

		e.Body = []byte(html)

		return
	}

	if !opt.InsertIfMissing {
		return
	}

	//--------------------------------------------------------
	// Before </body>
	//--------------------------------------------------------

	if opt.InsertHTMLBeforeBodyEnd {

		lower := strings.ToLower(html)

		idx := strings.LastIndex(
			lower,
			"</body>",
		)

		if idx >= 0 {

			html =
				html[:idx] +
					opt.HTMLSignature +
					html[idx:]

			e.Body = []byte(html)

			return
		}
	}

	//--------------------------------------------------------
	// Before </html>
	//--------------------------------------------------------

	lower := strings.ToLower(html)

	idx := strings.LastIndex(
		lower,
		"</html>",
	)

	if idx >= 0 {

		html =
			html[:idx] +
				opt.HTMLSignature +
				html[idx:]

		e.Body = []byte(html)

		return
	}

	//--------------------------------------------------------
	// Append
	//--------------------------------------------------------

	html += opt.HTMLSignature

	e.Body = []byte(html)
}

// ----------------------------------------------------------------------
// Plain Text
// ----------------------------------------------------------------------

func replaceText(
	e *Entity,
	opt ReplaceOptions,
) {

	if e == nil {
		return
	}

	text := string(e.Body)

	if strings.Contains(
		text,
		opt.Placeholder,
	) {

		text = strings.ReplaceAll(
			text,
			opt.Placeholder,
			opt.TextSignature,
		)

		e.Body = []byte(text)

		return
	}

	if !opt.InsertIfMissing {
		return
	}

	if opt.AppendText {

		if !strings.HasSuffix(
			text,
			"\n",
		) {

			text += "\r\n"
		}

		text += opt.TextSignature

		e.Body = []byte(text)
	}
}

// ----------------------------------------------------------------------
// Convenience Helpers
// ----------------------------------------------------------------------

func ReplaceHTML(
	root *Entity,
	signature string,
) {

	Replace(root, ReplaceOptions{
		Placeholder:             "%%SIGN%%",
		HTMLSignature:           signature,
		InsertIfMissing:         false,
		InsertHTMLBeforeBodyEnd: true,
	})
}

func ReplaceText(
	root *Entity,
	signature string,
) {

	Replace(root, ReplaceOptions{
		Placeholder:     "%%SIGN%%",
		TextSignature:   signature,
		InsertIfMissing: false,
		AppendText:      true,
	})
}

// ----------------------------------------------------------------------
// Raw Replace
// ----------------------------------------------------------------------

func ReplaceBytes(
	body []byte,
	placeholder string,
	value string,
) []byte {

	return bytes.ReplaceAll(
		body,
		[]byte(placeholder),
		[]byte(value),
	)
}
