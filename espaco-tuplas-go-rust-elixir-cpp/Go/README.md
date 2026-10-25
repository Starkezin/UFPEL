# Servidor de Espaço de Tuplas (Linda) - Implementação em Go

Este projeto consiste na implementação de um servidor concorrente de espaço de tuplas (Tuple Space) inspirado no modelo **Linda**, desenvolvido para a disciplina de **Programação Concorrente**.

A solução foi construída utilizando a linguagem **Go (Golang)**, aproveitando suas primitivas nativas de concorrência. O gerenciamento de estado e a sincronização foram implementados utilizando **Goroutines**, **Canais (Channels)** para o bloqueio eficiente de operações e `sync.RWMutex` para garantir o acesso seguro e sem *busy-waiting* às estruturas de dados em memória.

## 📌 Requisitos

* **Go 1.18** ou superior.
* Compilador `g++` (ambiente Linux/WSL recomendado) para executar o programa cliente de teste em C++.

## 🚀 Compilação e Execução

Para testar o projeto, você precisará de **dois terminais** abertos na raiz da pasta `go`.

**1. Iniciar o servidor (Terminal 1):**
Execute o comando abaixo para subir o servidor:

```bash
go run main.go
```

O servidor iniciará automaticamente e ficará ouvindo conexões na porta **54321**. *(Para encerrar graciosamente e liberar a porta, pressione `Ctrl + C` no terminal).*

**2. Executar o cliente de teste (Terminal 2):**
Com o servidor rodando no primeiro terminal, abra um segundo terminal (recomendado usar WSL/Linux para compatibilidade das bibliotecas C++), compile o testador e execute:

```bash
g++ -std=c++17 tester_linda.cpp -o tester_linda
./tester_linda 127.0.0.1 54321
```

## ⚙️ Configuração

* **Porta TCP:** `54321`.
* **Host:** `127.0.0.1` (Aceita conexões de múltiplos clientes simultaneamente via Goroutines).

## 📂 Estrutura de Arquivos

* `main.go`: Contém toda a lógica do servidor, incluindo o listener TCP, o protocolo de comunicação, a tabela estática de serviços e a estrutura segura do *Tuple Space*.
* `go.mod`: Arquivo de definição do módulo Go.
* `tester_linda.cpp`: Programa cliente de teste fornecido para validação das operações.

## 📡 Protocolo de Comunicação

O servidor opera sobre **TCP** e recebe comandos textuais. Os comandos devem ser separados por espaço e enviados com quebra de linha.

| Comando | Descrição | Bloqueante? | Resposta Sucesso |
| :--- | :--- | :---: | :--- |
| **`WR chave valor`** | Insere a tupla `(chave, valor)` no espaço. Nunca falha. | Não | `OK` |
| **`RD chave`** | Lê o valor da tupla mais antiga sem removê-la. | Sim | `OK valor` |
| **`IN chave`** | Lê e remove a tupla mais antiga do espaço. | Sim | `OK valor` |
| **`EX k_in k_out id`** | Consome `k_in` (bloqueante), aplica o serviço `id` e escreve o resultado em `k_out`. | Sim | `OK` |

> **Observações de Erro:**
> * Se o serviço solicitado no comando `EX` não existir, o servidor não altera o espaço e retorna `NO-SERVICE`.
> * Se o comando for mal formatado ou inválido, o servidor retorna `ERROR`.

## 🛠 Serviços Implementados (Operação EX)

Conforme a especificação do projeto, o servidor mantém uma tabela estática com os seguintes serviços de transformação de strings:

| ID (`svc_id`) | Serviço | Descrição | Exemplo |
| :---: | :--- | :--- | :--- |
| **`1`** | **To Upper** | Converte a string de entrada para letras maiúsculas. | `"abcdef"` → `"ABCDEF"` |
| **`2`** | **Reverse** | Inverte a ordem dos caracteres da string. | `"ghijkl"` → `"lkjihg"` |
| **`3`** | **Length** | Retorna o tamanho (quantidade de caracteres) da string em formato textual. | `"xyz"` → `"3"` |