package tracking

import (
	"math"

	"rei-julian/internal/domain"
	"rei-julian/internal/mapdata"
	"rei-julian/internal/routing"
)

type Service struct {
	repo      *mapdata.Repository
	routing   *routing.Service
	threshold float64
}

type UpdateResult struct {
	Status       string        `json:"status"`
	Recalculated bool          `json:"recalculated"`
	NearestNode  domain.Node   `json:"nearest_node"`
	NewRoute     *domain.Route `json:"new_route,omitempty"`
}

func NewService(repo *mapdata.Repository, routingService *routing.Service) *Service {
	return &Service{repo: repo, routing: routingService, threshold: 15}
}

func (s *Service) Update(pos domain.CyclistPosition, current domain.Route) (UpdateResult, error) {
	nearest := s.repo.NearestNode(pos)
	if distanceToRoute(pos, current) <= s.threshold {
		return UpdateResult{Status: "on_route", Recalculated: false, NearestNode: nearest}, nil
	}

	newRoute, err := s.routing.Calculate(nearest.ID, current.DestinationID, domain.RouteProfile(current.Profile))
	if err != nil {
		return UpdateResult{}, err
	}
	return UpdateResult{Status: "off_route", Recalculated: true, NearestNode: nearest, NewRoute: &newRoute}, nil
}

func distanceToRoute(pos domain.CyclistPosition, route domain.Route) float64 {
	if len(route.Nodes) == 0 {
		return math.Inf(1)
	}
	best := math.Inf(1)
	for i := 0; i < len(route.Nodes)-1; i++ {
		a := route.Nodes[i]
		b := route.Nodes[i+1]
		best = math.Min(best, distancePointToSegment(pos.X, pos.Y, a.X, a.Y, b.X, b.Y))
	}
	return best
}

func distancePointToSegment(px, py, ax, ay, bx, by float64) float64 {
	dx, dy := bx-ax, by-ay
	if dx == 0 && dy == 0 {
		return math.Hypot(px-ax, py-ay)
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(px-(ax+t*dx), py-(ay+t*dy))
}
