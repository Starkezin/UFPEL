# Processo de Desenvolvimento com IA - LibGeo

Este documento detalha o processo de obtenção do projeto **LibGeo**, destacando a interação com a Inteligência Artificial (IA) desde a análise da especificação básica até a entrega da solução funcional.

---

## 📅 Histórico de Interações e Etapas do Desenvolvimento

### 1. Análise da Especificação e Planejamento
* **Interação Inicial:** A IA recebeu a especificação detalhada exigindo uma aplicação React + Vite com APIs Open Library e REST Countries, usando Leaflet/React Leaflet e estruturação de componentes.
* **Ações da IA:** A IA analisou o escopo e propôs um **Plano de Implementação** estruturado:
  - Identificou que a Open Library retorna códigos de idioma de 3 letras (ISO 639-2) e verificou que a REST Countries os aceita diretamente, mitigando problemas de tradução direta.
  - Planejou um sistema de design premium em Vanilla CSS com Glassmorphic cards e tema escuro.
  - Identificou uma limitação comum no Leaflet ao rodar em ambientes Vite (ícones de marcador quebrados) e planejou a criação de **marcadores personalizados com bandeiras dos países**, evitando o bug padrão e elevando o nível visual do projeto.

---

### 2. Resolução de Problemas no Setup Inicial
* **O Problema:** Durante a execução do comando `npx -y create-vite@latest --help` para pesquisar os parâmetros de inicialização, a ferramenta `create-vite@latest` falhou com um erro de sintaxe (`SyntaxError`) devido à indisponibilidade de módulos modernos de Node (`node:util` - `styleText`) no ambiente local executando **Node v18.20.8**.
* **A Solução da IA:** A IA detectou rapidamente a incompatibilidade de versão do Node e alterou o comando para usar o gerador compatível com Node v18:
  ```bash
  npx -y create-vite@5 ./ --template react
  ```
  O comando foi executado com sucesso, inicializando o scaffolding do projeto na pasta raiz.

---

### 3. Instalação e Desenvolvimento de Código
* **Instalação:** A IA propôs a instalação do núcleo do projeto e das dependências extras de mapas e ícones:
  ```bash
  npm install && npm install leaflet react-leaflet lucide-react
  ```
* **Desenvolvimento Estruturado:**
  1. **Tradução de Idiomas (`src/utils/languageMap.js`):** Criou um dicionário robusto de mapeamento de códigos de 2 e 3 letras para português (ex: `por` -> Português, `eng` -> Inglês), garantindo a legibilidade dos dados.
  2. **Design System (`src/index.css`):** Formulou um arquivo de folha de estilos global contendo as variáveis HSL do tema escuro, efeitos de blur para glassmorphic cards, e as regras do mapa responsivo.
  3. **Barra de Pesquisa (`src/components/SearchBar.jsx`):** Escreveu o formulário de busca com validações, botão de limpeza e spinner para estados de carregamento.
  4. **Carregamento Visual (`src/components/LoadingState.jsx`):** Implementou skeleton loaders com animação pulsante para representar livros enquanto as requisições estão em trânsito.
  5. **Lista de Resultados (`src/components/BookList.jsx`):** Criou cartões interativos de exibição com capas reais da Open Library e suporte a covers de fallback em CSS para livros sem imagem cadastrada.
  6. **Painel de Detalhes (`src/components/BookDetails.jsx`):** Componentizou o detalhamento do livro contendo as informações básicas e os alertas dinâmicos de busca de países.
  7. **Mapa Geográfico (`src/components/CountryMap.jsx`):** Criou a interface com o Leaflet, adicionando marcadores circulares estilizados com a imagem da bandeira e um sub-componente `MapController` que escuta mudanças na lista de países e reajusta o enquadramento do mapa (`map.fitBounds()`).
  8. **Orquestração (`src/App.jsx`):** Consolidou as requisições assíncronas usando `fetch`, gerenciamento de erros (como 404 de países sem registros no banco) e controle de fluxos de tela.

---

### 4. Resolução de Erro em Tempo de Execução (React Leaflet)
* **O Problema:** Ao rodar e carregar a aplicação no navegador, a tela ficava totalmente preta/vazia (em branco), com o console acusando o erro `Uncaught TypeError: render2 is not a function` sob o componente `<Context.Consumer>` dentro do script do `react-leaflet.js`.
* **A Análise:** A IA identificou que a versão padrão instalada de `react-leaflet` foi a **v5.0.0**, desenvolvida exclusivamente para o **React 19**. Como o projeto foi construído no **React 18.3.1**, o React Leaflet tentava renderizar um Contexto diretamente como componente JSX (suportado apenas na versão 19), gerando uma falha grave na árvore de renderização.
* **A Solução da IA:** A IA procedeu com o downgrade do pacote para a última versão estável do ecossistema React 18:
  ```bash
  npm install react-leaflet@^4.2.1
  ```
  Com essa alteração no lockfile e no `package.json`, a aplicação voltou a renderizar a árvore de componentes perfeitamente.

---

## 🎨 Decisões de Design Alinhadas com a IA
A interação contínua permitiu adotar soluções que agregam valor estético e funcional:
- **Tonalidades do Mapa:** Adicionado um alternador de tema flutuante no mapa (CartoDB Dark Matter / CartoDB Voyager) para que o usuário possa escolher a exibição clara ou escura de forma contextualizada.
- **Micro-animações de Foco:** O marcador de bandeira possui uma animação em onda (pulse) que atrai a atenção visual do usuário para os pontos destacados.
- **Resiliência a Falhas:** Caso a API REST Countries retorne 404 (para línguas fictícias/indefinidas da Open Library como `sjn` ou `und`), o painel de detalhes notifica o usuário sem estragar a experiência de navegação.
