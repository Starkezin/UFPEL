// src/seq/main.c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

// Função auxiliar para tempo em segundos
double get_time() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts.tv_sec + ts.tv_nsec / 1e9;
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Uso: %s <N>\n", argv[0]);
        return 1;
    }

    long long N = atoll(argv[1]);
    int repeticoes = 20; // Deve ser igual ao da versao OMP

    // Alocação
    double *A = (double*) malloc(N * sizeof(double));
    double *B = (double*) malloc(N * sizeof(double));

    if (!A || !B) {
        fprintf(stderr, "Erro de alocacao.\n");
        return 1;
    }

    // Inicialização
    for(long long i=0; i<N; i++) {
        A[i] = i * 0.1;
        B[i] = i * 0.2;
    }

    double start = get_time();

    // Loop de repetições (Carga de trabalho)
    for (int r = 0; r < repeticoes; r++) {
        // Parte 1
        for (long long i = 0; i < N; i++) {
            A[i] = B[i] * 1.01;
        }
        
        // Parte 2
        for (long long i = 0; i < N; i++) {
            B[i] = A[i] + 0.5;
        }
    }

    double end = get_time();

    // Saída CSV: N, Threads, Versao, TempoTotal
    // Threads = 1, Versao = 0 (Sequencial)
    printf("%lld,1,0,%f\n", N, end - start);

    free(A);
    free(B);
    return 0;
}