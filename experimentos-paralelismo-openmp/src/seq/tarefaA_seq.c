#include <stdio.h>
#include <stdlib.h>
#include <time.h> // Usamos time.h ao invés de omp.h

// Função auxiliar para medir tempo em segundos (POSIX)
double get_time() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts.tv_sec + ts.tv_nsec / 1e9;
}

// Fibonacci recursivo (carga de trabalho)
long long fib(int n) {
    if (n < 2) return n;
    return fib(n - 1) + fib(n - 2);
}

int main(int argc, char *argv[]) {
    // Validação dos argumentos
    // Espera: <N> <K> (Threads, Schedule e Chunk são ignorados no seq)
    if (argc < 3) {
        fprintf(stderr, "Uso: %s <N> <K>\n", argv[0]);
        return 1;
    }

    int N = atoi(argv[1]);
    int K = atoi(argv[2]);

    // Alocação
    long long *v = (long long *)malloc(N * sizeof(long long));
    if (v == NULL) {
        perror("Erro ao alocar memoria");
        return 1;
    }

    // Medição de Tempo (Nativa do C)
    double start_time = get_time();

    // Loop Sequencial
    for (int i = 0; i < N; i++) {
        v[i] = fib(i % K);
    }

    double end_time = get_time();
    double time_taken = end_time - start_time;

    // Formato de saída compatível com o CSV:
    // N, K, Threads=1, Schedule=SEQ, Chunk=0, Tempo
    printf("%d,%d,1,SEQ,0,%.6f\n", N, K, time_taken);

    free(v);
    return 0;
}