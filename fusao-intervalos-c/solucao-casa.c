int comparator (const void *a, const void *b){
    return ( *(int*)a - *(int*)b );
}

int** merge(int** intervals, int intervalsSize, int* intervalsColSize, int* returnSize, int** returnColumnSizes){
    int* aux = NULL;
    int count = 0;

    if(intervals == 0) return NULL;

    qsort(intervals, intervalsSize, sizeof(int*), comparator);

    aux = intervals[0];
    for( int i = 1; i < intervalsSize; i++){
        if ( aux[1] >= intervals[i][0] ){
            aux[1] = (aux[1] > intervals[i][1])?aux[1]:intervals[i][1];
        }
        else{
            intervals[count][0] = aux[0];
            intervals[count][1] = aux[1];
            count++;
            aux = intervals[i];
        }
    }
    intervals[count][0] = aux[0];
    intervals[count][1] = aux[1];
    count++;

    *returnSize = count;
    (*returnColumnSizes) = (int*)malloc(count*sizeof(int));
    for(int i = 0; i < count; i++){
        (*returnColumnSizes)[i] = 2;
    }

    return intervals;
}