#!/bin/bash

# Garante que a pasta de resultados existe
mkdir -p results

# ==============================================================================
#                                  TAREFA A
# ==============================================================================

# Ajuste de caminhos (Fundamental para execução com o novo Makefile)
OUTPUT="results/resultados_tarefaA.csv"
EXEC_OMP="./bin/tarefaA"
EXEC_SEQ="./bin/tarefaA_seq"

echo "N,K,Threads,Schedule,Chunk,Tempo" > $OUTPUT

LISTA_N=(100000 500000 1000000)

LISTA_K=(20 24 28)

LISTA_THREADS=(1 2 4 8 16)

LISTA_SCHED=("static" "dynamic" "guided")

LISTA_CHUNKS=(1 4 16 64)

REPETICOES=5

echo "Iniciando testes..."

echo "=== 1. Rodando Sequencial (Baseline) ==="
for n in "${LISTA_N[@]}"; do
    for k in "${LISTA_K[@]}"; do
        echo "Sequencial: N=$n K=$k"
        for ((i=1; i<=REPETICOES; i++)); do
            $EXEC_SEQ $n $k >> $OUTPUT
        done
    done
done

echo "=== 2. Rodando OpenMP ==="
echo "Saída será gravada em: $OUTPUT"

for n in "${LISTA_N[@]}"; do
    for k in "${LISTA_K[@]}"; do
        for th in "${LISTA_THREADS[@]}"; do
            for sched in "${LISTA_SCHED[@]}"; do
                for chunk in "${LISTA_CHUNKS[@]}"; do
                    
                    echo "Rodando: N=$n K=$k Threads=$th Sched=$sched Chunk=$chunk"
                    
                    for ((i=1; i<=REPETICOES; i++)); do
                        $EXEC_OMP $n $k $th $sched $chunk >> $OUTPUT
                    done
                    
                done
            done
        done
    done
done

# ==============================================================================
#                                  TAREFA D
# ==============================================================================

OUTPUT_D="results/resultados_tarefaD.csv"
EXEC_OMP_D="./bin/tarefaD"
EXEC_SEQ_D="./bin/tarefaD_seq"

# Cabeçalho CSV da Tarefa D
echo "N,Threads,Version,Time" > $OUTPUT_D

# Listas específicas para D (Overhead vs Throughput)
LISTA_N_D=(100000 1000000 10000000 50000000)
LISTA_THREADS_D=(1 2 4 8 16)
LISTA_VERSOES=(1 2) # 1=Ingenua, 2=Arrumada

echo ""
echo "=== 3. Iniciando Tarefa D ==="

echo "--- D: Rodando Sequencial ---"
for n in "${LISTA_N_D[@]}"; do
    echo "D Sequencial: N=$n"
    for ((i=1; i<=REPETICOES; i++)); do
        $EXEC_SEQ_D $n >> $OUTPUT_D
    done
done

echo "--- D: Rodando OpenMP ---"
echo "Saída será gravada em: $OUTPUT_D"

for n in "${LISTA_N_D[@]}"; do
    for th in "${LISTA_THREADS_D[@]}"; do
        # Precisamos exportar a var de ambiente para o código pegar o num_threads
        export OMP_NUM_THREADS=$th
        
        for v in "${LISTA_VERSOES[@]}"; do
            echo "D Rodando: N=$n Threads=$th Versao=$v"
            for ((i=1; i<=REPETICOES; i++)); do
                $EXEC_OMP_D $n $v >> $OUTPUT_D
            done
        done
    done
done

echo "Concluído Geral! Dados salvos em results/, use 'make plot' para gerar os gráficos"