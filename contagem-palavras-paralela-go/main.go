package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <caminho-do-arquivo> [num-workers]")
		os.Exit(1)
	}

	filePath := os.Args[1]

	numWorkers := runtime.NumCPU()
	if len(os.Args) > 2 {
		val, err := strconv.Atoi(os.Args[2])
		if err == nil && val > 0 {
			numWorkers = val
		} else {
			fmt.Printf("Aviso: Número de workers inválido '%s'. Usando valor padrão: %d\n", os.Args[2], numWorkers)
		}
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read file contents into memory
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	// 1. Measure sequential execution time
	startSeq := time.Now()
	seqCounts := CountSequential(content)
	elapsedSeq := time.Since(startSeq)

	// 2. Measure concurrent execution time
	startConc := time.Now()
	concCounts := CountConcurrent(content, numWorkers)
	elapsedConc := time.Since(startConc)

	// 3. Compare results
	equal := CompareMaps(seqCounts, concCounts)
	equalStr := "não"
	if equal {
		equalStr = "sim"
	}

	// 4. Sort and get top 20 words
	top20 := GetTop20(seqCounts)

	// 5. Print results matching format requirements
	fmt.Printf("Tempo sequencial: %v\n", elapsedSeq)
	fmt.Printf("Tempo concorrente: %v (workers: %d)\n", elapsedConc, numWorkers)
	fmt.Printf("Resultados iguais: %s\n", equalStr)
	fmt.Println("Top 20 palavras:")
	for _, wc := range top20 {
		fmt.Printf("- %s: %d\n", wc.Word, wc.Count)
	}
}
