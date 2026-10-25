# Lucas ludwig e Samuel Starle

## Descrição
Consiste em um simulador de partículas 2D que utiliza uma combinação de C e Python. A física de interação das partículas e aplicação da força gravitacional é implementada em C,
enquanto a visualização e a animação é realizada em python com o auxilio de bibliotecas.

## Arquivos
particles.h > Cabeçalho com as definições de 'particles.c'.
particles.c > Implementação lógica das partículas.
particles_app.py > Script python que chama o código C e simula a dinâmica das partículas em uma animação 2D.
Makefile: Automação da compilação e execução do projeto.

## Dependências 
Compilador GCC
Python 3.x
biblioteca 'numpy' instalado via pip 
biblioteca 'matplotlib' instalado via pip
Codigo para instalação das bibliotecas: primeiro 'sudo apt install python3-pip', e depois 'pip install numpy matplotlib'
usando um simulador WSL para executar o código, talvez seja necessário a instalação de um visualizador gráfico (sudo apt-get install x11-apps)
e também a instalação da biblioteca tk do python (sudo apt-get install python3-tk).

## Compilação e Execução
para compilar o código C: make
para executar o simulador: make run
para limpar os arquivos gerados: make clean (Apenas caso queira compilar o código C novamente)
A compilação do arquivo C irá gerar um arquivo .so que pode ser lido apenas pelo linux.

Caso não seja possivel a compilação e execução do código pelo comando make, primeiro compile o arquivo C (gcc -shared -o particles.so -fPIC particles.),
após execute a interface de python (python3 particles_app.py).

## Documentação 
A documentação completa do projeto pode ser encontrada em documentacao.pdf.



