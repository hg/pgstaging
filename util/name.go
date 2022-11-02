package util

import (
	"regexp"
	"strings"
)

var reNonAlnum = regexp.MustCompile(`[^a-z0-9_]]`)

func NormalizeName(name string) string {
	name = strings.ToLower(name)
	return reNonAlnum.ReplaceAllString(name, "")
}
