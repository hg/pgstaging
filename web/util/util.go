package util

import (
	"regexp"
	"strings"
)

const prefix = "dev_"

var reNonAlnum = regexp.MustCompile(`[^a-z0-9_]`)
var reAlnum = regexp.MustCompile(`^[a-zA-Z0-9_]{8,}$`)

func AddPrefix(name string) string {
	if !strings.HasPrefix(name, prefix) {
		name = prefix + name
	}
	return name
}

func IsOkPassword(text string) bool {
	return reAlnum.MatchString(text)
}

func NormalizeName(name string) string {
	name = strings.ToLower(name)
	return reNonAlnum.ReplaceAllString(name, "")
}

func IsDevName(name string) bool {
	return strings.HasPrefix(name, prefix) &&
		len(name) > len(prefix)
}
