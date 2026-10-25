# Relatorio do Projeto Rei-Julian

## Identificacao

**Projeto:** Rei-Julian  
**Tema:** Aplicacao de rotas seguras para ciclistas em Pelotas/RS  
**Tecnologias:** Go, API HTTP, HTML Canvas, JavaScript, OpenStreetMap  
**Objetivo final:** transformar o MVP inicial em uma aplicacao funcional com mapa realista, roteamento por seguranca, tracking, simulacao, busca, comparacao de rotas e visualizacao otimizada.

---

## Prompts Utilizados na Conversa

Abaixo estao os principais prompts enviados durante o desenvolvimento e o resultado pratico de cada etapa.

### 1. Gerar mapa procedural baseado em Pelotas

> Quero que voce use o contexto do rei-julian, para gerar um mapa procedural completo baseado na cidade de Pelotas, com as ruas reais, avenidas, distancia, locais aonde tem ciclovias e tudo mais. Pode usar quaisquer informacoes que voce tiver, quero que voce gere o mapa o mais parecido com Pelotas possivel.

**O que foi feito:**

- Analisei a estrutura inicial do projeto Rei-Julian.
- Identifiquei o contrato do mapa em `maps/city_01.json`.
- Substitui o mapa ficticio inicial por um mapa procedural inspirado em Pelotas.
- Adicionei regioes como Centro, Fragata, Tres Vendas, Areal, Porto e Laranjal.
- Inclui avenidas, rodovias, ciclovias e ciclofaixas em formato compativel com o grafo usado pelo roteador.
- Atualizei testes para aceitar tanto `cycle_track` quanto `bike_lane` como infraestrutura cicloviaria.

### 2. Finalizar o projeto, nao apenas manter como MVP

> Nao quero mais um mapa para MVP, quero agora ter o projeto finalizado por inteiro.

**O que foi feito:**

- Refatorei o projeto para se apresentar como aplicacao final, nao apenas MVP.
- Adicionei validacoes mais rigidas no carregamento do mapa.
- Adicionei endpoint `GET /stats`.
- Atualizei README e arquitetura.
- Melhorei a interface Canvas com estatisticas, instrucoes de rota, selecao pelo mapa e simulacao de desvio.
- Adicionei novos testes de mapa, roteamento e tracking.

### 3. Atualizar o mapa com ruas reais e conexoes completas

> Atualize esse mapa de pelotas, deixando as conexoes de ruas completas, estao faltando muitas ruas que nao tem ciclovias. tem conexoes que nao fazem sentido. Procure mapas de pelotas em grafos, ou estruturas de dados, existentes no github ou na internet.

**O que foi feito:**

- Pesquisei por bases prontas de grafo de Pelotas no GitHub e na internet.
- Nao encontrei uma base especifica pronta e reutilizavel.
- Usei dados reais do OpenStreetMap via recortes XML.
- Criei o conversor `tools/osm_to_reijulian.py`.
- Converti dados OSM para o formato `Node`/`Edge` do Rei-Julian.
- O mapa passou de uma malha manual pequena para um grafo real com:
  - 19.780 nos.
  - 23.316 conexoes.
  - ruas residenciais, vias de servico, avenidas, rodovias, ciclovias e ciclofaixas.
- Mantive apenas o maior componente conectado da malha urbana para evitar ilhas desconectadas.

### 4. Continuar a opcao correta apos interrupcao acidental

> Faca a opcao do 1, cliquei errado.

**O que foi feito:**

- Continuei a abordagem baseada em OpenStreetMap.
- Finalizei a geracao do mapa OSM.
- Atualizei testes para usar IDs reais do grafo OSM.
- Atualizei README e arquitetura com a origem dos dados e licenca OSM/ODbL.

### 5. Adicionar zoom no mapa

> Adicione a funcionalidade de dar zoom no mapa.

**O que foi feito:**

- Adicionei zoom com roda do mouse.
- Adicionei botoes `+`, `-` e `Reset`.
- Adicionei pan/arraste no Canvas.
- Mantive clique e hover funcionando corretamente com zoom aplicado.
- Adicionei indicador de percentual de zoom.

### 6. Identificar o que ainda faltava na arquitetura

> Dentro da arquitetura inicial, o que ainda nao foi implementado? e o que poderia ser implementado para melhorar o projeto

**O que foi feito:**

- Comparei o documento `ARQUITETURA_GO_REI_JULIAN.md` com o codigo implementado.
- Identifiquei funcionalidades faltantes:
  - Simulation Service real.
  - GPS real.
  - sessoes independentes.
  - busca por rua/local.
  - comparacao de perfis.
  - persistencia.
  - otimizacao visual.
- Sugeri ordem de implementacao priorizando impacto tecnico.

### 7. Implementar melhorias, otimizar visualizacao e remover excesso de pontos genericos

> Adicione esse itens no projeto. Tambem otimize a visualizacao do mapa. Esta mostrando muitos detalhes mesmo com o zoom longe, mostre mais detalhe apenas quando o zoom aproximar. Tambem existem mts pontos com "Acesso local sem nome" Retire esses pontos que estao sobre outras ruas ou rodovias. Mantenha o desenho das ruas normalmente, apenas retire esses pontos arredondados de acessos em execesso

**O que foi feito:**

- Implementei sessoes por usuario com `session_id`.
- Adicionei busca por rua/local: `GET /search`.
- Adicionei comparacao de perfis: `POST /routes/compare`.
- Adicionei Simulation Service real em `internal/simulation`.
- Adicionei endpoints:
  - `POST /simulation/start`
  - `POST /simulation/step`
- Adicionei suporte a GPS no frontend usando `navigator.geolocation`.
- Adicionei `lat` e `lon` nos nodes do mapa para relacionar GPS com o grafo.
- Otimizei o desenho do mapa por nivel de zoom:
  - zoom distante mostra vias principais, ciclovias e ciclofaixas;
  - zoom medio mostra mais ruas;
  - zoom aproximado mostra mais detalhes.
- Removi pontos genericos como `Acesso local sem nome` da selecao e exibicao em excesso, mantendo as ruas desenhadas normalmente.

### 8. Gerar este relatorio

> Agora quero que voce monte um arquivo .md e coloque todos os prompts utilizados nessa conversa, e tambem tudo que foi feito por voce, no estilo de um relatorio para apresentar ao meu professor. Nao me peca permissao para nada, apenas faca ate o final. quero sair e ir no mercado enquanto tu trabalha

**O que foi feito:**

- Criei este relatorio em Markdown com os prompts, etapas realizadas, arquivos modificados, funcionalidades implementadas, validacoes e proximos passos.

---

## Estado Inicial do Projeto

O projeto comecou como um MVP em Go para recomendacao de rotas de bicicleta em um mapa ficticio 2D. A estrutura inicial continha:

- `cmd/api/main.go`: servidor HTTP.
- `internal/domain/types.go`: modelos de dominio.
- `internal/mapdata/repository.go`: carregamento e validacao do mapa.
- `internal/routing/service.go`: Dijkstra e calculo de custo.
- `internal/tracking/service.go`: deteccao de desvio.
- `internal/transport/httpapi/handlers.go`: endpoints HTTP.
- `web/index.html` e `web/app.js`: visualizacao Canvas simples.
- `maps/city_01.json`: mapa pequeno e ficticio.

---

## Principais Implementacoes Realizadas

### 1. Mapa realista de Pelotas com OpenStreetMap

O mapa deixou de ser manual/procedural e passou a ser gerado a partir de dados reais do OpenStreetMap.

Foi criado o arquivo:

```text
tools/osm_to_reijulian.py
```

Esse script converte recortes XML do OSM para o formato do Rei-Julian:

```json
{
  "nodes": [],
  "edges": []
}
```

O conversor:

- le `highway=*` do OSM;
- reconhece ciclovias e ciclofaixas por tags como `highway=cycleway`, `cycleway=*` e `bicycle=designated`;
- calcula distancia por Haversine;
- classifica vias em `cycle_track`, `bike_lane`, `residential`, `common_street`, `busy_avenue`, `highway` e `unsafe`;
- mantem o maior componente conectado;
- adiciona coordenadas normalizadas `x/y` para o Canvas;
- adiciona `lat/lon` para suporte a GPS.

Resultado atual do mapa:

| Metrica | Valor |
|---|---:|
| Nos | 19.780 |
| Conexoes | 23.316 |
| Ciclovias (`cycle_track`) | 946 |
| Ciclofaixas (`bike_lane`) | 452 |
| Ruas residenciais | 10.527 |
| Ruas comuns | 4.296 |
| Avenidas movimentadas | 6.613 |
| Rodovias | 482 |

### 2. Backend HTTP mais completo

Foram adicionados ou atualizados os endpoints:

```text
GET  /health
GET  /map
GET  /nodes
GET  /stats
GET  /search?q=...
POST /routes
POST /routes/compare
POST /tracking/update
POST /simulation/start
POST /simulation/step
```

### 3. Sessoes por usuario

Antes, o servidor guardava apenas uma rota global em memoria. Isso causaria problema com varios usuarios, pois um usuario sobrescreveria a rota do outro.

Foi implementado `session_id` para manter estados separados:

- rota atual da sessao;
- simulacao atual da sessao;
- tracking baseado na rota da sessao.

### 4. Busca por ruas e locais

Foi criada busca por nomes de nodes do grafo:

```text
GET /search?q=ferreira%20viana
```

A busca remove acentos para facilitar consultas e filtra resultados genericos demais, como muitos pontos de `Acesso local sem nome`.

### 5. Comparacao de perfis de rota

Foi criado:

```text
POST /routes/compare
```

Ele calcula e compara:

- `fastest`
- `balanced`
- `safest`

A interface mostra distancia, tempo e seguranca de cada perfil.

### 6. Simulation Service real

Foi criado o pacote:

```text
internal/simulation
```

Ele simula o deslocamento ao longo da rota, avancando por metros a cada passo.

Endpoints:

```text
POST /simulation/start
POST /simulation/step
```

A interface usa esses endpoints para animar a posicao simulada no mapa.

### 7. Tracking com recalculo

O projeto ja tinha tracking inicial. Ele foi mantido e integrado a sessoes.

Fluxo atual:

1. usuario calcula rota;
2. backend cria ou reutiliza uma sessao;
3. frontend envia posicao atual com `session_id`;
4. backend detecta se saiu da rota;
5. se saiu, calcula nova rota a partir do no mais proximo.

### 8. Suporte a GPS no frontend

Foi adicionado uso de:

```javascript
navigator.geolocation.getCurrentPosition(...)
```

O GPS procura o node real mais proximo usando `lat/lon` e associa a posicao do usuario ao grafo.

### 9. Interface Canvas melhorada

A interface agora possui:

- mapa clicavel;
- zoom com roda do mouse;
- botoes `+`, `-` e `Reset`;
- pan/arraste;
- busca de origem e destino;
- comparacao de perfis;
- simulacao de movimento;
- botao para GPS;
- visualizacao de rota em destaque;
- marcador de posicao simulada;
- menos detalhes quando o zoom esta distante.

### 10. Otimizacao visual por nivel de zoom

Foi implementado filtro de detalhe:

- zoom distante: mostra rodovias, avenidas, ciclovias e ciclofaixas;
- zoom intermediario: inclui mais ruas importantes;
- zoom aproximado: mostra ruas menores e detalhes locais.

Tambem foi reduzido o ruido visual causado por pontos genericos:

- `Acesso local sem nome`;
- `Rua sem nome`;
- pontos OSM genericos.

Esses pontos continuam existindo no grafo quando necessarios para conexao, mas nao aparecem como marcadores arredondados em excesso.

---

## Arquivos Criados ou Modificados

### Criados

```text
internal/simulation/service.go
internal/simulation/service_test.go
tools/osm_to_reijulian.py
RELATORIO_REI_JULIAN.md
```

### Modificados

```text
maps/city_01.json
README.md
ARQUITETURA_GO_REI_JULIAN.md
internal/domain/types.go
internal/mapdata/repository.go
internal/mapdata/repository_test.go
internal/routing/service.go
internal/routing/service_test.go
internal/tracking/service.go
internal/tracking/service_test.go
internal/transport/httpapi/handlers.go
web/index.html
web/app.js
```

---

## Validacoes Realizadas

Foram executados os seguintes comandos de validacao:

```bash
go test ./...
node --check web/app.js
```

Resultado: todos os testes Go passaram e o JavaScript foi validado sintaticamente.

Tambem foram testados manualmente via API:

```text
GET  /health
GET  /stats
GET  /search
POST /routes
POST /routes/compare
POST /simulation/start
POST /simulation/step
```

---

## Fontes e Dados Externos

A malha viaria foi baseada em dados publicos do OpenStreetMap.

Referencias:

- OpenStreetMap API v0.6: https://wiki.openstreetmap.org/wiki/Api06
- Licenca e atribuicao OSM/ODbL: https://www.openstreetmap.org/copyright
- OpenStreetMap Pelotas: https://www.openstreetmap.org/search?query=Pelotas%20RS

Observacao: os XMLs brutos do OSM nao foram versionados no projeto; o projeto mantem o JSON convertido em `maps/city_01.json` e o conversor em `tools/osm_to_reijulian.py`.

---

## Funcionalidades Atuais do Projeto

Atualmente o Rei-Julian possui:

- mapa realista de Pelotas baseado em OSM;
- roteamento com Dijkstra;
- perfis `fastest`, `balanced` e `safest`;
- custo ponderado por distancia, risco, tipo de via, trafego e seguranca;
- tracking com deteccao de desvio;
- recalculo automatico;
- sessoes por usuario;
- busca por ruas/locais;
- comparacao entre perfis;
- simulacao de movimento;
- suporte inicial a GPS;
- visualizacao Canvas com zoom, pan e nivel de detalhe;
- testes automatizados.

---

## Limitacoes Atuais

Apesar da evolucao, ainda existem pontos que podem ser melhorados:

1. **Persistencia**  
   O estado ainda fica em memoria. Se o servidor reiniciar, sessoes e rotas sao perdidas.

2. **Sem banco de dados**  
   Ainda nao ha SQLite ou PostgreSQL para historico de rotas, usuarios ou avaliacoes.

3. **Busca simples**  
   A busca funciona por nome de node, mas ainda nao e uma busca geocodificada completa por endereco.

4. **Sem A***  
   O roteamento usa Dijkstra. Funciona, mas A* poderia melhorar performance em grafos maiores.

5. **Sem restricoes avancadas de trafego**  
   O projeto ainda nao considera restricoes complexas como conversoes proibidas, horario, iluminacao, acidentes ou fluxo real.

6. **GPS inicial**  
   O GPS associa a posicao ao no mais proximo, mas ainda nao faz acompanhamento continuo com `watchPosition`.

---

## Melhorias Futuras Recomendadas

Para evoluir ainda mais o projeto, eu recomendaria:

1. Implementar SQLite para persistir sessoes, rotas e historico.
2. Trocar Dijkstra por A* usando heuristica espacial.
3. Criar geocodificacao local por rua, bairro e POI.
4. Adicionar `watchPosition` para tracking em tempo real.
5. Criar cache de tiles ou renderizacao em camadas para melhorar performance visual.
6. Adicionar avaliacoes de usuarios sobre seguranca dos trechos.
7. Considerar iluminacao, acidentes, horario e superficie da via no custo.
8. Criar tela de historico de rotas.
9. Criar exportacao de rota em GeoJSON ou GPX.
10. Automatizar download dos recortes OSM com metadados de data e licenca.

---

## Conclusao

O Rei-Julian evoluiu de um MVP com mapa ficticio para uma aplicacao funcional de rotas para bicicleta em Pelotas/RS. O projeto agora usa dados reais do OpenStreetMap, possui uma API mais completa, trabalha com sessoes independentes, calcula rotas por diferentes perfis, simula movimento, permite busca por ruas e apresenta o mapa em Canvas com zoom, pan e controle de detalhe.

A arquitetura inicial foi respeitada em seu nucleo: cidade como grafo, ruas como arestas ponderadas, rota calculada por custo e tracking com recalculo. Alem disso, foram implementadas varias evolucoes que tornam o projeto mais proximo de uma aplicacao real.
