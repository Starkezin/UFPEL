#include <stdio.h>

int fillCups(int* amount, int amountSize){
    int segundos = 0;
    while(amount[0] > 0 || amount[1] > 0 || amount[2] > 0){
        if(amount[0] >= amount[1] && amount[1] >= amount[2]){
            amount[0] = amount[0] - 1;
            amount[1] = amount[1] - 1;
            segundos++;
        }
        else if(amount[0] >= amount[1] && amount[2] >= amount[1]){
            amount[0] = amount[0] - 1;
            amount[2] = amount[2] - 1;
            segundos++;
        }
        else {
            amount[1] = amount[1] - 1;
            amount[2] = amount[2] - 1;
            segundos++;
        }
    }
    return segundos;
}

int main(){
    int amount[] = {7,5,0};
    int amountSize = sizeof(amount) / sizeof(amount[0]);
    int segundos = fillCups(amount, amountSize);
    printf("Demora %d segundos", segundos);
}