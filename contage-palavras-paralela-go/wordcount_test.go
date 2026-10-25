package main

import (
	"testing"
)

func TestCleanWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Scrooge", "scrooge"},
		{"Scrooge's", "scrooges"},
		{"hello!", "hello"},
		{"it", ""}, // under 3 chars
		{"a", ""},  // under 3 chars
		{"", ""},
		{",.;!?", ""},
		{"coffin-nail", "coffinnail"},
		{"didn't", "didnt"},
		{"Self-contained", "selfcontained"},
	}

	for _, tt := range tests {
		got := CleanWord(tt.input)
		if got != tt.expected {
			t.Errorf("CleanWord(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestCompareMaps(t *testing.T) {
	tests := []struct {
		name     string
		m1       map[string]int
		m2       map[string]int
		expected bool
	}{
		{
			name:     "identical maps",
			m1:       map[string]int{"hello": 2, "world": 1},
			m2:       map[string]int{"hello": 2, "world": 1},
			expected: true,
		},
		{
			name:     "different values",
			m1:       map[string]int{"hello": 2, "world": 1},
			m2:       map[string]int{"hello": 2, "world": 2},
			expected: false,
		},
		{
			name:     "different keys",
			m1:       map[string]int{"hello": 2, "world": 1},
			m2:       map[string]int{"hello": 2, "there": 1},
			expected: false,
		},
		{
			name:     "different lengths",
			m1:       map[string]int{"hello": 2, "world": 1},
			m2:       map[string]int{"hello": 2},
			expected: false,
		},
		{
			name:     "both empty",
			m1:       map[string]int{},
			m2:       map[string]int{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareMaps(tt.m1, tt.m2)
			if got != tt.expected {
				t.Errorf("CompareMaps() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestGetTop20(t *testing.T) {
	// Map with 25 distinct words with varying frequencies
	input := map[string]int{
		"apple": 5, "banana": 5, "cherry": 10, "date": 1, "elderberry": 8,
		"fig": 2, "grape": 2, "honeydew": 12, "kiwi": 3, "lemon": 4,
		"mango": 6, "nectarine": 7, "orange": 9, "papaya": 11, "quince": 15,
		"raspberry": 14, "strawberry": 13, "tangerine": 4, "ugli": 1, "victoria": 3,
		"watermelon": 20, "xigua": 2, "yam": 1, "zucchini": 2, "extra": 1,
	}

	top20 := GetTop20(input)
	if len(top20) != 20 {
		t.Fatalf("GetTop20() returned %d elements, want 20", len(top20))
	}

	// Verify sorting (descending count, alphabetical fallback for ties)
	// First should be watermelon (20)
	if top20[0].Word != "watermelon" || top20[0].Count != 20 {
		t.Errorf("First element is %+v, want watermelon: 20", top20[0])
	}

	// Compare apple (5) and banana (5).
	// Since counts are equal, apple should precede banana alphabetically.
	foundApple, foundBanana := -1, -1
	for idx, wc := range top20 {
		if wc.Word == "apple" {
			foundApple = idx
		}
		if wc.Word == "banana" {
			foundBanana = idx
		}
	}
	if foundApple != -1 && foundBanana != -1 && foundApple > foundBanana {
		t.Errorf("apple should precede banana in ties, apple idx: %d, banana idx: %d", foundApple, foundBanana)
	}
}

func TestEquivalence(t *testing.T) {
	text := []byte(`
		Marley was dead: to begin with. There is no doubt whatever about that.
		The register of his burial was signed by the clergyman, the clerk, the
		undertaker, and the chief mourner. Scrooge signed it: and Scrooge's name
		was good upon 'Change, for anything he chose to put his hand to.
		Old Marley was as dead as a door-nail.
		Mind! I don't mean to say that I know, of my own knowledge, what there
		is particularly dead about a door-nail. I might have been inclined,
		myself, to regard a coffin-nail as the deadest piece of ironmongery in
		the trade. But the wisdom of our ancestors is in the simile; and my
		unhallowed hands shall not disturb it, or the Country's done for.
		You will therefore permit me to repeat, emphatically, that Marley was
		as dead as a door-nail--this is a great example.
	`)

	seqCounts := CountSequential(text)

	// Test concurrent with various worker counts
	workerCounts := []int{1, 2, 3, 4, 8, 16, 32}
	for _, w := range workerCounts {
		concCounts := CountConcurrent(text, w)
		if !CompareMaps(seqCounts, concCounts) {
			t.Errorf("Mismatched counts with %d workers. Seq len: %d, Conc len: %d", w, len(seqCounts), len(concCounts))
			// Print mismatches
			for k, vSeq := range seqCounts {
				if vConc, ok := concCounts[k]; !ok || vSeq != vConc {
					t.Errorf("Key %q: seq=%d, conc=%d (exists=%t)", k, vSeq, vConc, ok)
				}
			}
			for k, vConc := range concCounts {
				if _, ok := seqCounts[k]; !ok {
					t.Errorf("Key %q: only in conc=%d", k, vConc)
				}
			}
		}
	}
}
