package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type WordCount struct {
	Word  string
	Count int
}

func main() {
	mode := flag.String("mode", "both", "Modo de execução: 'seq', 'conc' ou 'both'")
	workers := flag.Int("workers", runtime.NumCPU(), "Número de worker threads para versão concorrente")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Erro: O caminho do arquivo é obrigatório como argumento posicional.\nUso: go run . [flags] <caminho_do_arquivo> [workers]")
	}
	filePath := args[0]

	numWorkers := *workers
	if len(args) >= 2 {
		val, err := strconv.Atoi(args[1])
		if err == nil && val > 0 {
			numWorkers = val
		}
	}

	if *mode != "seq" && *mode != "conc" && *mode != "both" {
		log.Fatalf("Erro: Modo '%s' inválido. Use 'seq', 'conc' ou 'both'.", *mode)
	}

	fmt.Printf("Configuração:\n")
	fmt.Printf(" - Arquivo: %s\n", filePath)
	fmt.Printf(" - Modo: %s\n", *mode)
	if *mode == "conc" || *mode == "both" {
		fmt.Printf(" - Workers: %d\n", numWorkers)
	}
	fmt.Println(stringsRepeat("-", 40))

	var seqCounts map[string]int
	var concCounts map[string]int
	var seqDuration time.Duration
	var concDuration time.Duration

	// Run Sequential
	if *mode == "seq" || *mode == "both" {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Erro ao abrir arquivo para versão sequencial: %v", err)
		}

		start := time.Now()
		seqCounts, err = CountSequential(f)
		seqDuration = time.Since(start)
		f.Close()

		if err != nil {
			log.Fatalf("Erro na versão sequencial: %v", err)
		}

		if *mode == "seq" {
			fmt.Printf("Tempo sequencial: %v\n", seqDuration)
		}
	}

	// Run Concurrent
	if *mode == "conc" || *mode == "both" {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Erro ao abrir arquivo para versão concorrente: %v", err)
		}

		start := time.Now()
		concCounts, err = CountConcurrent(f, numWorkers)
		concDuration = time.Since(start)
		f.Close()

		if err != nil {
			log.Fatalf("Erro na versão concorrente: %v", err)
		}

		if *mode == "conc" {
			fmt.Printf("Tempo concorrente: %v\n", concDuration)
		}
	}

	// Comparison and validation if both ran
	if *mode == "both" {
		fmt.Printf("Tempo sequencial: %v\n", seqDuration)
		fmt.Printf("Tempo concorrente: %v\n", concDuration)

		iguais := "não"
		if compareMaps(seqCounts, concCounts) {
			iguais = "sim"
		}
		fmt.Printf("Resultados iguais: %s\n", iguais)
	}

	// Print top 20 words
	var sampleMap map[string]int
	if seqCounts != nil {
		sampleMap = seqCounts
	} else {
		sampleMap = concCounts
	}

	top20 := getTopWords(sampleMap, 20)
	fmt.Println("Top 20 palavras:")
	for i, wc := range top20 {
		fmt.Printf("%2d. %s: %d\n", i+1, wc.Word, wc.Count)
	}
}

// stringsRepeat is a helper because we don't want to import "strings" in main.go if we can avoid it.
// Wait, we can import "strings" or write a tiny helper loop. Let's just use strings.Repeat.
// Oh, strings is already imported in counter.go, so importing it here is clean.
func stringsRepeat(s string, count int) string {
	var result string
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func compareMaps(m1, m2 map[string]int) bool {
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

func getTopWords(counts map[string]int, limit int) []WordCount {
	var list []WordCount
	for w, c := range counts {
		list = append(list, WordCount{Word: w, Count: c})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Count == list[j].Count {
			return list[i].Word < list[j].Word
		}
		return list[i].Count > list[j].Count
	})

	if len(list) > limit {
		return list[:limit]
	}
	return list
}
