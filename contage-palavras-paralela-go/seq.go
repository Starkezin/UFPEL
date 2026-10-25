package main

import "strings"

// CountSequential counts the word frequencies in the given content sequentially.
func CountSequential(content []byte) map[string]int {
	counts := make(map[string]int)

	// Preprocess text to replace double hyphens and em-dashes/en-dashes with space.
	// This prevents words on either side of a dash from being concatenated.
	text := string(content)
	text = strings.ReplaceAll(text, "--", " ")
	text = strings.ReplaceAll(text, "—", " ")
	text = strings.ReplaceAll(text, "–", " ")

	// Split the text into tokens by any whitespace.
	words := strings.Fields(text)

	// Clean and count each word.
	for _, word := range words {
		cleaned := CleanWord(word)
		if cleaned != "" {
			counts[cleaned]++
		}
	}

	return counts
}
