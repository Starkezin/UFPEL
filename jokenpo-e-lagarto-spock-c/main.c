#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main(){
int resp, maquina, pessoa;
do {
    srand(time(NULL));
    maquina = rand() % 3 + 1;
// Empate
    printf("Digite sua escolha:\n 1 Para Pedra\n 2 Para Papel\n 3 Para Tesoura\n");
    scanf("%d", &pessoa);
    if (pessoa == 1 && maquina == 1) {
        printf("Empate\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 2){
        printf("Empate\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 3){
        printf("Empate\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
// Empate 
// Win Maquina
    else if (pessoa == 1 && maquina == 2){
        printf("Maquina venceu\n Papel cobre pedra\n Pessoa:%d\n Maquina:%d", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 3){
        printf("Maquina venceu\n Tesoura picota o papel\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 1){
        printf("Maquina venceu\n Pedra quebra a tesoura\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
// Win Maquina
// Win Pessoa
    else if (pessoa == 1 && maquina == 3){
        printf("Voce venceu\n Pedra quebra a tesoura\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 1){
        printf("Voce venceu\n Papel cobre a pedra\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 2){
        printf("Voce venceu\n Tesoura picota o papel\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    }while (resp == 1);
    return 0;
}