import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os
import sys

# --- Configuração de Diretórios ---
# Garante que funciona rodando da raiz ou da pasta scripts
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
RES_DIR = os.path.join(SCRIPT_DIR, '../results')
PLOT_DIR = os.path.join(SCRIPT_DIR, '../plots')

# Nome do arquivo gerado pelo run.sh
CSV_NAME = 'resultados_tarefaD.csv'

# Cria a pasta de plots se ela não existir
os.makedirs(PLOT_DIR, exist_ok=True)

# Define o estilo dos gráficos
plt.style.use('bmh')
sns.set_palette("viridis")

def plot_tarefa_D():
    csv_path = os.path.join(RES_DIR, CSV_NAME)
    
    print(f"Lendo dados de: {csv_path}")

    if not os.path.exists(csv_path):
        print(f"ERRO: O arquivo '{csv_path}' não foi encontrado.")
        print("Certifique-se de que o 'run.sh' foi executado e gerou o arquivo na pasta 'results'.")
        return

    try:
        # Carregar dados
        df = pd.read_csv(csv_path)
        
        # Limpar espaços em branco nos nomes das colunas
        df.columns = df.columns.str.strip()

        # --- CORREÇÃO: Padronizar nomes das colunas ---
        # O CSV vem como: N, Threads, Version, Time
        # O script usa:   N, Threads, Versao, Tempo
        mapa_colunas = {
            'Version': 'Versao',
            'Time': 'Tempo'
        }
        df.rename(columns=mapa_colunas, inplace=True)

        # Verifica se as colunas esperadas existem após a renomeação
        required_cols = {'N', 'Threads', 'Versao', 'Tempo'}
        if not required_cols.issubset(df.columns):
            print(f"Erro: O CSV não contém as colunas necessárias.")
            print(f"Esperado: {required_cols}")
            print(f"Encontrado: {set(df.columns)}")
            return

        # Calcular média das repetições
        df_mean = df.groupby(['N', 'Threads', 'Versao'])['Tempo'].mean().reset_index()

        # Pivotar a tabela: Versao vira coluna (0, 1, 2)
        df_pivot = df_mean.pivot(index=['N', 'Threads'], columns='Versao', values='Tempo').reset_index()
        
        # --- VERIFICAÇÃO DE DADOS ---
        # Verifica quais versões estão presentes no CSV
        versoes_presentes = [c for c in df_pivot.columns if isinstance(c, int)]
        print(f"Versões encontradas no CSV: {versoes_presentes}")

        if 1 not in df_pivot.columns or 2 not in df_pivot.columns:
            print("\nAVISO CRÍTICO: Dados insuficientes para comparação V1 vs V2.")
            print("Certifique-se de que o 'run.sh' rodou a Tarefa D completa (Versões 1 e 2).")
            print("Se você rodou apenas o Sequencial (0), não há comparação de overhead OpenMP para plotar.")
            return

        # Renomear colunas para facilitar acesso
        df_pivot.rename(columns={1: 'Tempo_V1', 2: 'Tempo_V2'}, inplace=True)
        
        # Se existir versão 0 (Sequencial), renomeia também para uso futuro
        if 0 in df_pivot.columns:
            df_pivot.rename(columns={0: 'Tempo_Seq'}, inplace=True)

        # Calcular Métricas (Foco em V1 vs V2 conforme enunciado)
        df_pivot['Speedup_Otimizacao'] = df_pivot['Tempo_V1'] / df_pivot['Tempo_V2']
        
        # --- GRÁFICO 1: Comparação de Tempos (Foco em Overhead - Menor N) ---
        n_min = df_pivot['N'].min()
        subset_min = df_pivot[df_pivot['N'] == n_min]

        plt.figure(figsize=(10, 6))
        
        plt.plot(subset_min['Threads'], subset_min['Tempo_V1'], marker='o', label='V1: Ingênua', color='red', linestyle='--')
        plt.plot(subset_min['Threads'], subset_min['Tempo_V2'], marker='s', label='V2: Arrumada', color='green', linestyle='-')
        
        # Se tiver sequencial, plota também para referência
        if 'Tempo_Seq' in subset_min.columns:
             # O sequencial tem Threads=1, pegamos o valor base
             tempo_seq = subset_min['Tempo_Seq'].iloc[0] 
             if pd.notna(tempo_seq):
                plt.axhline(y=tempo_seq, color='blue', linestyle=':', label='Sequencial (Ref)')

        plt.title(f'Tarefa D: Comparação de Tempo (N={n_min}) - Overhead')
        plt.xlabel('Número de Threads')
        plt.ylabel('Tempo Total (s)')
        plt.legend()
        plt.grid(True, alpha=0.3)
        if len(subset_min['Threads'].unique()) > 1:
            plt.xticks(subset_min['Threads'].unique())
        
        out_path1 = os.path.join(PLOT_DIR, 'tarefaD_overhead_tempos.png')
        plt.savefig(out_path1)
        plt.close()
        print(f"-> Gráfico gerado: {out_path1}")

        # --- GRÁFICO 2: Speedup da Otimização por Tamanho de N ---
        plt.figure(figsize=(10, 6))
        
        sns.lineplot(data=df_pivot, x='Threads', y='Speedup_Otimizacao', hue='N', marker='D', palette='viridis', style='N')
        
        plt.axhline(1.0, color='gray', linestyle='--', linewidth=1, label='Baseline (Igual)')
        
        plt.title('Tarefa D: Ganho da Versão Arrumada sobre a Ingênua')
        plt.ylabel('Speedup (V1 / V2)')
        plt.xlabel('Número de Threads')
        plt.legend(title='Tamanho N')
        plt.grid(True, alpha=0.3)
        if len(df_pivot['Threads'].unique()) > 1:
            plt.xticks(df_pivot['Threads'].unique())
        
        out_path2 = os.path.join(PLOT_DIR, 'tarefaD_speedup_otimizacao.png')
        plt.savefig(out_path2)
        plt.close()
        print(f"-> Gráfico gerado: {out_path2}")

        # --- TABELA RESUMO ---
        print("\n=== Resumo dos Ganhos (Speedup V2 vs V1) ===")
        cols_to_show = ['N', 'Threads', 'Tempo_V1', 'Tempo_V2', 'Speedup_Otimizacao']
        if 'Tempo_Seq' in df_pivot.columns:
            cols_to_show.insert(2, 'Tempo_Seq')
            
        print(df_pivot[cols_to_show].to_string(index=False))

    except Exception as e:
        print(f"Erro ao processar Tarefa D: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    print(f"Iniciando plotagem da Tarefa D...")
    plot_tarefa_D()
    print("Concluído.")