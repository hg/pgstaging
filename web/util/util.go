package util

import (
	"regexp"
	"strings"
)

const prefix = "dev_"

var reNonAlnum = regexp.MustCompile(`[^a-z0-9_]`)
var rePasswd = regexp.MustCompile(`^[a-zA-Z0-9_]{8,}$`)
var reName = regexp.MustCompile(`^dev_[a-z0-9_]+$`)

func AddPrefix(name string) string {
	if !strings.HasPrefix(name, prefix) {
		name = prefix + name
	}
	return name
}

func IsOkPassword(text string) bool {
	return rePasswd.MatchString(text)
}

func NormalizeName(name string) string {
	name = strings.ToLower(name)
	return reNonAlnum.ReplaceAllString(name, "")
}

func IsDevName(name string) bool {
	return reName.MatchString(name)
}
