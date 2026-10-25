int fillCups(int* amount, int amountSize){
    int c; 
    float i = 0;
    int tipos = 0;


    for(int count = 0; count < amountSize; count++){
        if(amount[count] > 0){
            tipos++;
        }
    }

    for(c = 0;c < amountSize  ; c++){

        while(amount[c] > 0){
            amount[c] = amount[c] - 1;
            i++;
        }
    }
    return i;
}

int main(){
    int amount[] = {7,5,0};
    int amountSize = sizeof(amount) / sizeof(amount[0]);
    int segundos = fillCups(amount, amountSize);
    printf("Demora %d segundos", segundos);
}
