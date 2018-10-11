package tfidf

import "bytes"

// ReplaceSmartQuotes replaces non-ASCII smart quotes that come from
// applications like MS Word.
func ReplaceSmartQuotes(text string) string {
	buf := bytes.Buffer{}
	for _, c := range text {
		switch c {
		case '\u2013', '\u2014', '\u2015':
			buf.WriteRune('-')
		case '\u2017':
			buf.WriteRune('_')
		case '\u2018', '\u2019', '\u201b', '\u2032', '\u2035':
			buf.WriteRune('\'')
		case '\u201a':
			buf.WriteRune(',')
		case '\u201c', '\u201d', '\u201e', '\u201f', '\u2033', '\u2036':
			buf.WriteRune('"')
		case '\u2026':
			buf.WriteString("...")
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
