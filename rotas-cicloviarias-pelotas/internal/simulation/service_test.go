package simulation

import (
	"testing"

	"rei-julian/internal/domain"
)

func TestStepMovesAlongRoute(t *testing.T) {
	route := domain.Route{
		Nodes: []domain.Node{{ID: "A", X: 0, Y: 0}, {ID: "B", X: 10, Y: 0}},
		Edges: []domain.Edge{{ID: "E1", From: "A", To: "B", Distance: 100}},
	}
	state := NewState(25)
	result := Step(route, state)

	if result.State.Done {
		t.Fatalf("simulacao nao deveria terminar no primeiro passo")
	}
	if result.Position.X <= 0 || result.Position.X >= 10 {
		t.Fatalf("posicao deveria avancar no segmento: %+v", result.Position)
	}
}
