#include <stdio.h>
#include <string.h>
#include <ctype.h>

#define TAMANHO_ALFABETO 26

typedef struct {
    int count[TAMANHO_ALFABETO];
} contagemLetras;

void iniContadorLetras(contagemLetras *lc) {
    for (int i = 0; i < TAMANHO_ALFABETO; i++) {
        lc->count[i] = 0;
    }
}

void attContagemLetras(contagemLetras *lc, char *str) {
    for (int i = 0; str[i] != '\0'; i++) {
        char c = tolower(str[i]);
        if (isalpha(c)) {
            lc->count[c - 'a']++;
        }
    }
}

int palavraCompleta(contagemLetras *licenseCount, contagemLetras *qtdLetras) {
    for (int i = 0; i < TAMANHO_ALFABETO; i++) {
        if (licenseCount->count[i] > qtdLetras->count[i]) {
            return 0;
        }
    }
    return 1;
}

char *shortestCompletingWord(char *licensePlate, char **words, int wordsSize) {
    contagemLetras licenseCount;
    iniContadorLetras(&licenseCount);
    attContagemLetras(&licenseCount, licensePlate);

    char *result = NULL;
    int maisCurta = -1;

    for (int i = 0; i < wordsSize; i++) {
        contagemLetras qtdLetras;
        iniContadorLetras(&qtdLetras);
        attContagemLetras(&qtdLetras, words[i]);

        if (palavraCompleta(&licenseCount, &qtdLetras)) {
            int tamPalavra = strlen(words[i]);
            if (maisCurta == -1 || tamPalavra < maisCurta) {
                result = words[i];
                maisCurta = tamPalavra;
            }
        }
    }

    return result;
}

int main() {
    char *licensePlate = "aBc 12c";
    char *words[] = {"abccdef", "caaacab", "cbca", "abcde", "abcdefg"};
    int wordsSize = sizeof(words) / sizeof(words[0]);

    char *resultado = shortestCompletingWord(licensePlate, words, wordsSize);
    printf("Shortest completing word: %s\n", resultado);

    return 0;
}