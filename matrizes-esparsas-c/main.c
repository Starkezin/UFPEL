//SAMUEL STARKE E VITOR MONTEIRO COLOMBO
//TRABALHO AED

#include <stdio.h>
#include <stdlib.h>
#include <time.h>

typedef struct t_celula{
    struct t_celula *right, *below;
    int line, column;
    float valor;
}celula;

typedef struct{
    int m, n;
    celula *inicio, *endline, *endcolumn;
}matriz;

typedef matriz Matrix;

Matrix *matrix_create();
Matrix *create_esparse_matrix(int m, int n);
int startHeads(Matrix *mat);
int insertHeadLine(Matrix *mat);
int insertHeadColumn(Matrix *mat);
int matrix_setelem(Matrix *mat, int line, int column, float valor);
int matrix_print(Matrix *mat);
int matrix_destroy(Matrix *mat);
Matrix *matrix_add(Matrix *matA, Matrix *matB);
Matrix *matrix_multiply(Matrix *matA, Matrix *matB);
Matrix *matrix_transpose(Matrix *m);
float matrix_getelem(Matrix *mat, int line, int column);
Matrix* create_gigaMatrix(unsigned int tam);

int main( void ) {
    /* Inicializacao da aplicacao ... */
    Matrix *A = matrix_create();
    Matrix *B = matrix_create(); 
    matrix_print( A );
    matrix_print( B );
    Matrix *C = matrix_add( A, B ); 
    matrix_print( C );
    matrix_destroy( C );
    C = matrix_multiply( A, B ); 
    matrix_print( C );
    matrix_destroy( C );
    C = matrix_transpose( A ); 
    matrix_print( C );
    matrix_destroy( C );
    matrix_destroy( A );
    matrix_destroy( B );
    
    /*
    Matrix* A = create_gigaMatrix(2800);
    matrix_print( A );
    */
    return 0;
}

Matrix *matrix_create(){
    Matrix *matrix;
	matrix = (Matrix *)malloc(sizeof(Matrix));

    int m, n;
    printf("Digite a quantidade de linhas e de columns respectivamente: ");
    scanf("%d %d", &m, &n);

	if (!matrix || m<=0 || n<=0){
        return 0;
	}
	matrix->inicio = NULL;
	matrix->endline = NULL;
	matrix->endcolumn = NULL;
	matrix->m = m;
	matrix->n = n;

    int line, column;
    float valor;

    startHeads(matrix); 

    printf("Digite em qual linha, coluna e qual o valor deseja inserir respectivamente: ");
    scanf("%d %d %f", &line, &column, &valor);
    matrix_setelem(matrix, line, column, valor);
    while(valor != 0){
        printf("Digite em qual linha, coluna e qual o valor deseja inserir respectivamente: ");
        scanf("%d %d %f", &line, &column, &valor);
        matrix_setelem(matrix, line, column, valor);
    }
	return matrix;
}

Matrix *create_esparse_matrix(int m, int n){
    Matrix *matrix;
	matrix = (Matrix *)malloc(sizeof(Matrix));

	if (!matrix|| m<=0 || n<=0){
        return 0;
	}
	matrix->inicio = NULL;
	matrix->endline = NULL;
	matrix->endcolumn = NULL;
	matrix->m = m;
	matrix->n = n;

	startHeads(matrix);

	return matrix;
}

Matrix *matrix_add(Matrix *matA, Matrix *matB){
    int i, j;
    float soma;
    Matrix *matC;
    celula *pCelulaA, *pCelulaB;

    if (matA->m != matB->m || matA->n != matB->n){
        printf("\n Tamanhos incompativeis");
        return NULL;
    }

    if (!matA || !matB || !matA->m || !matA->n){
        return NULL;
    }

    matC = create_esparse_matrix(matA->m, matB->n);

    pCelulaA = matA->inicio->below;
    pCelulaB = matB->inicio->below;

    for (i = 1; i <= matA->m; i++){
        for (j = 1; j <= matA->n; j++){
            if (j == pCelulaA->right->column && j == pCelulaB->right->column){
                soma = pCelulaA->right->valor + pCelulaB->right->valor;
                if (soma){
                    matrix_setelem(matC, i, j, soma);
                }
                pCelulaA = pCelulaA->right;
                pCelulaB = pCelulaB->right;
            }
            else if (j == pCelulaA->right->column){
                matrix_setelem(matC, i, j, pCelulaA->right->valor);
                pCelulaA = pCelulaA->right;
            }
            else if (j == pCelulaB->right->column){
                matrix_setelem(matC, i, j, pCelulaB->right->valor);
                pCelulaB = pCelulaB->right;
            }
        }
        pCelulaA = pCelulaA->right->below;
        pCelulaB = pCelulaB->right->below;
    }

    return matC;
}

int startHeads(Matrix *mat){
    int i=0;

    celula *head;
    head = (celula*)malloc(sizeof(celula));

    if (!head){
        return 0; //erro
    }

    //celula principal
    head->column = -1;
    head->line = -1;

    mat->inicio = head;
    mat->endline = head;
    mat->endcolumn = head;

    for (i = 1; i <= mat->n; i++){ 
        insertHeadColumn(mat);
    }
    for (i = 1; i <= mat->m; i++){ 
        insertHeadLine(mat);
    }
    return 1;
}
int insertHeadColumn(Matrix *mat){ 
    celula *head;
    head = (celula*)malloc(sizeof(celula));

    if (!head){
        return 0; 
    }

    head->column = -1;
    head->line = 0;

    mat->endcolumn->right = head;
    mat->endcolumn = head;

    head->right = mat->inicio;
    head->below = head;
    return 1;
}

int insertHeadLine(Matrix *mat){
    celula *head;
    head = (celula*)malloc(sizeof(celula));

    if (!head){
        return 0;
    }

    head->column = 0;
    head->line = -1;

    mat->endline->below = head;
    mat->endline = head;

    head->below = mat->inicio;
    head->right = head;

    return 1;
}

int matrix_setelem(Matrix *mat, int line, int column, float valor){
    int i;

    if (!mat || mat->m <= 0 || mat->n <= 0 || !valor){ 
        return 0;
    }
    if (line>mat->m || column>mat->n || !valor || line < 1 || column < 1){
        return 0;
    }

    celula *pCelula;
    celula *pCelulacolumn; 
    celula *pCelulaline; 

    pCelula = (celula*)malloc(sizeof(celula));

    if (!pCelula){ 
        return 0;
    }

    pCelula->line = line;
    pCelula->column = column;
    pCelula->valor = valor;

    pCelulaline = mat->inicio->below; //pCelulaline aponta para cabeça das linhas
    pCelulacolumn = mat->inicio->right; //pCelulacolumn aponta para cabeça das colunas

    //Vai até a linha desejada
    for (i=0; i<line-1; i++){ 
        pCelulaline = pCelulaline->below; 
    }
    i=0;
    while (i<column && pCelulaline->right->line != -1){
        if (pCelulaline->right->column > pCelula->column){ //Verifica se a celula que devemos inserir esta a esquerda
            pCelula->right = pCelulaline->right;
            pCelulaline->right = pCelula;
        }
        else{
            pCelulaline = pCelulaline->right; //Pula uma celula para a direita
        }
        i++;
    }
    if (pCelulaline->right->line == -1){ 
        pCelula->right = pCelulaline->right; //Aponta para a cabeça
        pCelulaline->right = pCelula; 
    }

    //Vai ate a coluna desejada
    for (i = 0; i < column-1; i++){
        pCelulacolumn = pCelulacolumn->right; 
    }
    i=0;
    while (i<line && pCelulacolumn->below->column != -1){ //Busca pelas celulas a posição
        if (pCelulacolumn->below->line > pCelula->line){
            pCelula->below = pCelulacolumn->below;
            pCelulacolumn->below = pCelula;
        }
        else{
            pCelulacolumn = pCelulacolumn->below; 
        }
        i++;
    }
    if (pCelulacolumn->below->column == -1){
        pCelula->below = pCelulacolumn->below; //Aponta para a cabeça
        pCelulacolumn->below = pCelula; 
    }
    return 1;
}


int matrix_print(Matrix *mat){
    int i, j;
    celula *pCelula;

    if (!mat || !mat->m || !mat->n){
        return 0;
    }

    pCelula = mat->inicio->below;

    printf("\n");

    for (i = 1; i <= mat->m; i++){
        for (j = 1; j <= mat->n; j++){
            if (pCelula->right->line == i && pCelula->right->column == j){
                pCelula = pCelula->right;
                printf("  \t%0.2f   ", pCelula->valor);
            }
            else{
                printf("  \t%0.2f   ", 0.0F);
            }
        }
        printf("\n");
        pCelula = pCelula->right->below;
    }

    return 1;
}


int matrix_destroy(Matrix *mat){
    int i = 1;
    int j = 1;
    celula *pCelula, *aux;

    if (!mat || !mat->m || !mat->n){
        return 0;
    }

    pCelula = mat->inicio->below;
    pCelula = pCelula->right;

    for (i = 1; i <= mat->m; i++){
        for (j = 1; j <= mat->n; j++){
            if (pCelula->line == i && pCelula->column == j){
                aux = pCelula;
                pCelula = pCelula->right;
                free(aux);
            }
        }
        pCelula = pCelula->below->right;
    }

    
    pCelula = mat->inicio->right;

    for (i = 0; i < mat->n; i++){   //Apaga cabeças da coluna
        aux = pCelula;              //aux: variavel a ser destruida
        pCelula = pCelula->right;   //passa pCelula para a direita
        free(aux);
    }

    pCelula = mat->inicio->below;
    for (i = 0; i < mat->m; i++){ //Apaga cabeças da linha
        aux = pCelula;            //aux: variavel a ser destruida
        pCelula = pCelula->below; //passa a pCelula para baixo
        free(aux);
    }

    pCelula = mat->inicio;
    free(pCelula);
    
    //seta os ponteiros para nullo e as linhas e colunas em 0
    mat->endcolumn = mat->endline = mat->inicio = NULL;
    mat->m = mat->n = 0;
    mat = NULL;

    return 1;
}

Matrix *matrix_multiply(Matrix *matA, Matrix *matB){
    int i=0, j=0;
    float total;
    Matrix *matC;

    if (matA->n != matB->m){ //so pode multiplicar se o numero de colunas de A eh igual ao numero de linhas de B
        printf("\nTamanhos incompativeis");
        return NULL;
    }

    if (!matA || !matB || !matA->m || !matA->n || !matB->n){
        return NULL;
    }

    matC = create_esparse_matrix(matA->m, matB->n); //C é formada pelo numero de linhas de A e de colunas de B

    for (i = 1; i <= matA->m; i++){
        for (j = 1; j <= matB->n; j++){
            total = 0;
            total += matrix_getelem(matA,i,j) * matrix_getelem(matB,i,j);
            if (total != 0){
                matrix_setelem(matC,i,j,total);
            }
        }
    }
    return matC;
}

float matrix_getelem(Matrix *mat, int line, int column){
    celula *pCelula;
    int i = 0;

    pCelula = mat->inicio->right;

    for (i = 0; i < column-1; i++){ //Vai ate a coluna do elemento a ser buscado
        pCelula = pCelula->right;
    }

    do{
        pCelula = pCelula->below;
        if (pCelula->line == line){
            return pCelula->valor;
        }
    }while(pCelula->column != -1);

    return 0; 
}
Matrix* matrix_transpose(Matrix* m) {
    int i=0, j=0;
    float aux;
    Matrix *matC;

    if (!m || !m->m || !m->n){
        return NULL;
    }

    matC = create_esparse_matrix(m->m,m->n); //Cria matriz que recebera a matriz transposta
    for (i = 1; i <= m->m; i++){
        for (j = 1; j <= m->n; j++){
            aux = 0; 
            aux = matrix_getelem(m,i,j); //Passa o valor da linha I e coluna J para a variavel aux;
            if (aux != 0){ 
                matrix_setelem(matC,j,i,aux); //Recebe o valor de aux na linha J e coluna I;
            }
        }
    }
    return matC;
}

Matrix* create_gigaMatrix(unsigned int tam){
    Matrix *matrix;
	matrix = (Matrix *)malloc(sizeof(Matrix));
    unsigned int m, n;
    m = tam;
    n = tam;

	if (!matrix|| m<=0 || n<=0){
        return 0;
	}

	matrix->inicio = NULL;
	matrix->endline = NULL;
	matrix->endcolumn = NULL;
	matrix->m = m;
	matrix->n = n;

	startHeads(matrix);

    int aux;
    srand(time(NULL));
    for(unsigned int i = 0; i < tam; i++){
        for(unsigned int j = 0; j < tam; j++){
            aux = rand() % 11;
            if(aux < 6) matrix_setelem(matrix, i, j, 0);
            else matrix_setelem(matrix, i, j, 1 + (rand() % 9));
        }
    }

	return matrix;
}