#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <omp.h>

// Fibonacci, onde o custo cresce exponencialmente com n
long long fib(int n) {
    if (n < 2) return n;
    return fib(n - 1) + fib(n - 2);
}

int main(int argc, char *argv[]) {
    // Validação dos argumentos
    if (argc < 6) {
        printf("Uso: %s <N> <K> <threads> <sched_type> <chunk>\n", argv[0]);
        printf("  sched_type: static, dynamic, guided\n");
        return 1;
    }

    int N = atoi(argv[1]);           // Tamanho do vetor
    int K = atoi(argv[2]);           // Parâmetro de carga (fib(i % K))
    int num_threads = atoi(argv[3]); // Número de threads
    char *sched_str = argv[4];       // Tipo de schedule (texto)
    int chunk_size = atoi(argv[5]);  // Tamanho do chunk

    // Alocação do vetor de resultados
    long long *v = (long long *)malloc(N * sizeof(long long));
    if (v == NULL) {
        perror("Erro ao alocar memoria");
        return 1;
    }

    // Configuração do Ambiente OpenMP
    omp_set_num_threads(num_threads);

    omp_sched_t kind;
    if (strcmp(sched_str, "static") == 0) {
        kind = omp_sched_static;
    } else if (strcmp(sched_str, "dynamic") == 0) {
        kind = omp_sched_dynamic;
    } else if (strcmp(sched_str, "guided") == 0) {
        kind = omp_sched_guided;
    } else {

        kind = omp_sched_auto; 
    }

    // Define o schedule que será usado na cláusula 'runtime'
    omp_set_schedule(kind, chunk_size);

    // Execução do Kernel (Medição de Tempo)
    double start_time = omp_get_wtime();

    #pragma omp parallel for schedule(runtime)
    for (int i = 0; i < N; i++) {
        // i % K garante que o valor passado para fib esteja entre 0 e K-1
        v[i] = fib(i % K);
    }

    double end_time = omp_get_wtime();
    double time_taken = end_time - start_time;

    // Formato: N, K, Threads, Schedule, Chunk, Tempo
    printf("%d,%d,%d,%s,%d,%.6f\n", N, K, num_threads, sched_str, chunk_size, time_taken);

    free(v);

    return 0;
}