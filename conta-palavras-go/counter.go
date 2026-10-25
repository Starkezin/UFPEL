package main

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// CleanAndSplit splits a text string into lowercase alphanumeric tokens.
// It ignores tokens with fewer than 3 characters (runes).
// It uses unicode definitions of letters and numbers to support accented and multi-byte characters.
func CleanAndSplit(text string) []string {
	// FieldsFunc splits the string at each rune where isNotAlphanumeric returns true.
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

// CountSequential counts the word frequencies of r sequentially.
func CountSequential(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		words := CleanAndSplit(scanner.Text())
		for _, w := range words {
			counts[w]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

// isAlphanumericASCII checks if a byte is a standard ASCII alphanumeric character.
func isAlphanumericASCII(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// CountConcurrent counts the word frequencies of r concurrently using a worker pool.
// It reads the reader in chunks, ensures no word is cut in half at chunk boundaries,
// processes chunks in parallel, and merges the local maps.
func CountConcurrent(r io.Reader, numWorkers int) (map[string]int, error) {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	jobs := make(chan string, numWorkers*2)
	results := make(chan map[string]int, numWorkers)
	var wg sync.WaitGroup

	// Start the worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localCounts := make(map[string]int)
			for chunk := range jobs {
				words := CleanAndSplit(chunk)
				for _, w := range words {
					localCounts[w]++
				}
			}
			results <- localCounts
		}()
	}

	// Read the file in chunks and send to jobs channel
	const chunkSize = 64 * 1024 // 64 KB chunks
	go func() {
		defer close(jobs)

		buf := make([]byte, chunkSize)
		var leftover []byte

		for {
			n, err := r.Read(buf)
			if n > 0 {
				chunk := append(leftover, buf[:n]...)
				leftover = nil

				// Find the last word separator boundary in the chunk.
				// Any byte < 128 that is not alphanumeric is a safe ASCII separator.
				// This ensures we never split a multi-byte UTF-8 rune in half (bytes >= 128).
				lastBoundary := len(chunk)
				for lastBoundary > 0 {
					b := chunk[lastBoundary-1]
					if b < 128 && !isAlphanumericASCII(b) {
						break
					}
					lastBoundary--
				}

				if lastBoundary == 0 {
					// No separator found; the entire chunk is part of a single word.
					// Keep it in leftover to read more.
					leftover = chunk
				} else {
					leftover = make([]byte, len(chunk)-lastBoundary)
					copy(leftover, chunk[lastBoundary:])
					jobs <- string(chunk[:lastBoundary])
				}
			}

			if err != nil {
				if err == io.EOF {
					break
				}
				// In case of read error, we break. The parent function will handle
				// completing workers and returning the result.
				break
			}
		}

		// Send any remaining leftover text
		if len(leftover) > 0 {
			jobs <- string(leftover)
		}
	}()

	// Wait for workers in a separate goroutine and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Merge the local maps from workers
	finalCounts := make(map[string]int)
	for localCounts := range results {
		for w, count := range localCounts {
			finalCounts[w] += count
		}
	}

	return finalCounts, nil
}
