package textvector

import (
	"strings"
	"unicode"
)

// it accepts an input which returns all the words in that input as an array
func Tokenize(input string) []string {
	lowercase := strings.ToLower(input)

	words := strings.FieldsFunc(lowercase, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	return words
}
