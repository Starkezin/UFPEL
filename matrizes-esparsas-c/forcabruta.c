#include <stdio.h>
#include <unistd.h>
#include <time.h>
#include <malloc.h>
#define TAM 5120

int main(){
    int **mat = malloc(TAM  * sizeof(int*));
    int **mat2 = malloc(TAM * sizeof(int*));
    int **mat3 = malloc(TAM * sizeof(int*));
    for(int i = 0 ; i < TAM ; i++){
        mat[i] = malloc(TAM * sizeof(int));
    }
    
    for(int i = 0 ; i < TAM ; i++){
        mat2[i] = malloc(TAM * sizeof(int));
    }
    for(int i = 0 ; i < TAM ; i++){
        mat3[i] = malloc(TAM * sizeof(int));
    }
    srand(time(NULL));
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){
            if(rand() % 10 < 6)
            mat[i][j] = 0;
            else
            mat[i][j] = rand() % 10;
        }
    }
    /*
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){
            printf("%d\t",mat[i][j]);
        }
        printf("\n");
    }
    printf("\n");
    */
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){
            if(rand() % 10 < 8)
            mat2[i][j] = 0;
            else
            mat2[i][j] = rand() % 10;
        }
    }
    /*
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){
            printf("%d\t",mat2[i][j]);
        }
        printf("\n");
    }
    printf("\n");
    */
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){  
            mat3[i][j] = mat[i][j] * mat2[i][j];
        }
    }
    sleep(10);
    for (int i = 0; i < TAM; i++) {
        free(mat[i]);
    }
    sleep(4);
    for (int i = 0; i < TAM; i++) {
        free(mat2[i]);
    }
    sleep(4);
    for (int i = 0; i < TAM; i++) {
        free(mat3[i]);
    }
    sleep(4);
    free(mat);
    free(mat2);
    free(mat3);
    /*
    for(int i = 0; i < TAM ; i++){
        for(int j = 0 ; j < TAM ; j++){
            printf("%d\t",mat3[i][j]);
        }
        printf("\n");
    }
    */ 
    return 0;

}