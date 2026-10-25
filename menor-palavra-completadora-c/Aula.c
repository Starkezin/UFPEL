void getString(int pos, int len, int c, char string[]){
    char substring;
    while (c < len) {
        substring = string[pos + c - 1];
        c++;
    }
    substring = '\0';
}
char * shortestCompletingWord(char * licensePlate, char ** words, int wordsSize){
    int qtdIgual, qtdAnt;
    for(int i=0; i < sizeof(licensePlate); i++){
        char tmp;
        int zero = 0;
        tmp = getString(i, 1, zero, licensePlate);
        for(int j=0; j < wordsSize; j++){
            for(int c=0; c < sizeof(words[j]); c++){
                char temp;
                temp = getString(c, 1, zero, words[j]);
                if(temp == tmp){
                    qtdIgual = qtdIgual + 1;
                }
            }
        qtdAnt = qtdIgual;
        }     
    }
    return qtdIgual;
}