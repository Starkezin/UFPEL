## Etapa 1

O agente implementou as funções solicitadas no Prompt 01, não fez questionamentos ou nada do tipo. Além disso, o agente também implementou os códigos no diretório errado.

## Etapa 2

O agente informou que o código já estava implementado e informou que possuia os seguintes detalhes:
    - Validação de tamanho utilizando utf8.RuneCountInString invés de len(). 
    - Separação de palavras e símbolos utilizando regex e FieldsFunc.
    - Conversão para minúsculas utilizando strings.ToLower().

# Respostas:
o código compila?
- R: Sim
a função tem a assinatura esperada?
- R: Sim
a normalização está correta?
- R: Sim
palavras curtas são ignoradas?
- R: Sim
acentos são preservados?
- R: Sim

    
## Etapa 3

Foram gerados dois testes, um para limpar e separar palavras e outro para contar as palavras, ambos utilizando o pacote testing.

PASS
ok  	testes-automaticos-go	0.003s

Nenhum erro foi encontrado.

## Etapa 4

A revisão aponta que ambos os testes estão corretamente implementados, utilizando os recursos do pacote testing de forma apropriada.

## Etapa 5

Novamente o agente implementou por conta própria o código, invés de apenas sugerir novos testes para inclusão.

Os testes implementados foram:
    - Texto Vazio e Delimitadores
    - Apenas Palavras Curtas (< 3 caracteres)
    - Combinações Extremas de Caixa e Pontuação
    
## Etapa 6

Testes não cobertos:
    - Pontuação interna ("guarda-chuva", "d'água").
    - Números ("123"), são sempre tratados como string.
    - Alfanuméricos ("Python3", "IPv6"). 


    
## Análise do estudante:

Acredito que em termos de implementação, a qualidade do código e dos testes tem um nível muito acima do que o uma pessoa entregaria no mesmo tempo que levou para ser feito por uma IA. Contudo, especificamente nas etapas de implementação utilizando o Antigravity, o agente não faz perguntas críticas ou um plano de execução para implementação, ele apenas recebe a instrução e implementa, algo que não acontece quando utilizo outros agentes, como claude e codex.
