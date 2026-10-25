# Rei-Julian - Arquitetura em Go para rotas seguras de bicicleta

## 1. Visão geral

O Rei-Julian é uma aplicação de rastreamento e recomendação de rotas para ciclistas. O objetivo principal não é apenas encontrar o caminho mais curto, mas escolher uma rota equilibrada entre rapidez e segurança.

O sistema usa um grafo 2D derivado de dados OpenStreetMap de Pelotas/RS. O mapa representa ruas residenciais, vias de servico, avenidas, rodovias, ciclovias, ciclofaixas e trechos de maior risco para bicicleta, mantendo o contrato `Node`/`Edge` simples para roteamento local e testes automatizados.

## 2. Objetivos do projeto

- Representar Pelotas/RS em um grafo 2D derivado de OpenStreetMap.
- Permitir informar uma origem e um destino.
- Calcular uma rota considerando distância e segurança.
- Priorizar ciclovias, ciclofaixas e ruas tranquilas.
- Penalizar avenidas movimentadas e vias ruins para ciclistas.
- Simular a posição do ciclista se movendo pela rota.
- Detectar quando o ciclista saiu da rota planejada.
- Recalcular uma nova rota a partir da posição atual.
- Expor a lógica por uma API HTTP em Go.
- Servir uma interface web final com mapa clicavel, instrucoes de rota e simulacao de tracking.

## 3. Decisão principal: grafo urbano 2D

Para o primeiro desenvolvimento, o mapa pode ser uma grade ou um grafo 2D.

A melhor escolha para este projeto é usar um grafo.

Cada ponto importante do mapa será um `Node`, e cada rua ou trecho entre dois pontos será uma `Edge`. A rota será calculada sobre esse grafo.

Exemplo conceitual:

```text
Node A ---- rua residencial ---- Node B ---- ciclovia ---- Node C
  |                                                   |
avenida                                             ciclofaixa
  |                                                   |
Node D ---------------- rua comum ---------------- Node E
```

Essa estrutura preserva a topologia urbana necessaria para o roteamento e permite converter dados OSM para um JSON simples consumido pela API.

## 4. Modelo de domínio

### 4.1 Node

Representa um ponto do mapa.

Campos sugeridos:

```go
type Node struct {
    ID   string  `json:"id"`
    Name string  `json:"name"`
    X    float64 `json:"x"`
    Y    float64 `json:"y"`
}
```

### 4.2 Edge

Representa uma rua, avenida, ciclovia ou trecho entre dois pontos.

```go
type Edge struct {
    ID           string   `json:"id"`
    From         string   `json:"from"`
    To           string   `json:"to"`
    Distance     float64  `json:"distance"`
    StreetType   RoadType `json:"street_type"`
    SafetyScore  float64  `json:"safety_score"`
    SpeedLimit   int      `json:"speed_limit"`
    TrafficLevel int      `json:"traffic_level"`
    Bidirectional bool    `json:"bidirectional"`
}
```

### 4.3 RoadType

Tipos de trecho simulados:

```go
type RoadType string

const (
    RoadBikeLane       RoadType = "bike_lane"
    RoadCycleTrack     RoadType = "cycle_track"
    RoadResidential    RoadType = "residential"
    RoadCommonStreet   RoadType = "common_street"
    RoadBusyAvenue     RoadType = "busy_avenue"
    RoadUnsafe         RoadType = "unsafe"
)
```

Sugestão de pesos:

| Tipo de via | Descrição | Prioridade |
|---|---|---|
| `cycle_track` | Ciclovia protegida | Muito alta |
| `bike_lane` | Ciclofaixa | Alta |
| `residential` | Rua residencial tranquila | Boa |
| `common_street` | Rua comum | Média |
| `busy_avenue` | Avenida movimentada | Baixa |
| `unsafe` | Trecho inadequado | Evitar |

### 4.4 Route

Resultado do cálculo de rota.

```go
type Route struct {
    OriginID       string        `json:"origin_id"`
    DestinationID  string        `json:"destination_id"`
    Nodes          []Node        `json:"nodes"`
    Edges          []Edge        `json:"edges"`
    TotalDistance  float64       `json:"total_distance"`
    EstimatedTime  float64       `json:"estimated_time"`
    SafetyScore    float64       `json:"safety_score"`
}
```

### 4.5 CyclistPosition

Representa a posição atual do ciclista.

```go
type CyclistPosition struct {
    X        float64 `json:"x"`
    Y        float64 `json:"y"`
    NodeHint string  `json:"node_hint,omitempty"`
}
```

## 5. Cálculo de rota

O cálculo pode usar Dijkstra ou A*.

Dijkstra e suficiente para o tamanho atual do grafo e mantem o comportamento previsivel nos testes. A* pode ser introduzido depois se o mapa crescer muito.

### 5.1 Custo do trecho

Cada trecho deve ter um custo composto por distância e risco.

Exemplo:

```text
cost = distance * distanceWeight + riskPenalty * safetyWeight
```

Sugestão:

```text
distanceWeight = 1.0
safetyWeight   = 2.0
```

Penalidades por tipo:

| Tipo | Penalidade |
|---|---:|
| `cycle_track` | 0.2 |
| `bike_lane` | 0.4 |
| `residential` | 0.8 |
| `common_street` | 1.2 |
| `busy_avenue` | 2.5 |
| `unsafe` | 8.0 |

Com isso, uma ciclovia um pouco mais longa pode ser escolhida no lugar de uma avenida perigosa mais curta.

### 5.2 Perfis de rota

O sistema pode aceitar diferentes perfis:

```go
type RouteProfile string

const (
    ProfileBalanced RouteProfile = "balanced"
    ProfileSafest   RouteProfile = "safest"
    ProfileFastest  RouteProfile = "fastest"
)
```

Exemplo de pesos por perfil:

| Perfil | Distância | Segurança |
|---|---:|---:|
| `fastest` | 1.5 | 0.7 |
| `balanced` | 1.0 | 2.0 |
| `safest` | 0.8 | 4.0 |

## 6. Dados do mapa

O mapa fica versionado em `maps/city_01.json` e e carregado na inicializacao da API. O arquivo e gerado por `tools/osm_to_reijulian.py` a partir de recortes XML do OpenStreetMap.

Exemplo:

```json
{
  "nodes": [
    { "id": "A", "name": "Praça Central", "x": 10, "y": 10 },
    { "id": "B", "name": "Rua das Flores", "x": 30, "y": 10 },
    { "id": "C", "name": "Universidade", "x": 50, "y": 20 }
  ],
  "edges": [
    {
      "id": "E1",
      "from": "A",
      "to": "B",
      "distance": 20,
      "street_type": "cycle_track",
      "safety_score": 0.95,
      "speed_limit": 30,
      "traffic_level": 1,
      "bidirectional": true
    }
  ]
}
```

Esse JSON pode ser carregado na inicialização da aplicação.

## 7. Evolucao do mapa

Existem duas opções.

### 7.1 Mapa manual

Criar um JSON pequeno com 20 a 50 nós.

Vantagens:

- Mais previsível para testes.
- Fácil de demonstrar.
- Bom para validar o algoritmo de rota.

### 7.2 Gerador procedural

Criar um módulo em Go que gera uma cidade fictícia automaticamente.

Regras possíveis:

- Criar uma grade de pontos.
- Conectar pontos vizinhos.
- Marcar algumas linhas horizontais e verticais como avenidas.
- Criar corredores de ciclovia atravessando bairros.
- Marcar áreas residenciais com menor tráfego.
- Sortear alguns trechos perigosos.

A versao atual usa recortes OSM convertidos para o contrato de `Node` e `Edge`. Uma evolucao natural e automatizar o download dos recortes, armazenar metadados de data/licenca da extracao e enriquecer o grafo com POIs e restricoes de acesso especificas para bicicleta.

## 8. Recálculo de rota

Durante o deslocamento, o sistema compara a posição atual do ciclista com a rota planejada.

Fluxo:

1. Receber nova posição.
2. Calcular a distância da posição até os trechos da rota atual.
3. Se a distância for menor que o limite, manter a rota.
4. Se a distância for maior que o limite, considerar que houve desvio.
5. Encontrar o nó mais próximo da posição atual.
6. Calcular nova rota desse nó até o destino.

Exemplo de limite:

```text
offRouteThreshold = 15 unidades do mapa
```

No mapa fictício, essas unidades podem representar metros de forma simplificada.

## 9. Componentes da aplicação

### 9.1 API HTTP

Responsável por receber requisições e devolver respostas JSON.

Endpoints sugeridos:

```text
GET  /health
GET  /map
GET  /nodes
POST /routes
POST /tracking/update
POST /simulation/start
POST /simulation/step
```

### 9.2 Map Service

Responsável por carregar e consultar o mapa.

Funções:

- Carregar mapa do JSON.
- Retornar nós e trechos.
- Encontrar nó mais próximo de uma posição.
- Retornar vizinhos de um nó.

### 9.3 Routing Service

Responsável por calcular rotas.

Funções:

- Calcular rota entre origem e destino.
- Aplicar pesos por tipo de rua.
- Suportar perfis de rota.
- Retornar distância, tempo estimado e pontuação de segurança.

### 9.4 Tracking Service

Responsável por acompanhar a posição do ciclista.

Funções:

- Receber posição atual.
- Verificar se o ciclista está na rota.
- Detectar desvio.
- Solicitar recálculo da rota.

### 9.5 Simulation Service

Responsável por simular movimento sem precisar de GPS real.

Funções:

- Mover o ciclista ao longo de uma rota.
- Gerar atualizações periódicas de posição.
- Simular desvios.
- Testar recálculo de rota.

## 10. Estrutura de pastas sugerida

```text
rei-julian/
  cmd/
    api/
      main.go
  internal/
    domain/
      node.go
      edge.go
      route.go
      position.go
    mapdata/
      loader.go
      repository.go
      nearest.go
    routing/
      cost.go
      dijkstra.go
      service.go
    tracking/
      service.go
      off_route.go
    simulation/
      service.go
      movement.go
    transport/
      http/
        handlers.go
        routes.go
        dto.go
  maps/
    city_01.json
  tests/
    routing_test.go
    tracking_test.go
  go.mod
  README.md
```

## 11. Fluxo principal

### 11.1 Calcular rota

```text
Cliente -> POST /routes
API -> RoutingService.CalculateRoute(origin, destination, profile)
RoutingService -> MapRepository.GetNeighbors(node)
RoutingService -> CostCalculator.EdgeCost(edge, profile)
API -> retorna Route em JSON
```

### 11.2 Atualizar posição

```text
Cliente -> POST /tracking/update
API -> TrackingService.UpdatePosition(position)
TrackingService -> verifica distância até a rota atual
TrackingService -> se saiu da rota, chama RoutingService.CalculateRoute(...)
API -> retorna status e rota atualizada se necessário
```

## 12. Contratos de API

### 12.1 Criar rota

Request:

```json
{
  "origin_id": "A",
  "destination_id": "C",
  "profile": "balanced"
}
```

Response:

```json
{
  "route": {
    "origin_id": "A",
    "destination_id": "C",
    "total_distance": 52.4,
    "estimated_time": 14.2,
    "safety_score": 0.87,
    "nodes": [],
    "edges": []
  }
}
```

### 12.2 Atualizar posição

Request:

```json
{
  "x": 24.5,
  "y": 11.2
}
```

Response sem desvio:

```json
{
  "status": "on_route",
  "recalculated": false
}
```

Response com desvio:

```json
{
  "status": "off_route",
  "recalculated": true,
  "new_route": {}
}
```

## 13. Visualização 2D

Mesmo que o backend seja em Go, a visualização pode ser feita com uma página simples usando HTML Canvas.

O frontend pode desenhar:

- Nós como pontos.
- Trechos como linhas.
- Cores diferentes por tipo de via.
- Rota selecionada em destaque.
- Ciclista como um marcador se movendo.

Sugestão de cores:

| Tipo | Cor |
|---|---|
| Ciclovia | Verde |
| Ciclofaixa | Azul |
| Rua residencial | Cinza claro |
| Rua comum | Cinza |
| Avenida movimentada | Laranja |
| Via insegura | Vermelho |
| Rota atual | Roxo ou preto destacado |

## 14. Banco de dados

Para o MVP, não é obrigatório usar banco de dados.

Recomendação inicial:

- Carregar `maps/city_01.json` em memória.
- Manter rota atual e simulação em memória.
- Usar testes automatizados para validar os algoritmos.

Depois, se necessário:

- SQLite para persistir mapas, rotas e sessões.
- PostgreSQL se o projeto evoluir para múltiplos usuários.

## 15. Testes recomendados

Testes de routing:

- Deve escolher ciclovia mesmo se for um pouco mais longa.
- Deve evitar trecho `unsafe` quando existir alternativa.
- Deve aceitar perfil `fastest` escolhendo rota mais curta.
- Deve aceitar perfil `safest` escolhendo rota mais segura.

Testes de tracking:

- Deve detectar posição dentro da rota.
- Deve detectar desvio quando a posição estiver longe.
- Deve recalcular rota a partir do nó mais próximo.

Testes de mapa:

- Deve carregar JSON válido.
- Deve rejeitar edge apontando para node inexistente.
- Deve calcular vizinhos corretamente.

## 16. Ordem de implementação sugerida

1. Criar modelos em `internal/domain`.
2. Criar `maps/city_01.json` com mapa pequeno.
3. Implementar carregador do mapa.
4. Implementar cálculo de custo por trecho.
5. Implementar Dijkstra.
6. Criar endpoint `POST /routes`.
7. Criar tracking com detecção de desvio.
8. Criar recálculo de rota.
9. Criar simulação de movimento.
10. Criar visualização 2D simples com Canvas.

## 17. Evolução futura

Depois que o MVP funcionar com mapa fictício, o projeto pode evoluir para:

- Importar dados reais de OpenStreetMap.
- Usar GPS real no cliente.
- Salvar histórico de rotas.
- Permitir avaliação manual de trechos por usuários.
- Ajustar pesos com base em acidentes, iluminação ou horário.
- Criar múltiplas cidades fictícias para testes.
- Comparar rota mais rápida, mais segura e balanceada lado a lado.

## 18. Resumo da arquitetura

O núcleo do Rei-Julian deve ser independente de mapas reais. A parte mais importante é modelar a cidade como um grafo ponderado, onde cada trecho possui distância e nível de segurança. Com isso, o backend em Go consegue calcular rotas seguras, simular deslocamento e recalcular caminhos quando o ciclista sai da rota.

Essa abordagem reduz dependências externas, facilita testes e permite demonstrar o funcionamento da ideia antes de integrar APIs reais ou dados de GPS.
