#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main(){
int resp, maquina, pessoa;
do {
    srand(time(NULL));
    maquina = rand() % 3 + 1;
// Empate
    printf("Digite sua escolha:\n 1 Para Pedra\n 2 Para Papel\n 3 Para Tesoura\n 4 Para Lagarto\n 5 Para Spock\n");
    scanf(" %d", &pessoa);
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
    else if (pessoa == 4 && maquina == 4){
        printf("Empate\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 5 && maquina == 5){
        printf("Empate\n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); } 
// Empate 

// Win Maquina
    else if (pessoa == 1 && maquina == 2){
        printf("Maquina venceu\n Papel cobre a pedra \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 3){
        printf("Maquina venceu\n Tesoura picota o papel \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 1){
        printf("Maquina venceu\n Pedra quebra a tesoura \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 4 && maquina == 1){
        printf("Maquina venceu\n Pedra esmaga o lagarto \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 5 && maquina == 4){
        printf("Maquina venceu\n Lagarto envenena o spock \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 5){
        printf("Maquina venceu\n Spock derrete a tesoura \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 5 && maquina == 2){
        printf("Maquina venceu\n Papel refuta o spock \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n"); 
        scanf("%d", &resp); }
    else if (pessoa == 1 && maquina == 5){
        printf("Maquina venceu\n Spock vaporiza a pedra \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 4){
        printf("Maquina venceu\n Lagarto come o papel \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 4 && maquina == 3){
        printf("Maquina venceu\n Tesoura decapita o lagarto \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina);
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
// Win Maquina

// Win Pessoa
    else if (pessoa == 1 && maquina == 3){
        printf("Voce venceu\n Pedra quebra a tesoura  \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 1){
        printf("Voce venceu\n Papel cobre a pedra  \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 2){
        printf("Voce venceu\n Tesoura picota o papel \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 1 && maquina == 4){
        printf("Voce venceu\n Pedra esmaga o lagarto \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 4 && maquina == 5){
        printf("Voce venceu\n Lagarto envenena o spock \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 5 && maquina == 3){
        printf("Voce venceu\n spock derrete a tesoura \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 3 && maquina == 4){
        printf("Voce venceu\n Tesoura decapita o lagarto \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 4 && maquina == 2){
        printf("Voce venceu\n Lagarto come o papel \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 2 && maquina == 5){
        printf("Voce venceu\n Papel refuta spock \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    else if (pessoa == 5 && maquina == 1){
        printf("Voce venceu\n Spock vaporiza pedra \n Pessoa:%d\n Maquina:%d\n", pessoa, maquina); 
        printf("Deseja continuar? 1 para cotinuar e 2 parar\n");
        scanf("%d", &resp); }
    }while (resp == 1);
    return 0;
} 