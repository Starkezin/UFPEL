# TSP (Brute Force & Nearest Neighbor) em C++

- **Brute force** (força bruta) — gera todas as permutações para encontrar a solução ótima (usado apenas em instâncias muito pequenas).
- **Nearest Neighbor** (Vizinho Mais Próximo) — heurística gulosa que encontra uma solução rápida e aproximada.
- **Double Tree** - Algoritmo aproximativo para o TSP euclidiano que gera um ciclo hamiltoniano a partir da MST do grafo.
Garante solução com custo no máximo 2× o do ciclo ótimo.

---
## Estrutura e arquivos

- `tsp_bruteforce.cpp` — implementação por força bruta (testa todas as permutações). 
- `tsp_nearest_neighbor.cpp` — implementação da heurística Nearest Neighbor.
- `tsp_double_tree.cpp` - Implementação do método aproximativo double tree.

---
## Formato do arquivo de entrada

O arquivo deve conter a **matriz de adjacência NxN** de distâncias, uma linha por linha, valores inteiros separados por espaços.

--- 

## Compilando (exemplos com g++)

```bash
# Compilar brute force
g++ tsp_bruteforce.cpp -o tsp_bf

# Compilar nearest neighbor
g++ tsp_nearest_neighbor.cpp -o tsp_nn

#Compilar aproximativo (double-tree)
g++ tsp_double_tree.cpp -o tsp_dt
```

---

## Executando
**IMPORTANTE:** passe o arquivo de entrada como argumento da linha de comando.
Exemplos:

```bash
# Executar brute force
./tsp_bf data/tsp1_253.txt

# Executar Nearest Neighbor
./tsp_nn data/tsp1_253.txt

#Executar aproximativo (double-tree)
./tsp_dt data/tsp1_253.txt
```
---

## Saídas esperadas

- **Brute force:** imprimirá a menor distância encontrada (o custo ótimo) e o tempo de execução (ms).
- **Nearest Neighbor:** imprimirá a rota encontrada (começando pela cidade inicial) e a distância total aproximada, além do tempo de execução (ms).
- **Double tree:** Imprimirá o tempo de execução (ms), o caminho aproximado e o custo total aproximado.

Exemplo de saída (nearest neighbor):
```
Tempo de execucao: 0.004 ms
Rota encontrada (Nearest Neighbor): 0 8 7 2
Distancia total aproximada: 299
```

Exemplo de saída (aproximativo):
```
Tempo de execucao: 0.104 ms
Rota encontrada (Nearest Neighbor): 0 1 5 4
Distancia total aproximada: 34902
```
---


