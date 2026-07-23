package smtp

import (
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// ----------------------------------------------------------------------
// ParseBody
// ----------------------------------------------------------------------
// Replaces %%SIGN%% with the rendered signature.
//
// Rules:
//
// 1. Replace %%SIGN%% if present (case-insensitive)
// 2. Otherwise insert before </body>
// 3. Otherwise append to the end
// ----------------------------------------------------------------------

func ParseBody(body, sender string) string {

	if strings.TrimSpace(body) == "" {
		return body
	}

	//------------------------------------------------------------------
	// LDAP User
	//------------------------------------------------------------------

	user, err := GetUser(sender)
	if err != nil {
		return body
	}

	//------------------------------------------------------------------
	// Load Template
	//------------------------------------------------------------------

	tmpl := config.Get(user.Office)

	if tmpl == nil {
		tmpl = config.Get("default")
	}

	if tmpl == nil {
		return body
	}

	//------------------------------------------------------------------
	// Render Signature
	//------------------------------------------------------------------

	signature := Render(tmpl.HTML, user)

	return InjectSignature(body, signature)
}

// ----------------------------------------------------------------------
// InjectSignature
// ----------------------------------------------------------------------

func InjectSignature(body, signature string) string {

	if body == "" {
		return signature
	}

	upper := strings.ToUpper(body)

	//------------------------------------------------------------------
	// Replace %%SIGN%%
	//------------------------------------------------------------------

	if idx := strings.Index(upper, "%%SIGN%%"); idx >= 0 {

		return body[:idx] +
			signature +
			body[idx+len("%%SIGN%%"):]
	}

	//------------------------------------------------------------------
	// Insert before </body>
	//------------------------------------------------------------------

	lower := strings.ToLower(body)

	if idx := strings.LastIndex(lower, "</body>"); idx >= 0 {

		return body[:idx] +
			"\n" +
			signature +
			"\n" +
			body[idx:]
	}

	//------------------------------------------------------------------
	// HTML Message
	//------------------------------------------------------------------

	if strings.Contains(lower, "<html") ||
		strings.Contains(lower, "<div") ||
		strings.Contains(lower, "<table") ||
		strings.Contains(lower, "<p") {

		return body + "<br><br>" + signature
	}

	//------------------------------------------------------------------
	// Plain Text
	//------------------------------------------------------------------

	text := HTMLToText(signature)

	return body + "\r\n\r\n" + text
}

// ----------------------------------------------------------------------
// HTMLToText
// ----------------------------------------------------------------------

func HTMLToText(html string) string {

	r := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n",
		"<p>", "",
		"</div>", "\n",
		"<div>", "",
		"</td>", "\t",
		"<td>", "",
		"</tr>", "\n",
		"<tr>", "",
		"</table>", "\n",
		"<table>", "",
	)

	text := r.Replace(html)

	// Remove remaining tags
	for {

		start := strings.Index(text, "<")
		if start < 0 {
			break
		}

		end := strings.Index(text[start:], ">")
		if end < 0 {
			break
		}

		text = text[:start] + text[start+end+1:]
	}

	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")

	return strings.TrimSpace(text)
}
