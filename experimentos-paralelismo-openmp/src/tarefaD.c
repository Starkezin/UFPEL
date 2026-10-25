#include <stdio.h>
#include <stdlib.h>
#include <omp.h>

int main(int argc, char *argv[]) {
    if (argc < 3) {
        fprintf(stderr, "Uso: %s <N> <Versao [1=Ingenua, 2=Arrumada]>\n", argv[0]);
        return 1;
    }
    
    int repeticoes = 50;

    int N = atoi(argv[1]);
    int version = atoi(argv[2]);
    int n_threads = omp_get_max_threads(); // Pega o que foi setado no export OMP_NUM_THREADS

    // Alocação
    double *A = (double*) malloc(N * sizeof(double));
    double *B = (double*) malloc(N * sizeof(double));

    // Inicialização (fora da medição)
    #pragma omp parallel for
    for(int i=0; i<N; i++) {
        A[i] = i * 0.1;
        B[i] = i * 0.2;
    }

    double start, end;
    
    // --- INICIO DA MEDIÇÃO ---
    start = omp_get_wtime();

    if (version == 1) {
        // --- V1: INGÊNUA (Fork-Join repetitivo) ---
        // Cria e destroi threads 2 * repeticoes vezes
        for (int r = 0; r < repeticoes; r++) {
            
            // Região 1: Cria threads -> Executa -> Destroi (ou estaciona)
            #pragma omp parallel for schedule(static)
            for (int i = 0; i < N; i++) {
                A[i] = B[i] * 1.01; // Operação leve
            }

            // Região 2: Cria threads DE NOVO -> Executa -> Destroi
            #pragma omp parallel for schedule(static)
            for (int i = 0; i < N; i++) {
                B[i] = A[i] + 0.5;
            }
        }

    } else {
        // --- V2: ARRUMADA (Uma região, barreira interna) ---
        // Cria threads APENAS 1 * repeticoes vezes
        for (int r = 0; r < repeticoes; r++) {
            
            #pragma omp parallel
            {
                // As threads nascem aqui e ficam vivas
                
                #pragma omp for schedule(static)
                for (int i = 0; i < N; i++) {
                    A[i] = B[i] * 1.01;
                }
                // Barreira implícita aqui (threads esperam, mas não morrem)

                // Reutilização de cache: a mesma thread provavelmente pega o mesmo pedaço
                #pragma omp for schedule(static)
                for (int i = 0; i < N; i++) {
                    B[i] = A[i] + 0.5;
                }
            } // Threads morrem aqui
        }
    }

    end = omp_get_wtime();
    // --- FIM DA MEDIÇÃO ---

    // Saída CSV: N, Threads, Versao, TempoTotal
    printf("%d,%d,%d,%f\n", N, n_threads, version, end - start);

    free(A);
    free(B);
    return 0;
}