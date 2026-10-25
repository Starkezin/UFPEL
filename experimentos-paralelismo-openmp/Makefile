CC = gcc
# Flags
CFLAGS_OMP = -O3 -Wall -fopenmp
CFLAGS_SEQ = -O3 -Wall

# Diretórios
SRC_OMP = src/omp
SRC_SEQ = src/seq
BIN_DIR = bin
SCRIPTS_DIR = scripts
RESULTS_DIR = results

# --- ALVOS GERAIS ---
all: directories seq omp

# Cria os diretórios bin e results se não existirem
directories:
	mkdir -p $(BIN_DIR) $(RESULTS_DIR)

# Compila versões paralelas
omp: tarefaA tarefaD

# Compila versões sequenciais
seq: tarefaA_seq tarefaD_seq

# --- TAREFA A ---
tarefaA: directories $(SRC_OMP)/tarefaA.c
	$(CC) $(CFLAGS_OMP) -o $(BIN_DIR)/tarefaA $(SRC_OMP)/tarefaA.c

tarefaA_seq: directories $(SRC_SEQ)/tarefaA_seq.c
	$(CC) $(CFLAGS_SEQ) -o $(BIN_DIR)/tarefaA_seq $(SRC_SEQ)/tarefaA_seq.c

# --- TAREFA D ---
tarefaD: directories $(SRC_OMP)/tarefaD.c
	$(CC) $(CFLAGS_OMP) -o $(BIN_DIR)/tarefaD $(SRC_OMP)/tarefaD.c

tarefaD_seq: directories $(SRC_SEQ)/tarefaD_seq.c
	$(CC) $(CFLAGS_SEQ) -o $(BIN_DIR)/tarefaD_seq $(SRC_SEQ)/tarefaD_seq.c

# --- EXECUÇÃO (RUN) ---
# Executa o script único que roda A e D sequencialmente
run: all
	chmod +x scripts/run.sh
	./scripts/run.sh

# --- PLOTAGEM ---
# Gera os gráficos para ambas as tarefas
plot: plot_A plot_D

plot_A:
	python3 $(SCRIPTS_DIR)/plot_tarefaA.py

plot_D:
	python3 $(SCRIPTS_DIR)/plot_tarefaD.py

# --- LIMPEZA ---
clean:
	rm -rf $(BIN_DIR)/*
	rm -rf $(RESULTS_DIR)/*
	rm -f *.png plots/*.png