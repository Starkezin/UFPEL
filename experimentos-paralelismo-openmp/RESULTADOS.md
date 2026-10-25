# TAREFA A

Este trabalho avalia o impacto de diferentes políticas de escalonamento (static, dynamic, guided) em um laço irregular paralelizado com OpenMP.

O kernel estudado consiste no cálculo de fib(i % K) sem memorização, o que gera custo computacional irregular e crescente conforme o parâmetro K.

O kernel da Tarefa A não requer sincronização explícita, pois cada thread escreve em posições distintas do vetor de saída (v[i]). Assim, não há uso de atomic, critical ou agregação local nesse kernel.

Para o SIMD funcionar o kernel não pode conter recursão, mas o v[i] = fib(i % K); contém fib que apresentra recursão:
```
long long fib(int n) {
    if (n < 2) return n;
    return fib(n - 1) + fib(n - 2);
}
```

Variáveis:
```
Tamanho do problema (N)

Granularidade da carga (K)

Número de threads

Política de escalonamento

Tamanho de chunk
```
Cada configuração foi executada 5 vezes.

## Irregularidade do Kernel

O kernel principal é definido por:
```
v[i] = fib(i % K);
```
O custo cresce exponencialmente com o valor de i % K, gerando um desbalanceamento de carga, ou seja, alta variação no tempo de execução das iterações.

## Paralelização

O laço foi paralelizado utilizando:
```
#pragma omp parallel for schedule(runtime)
```

A política de escalonamento (static, dynamic, guided) e o tamanho de chunk foram configurados dinamicamente via omp_set_schedule. O número de threads foi controlado por omp_set_num_threads.

## Speedup
### Carga leve (K = 20)

Figura: speedup_N1000000_K20.png

Este caso foi escolhido pois foi o que demonstrou de forma mais clara a vantagem do dynamic e o guided em relação ao static.

Para cargas leves e granularidade fina, a redistribuição dinâmica das iterações compensa o overhead adicional.

O speedup cresce até aproximadamente 8 threads e depois satura. Esse desvio em relação ao speedup ideal indica um limite do sistema, pois o processador atingiu seu limite físico de threads e sofre de sobrecarga.

### Carga pesada (K = 28)

Figura: speedup_N1000000_K28.png

No caso mais pesado, as políticas apresentam um desempenho muito semelhante. Isso nos leva a concluir que o custo computacional por iteração domina o overhead de escalonamento, tornando as vantagens de dynamic e guided ínfimas.

## Tempo Absoluto
A figura tempo_fixo_chunk16_N1000000_K28.png mostra o tempo absoluto com chunk = 16 pois é o chunk com melhor precisão (com chunk = 8 também serve). Novamente obteve-se o resultado esperado, onde todos os schedules convergem para tempos semelhantes com uma limitação após 8 threads. Além disso, foi adicionado o tempo sequencial para mostrar o benefício do paralelismo em relação ao sequencial.

## Tamanho do Chunk
Figura: impacto_chunk_N1000000_K28.png

A figura mostra que o aumento do tamanho do chunk tem impacto secundário para os escalonamentos dynamic e guided, enquanto exerce influência perceptível para static ao variar de 1 para 4. Isso ocorre porque o static realiza a distribuição das iterações de forma fixa e pré-determinada, de modo que chunks muito pequenos aumentam o overhead de controle do laço. Ao aumentar o tamanho do chunk, esse overhead é reduzido. A partir de chunk = 4, novos aumentos não produzem ganhos relevantes, pois o overhead já se torna desprezível e o desempenho passa a ser dominado pelo custo computacional do kernel.

Tanto dynamic quanto guided realizam a distribuição das iterações em tempo de execução, o que torna o desempenho menos sensível ao tamanho do chunk. No caso do guided, a redução progressiva do chunk contribui adicionalmente para limitar o overhead inicial.


## TAREFA D: Organização de Região Paralela (Overhead)

Esta tarefa analisa o impacto da criação e destruição de threads (*Fork-Join*) versus o reaproveitamento de threads em uma região paralela única.

Foram implementadas duas versões para processar dois laços consecutivos de operações vetoriais:
1.  **V1 (Ingênua):** Abre e fecha regiões paralelas (`#pragma omp parallel for`) múltiplas vezes dentro de um loop de repetição.
2.  **V2 (Arrumada):** Abre a região paralela (`#pragma omp parallel`) apenas uma vez e utiliza diretivas de distribuição de trabalho internas (`#pragma omp for`), mantendo as threads vivas entre as iterações através de barreiras implícitas.

### Análise de Overhead (N Pequeno)

**Gráfico de Referência:** `tarefaD_overhead_tempos.png` (N=100.000)

Para cargas de trabalho pequenas ($N=10^5$), o tempo de computação útil é ínfimo, tornando o custo de gerenciamento de threads o fator dominante no tempo total.

1.  **Escalabilidade Negativa:** Observa-se que aumentar o número de threads (de 1 para 16) **aumenta** o tempo de execução em ambas as versões. Isso ocorre porque o custo de coordenar 16 threads supera o ganho de paralelizar um vetor tão pequeno.
2.  **Impacto do Fork/Join na V1 (Ingênua):** A linha vermelha tracejada (V1) apresenta um aumento drástico de tempo com 16 threads (chegando a ~0.008s).
    * **Causa:** A V1 realiza operações de *fork* (criar/acordar threads) e *join* (destruir/dormir threads) repetidamente a cada iteração externa. Com 16 threads, a disputa pelo escalonador do Sistema Operacional e o custo de troca de contexto tornam-se proibitivos.
3.  **Vantagem da V2 (Arrumada):** A linha verde (V2) mantém-se consistentemente abaixo da V1. Embora também sofra com overhead ao aumentar threads, ela é mais eficiente por realizar o *fork* apenas uma vez. As sincronizações internas são muito mais leves que a recriação do time de threads.
4.  **Sequencial vs Paralelo:** Para este tamanho de problema, a versão **Sequencial (linha azul)** é mais rápida que qualquer versão paralela com múltiplas threads, comprovando que paralelizar tarefas triviais sem carga suficiente introduz prejuízo de desempenho.

### Análise de Speedup da Otimização (V1 vs V2)

**Gráfico de Referência:** `tarefaD_speedup_otimizacao.png`

Este gráfico mostra a razão $Tempo(V1) / Tempo(V2)$. Valores acima de 1.0 indicam que a versão otimizada (V2) foi mais rápida.

1.  **Alto Ganho em N Pequeno ($10^5$ - Roxo):**
    * O speedup atinge picos de 1.6x (60% de ganho). Isso confirma que, quando o processamento é leve, a escolha correta da região paralela (V2) oferece ganhos expressivos apenas por eliminar chamadas de sistema desnecessárias.
2.  **Diluição em N Grande ($5 \times 10^7$ - Amarelo):**
    * Para vetores muito grandes, a linha amarela aproxima-se de 1.0 (linha horizontal cinza).
    * **Causa:** Conforme prevê a Lei de Amdahl, quando $N$ é grande, o tempo gasto calculando as operações matemáticas (tempo útil) é de vários segundos. O tempo economizado não recriando threads (milissegundos) torna-se uma fração desprezível do total.
3.  **Conclusão de Design:**
    * A otimização de "Região Paralela Única" é **crítica para tarefas de granulação fina** ou arquiteturas onde o *launch overhead* é alto.
    * A V2 também favorece a localidade de cache, pois as mesmas threads tendem a processar os mesmos blocos de dados repetidamente, diferentemente da V1 onde o SO pode alocar novas threads em núcleos diferentes a cada *fork*.