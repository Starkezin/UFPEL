# Compilador e ferramentas
JAVAC = javac
JAVA = java
JAVACC = javacc

# Arquivo da gramática
GRAMMAR = Parser.jj

# Diretórios
SRC_DIR = javaFiles
BUILD_DIR = build

# Nome da classe principal
MAIN_CLASS = Parser

# Arquivo de entrada (pode ser alterado ao rodar)
INPUT_FILE = exemplo.lov

# Regra padrão
all: $(BUILD_DIR)/$(MAIN_CLASS).class

# 1. Gera os arquivos Java a partir da gramática
ParserFiles: $(GRAMMAR)
	@echo ">> Gerando arquivos Java a partir da gramática..."
	$(JAVACC) $(GRAMMAR)

# 2. Compila todos os .java da pasta files/
$(BUILD_DIR)/$(MAIN_CLASS).class: ParserFiles
	@echo ">> Preparando diretorio de build..."
	@mkdir -p $(BUILD_DIR)
	@echo ">> Compilando código-fonte..."
	$(JAVAC) -d $(BUILD_DIR) $(SRC_DIR)/*.java

# 3. Executa o parser (pode sobrescrever o arquivo na linha de comando)
run:
	@echo ">> Executando o parser com o arquivo $(INPUT_FILE)..."
	$(JAVA) -cp $(BUILD_DIR) $(MAIN_CLASS) $(INPUT_FILE)

# 4. Limpa arquivos gerados
clean:
	@echo ">> Limpando arquivos gerados..."
	if exist build rmdir /S /Q build
	if exist javaFiles rmdir /S /Q javaFiles
	if exist -p rmdir /S /Q -p
	@echo ">> Limpeza concluida."
