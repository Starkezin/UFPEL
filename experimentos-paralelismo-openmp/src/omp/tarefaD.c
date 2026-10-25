// src/omp/main.c
#include <stdio.h>
#include <stdlib.h>
#include <omp.h>

int main(int argc, char *argv[]) {
    if (argc < 3) {
        fprintf(stderr, "Uso: %s <N> <Versao [1=Ingenua, 2=Arrumada]>\n", argv[0]);
        return 1;
    }
    
    long long N = atoll(argv[1]);
    int version = atoi(argv[2]);
    int repeticoes = 20; // Mesmo valor do sequencial
    
    // Define numero de threads
    int n_threads = omp_get_max_threads();

    // Alocação
    double *A = (double*) malloc(N * sizeof(double));
    double *B = (double*) malloc(N * sizeof(double));

    if (!A || !B) {
        fprintf(stderr, "Erro de alocacao.\n");
        return 1;
    }

    // Inicialização (First touch para correta alocação de páginas na NUMA)
    #pragma omp parallel for schedule(static)
    for(long long i=0; i<N; i++) {
        A[i] = i * 0.1;
        B[i] = i * 0.2;
    }

    double start = omp_get_wtime();

    if (version == 1) {
        // --- V1: INGÊNUA ---
        // Cria e destroi threads repetidamente (Alto Overhead)
        for (int r = 0; r < repeticoes; r++) {
            
            #pragma omp parallel for schedule(static)
            for (long long i = 0; i < N; i++) {
                A[i] = B[i] * 1.01;
            }

            #pragma omp parallel for schedule(static)
            for (long long i = 0; i < N; i++) {
                B[i] = A[i] + 0.5;
            }
        }

    } else {
        // --- V2: ARRUMADA ---
        // Cria threads apenas uma vez (Baixo Overhead)
        for (int r = 0; r < repeticoes; r++) {
            
            #pragma omp parallel
            {
                // Threads já existem aqui dentro
                
                #pragma omp for schedule(static)
                for (long long i = 0; i < N; i++) {
                    A[i] = B[i] * 1.01;
                }
                // O 'omp for' sem 'nowait' possui uma barreira implícita aqui.
                // Isso garante que todo A[] esteja atualizado antes de B[] ser calculado.

                #pragma omp for schedule(static)
                for (long long i = 0; i < N; i++) {
                    B[i] = A[i] + 0.5;
                }
            } // Fim da região paralela (join apenas no final do loop interno)
        }
    }

    double end = omp_get_wtime();

    // Saída CSV: N, Threads, Versao, TempoTotal
    printf("%lld,%d,%d,%f\n", N, n_threads, version, end - start);

    free(A);
    free(B);
    return 0;
}