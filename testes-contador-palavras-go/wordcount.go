package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// CleanAndSplit splits a text string into lowercase alphanumeric tokens.
// It ignores tokens with fewer than 3 characters (runes).
// It uses unicode definitions of letters and numbers to support accented and multi-byte characters.
func CleanAndSplit(text string) []string {
	// FieldsFunc splits the string at each rune where the function returns true.
	// We split on any character that is not a letter or a number (punctuation, space, etc.)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var filtered []string
	for _, w := range words {
		wLower := strings.ToLower(w)
		if utf8.RuneCountInString(wLower) >= 3 {
			filtered = append(filtered, wLower)
		}
	}

	return filtered
}

// CountWords counts the frequency of words in a text.
// It converts words to lowercase, removes simple punctuation, and ignores words with fewer than 3 characters.
func CountWords(text string) map[string]int {
	counts := make(map[string]int)
	words := CleanAndSplit(text)
	for _, w := range words {
		counts[w]++
	}
	return counts
}
