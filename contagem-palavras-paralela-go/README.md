# Conta-palavras Concorrente em Go

Este repositório contém uma solução concorrente e sequencial desenvolvida em Go para contar a frequência de palavras em um arquivo texto, comparando a correção e o tempo de execução entre ambas as abordagens.

---

## 🚀 Como Executar o Programa

### Pré-requisitos
- Ter o [Go](https://go.dev/) instalado (versão 1.18 ou superior recomendada).

### Download do Dataset
O arquivo de texto *A Christmas Carol* de Charles Dickens foi baixado no diretório raiz do projeto. Caso necessite baixar novamente:
```bash
curl -L -o AChristmasCarol_CharlesDickens_English.txt https://raw.githubusercontent.com/FilipeLopesPires/LargeText-WordCount/main/datasets/AChristmasCarol_CharlesDickens/AChristmasCarol_CharlesDickens_English.txt
```

### Execução da Aplicação
Para executar o programa com a contagem sequencial e concorrente, passe o caminho do arquivo como parâmetro:
```bash
go run . AChristmasCarol_CharlesDickens_English.txt
```

Para definir um número específico de workers/goroutines para a versão concorrente (por exemplo, 4 workers):
```bash
go run . AChristmasCarol_CharlesDickens_English.txt 4
```

Se nenhum número de workers for fornecido, a aplicação utilizará a quantidade padrão de CPUs lógicas disponíveis no sistema (`runtime.NumCPU()`).

### Execução dos Testes Automatizados
Para executar os testes unitários e validar a corretude das funções e a equivalência das abordagens:
```bash
go test -v ./...
```

---

## 🛠️ Estratégia Concorrente

A versão concorrente foi implementada com foco em alto desempenho e segurança de threads (thread-safety), evitando contenção por travas (mutexes) durante a fase de processamento paralelo.

### 1. Divisão de Trabalho
1. O conteúdo do arquivo texto é lido inteiramente na memória como um buffer de bytes (`[]byte`).
2. O buffer de bytes é dividido logicamente em `N` fatias (onde `N` é o número configurado de workers).
3. **Evitando corte de palavras**: Para garantir que nenhuma palavra seja cortada ao meio na fronteira das fatias (por exemplo, metade em um worker e metade em outro), o ponto de término de cada fatia é estendido incrementalmente até encontrar o próximo caractere ASCII de espaço em branco (`' '`, `'\n'`, `'\r'`, `'\t'`). 

### 2. Processamento em Goroutines
Cada fatia ajustada é passada para uma goroutine dedicada. Essa goroutine executa de forma independente o processo de:
- Substituir separadores sintáticos especiais (como travessões e hífens duplos) por espaços.
- Separar o texto por espaços em branco para extrair tokens.
- Converter letras para minúsculas e limpar pontuações simples.
- Ignorar tokens com menos de 3 caracteres.
- Registrar as frequências em um mapa local (`map[string]int`).

### 3. Combinação de Resultados (Reduce)
Cada goroutine envia seu mapa de frequências local para a thread principal por meio de um canal seguro (`chan`). A thread principal aguarda o término das goroutines e mescla de forma segura os resultados parciais em um mapa de contagem global definitivo (`finalCounts`).

---

## ✅ Verificação de Correção

A consistência e correção das implementações foram verificadas de duas formas:

1. **Testes Unitários Automatizados (`wordcount_test.go`)**:
   - `TestCleanWord`: valida regras de conversão para minúsculas, remoção de caracteres de pontuação específicos e filtragem por comprimento mínimo.
   - `TestCompareMaps`: valida se a rotina de comparação detecta divergências de chaves, valores e tamanhos em mapas.
   - `TestGetTop20`: valida se a lista final das 20 palavras mais frequentes está classificada corretamente em ordem decrescente de frequência e alfabética secundária.
   - `TestEquivalence`: processa um texto de teste complexo usando a versão sequencial e a versão concorrente sob diferentes cenários de paralelismo (1, 2, 3, 4, 8, 16, 32 workers) e exige que os mapas resultantes sejam 100% idênticos.

2. **Comparação Dinâmica em Tempo de Execução**:
   - Toda execução da aplicação compara o mapa sequencial e concorrente por completo (chaves e valores) e imprime na tela `Resultados iguais: sim` caso a equivalência matemática seja confirmada.

---

## ⏱️ Medição de Tempo e Resultados

### Método de Medição
A medição de tempo utiliza o pacote padrão do Go (`time.Now()` e `time.Since()`). O trecho medido em cada versão é apenas o processamento e contagem do buffer de bytes na memória, garantindo que a latência de entrada e saída (I/O) de arquivos em disco não distorça as métricas comparativas.

### Resultados Obtidos (Dataset de 159 KB)
Os tempos foram coletados rodando o projeto em um ambiente de CPU multi-core (Linux):

| Versão / Configuração | Tempo de Execução | Resultados Iguais? |
| :--- | :--- | :--- |
| **Sequencial** (Referência) | ~7.29 ms | - |
| **Concorrente (1 worker)** | ~4.79 ms | Sim |
| **Concorrente (2 workers)** | ~4.17 ms | Sim |
| **Concorrente (4 workers)** | ~4.23 ms | Sim |
| **Concorrente (8 workers)** | ~4.51 ms | Sim |
| **Concorrente (16 workers)** | ~2.97 ms | Sim |

*Nota: A versão concorrente obteve melhor desempenho geral do que a sequencial (cerca de 2.5x mais rápida com 16 workers). A divisão das tarefas em fatias e o uso de mapas locais permitiram uma aceleração significativa mesmo com arquivos pequenos, superando com facilidade o custo de gerenciar goroutines e canais em Go.*

---

## 🧠 Dificuldades Encontradas

1. **Garantia de Limites de Palavras na Divisão (Word Chunks)**:
   - Dividir o arquivo de forma estrita em partes de tamanhos iguais (`tamanho / N`) causaria o corte de palavras no meio. A solução foi ler o buffer de bytes e deslizar o marcador de término até encontrar o próximo espaço em branco, garantindo a integridade dos limites lexicais.
2. **Divisões de Frase Sem Espaço (`--` e `—`)**:
   - No texto clássico de Dickens, há muitos travessões colados às palavras (ex: `desired--a`). Apenas separar por espaço em branco uniria as duas palavras em uma só (`desireda`). A solução foi realizar a substituição desses caracteres por espaços em branco antes de quebrar o texto em tokens, isolando corretamente as palavras de forma idêntica nas duas versões.

---

## 🤖 Relatório de Uso de IA

Conforme exigido pelas instruções da atividade, este é o relatório detalhado do uso de Inteligência Artificial para a execução do trabalho:

- **Ferramenta Utilizada**: Antigravity, assistente de programação agentivo projetado pela equipe do Google DeepMind.
- **Ambiente e Modelo**: Executado em um ambiente sandbox integrado no ecossistema de desenvolvimento Linux local, utilizando o modelo de linguagem **Gemini 3.5 Flash (Medium)**.
- **Passos Vencidos com a IA**:
  1. *Concepção e Arquitetura*: Definição das estruturas compartilhadas (`WordCount`), tratamento lexical e estratégia de merge de mapas locais.
  2. *Escrita do Código em Go*: Implementação assistida das rotinas sequenciais (`seq.go`), concorrentes (`concurrent.go`), utilitários (`utils.go`) e ponto de entrada (`main.go`).
  3. *Tratamento de Corner Cases*: Ajuste do loop de fatiamento no limite do buffer para evitar vazamento de índices (out-of-bounds) e ajuste da tokenização para tratar adequadamente pontuações sem espaço (como os travessões `--` e `—`).
  4. *Desenvolvimento de Testes Unitários*: Escrita automatizada de testes com foco na equivalência das abordagens e validação determinística de ordenação alfabética em caso de empates.
  5. *Documentação*: Estruturação deste documento README com análise crítica de performance.
