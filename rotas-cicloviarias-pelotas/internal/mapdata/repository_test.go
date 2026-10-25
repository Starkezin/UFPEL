package mapdata

import (
	"testing"

	"rei-julian/internal/domain"
)

func TestCityMapRespectsUrbanRules(t *testing.T) {
	repo, err := Load("../../maps/city_01.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if err := repo.ValidateUrbanRules(); err != nil {
		t.Fatalf("ValidateUrbanRules() error = %v", err)
	}
}

func TestPelotasMapHasProductionCoverage(t *testing.T) {
	repo, err := Load("../../maps/city_01.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	stats := repo.Stats()

	if stats.NodeCount < 15000 {
		t.Fatalf("mapa final deveria ter cobertura urbana ampla, recebeu %d nodes", stats.NodeCount)
	}
	if stats.EdgeCount < 20000 {
		t.Fatalf("mapa final deveria ter malha viaria ampla, recebeu %d edges", stats.EdgeCount)
	}
	if stats.ByStreetType[domain.RoadCycleTrack] < 500 {
		t.Fatalf("mapa deveria representar varias ciclovias, recebeu %+v", stats.ByStreetType)
	}
	if stats.ByStreetType[domain.RoadBikeLane] < 250 {
		t.Fatalf("mapa deveria representar varias ciclofaixas, recebeu %+v", stats.ByStreetType)
	}
	if stats.CyclingDistance < 50000 {
		t.Fatalf("infraestrutura cicloviaria muito pequena: %.2f", stats.CyclingDistance)
	}
}
