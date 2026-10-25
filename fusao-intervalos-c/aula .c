int** merge(int** intervals, int intervalsSize, int* intervalsColSize, int* returnSize, int** returnColumnSizes){
    int** result[returnSize][**returnColumnSizes];
    int aux;
    int countl = 0, countc = 0;
    for(int i = 0; i < intervalsSize; i++){
        for(int c = 1; i < *intervalsColSize; c++){
            /*
            Acredito que usando um algoritmo semelhante ou ao InsertionSort, onde eu percorro todos os elemntos
            e verifico se é menor que os ultimos ja apresentados.
            Mas nesse caso invés de ordenar ele, eu apenas colocaria o valor em uma nova matriz.

			Não consegui desenvolver nenhum codigo na aula pois tive dificuldades com os ponteiros duplos;
            */
        }
    }

   return intervals;
}