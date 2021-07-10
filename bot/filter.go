package bot

import (
	"strings"
	"unicode"

	"github.com/mtibben/confusables"
)

type filter func(string) string

var filters = []filter{
	CleanConfusables,
	LowerLettersOnly,
}

func Filter(s string) string {
	for _, f := range filters {
		s = f(s)
	}

	return s
}

func LowerLettersOnly(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}

		return -1
	}, s)
}

func CleanConfusables(s string) string {
	return confusables.Skeleton(s)
}
