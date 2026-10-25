import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os
import numpy as np

# ALTERAÇÃO AQUI: Apontando para a pasta results
FILE_NAME = 'results/resultados_tarefaA.csv'

if not os.path.exists(FILE_NAME):
    print(f"Erro: '{FILE_NAME}' não encontrado. Rode 'make run' primeiro.")
    exit()

try:
    df = pd.read_csv(FILE_NAME)
    df.columns = df.columns.str.strip()
    df['Schedule'] = df['Schedule'].astype(str).str.strip()
except Exception as e:
    print(f"Erro ao ler CSV: {e}")
    exit()

sns.set_theme(style="whitegrid")

# Tratamento para identificar o sequencial
# Certifique-se que seu C sequencial imprime "SEQ" na coluna Schedule
# ou ajuste a lógica abaixo caso ele imprima outra coisa (ex: "static" com threads=1)
df_seq = df[df['Schedule'] == 'SEQ'][['N', 'K', 'Tempo']].rename(columns={'Tempo': 'Tempo_Seq'})

df_seq = df_seq.groupby(['N', 'K']).mean().reset_index()

df = df.merge(df_seq, on=['N', 'K'], how='left')

df['Speedup'] = df['Tempo_Seq'] / df['Tempo']

df_plot = df[df['Schedule'] != 'SEQ'].copy()


def plot_speedup(n_val, k_val, sufixo_nome):

    subset = df_plot[(df_plot['N'] == n_val) & (df_plot['K'] == k_val)]
    
    if subset.empty:
        return

    resumo = subset.groupby(['Schedule', 'Threads', 'Chunk'])['Speedup'].mean().reset_index()
    # Pega a melhor configuração de Chunk para cada par Schedule/Thread
    best_configs = resumo.loc[resumo.groupby(['Schedule', 'Threads'])['Speedup'].idxmax()]

    plt.figure(figsize=(10, 6))
    
    sns.lineplot(data=best_configs, x='Threads', y='Speedup', hue='Schedule', style='Schedule', 
                 markers=True, dashes=False, linewidth=2.5, markersize=8)
    
    max_threads = subset['Threads'].max()
    plt.plot([1, max_threads], [1, max_threads], '--', color='grey', label='Ideal (Linear)')

    plt.title(f'Speedup (Melhor Chunk): N={n_val}, K={k_val} ({sufixo_nome})', fontsize=14)
    plt.ylabel('Speedup (x vezes mais rápido)', fontsize=12)
    plt.xlabel('Número de Threads', fontsize=12)
    plt.legend(title='Schedule')
    plt.xticks(sorted(subset['Threads'].unique()))
    
    nome_arquivo = f'speedup_N{n_val}_K{k_val}.png'
    plt.savefig("plots/" + nome_arquivo)
    plt.close()
    print(f"Gerado: {nome_arquivo}")

def plot_tempo_chunk_fixo(n_val, k_val, chunk_foco=16):

    subset = df_plot[(df_plot['N'] == n_val) & (df_plot['K'] == k_val) & (df_plot['Chunk'] == chunk_foco)]
    
    if subset.empty:
        return

    plt.figure(figsize=(10, 6))
    sns.lineplot(data=subset, x='Threads', y='Tempo', hue='Schedule', style='Schedule', markers=True, dashes=False)
    
    # Tenta pegar o tempo sequencial, se existir para esses parâmetros
    if not df_seq[(df_seq['N'] == n_val) & (df_seq['K'] == k_val)].empty:
        tempo_seq = df_seq[(df_seq['N'] == n_val) & (df_seq['K'] == k_val)]['Tempo_Seq'].values[0]
        plt.axhline(y=tempo_seq, color='red', linestyle='--', label=f'Sequencial ({tempo_seq:.2f}s)')

    plt.title(f'Tempo de Execução (Chunk={chunk_foco}): N={n_val}, K={k_val}', fontsize=14)
    plt.ylabel('Tempo (s)')
    plt.xlabel('Threads')
    plt.legend()
    
    nome_arquivo = f'tempo_fixo_chunk{chunk_foco}_N{n_val}_K{k_val}.png'
    plt.savefig("plots/" + nome_arquivo)
    plt.close()
    print(f"Gerado: {nome_arquivo}")

def plot_impacto_chunk(n_val, k_val, thread_alvo=16):

    subset = df_plot[(df_plot['N'] == n_val) & (df_plot['K'] == k_val) & (df_plot['Threads'] == thread_alvo)]
    
    subset = subset[subset['Schedule'].isin(['static','dynamic', 'guided'])]

    if subset.empty:
        return

    plt.figure(figsize=(10, 6))
    sns.lineplot(data=subset, x='Chunk', y='Tempo', hue='Schedule', marker='o')
    
    plt.title(f'Impacto do Chunk (Threads={thread_alvo}): N={n_val}, K={k_val}', fontsize=14)
    plt.xscale('log', base=2)
    plt.ylabel('Tempo (s)')
    plt.xlabel('Tamanho do Chunk')
    plt.xticks(sorted(subset['Chunk'].unique()))
    plt.grid(True, which="both", ls="-", alpha=0.5)
    
    nome_arquivo = f'impacto_chunk_N{n_val}_K{k_val}.png'
    plt.savefig("plots/" + nome_arquivo)
    plt.close()
    print(f"Gerado: {nome_arquivo}")

    
print("=== Gerando Gráficos de Speedup ===")
plot_speedup(1000000, 20, "Carga Leve/Granularidade Fina") 
plot_speedup(1000000, 28, "Carga Pesada/Irregular")

print("\n=== Gerando Gráficos de Tempo (Chunk Fixo 16) ===")
plot_tempo_chunk_fixo(1000000, 28, chunk_foco=16)

print("\n=== Gerando Análise de Chunk (Ajuste Fino) ===")
plot_impacto_chunk(1000000, 28, thread_alvo=16)


print("\nConcluído!")