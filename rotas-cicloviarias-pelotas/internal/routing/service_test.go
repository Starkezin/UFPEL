package routing

import (
	"testing"

	"rei-julian/internal/domain"
	"rei-julian/internal/mapdata"
)

const testOrigin = "N1675175053"      // Rua Andrade Neves / Rua Coronel Gomes Carneiro
const testDestination = "N1241501806" // Avenida Doutor Antonio Augusto de Assumpcao Junior / Avenida Rio Grande do Sul

func TestBalancedRouteUsesCyclingInfraAndAvoidsHighway(t *testing.T) {
	repo, err := mapdata.Load("../../maps/city_01.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	service := NewService(repo)

	route, err := service.Calculate(testOrigin, testDestination, domain.ProfileBalanced)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	hasCyclingInfra := false
	for _, edge := range route.Edges {
		if edge.StreetType == domain.RoadHighway || edge.StreetType == domain.RoadUnsafe {
			t.Fatalf("rota balanceada nao deveria usar %s: edge %s", edge.StreetType, edge.ID)
		}
		if edge.StreetType == domain.RoadCycleTrack || edge.StreetType == domain.RoadBikeLane {
			hasCyclingInfra = true
		}
	}
	if !hasCyclingInfra {
		t.Fatalf("rota balanceada deveria usar infraestrutura cicloviaria: %+v", route.Edges)
	}
}

func TestInvalidProfileIsRejected(t *testing.T) {
	repo, err := mapdata.Load("../../maps/city_01.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	service := NewService(repo)

	_, err = service.Calculate(testOrigin, testDestination, domain.RouteProfile("turbo"))
	if err == nil {
		t.Fatalf("Calculate() deveria rejeitar perfil invalido")
	}
}
