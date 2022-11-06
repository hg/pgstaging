package util

import (
	"regexp"
	"strings"
)

const prefix = "dev_"

var reNonAlnum = regexp.MustCompile(`[^a-z0-9_]`)

func AddPrefix(name string) string {
	if !strings.HasPrefix(name, prefix) {
		name = prefix + name
	}
	return name
}

func NormalizeName(name string) string {
	name = strings.ToLower(name)
	return reNonAlnum.ReplaceAllString(name, "")
}

func IsDevName(name string) bool {
	return strings.HasPrefix(name, prefix) &&
		len(name) > len(prefix)
}
