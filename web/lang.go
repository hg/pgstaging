package web

import (
	"regexp"
	"strings"
)

var reLangPrefix = regexp.MustCompile(`^(\w+)`)

func toIsoCodes(value string) (out []string) {
	seen := make(map[string]bool)

	for _, tag := range strings.Split(value, ",") {
		tag = strings.ToLower(strings.TrimSpace(tag))
		lang := reLangPrefix.FindString(tag)

		if lang != "" && !seen[lang] {
			out = append(out, lang)
			seen[lang] = true
		}
	}

	return
}
