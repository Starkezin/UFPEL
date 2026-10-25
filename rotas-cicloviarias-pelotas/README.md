# Rei-Julian

Aplicacao em Go para recomendacao de rotas de bicicleta em Pelotas/RS, considerando distancia, seguranca, tipo de via e infraestrutura cicloviaria.

## Funcionalidades

- Mapa de Pelotas/RS em `maps/city_01.json`, gerado a partir de recortes OpenStreetMap, com ruas residenciais, vias de servico, avenidas, rodovias, ciclovias e ciclofaixas.
- API HTTP em Go com mapa, estatisticas, busca, calculo de rotas, comparacao de perfis, sessoes, tracking e simulacao.
- Calculo de rota com Dijkstra e perfis `balanced`, `safest` e `fastest`.
- Validacao estrutural do mapa ao carregar: IDs duplicados, referencias inexistentes, distancias invalidas, tipos de via invalidos, scores fora da faixa e componentes desconectados.
- Interface Canvas em `web/` com mapa clicavel, zoom/pan, busca por ruas, resumo da rota, instrucoes por trecho, comparacao de perfis, estatisticas da malha, GPS e simulacao de movimento.

## Dados do mapa

O mapa padrao agora e um grafo derivado de dados OpenStreetMap. Os recortes XML foram convertidos por `tools/osm_to_reijulian.py`, preservando conexoes entre segmentos reais de `highway=*`, nomes de ruas, distancias calculadas por Haversine e tags de infraestrutura cicloviaria (`highway=cycleway`, `cycleway=*`, `bicycle=designated`).

Para manter o roteamento consistente, o conversor descarta pequenas ilhas desconectadas e mantem o maior componente conectado da malha urbana. O JSON gerado contem 19.780 pontos e 23.316 conexoes. Dados OpenStreetMap estao sob ODbL; veja https://www.openstreetmap.org/copyright.

## Atualizar mapa a partir do OSM

Os XMLs brutos nao ficam versionados. Baixe recortes pelo endpoint de mapa do OpenStreetMap ou por outro extrator OSM e rode:

```bash
python3 tools/osm_to_reijulian.py /tmp/pelotas_core.osm /tmp/pelotas_laranjal.osm -o maps/city_01.json
```

Depois valide:

```bash
go test ./...
node --check web/app.js
```

## Rodar

```bash
go run ./cmd/api -addr :8080
```

Abra:

```text
http://localhost:8080
```

Verificar se o servico esta respondendo:

```bash
curl http://localhost:8080/health
```


## Desligar

Se o servico foi iniciado no terminal com `go run`, pressione:

```text
Ctrl+C
```

Se a porta `8080` continuar ocupada, encontre o processo:

```bash
lsof -i :8080
```

Depois finalize pelo PID mostrado:

```bash
kill <PID>
```

Se quiser iniciar em outra porta:

```bash
go run ./cmd/api -addr :8081
```

## Endpoints

```text
GET  /health
GET  /map
GET  /nodes
GET  /stats
GET  /search?q=ferreira%20viana
POST /routes
POST /routes/compare
POST /tracking/update
POST /simulation/start
POST /simulation/step
```

Exemplo de rota:

```bash
curl -X POST http://localhost:8080/routes \
  -H 'Content-Type: application/json' \
  -d '{"origin_id":"N1675175053","destination_id":"N1241501806","profile":"balanced"}'
```

A resposta inclui `session_id`; envie esse valor em tracking e simulacao para manter estados independentes por usuario.

## Testes

```bash
go test ./...
```

## Verificacao final

```bash
go test ./...
node --check web/app.js
```

Exemplo de estatisticas do mapa:

```bash
curl http://localhost:8080/stats
```
