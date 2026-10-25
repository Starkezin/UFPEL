package routing

import (
	"container/heap"
	"errors"
	"fmt"
	"math"

	"rei-julian/internal/domain"
	"rei-julian/internal/mapdata"
)

type Service struct {
	repo *mapdata.Repository
}

func NewService(repo *mapdata.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Calculate(originID, destinationID string, profile domain.RouteProfile) (domain.Route, error) {
	if profile == "" {
		profile = domain.ProfileBalanced
	}
	if !validProfile(profile) {
		return domain.Route{}, fmt.Errorf("perfil invalido: %s", profile)
	}
	if _, ok := s.repo.Node(originID); !ok {
		return domain.Route{}, fmt.Errorf("origem inexistente: %s", originID)
	}
	if _, ok := s.repo.Node(destinationID); !ok {
		return domain.Route{}, fmt.Errorf("destino inexistente: %s", destinationID)
	}

	dist := map[string]float64{}
	prevNode := map[string]string{}
	prevEdge := map[string]domain.Edge{}
	for _, node := range s.repo.Nodes() {
		dist[node.ID] = math.Inf(1)
	}
	dist[originID] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{nodeID: originID, cost: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*item)
		if current.cost > dist[current.nodeID] {
			continue
		}
		if current.nodeID == destinationID {
			break
		}
		for _, edge := range s.repo.Neighbors(current.nodeID) {
			next := edge.To
			cost := current.cost + EdgeCost(edge, profile)
			if cost < dist[next] {
				dist[next] = cost
				prevNode[next] = current.nodeID
				prevEdge[next] = edge
				heap.Push(pq, &item{nodeID: next, cost: cost})
			}
		}
	}

	if math.IsInf(dist[destinationID], 1) {
		return domain.Route{}, errors.New("rota nao encontrada")
	}
	return s.buildRoute(originID, destinationID, string(profile), prevNode, prevEdge)
}

func (s *Service) buildRoute(originID, destinationID, profile string, prevNode map[string]string, prevEdge map[string]domain.Edge) (domain.Route, error) {
	nodeIDs := []string{destinationID}
	edges := []domain.Edge{}
	for cursor := destinationID; cursor != originID; {
		edge, ok := prevEdge[cursor]
		if !ok {
			return domain.Route{}, errors.New("rota incompleta")
		}
		edges = append([]domain.Edge{edge}, edges...)
		cursor = prevNode[cursor]
		nodeIDs = append([]string{cursor}, nodeIDs...)
	}

	nodes := make([]domain.Node, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		node, _ := s.repo.Node(id)
		nodes = append(nodes, node)
	}

	var totalDistance, safetySum, speedSum float64
	for _, edge := range edges {
		totalDistance += edge.Distance
		safetySum += edge.SafetyScore
		if edge.SpeedLimit > 0 {
			speedSum += float64(edge.SpeedLimit)
		}
	}

	avgSafety := 1.0
	if len(edges) > 0 {
		avgSafety = safetySum / float64(len(edges))
	}
	avgSpeed := 16.0
	if len(edges) > 0 && speedSum > 0 {
		avgSpeed = math.Min(speedSum/float64(len(edges)), 22)
	}

	return domain.Route{
		OriginID:      originID,
		DestinationID: destinationID,
		Profile:       profile,
		Nodes:         nodes,
		Edges:         edges,
		TotalDistance: round(totalDistance),
		EstimatedTime: round((totalDistance / 1000) / avgSpeed * 60),
		SafetyScore:   round(avgSafety),
	}, nil
}

func validProfile(profile domain.RouteProfile) bool {
	switch profile {
	case domain.ProfileBalanced, domain.ProfileSafest, domain.ProfileFastest:
		return true
	default:
		return false
	}
}

func EdgeCost(edge domain.Edge, profile domain.RouteProfile) float64 {
	distanceWeight, safetyWeight := 1.0, 2.0
	switch profile {
	case domain.ProfileFastest:
		distanceWeight, safetyWeight = 1.5, 0.7
	case domain.ProfileSafest:
		distanceWeight, safetyWeight = 0.8, 4.0
	}
	return edge.Distance*distanceWeight + riskPenalty(edge)*safetyWeight
}

func riskPenalty(edge domain.Edge) float64 {
	byType := map[domain.RoadType]float64{
		domain.RoadCycleTrack:   0.2,
		domain.RoadBikeLane:     0.4,
		domain.RoadResidential:  0.8,
		domain.RoadCommonStreet: 1.2,
		domain.RoadBusyAvenue:   2.5,
		domain.RoadHighway:      6.0,
		domain.RoadUnsafe:       8.0,
	}
	return byType[edge.StreetType] + float64(edge.TrafficLevel)*0.15 + (1-edge.SafetyScore)*2
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}

type item struct {
	nodeID string
	cost   float64
	index  int
}

type priorityQueue []*item

func (p priorityQueue) Len() int           { return len(p) }
func (p priorityQueue) Less(i, j int) bool { return p[i].cost < p[j].cost }
func (p priorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}
func (p *priorityQueue) Push(x any) {
	it := x.(*item)
	it.index = len(*p)
	*p = append(*p, it)
}
func (p *priorityQueue) Pop() any {
	old := *p
	it := old[len(old)-1]
	*p = old[:len(old)-1]
	return it
}
