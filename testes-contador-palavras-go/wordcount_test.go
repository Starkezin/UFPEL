package main

import (
	"reflect"
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
		{"   ", nil},
		{"aaa...bbb...ccc", []string{"aaa", "bbb", "ccc"}},
	}

	for _, tc := range tests {
		got := CleanAndSplit(tc.input)
		if len(got) == 0 && len(tc.expected) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("CleanAndSplit(%q) = %v; expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name: "Caso Original - Especificação Base",
			input: `Casa, casa! A casa é azul.
Árvore; árvore? verde.
Go go Go. IA é útil, mas IA erra.`,
			expected: map[string]int{
				"casa":   3,
				"árvore": 2,
				"azul":   1,
				"verde":  1,
				"útil":   1,
				"mas":    1,
				"erra":   1,
			},
		},
		{
			name:     "Caso 1 - Texto Vazio e Delimitadores",
			input:    "   \n\t  .,!?;: ",
			expected: map[string]int{},
		},
		{
			name:     "Caso 2 - Apenas Palavras Curtas (< 3 caracteres)",
			input:    "Oi, ir ao ar de pó só ou é um as no.",
			expected: map[string]int{},
		},
		{
			name:  "Caso 3 - Combinações Extremas de Caixa e Pontuação",
			input: "Gol! gol... GOL, GoL: gol.",
			expected: map[string]int{
				"gol": 5,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CountWords(tc.input)
			if len(got) == 0 && len(tc.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("CountWords() = %v; expected %v", got, tc.expected)
			}
		})
	}
}
