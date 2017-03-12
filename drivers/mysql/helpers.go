package mysql

import "strings"

func wrapName(s string) string {
	return "`" + strings.Replace(s, "`", "\\`", -1) + "`"
}

func wrapNames(s []string) string {
	wrapped := make([]string, len(s))
	for key, value := range s {
		wrapped[key] = wrapName(value)
	}
	return strings.Join(wrapped, ",")
}
