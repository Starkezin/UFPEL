# Registro de Uso da Inteligência Artificial

Este documento registra o processo de uso da Inteligência Artificial no desenvolvimento do Contador de Frequência de Palavras em Go.

---

## 1. Metadados do Ambiente

* **Ferramenta de IA Utilizada**: Antigravity IDE (alimentada pelo modelo Gemini 3.5 Flash / Medium).
* **Ambiente de Desenvolvimento**: Linux (Ubuntu), Terminal Bash, VSCode/Antigravity editor.
* **Versão do Go**: Go 1.18+ (utilizando recursos nativos de `sync` e `unicode/utf8`).
* **Especificações do Hardware**: Processador Intel Core i5-13420H (13ª geração, 12 núcleos lógicos).

---

## 2. Relação de Prompts Utilizados

### Prompts Obrigatórios (Etapas de Desenvolvimento)
1. **Planejamento inicial**: *"Quero implementar em Go um programa de contagem de frequência de palavras em um arquivo texto. O programa deve ter uma versão sequencial e uma versão concorrente. Explique o problema, as principais decisões de implementação e os cuidados necessários..."*
2. **Detalhamento do fluxo**: *"como o arquivo será lido; como as palavras serão separadas; como as palavras serão normalizadas; como as frequências serão armazenadas; como a versão sequencial pode servir como referência; que problemas podem aparecer na versão concorrente"*
3. **Comparação de alternativas**: *"Em Go, quais são as principais formas de implementar uma contagem concorrente de frequência de palavras em um arquivo texto? Compare alternativas..."*
4. **Justificativa**: *"Considerando o problema de contar frequência de palavras em um arquivo texto, qual estratégia concorrente em Go parece mais adequada? Justifique considerando correção, simplicidade, risco de condição de corrida, custo de sincronização e desempenho."*
5. **Proposta de Estrutura**: *"Proponha a estrutura de um programa em Go para resolver esse problema. O programa deve conter uma versão sequencial, uma versão concorrente, medição de tempo... Não escreva o código ainda; descreva apenas as principais funções necessárias."*
6. **Revisão de I/O**: *"Nao tem nenhuma função para leitura de arquivos?
7. **Implementação Sequencial**: *"Implemente em Go a versão sequencial do conta-palavras. O programa deve receber o caminho do arquivo texto por argumento de linha de comando, contar a frequência das palavras, ignorar palavras com menos de 3 caracteres, converter para minúsculas, remover pontuação simples, medir o tempo de execução da versão sequencial e imprimir as 20 palavras mais frequentes."*
8. **Dataset de Validação**: *"Crie um pequeno arquivo de teste para validar a contagem de palavras antes de usar o dataset principal. O arquivo deve conter palavras repetidas, letras maiúsculas e minúsculas, pontuação e palavras com menos de 3 caracteres..."*
9. **Homologação Sequencial**: *"Adicione ao programa uma forma simples de testar a versão sequencial usando um arquivo pequeno com resultado esperado conhecido..."*
10. **Implementação Concorrente**: *"A partir da versão sequencial, implemente uma versão concorrente em Go para a contagem de frequência de palavras. Use a estratégia escolhida anteriormente..."*
11. **Verificação de Exatidão**: *"Adicione ao programa uma função para comparar o mapa de frequências produzido pela versão sequencial com o mapa produzido pela versão concorrente. A comparação deve garantir que os dois mapas tenham exatamente as mesmas palavras com as mesmas frequências..."*
12. **Análise de Concorrência**: *"Analise a versão concorrente do programa abaixo. Existe alguma condição de corrida? Alguma estrutura de dados é escrita simultaneamente por mais de uma goroutine? Há risco de deadlock? A sincronização usada é necessária? Existe contenção desnecessária? Cole o código completo."*
13. **Diferentes Workers**: *"Quando a estratégia permitir configurar o número de workers, tarefas ou goroutines, execute diferentes configurações... Registre os tempos obtidos."*
14. **Análise de Tempos**: *"Analise os tempos de execução obtidos. Explique por que a versão concorrente foi mais rápida, mais lenta ou equivalente à versão sequencial. Considere tamanho do arquivo, custo de criação de goroutines, custo de sincronização, custo de comunicação, divisão do trabalho, número de núcleos e overhead da combinação dos resultados."*

---

## 3. Resumo das Respostas da IA e Decisões Técnicas

* **Tokenização e Normalização**: Definiu-se que pontuações e espaços delimitam palavras via `unicode.IsLetter` e `unicode.IsNumber`. Tokens são normalizados em minúsculas via `strings.ToLower`. Palavras menores que 3 caracteres (runas UTF-8) são ignoradas com `utf8.RuneCountInString(w) >= 3` para garantir corretude em acentos.
* **Modelo Concorrente**: Adotou-se o **Pool de Workers com Mapas Locais e Redução Final**. O arquivo é lido sequencialmente em blocos de 64KB e estendido nas bordas de leitura para não fatiar palavras ou runas UTF-8. Goroutines trabalhadoras computam sub-tabelas `map[string]int` privadas (zero contenção de locks) e a thread principal as mescla sequencialmente no final.
* **Comparação**: O utilitário faz comparação profunda chave-por-chave (`compareMaps`) confirmando a igualdade absoluta.

---

## 4. Sugestões Aceitas e Rejeitadas

* **Sugestões Aceitas**:
  1. *Uso de mapas locais por worker* para eliminar o uso de exclusão mútua (`Mutex`) durante o processamento do texto.
  2. *Uso da interface `io.Reader`* para desacoplar a leitura física do arquivo da lógica algorítmica de contagem, tornando a ferramenta flexível e testável sem mockar arquivos.
  3. *Uso de buffers de 64KB* para transmitir pedaços maiores de dados pelo canal, minimizando o overhead de mensagens.
* **Sugestões Rejeitadas**:
  1. *Mapa Global com Mutex* ou uso de *`sync.Map`*: Rejeitados por causarem contenção excessiva de CPU e gerarem degradação severa de desempenho em escritas massivas e frequentes.
  2. *Canais de palavra individual*: Rejeitado por criar gargalo severo na goroutine consumidora central.

---

## 5. Erros Gerados pela IA e Correções Realizadas

1. **Variável `err` Declarada e Não Utilizada**:
   * *O que houve*: A IA declarou `var err error` no escopo superior da função `main()`, mas nos blocos internos `if` usou atribuição por operador de curto-circuito `:=` (ex: `f, err := os.Open(...)`), gerando shadowing local do erro e deixando a variável externa sem uso.
   * *Correção*: Removida a variável redundante de escopo externo e unificado o tratamento local.
2. **Falha de Regressão nos Testes Unitários**:
   * *O que houve*: O teste preexistente `TestCleanAndSplit` continha a expectativa `"a...b...c" -> ["a", "b", "c"]`. Quando implementamos o filtro de tamanho mínimo (>= 3 caracteres), esses tokens de 1 caractere passaram a ser descartados legítimamente, quebrando o teste.
   * *Correção*: Atualizou-se o teste para a string `"aaa...bbb...ccc"` para compatibilizar a suíte de testes com a nova regra de negócio.

---

## 6. Avaliação Crítica do Estudante

* **Sobre a análise de condições de corrida (Etapa 10)**:
  Acredito que ele deu uma resposta bem completa, obviamente bem caracteristica de uma IA, mas detalhando de forma clara os pontos sobre condições de corrida.

* **Sobre a análise de tempos e hardware (Etapa 11)**:
  Vejo uma análise bem completa, trazendo desde questões de Hardware para a análise, não apenas o código em si. Também análise e traz informações sobre o conjunto de entrada. 
  
  Acredito que a análise técnica foi correta, contudo como padrão a IA acaba trazendo muitas informações que por vezes poderiam ser resumidas e simplificadas, sem tanta enrolação.

