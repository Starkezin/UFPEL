package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCleanAndSplit(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"Hello, World! 123", []string{"hello", "world", "123"}},
		{"Café com pão: muito bom!", []string{"café", "com", "pão", "muito", "bom"}},
		{"Olá, como vai? Tudo bem! Olá", []string{"olá", "como", "vai", "tudo", "bem", "olá"}},
		{"   ", []string(nil)}, // strings.FieldsFunc might return nil or empty slice for only delimiters
		{"aaa...bbb...ccc", []string{"aaa", "bbb", "ccc"}},
	}

	for _, tc := range tests {
		got := CleanAndSplit(tc.input)
		// Normalize empty slices to compare with nil or empty safely
		if len(got) == 0 && len(tc.expected) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("CleanAndSplit(%q) = %v; expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestCountSequentialSmall(t *testing.T) {
	text := `Olá Mundo! O mundo é belo, o mundo é GRANDE.
Um teste simples, Simples e curto.
Oi, tudo bem? Tudo ótimo!
Ir ao parque.`

	expected := map[string]int{
		"mundo":   3,
		"simples": 2,
		"tudo":    2,
		"olá":     1,
		"belo":    1,
		"grande":  1,
		"teste":   1,
		"curto":   1,
		"bem":     1,
		"ótimo":   1,
		"parque":  1,
	}

	r := strings.NewReader(text)
	got, err := CountSequential(r)
	if err != nil {
		t.Fatalf("CountSequential falhou: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Divergência na contagem sequencial pequena.\nObtido: %v\nEsperado: %v", got, expected)
	}
}

func TestCountCompare(t *testing.T) {
	text := `Olá mundo! Este é um teste para o contador de palavras.
O contador deve ter uma versão sequencial e uma versão concorrente.
Tanto a versão sequencial quanto a versão concorrente devem produzir EXATAMENTE o mesmo resultado.
Café com pão é muito bom. Olá de novo!`

	rSeq := strings.NewReader(text)
	seqCounts, err := CountSequential(rSeq)
	if err != nil {
		t.Fatalf("CountSequential falhou: %v", err)
	}

	workersList := []int{1, 2, 4, 8}
	for _, w := range workersList {
		rConc := strings.NewReader(text)
		concCounts, err := CountConcurrent(rConc, w)
		if err != nil {
			t.Fatalf("CountConcurrent com %d workers falhou: %v", w, err)
		}

		if !reflect.DeepEqual(seqCounts, concCounts) {
			t.Errorf("Resultados diferem para %d workers!\nSeq: %v\nConc: %v", w, seqCounts, concCounts)
		}
	}
}

func BenchmarkCountSequential(b *testing.B) {
	text := strings.Repeat("O rato roeu a roupa do rei de Roma. Café com pão é bom!\n", 50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(text)
		_, _ = CountSequential(r)
	}
}

func BenchmarkCountConcurrent_2Workers(b *testing.B) {
	text := strings.Repeat("O rato roeu a roupa do rei de Roma. Café com pão é bom!\n", 50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(text)
		_, _ = CountConcurrent(r, 2)
	}
}

func BenchmarkCountConcurrent_4Workers(b *testing.B) {
	text := strings.Repeat("O rato roeu a roupa do rei de Roma. Café com pão é bom!\n", 50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(text)
		_, _ = CountConcurrent(r, 4)
	}
}
