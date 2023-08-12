package helpers

import (
	"regexp"
	"strings"
)

func CreateSlug(name string) string {
	String := strings.ToLower(name)

	regexChar := regexp.MustCompile("[^a-z0-9]+")
	String = regexChar.ReplaceAllString(String, "-")
	String = strings.Trim(String, "-")

	return String
}
