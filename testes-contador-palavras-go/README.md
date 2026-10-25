# Testes Automáticos Go - Contador de Palavras

## Descrição Breve da Atividade
Esta atividade consiste na implementação e no teste de uma função em Go, `CountWords(text string) map[string]int`, encarregada de contar a frequência de palavras em uma string de texto. A solução engloba a normalização de caixa (conversão para minúsculas), a remoção de pontuações simples e o descarte de termos curtos (menos de 3 caracteres/runas).

## Explicação da Função CountWords
A função [CountWords](./wordcount.go) recebe uma string de entrada e retorna um mapa contendo as palavras e suas respectivas frequências. As etapas do processamento são:
1. **Tokenização e Remoção de Pontuação**: Utiliza `strings.FieldsFunc` mapeando símbolos que não são letras ou números (como `,`, `!`, `.`, `;`, `?`) como delimitadores de palavras.
2. **Normalização**: Converte cada palavra identificada para minúsculas usando `strings.ToLower`.
3. **Filtro de Tamanho**: Utiliza `utf8.RuneCountInString` para validar que a palavra possui pelo menos 3 caracteres/runas. Isso impede que caracteres Unicode multibyte (com acento) sejam contados incorretamente por seu tamanho em bytes.
4. **Contagem**: Insere as palavras validadas no mapa de retorno incrementando o contador de frequência correspondente.

## Descrição dos Casos de Teste Implementados
A suíte de testes unitários em [wordcount_test.go](./wordcount_test.go) cobre os seguintes cenários:
1. **Caso Original - Especificação Base**: Valida a contagem de repetições, case-insensitivity, descarte de palavras curtas (como `"A"`, `"é"`, `"Go"`, `"IA"`) e pontuações variadas de uma frase complexa.
2. **Caso 1 - Texto Vazio e Delimitadores**: Verifica que entradas sem conteúdo alfanumérico (como espaços, quebras de linha e pontuações isoladas) retornam um mapa vazio e inicializado, sem gerar pânico de execução.
3. **Caso 2 - Apenas Palavras Curtas**: Garante que uma entrada contendo somente palavras com menos de 3 caracteres (como `"Oi, ir ao ar de pó só"`) seja totalmente filtrada, retornando um mapa vazio.
4. **Caso 3 - Combinações Extremas de Caixa e Pontuação**: Estressa a normalização e junção de repetições para a mesma palavra sob formatos diversos (ex: `"Gol! gol... GOL, GoL: gol."` deve acumular tudo para a chave `"gol"` com frequência 5).
5. **CleanAndSplit Unit Tests**: Casos adicionais específicos testando a limpeza de tokens e conversão correta de acentuações como `"Café com pão: muito bom!"`.

## Instruções para Executar os Testes
Para executar a suíte de testes automáticos, execute o comando abaixo no terminal dentro do diretório do projeto:

```bash
go test -v ./...
```

## Resultado Obtido ao Executar `go test`
A execução dos testes retorna o seguinte output de sucesso:

```text
=== RUN   TestCleanAndSplit
--- PASS: TestCleanAndSplit (0.00s)
=== RUN   TestCountWords
=== RUN   TestCountWords/Caso_Original_-_Especificação_Base
=== RUN   TestCountWords/Caso_1_-_Texto_Vazio_e_Delimitadores
=== RUN   TestCountWords/Caso_2_-_Apenas_Palavras_Curtas_(<_3_caracteres)
=== RUN   TestCountWords/Caso_3_-_Combinações_Extremas_de_Caixa_e_Pontuação
--- PASS: TestCountWords (0.00s)
PASS
ok      testes-automaticos-go   0.002s
```
