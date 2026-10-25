package main

import (
	"sort"
	"strings"
)

// WordCount holds a word and its corresponding frequency count.
type WordCount struct {
	Word  string
	Count int
}

// CleanWord converts a string to lowercase and removes simple punctuation.
// If the resulting word is shorter than 3 characters, it returns an empty string.
func CleanWord(w string) string {
	w = strings.ToLower(w)
	var sb strings.Builder
	for _, r := range w {
		if isSimplePunctuation(r) {
			continue
		}
		sb.WriteRune(r)
	}
	cleaned := sb.String()
	if len(cleaned) < 3 {
		return ""
	}
	return cleaned
}

// isSimplePunctuation returns true if the rune is a simple punctuation mark or symbol.
func isSimplePunctuation(r rune) bool {
	switch r {
	case '.', ',', ';', ':', '!', '?', '(', ')', '[', ']', '{', '}', '"', '\'', '`', '_', '-', '*', '/', '\\', '@', '#', '$', '%', '^', '&', '+', '=', '<', '>', '|', '~', '“', '”', '‘', '’', '—', '–':
		return true
	}
	return false
}

// CompareMaps compares two frequency maps. It returns true if they are identical,
// and false otherwise, detailing differences if any could be logged.
func CompareMaps(m1, m2 map[string]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// GetTop20 returns the 20 most frequent words from the frequency map.
// The sorting is deterministic: primary by frequency descending, secondary by word alphabetically.
func GetTop20(m map[string]int) []WordCount {
	pairs := make([]WordCount, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, WordCount{Word: k, Count: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].Word < pairs[j].Word
		}
		return pairs[i].Count > pairs[j].Count
	})

	if len(pairs) > 20 {
		return pairs[:20]
	}
	return pairs
}
