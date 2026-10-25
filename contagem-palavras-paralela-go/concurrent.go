package main

// CountConcurrent counts the word frequencies in the given content concurrently using numWorkers goroutines.
func CountConcurrent(content []byte, numWorkers int) map[string]int {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	totalLen := len(content)
	if totalLen == 0 {
		return make(map[string]int)
	}

	// Step 1: Divide content into chunks along whitespace boundaries.
	chunks := make([][]byte, 0, numWorkers)
	start := 0
	for i := 0; i < numWorkers; i++ {
		if start >= totalLen {
			break
		}
		// Calculate initial end index for this worker.
		end := (i + 1) * totalLen / numWorkers
		if i == numWorkers-1 || end > totalLen {
			end = totalLen
		} else {
			// Find the next whitespace character to avoid splitting a word.
			for end < totalLen && !isWhitespaceByte(content[end]) {
				end++
			}
		}
		chunks = append(chunks, content[start:end])
		start = end
	}

	// Step 2: Set up a channel to collect results.
	type result struct {
		counts map[string]int
	}
	ch := make(chan result, len(chunks))

	// Step 3: Spin up workers.
	for _, chunk := range chunks {
		go func(c []byte) {
			localCounts := CountSequential(c)
			ch <- result{counts: localCounts}
		}(chunk)
	}

	// Step 4: Merge results in the main thread.
	finalCounts := make(map[string]int)
	for i := 0; i < len(chunks); i++ {
		res := <-ch
		for k, v := range res.counts {
			finalCounts[k] += v
		}
	}

	return finalCounts
}

// isWhitespaceByte returns true if the byte is ASCII whitespace.
func isWhitespaceByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}
