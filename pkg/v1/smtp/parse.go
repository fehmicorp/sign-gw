package smtp

import (
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// ParseBody replaces %%SIGN%% or %%SIGNATURE%% with the rendered user signature.
func ParseBody(body, sender string) string {
	if strings.TrimSpace(body) == "" {
		return body
	}

	// 1. LDAP User lookup
	user, err := GetUser(sender)
	if err != nil {
		return body
	}

	// 2. Load Template
	tmpl := config.Get(user.Office)
	if tmpl == nil {
		tmpl = config.Get("default")
	}

	if tmpl == nil {
		return body
	}

	// 3. Render Signature
	signature := Render(tmpl.HTML, user)

	return InjectSignature(body, signature)
}

// InjectSignature handles token replacement and fallback insertion logic.
func InjectSignature(body, signature string) string {
	if body == "" {
		return signature
	}

	upper := strings.ToUpper(body)
	isHTML := strings.Contains(strings.ToLower(body), "<html") ||
		strings.Contains(strings.ToLower(body), "<div") ||
		strings.Contains(strings.ToLower(body), "<table") ||
		strings.Contains(strings.ToLower(body), "<p")

	// Prepare signature formats
	htmlSig := signature
	textSig := HTMLToText(signature)

	//------------------------------------------------------------------
	// 1. Replace %%SIGN%% or %%SIGNATURE%%
	//------------------------------------------------------------------
	for _, token := range []string{"%%SIGN%%", "%%SIGNATURE%%"} {
		if idx := strings.Index(upper, token); idx >= 0 {
			// If replacing inside plain text, use plain text signature
			replacement := htmlSig
			if !isHTML {
				replacement = textSig
			}

			return body[:idx] + replacement + body[idx+len(token):]
		}
	}

	//------------------------------------------------------------------
	// 2. Insert before </body> (HTML messages)
	//------------------------------------------------------------------
	lower := strings.ToLower(body)
	if idx := strings.LastIndex(lower, "</body>"); idx >= 0 {
		return body[:idx] + "\n" + htmlSig + "\n" + body[idx:]
	}

	//------------------------------------------------------------------
	// 3. Append to HTML
	//------------------------------------------------------------------
	if isHTML {
		return body + "<br><br>" + htmlSig
	}

	//------------------------------------------------------------------
	// 4. Append to Plain Text
	//------------------------------------------------------------------
	return body + "\r\n\r\n" + textSig
}

// HTMLToText safely strips HTML tags and converts common markup to plain text.
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

	// Safe HTML Tag Removal (prevents infinite loops on malformed HTML)
	var sb strings.Builder
	inTag := false

	for _, ch := range text {
		if ch == '<' {
			inTag = true
			continue
		}
		if ch == '>' {
			inTag = false
			continue
		}
		if !inTag {
			sb.WriteRune(ch)
		}
	}

	res := sb.String()
	res = strings.ReplaceAll(res, "&nbsp;", " ")
	res = strings.ReplaceAll(res, "&amp;", "&")
	res = strings.ReplaceAll(res, "&lt;", "<")
	res = strings.ReplaceAll(res, "&gt;", ">")

	return strings.TrimSpace(res)
}
