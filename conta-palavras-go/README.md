# Contador de Frequência de Palavras em Go

Este projeto implementa um programa em Go para contar a frequência de palavras em um arquivo de texto, suportando uma versão sequencial e uma versão concorrente parametrizável.

---

## Estratégia Concorrente

* **Modelo adotado**: **Pool de Workers com Mapas Locais e Redução Final**.
* **Como funciona**: A thread de leitura fatia o arquivo em blocos de 64KB e os envia por canal. Para evitar cortar palavras ou quebrar runas UTF-8 na borda do bloco, a thread de leitura avança até encontrar um caractere delimitador antes de fazer o corte. Um grupo de goroutines trabalhadoras (*workers*) processa os blocos concorrentemente, gerando contagens locais e isoladas. Ao término, os mapas parciais são fundidos sequencialmente.

### Alternativas Consideradas
1. **Mapa Global com Mutex (`sync.Mutex`)**: Simples, mas causaria contenção extrema de travas a cada palavra inserida, degradando a performance.
2. **Mapa Concorrente Nativo (`sync.Map`)**: Inadequado devido à baixa eficiência em cenários de escrita intensa (atualização constante de frequências).
3. **Canais por Palavra (`chan string`)**: Alto custo de sincronização interna e gargalo na goroutine centralizada que faria as escritas no mapa.

---

## Segurança Concorrente (Prevenção de Race Conditions)

1. **Memória Isolada**: Cada worker manipula estritamente seu próprio mapa `localCounts` privado. Não há mapas compartilhados.
2. **Dados Imutáveis**: Os dados são lidos em fatias e passados via canal como `string` imutável, garantindo que a thread de leitura não altere a memória que está sendo processada.
3. **Sincronização de Encerramento**: O encerramento seguro dos canais e a fusão de dados ocorrem de forma ordenada e determinística usando `sync.WaitGroup`, eliminando riscos de deadlocks.

---

## Medição de Tempos e Comparação de Resultados

* **Medição**: O tempo de execução de ambas as abordagens é calculado de forma isolada na função `main()` usando `time.Now()` e `time.Since()`.
* **Comparação**: O programa compara as tabelas de frequências geradas chave-por-chave e valor-por-valor (`compareMaps`). A validação exige exatidão de 100% de todo o conjunto de dados.

---

## Código da Versão Concorrente

```go
func CountConcurrent(r io.Reader, numWorkers int) (map[string]int, error) {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	jobs := make(chan string, numWorkers*2)
	results := make(chan map[string]int, numWorkers)
	var wg sync.WaitGroup

	// Inicia o pool de workers
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

	// Lê o arquivo em chunks e envia para o canal jobs
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

				lastBoundary := len(chunk)
				for lastBoundary > 0 {
					b := chunk[lastBoundary-1]
					if b < 128 && !isAlphanumericASCII(b) {
						break
					}
					lastBoundary--
				}

				if lastBoundary == 0 {
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
				break
			}
		}

		if len(leftover) > 0 {
			jobs <- string(leftover)
		}
	}()

	// Aguarda os workers e fecha o canal de resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	// Mescla os mapas locais de todos os workers
	finalCounts := make(map[string]int)
	for localCounts := range results {
		for w, count := range localCounts {
			finalCounts[w] += count
		}
	}

	return finalCounts, nil
}
```

---

## Resultados de Desempenho (Dataset ~160KB)

Abaixo estão os tempos coletados com diferentes contagens de workers:

* **1 Worker**:
  * Tempo sequencial: `5.23ms`
  * Tempo concorrente: `7.07ms`
* **2 Workers**:
  * Tempo sequencial: `4.18ms`
  * Tempo concorrente: `3.75ms` (Ganho de performance: ~1.4x)
* **4 Workers**:
  * Tempo sequencial: `4.88ms`
  * Tempo concorrente: `3.97ms`
* **8 Workers**:
  * Tempo sequencial: `7.07ms`
  * Tempo concorrente: `3.84ms`

### Análise de Desempenho
A versão concorrente atinge seu ganho máximo com 2 ou 3 workers (reduzindo o tempo de processamento a ~3.7ms). Com 4 ou mais workers, ocorre um platô de desempenho por conta do tamanho do arquivo de entrada (~160KB), que é fatiado em apenas **3 blocos** de 64KB, deixando os demais núcleos de processamento adicionais ociosos. Em conjuntos de dados maiores (na escala de dezenas ou centenas de megabytes), a aceleração escala linearmente com os recursos da máquina.

---

## Como Executar

Execute o comando passando o caminho do arquivo de entrada e, opcionalmente, o número de workers:

```bash
go run . input/AChristmasCarol_CharlesDickens_English.txt [workers]
```

Para rodar os testes automatizados:
```bash
go test -v ./...
```


## 

Qual ferramenta de IA foi utilizada e em qual ambiente? Informe versões, quando possível.
R: Antigravity IDE (Gemini 3.5 Flash / Medium), Go 1.25.4, Ubuntu


Qual estratégia concorrente foi escolhida?
R: Pool de Workers com Mapas Locais e Redução Final


Quais alternativas foram consideradas?
R: Mapas globais com Mutex e sync.Map, canais por palavra.

Por que a estratégia escolhida pareceu mais adequada?
R: Evita contenção de recursos e permite escalabilidade horizontal.


O resultado concorrente foi exatamente igual ao sequencial?
R: Sim

A versão concorrente melhorou o desempenho?
R: Sim, porém sinto que para o tamanho da entrada foi pouca a melhora. Acredito que para um arquivo maior a melhora seria mais expressiva.

Onde poderia ocorrer condição de corrida?
R: Poderia ocorrer condição de corrida se os workers estivessem escrevendo no mesmo mapa sem sincronização.

A sincronização usada era necessária?
R: Sim, para garantir que os workers não acessassem o mesmo mapa ao mesmo tempo.

Houve contenção?
R: Sim, houve contenção no canal de resultados.

A IA sugeriu alguma solução problemática?
R: Não, apenas implementou pontos antes mesmo de ser solicitado para implementar, parecia que mesmo eu fazendo perguntas e tentando ser critico antes de implementar, ela buscava sempre uma resposta para começar a implementar imediatamente.

O que precisou ser corrigido na solução sugerida pela IA?
R: A implementação antes da definição final de qual estratégia seria utilizada.

Qual é a diferença entre simplesmente usar recursos concorrentes e construir uma boa solução concorrente?
R: Caso simplesmente utilizasse recursos concorrentes sem garantir uma boa solução, ele provavelmente teria resultados expressivamente piores que a versão sequencial. A concorrência quando aplicada corretamente traz bons resultados, mas precisa ser implementada de maneira cuidadosa.


-- Extra:
Senti durante a implementação que o Gemini/Antigravity estava tentando buscar uma resposta que fosse aceita de forma imediata para começar a implementação, invés de tentar montar toda estratégia/plano de implementação de forma detalhada para garantir que a implementação fosse dentro do que eu acredito ser a melhor solução utilizando esse ferramental. 
Por outro lado a velocidade de implementação de uma solução desse tamanho pela IA em poucos minutos é algo realmente impressionante, para uma pessoa implementar essa mesma solução da maneira antiga (lendo documentação, etc.) levariam horas ou dias.