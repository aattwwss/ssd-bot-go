package search

import (
	"regexp"
	"strings"
)

var alphaNumRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func ReplaceSpecialChar(s string, replaceInto string) string {
	return alphaNumRegex.ReplaceAllString(s, replaceInto)
}

// tokenize replace all special characters with space and return a list of words
func Tokenize(s string) []string {
	return removeDuplicate(strings.Split(alphaNumRegex.ReplaceAllString(s, " "), " "))
}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
