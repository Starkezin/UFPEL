package tracking

import (
	"testing"

	"rei-julian/internal/domain"
	"rei-julian/internal/mapdata"
	"rei-julian/internal/routing"
)

func TestTrackingRecalculatesWhenCyclistLeavesRoute(t *testing.T) {
	repo, err := mapdata.Load("../../maps/city_01.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	routingService := routing.NewService(repo)
	current, err := routingService.Calculate("N1675175053", "N1241501806", domain.ProfileBalanced)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	service := NewService(repo, routingService)

	result, err := service.Update(domain.CyclistPosition{X: 18, Y: 18}, current)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if !result.Recalculated || result.NewRoute == nil {
		t.Fatalf("tracking deveria recalcular rota quando ciclista sai do trajeto: %+v", result)
	}
	if result.NewRoute.OriginID != result.NearestNode.ID {
		t.Fatalf("nova rota deveria partir do node mais proximo")
	}
}
