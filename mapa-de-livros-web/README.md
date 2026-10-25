# LibGeo - Mapeamento Geográfico de Obras Literárias

O **LibGeo** é uma aplicação web interativa desenvolvida em **React** + **Vite** que permite pesquisar livros, visualizar seus detalhes e mapear, de forma geográfica, os países que compartilham do idioma principal da obra selecionada.

A aplicação consome dados públicos da **Open Library API** para busca e metadados de livros, e da **REST Countries API** para obter os países e suas coordenadas a partir do código do idioma correspondente. A renderização do mapa dinâmico é feita com a biblioteca **Leaflet** e seu wrapper **React Leaflet**.

---

## 🚀 Como Executar o Projeto

Siga os passos abaixo para rodar o projeto localmente:

1. **Instalar as Dependências:**
   No diretório raiz do projeto, instale os pacotes necessários rodando:
   ```bash
   npm install
   ```

2. **Iniciar o Servidor de Desenvolvimento:**
   Execute a aplicação em modo de desenvolvimento local com:
   ```bash
   npm run dev
   ```

3. **Acessar a Aplicação:**
   Abra o seu navegador e acesse a URL que aparecerá no terminal (normalmente `http://localhost:5173`).

---

## 🛠️ Principais Decisões Técnicas

Durante o desenvolvimento do projeto, adotamos as seguintes premissas e tecnologias:

1. **React com Vite:** 
   Utilizado para uma compilação ultra-rápida (Hot Module Replacement) e arquitetura de componentes reutilizáveis baseados em estados e hooks eficientes.
   
2. **Design System com Vanilla CSS:**
   Optamos por estilizar o projeto puramente com CSS padrão estruturado em variáveis HSL. Isso permite controle máximo sobre os detalhes estéticos sem a necessidade de frameworks adicionais (como Tailwind), criando uma interface premium com o conceito de **Glassmorphism** (efeito de vidro fosco), tema escuro (Dark Mode por padrão), micro-animações nas interações e um layout totalmente responsivo para dispositivos móveis.

3. **Geolocalização com React Leaflet:**
   A visualização geográfica utiliza o Leaflet alimentado por camadas de mapa limpas e modernas (CartoDB Dark Matter / Voyager) com atalhos para alternar entre temas visuais diretamente no mapa. Os marcadores no mapa são renderizados usando um `divIcon` customizado que encapsula a imagem da bandeira de cada país com animações em onda (pulse).
   
4. **Mapeamento de Idiomas (ISO 639):**
   A Open Library retorna os códigos de idioma no padrão ISO 639-2 (3 letras). Criamos um módulo utilitário (`languageMap.js`) capaz de traduzir esses códigos para seus respectivos nomes legíveis em português e gerenciar a busca adequada na REST Countries API, que aceita esses códigos nativamente.

5. **Tratamento Robusto de Estados e Erros:**
   O fluxo de dados trata os seguintes cenários excepcionais:
   - **Busca sem resultados:** Exibição de painel amigável sem quebrar o layout.
   - **Obras sem idioma:** Notifica o usuário que a obra carece de dados de idioma na Open Library.
   - **Idiomas sem países mapeados:** Trata retornos vazios ou erros 404 da REST Countries de forma silenciosa e informativa.
   - **Erros de conexão/API:** Telas de falha com botões de "Tentar Novamente" integrados para ambas as APIs.
   - **Responsividade:** Layout em grid no desktop (barra lateral + mapa completo) que colapsa de forma empilhada e organizada em telas menores de celulares.
