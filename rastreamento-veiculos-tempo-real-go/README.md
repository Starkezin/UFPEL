# Monitoramento de Dois Veiculos

Aplicacao web em Go para simulacao de rastreamento de dois veiculos em tempo real.

## O que o projeto faz

- Simula a movimentacao de dois veiculos em um mapa.
- Atualiza as posicoes a cada poucos segundos.
- Exibe o rastro percorrido por cada veiculo.
- Envia atualizacoes em tempo real do backend para o frontend usando SSE.

## Tecnologias

- Go
- HTML, CSS e JavaScript
- Leaflet
- OpenStreetMap

## Como executar

1. Instale o Go 1.18 ou superior.
2. Abra um terminal na pasta do projeto.
3. Rode o comando:

```powershell
go run ./cmd/rei-julian
```

4. Abra o navegador em:

```text
http://localhost:8080
```

## Endpoints

- `GET /positions`: retorna as posicoes atuais e o historico dos veiculos em JSON.
- `GET /stream`: envia atualizacoes em tempo real via Server-Sent Events.

## Ideia de uso academico

Este projeto pode ser descrito como um sistema simples de monitoramento de dois veiculos em area urbana, com exibicao da posicao atual e do trajeto percorrido.
