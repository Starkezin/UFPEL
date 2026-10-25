## Atividade extra-classe: OpenMP na prática, com ênfase em laços, escalonamento e seções críticas

Alunos: Diego Aquino e Samuel Starke

A divisão das tarefas foi realizada de forma que Diego ficou responsável pelo desenvolvimento da base da tarefa A e Samuel da tarefa D. Após cada um desenvolver uma versão funcional de sua tarefa, testes e adaptações finais para cumprir com o máximo de requisitos foram feitos em conjunto.

## Estrutura do Projeto

O projeto está organizado da seguinte forma para manter a limpeza do repositório:

- `src/`: Códigos fonte (.c) separados em `seq` (sequencial) e `omp` (paralelo).
- `bin/`: Diretório onde os executáveis são gerados.
- `results/`: Diretório onde os arquivos `.csv` com os tempos são salvos.
- `plots/`: Diretório onde os gráficos gerados são salvos.
- `scripts/`: Scripts Python para geração dos gráficos.

---

## Ambiente e Requisitos

Este projeto foi configurado para rodar em ambiente **Linux**.

### Pré-requisitos do Sistema
Certifique-se de ter as seguintes ferramentas instaladas no seu sistema operacional:
```bash
sudo apt update
sudo apt install build-essential libomp-dev python3-venv python3-pip
``` 
___
### Configuração do Ambiente Python (Venv)

Para gerar os gráficos sem interferir no Python do seu sistema, recomendamos o uso de um ambiente virtual.

1. Crie o ambiente virtual (apenas na primeira vez):
```bash
python3 -m venv venv
```

2. Ative o ambiente:
```bash
source venv/bin/activate
```

(Você verá (venv) no início do seu terminal).

3. Instale as dependências:
```bash
pip install pandas matplotlib seaborn
```
Nota: Sempre que for rodar o comando make plot, certifique-se de que o ambiente virtual está ativo (source venv/bin/activate).
___

### Como Executar

O projeto é controlado inteiramente via Makefile. Na raiz do projeto, execute os comandos na ordem abaixo:

1. Compilar

Compila os códigos fonte (src/) e gera os binários otimizados na pasta bin/.
```bash
make
```
2. Rodar Experimentos

Executa o script de automação (scripts/run.sh). Isso rodará as Tarefas A e D com diferentes parâmetros (N, K, Threads, Schedules).
```bash
make run
```
Nota: Os resultados serão salvos automaticamente na pasta results/ em formato CSV. Este passo pode levar alguns minutos.

3. Gerar Gráficos

Executa o script de visualização (scripts/plot.py). Ele lê os arquivos CSV e gera gráficos de desempenho na pasta plots/.
```bash
# Lembre-se de ativar o venv antes: source venv/bin/activate
make plot
```
4. Limpar 

Para remover os binários compilados, arquivos de resultados (.csv) e gráficos (.png):
```bash
make clean
```