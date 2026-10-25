package mapdata

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"rei-julian/internal/domain"
)

type Repository struct {
	city      domain.CityMap
	nodesByID map[string]domain.Node
	edgesByID map[string]domain.Edge
	adjacency map[string][]domain.Edge
}

type SearchResult struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Zone  string  `json:"zone"`
	Score int     `json:"score"`
}

type Stats struct {
	Name            string                  `json:"name"`
	NodeCount       int                     `json:"node_count"`
	EdgeCount       int                     `json:"edge_count"`
	TotalDistance   float64                 `json:"total_distance"`
	CyclingDistance float64                 `json:"cycling_distance"`
	AverageSafety   float64                 `json:"average_safety"`
	ByStreetType    map[domain.RoadType]int `json:"by_street_type"`
	ByZone          map[string]int          `json:"by_zone"`
}

func Load(path string) (*Repository, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var city domain.CityMap
	if err := json.Unmarshal(raw, &city); err != nil {
		return nil, err
	}

	repo := &Repository{
		city:      city,
		nodesByID: map[string]domain.Node{},
		edgesByID: map[string]domain.Edge{},
		adjacency: map[string][]domain.Edge{},
	}

	for _, node := range city.Nodes {
		if node.ID == "" {
			return nil, fmt.Errorf("node sem id")
		}
		if _, exists := repo.nodesByID[node.ID]; exists {
			return nil, fmt.Errorf("node duplicado: %s", node.ID)
		}
		repo.nodesByID[node.ID] = node
	}

	for _, edge := range city.Edges {
		if edge.ID == "" {
			return nil, fmt.Errorf("edge sem id")
		}
		if _, exists := repo.edgesByID[edge.ID]; exists {
			return nil, fmt.Errorf("edge duplicada: %s", edge.ID)
		}
		if _, ok := repo.nodesByID[edge.From]; !ok {
			return nil, fmt.Errorf("edge %s aponta origem inexistente: %s", edge.ID, edge.From)
		}
		if _, ok := repo.nodesByID[edge.To]; !ok {
			return nil, fmt.Errorf("edge %s aponta destino inexistente: %s", edge.ID, edge.To)
		}
		repo.edgesByID[edge.ID] = edge
		repo.adjacency[edge.From] = append(repo.adjacency[edge.From], edge)
		if edge.Bidirectional {
			reverse := edge
			reverse.From, reverse.To = edge.To, edge.From
			repo.adjacency[edge.To] = append(repo.adjacency[edge.To], reverse)
		}
	}

	if err := repo.ValidateUrbanRules(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) City() domain.CityMap {
	return r.city
}

func (r *Repository) Nodes() []domain.Node {
	return r.city.Nodes
}

func (r *Repository) Stats() Stats {
	stats := Stats{
		Name:         r.city.Name,
		NodeCount:    len(r.city.Nodes),
		EdgeCount:    len(r.city.Edges),
		ByStreetType: map[domain.RoadType]int{},
		ByZone:       map[string]int{},
	}
	var safetySum float64
	for _, node := range r.city.Nodes {
		stats.ByZone[node.Zone]++
	}
	for _, edge := range r.city.Edges {
		stats.TotalDistance += edge.Distance
		safetySum += edge.SafetyScore
		stats.ByStreetType[edge.StreetType]++
		if edge.StreetType == domain.RoadCycleTrack || edge.StreetType == domain.RoadBikeLane {
			stats.CyclingDistance += edge.Distance
		}
	}
	if len(r.city.Edges) > 0 {
		stats.AverageSafety = math.Round((safetySum/float64(len(r.city.Edges)))*100) / 100
	}
	stats.TotalDistance = math.Round(stats.TotalDistance*100) / 100
	stats.CyclingDistance = math.Round(stats.CyclingDistance*100) / 100
	return stats
}

func (r *Repository) Node(id string) (domain.Node, bool) {
	node, ok := r.nodesByID[id]
	return node, ok
}

func (r *Repository) SearchNodes(query string, limit int) []SearchResult {
	query = normalize(query)
	if query == "" {
		return nil
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	results := []SearchResult{}
	for _, node := range r.city.Nodes {
		name := normalize(node.Name)
		if !strings.Contains(name, query) {
			continue
		}
		if strings.Contains(name, "acesso local sem nome") && len(r.adjacency[node.ID]) < 4 {
			continue
		}
		score := len(r.adjacency[node.ID])
		if strings.Contains(name, query) {
			score += 10
		}
		if strings.HasPrefix(name, query) {
			score += 20
		}
		if !strings.Contains(name, "sem nome") {
			score += 10
		}
		results = append(results, SearchResult{ID: node.ID, Name: node.Name, X: node.X, Y: node.Y, Zone: node.Zone, Score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Name < results[j].Name
		}
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}

func normalize(value string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a",
		"é", "e", "ê", "e", "è", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ç", "c",
	)
	return replacer.Replace(strings.ToLower(strings.TrimSpace(value)))
}

func (r *Repository) Neighbors(nodeID string) []domain.Edge {
	return r.adjacency[nodeID]
}

func (r *Repository) NearestNode(pos domain.CyclistPosition) domain.Node {
	var nearest domain.Node
	best := math.Inf(1)
	for _, node := range r.city.Nodes {
		d := math.Hypot(node.X-pos.X, node.Y-pos.Y)
		if d < best {
			best = d
			nearest = node
		}
	}
	return nearest
}

func (r *Repository) ValidateUrbanRules() error {
	if r.city.Name == "" {
		return fmt.Errorf("mapa sem nome")
	}
	if len(r.city.Nodes) == 0 {
		return fmt.Errorf("mapa sem nodes")
	}
	for _, node := range r.city.Nodes {
		if node.Name == "" {
			return fmt.Errorf("node %s sem nome", node.ID)
		}
		if node.Zone == "" {
			return fmt.Errorf("node %s sem zona", node.ID)
		}
	}
	for _, edge := range r.city.Edges {
		if edge.Name == "" {
			return fmt.Errorf("edge %s sem nome", edge.ID)
		}
		if edge.Distance <= 0 {
			return fmt.Errorf("edge %s com distancia invalida", edge.ID)
		}
		if edge.SafetyScore < 0 || edge.SafetyScore > 1 {
			return fmt.Errorf("edge %s com safety_score invalido", edge.ID)
		}
		if edge.SpeedLimit <= 0 {
			return fmt.Errorf("edge %s com speed_limit invalido", edge.ID)
		}
		if edge.TrafficLevel < 0 || edge.TrafficLevel > 5 {
			return fmt.Errorf("edge %s com traffic_level invalido", edge.ID)
		}
		if !validRoadType(edge.StreetType) {
			return fmt.Errorf("edge %s com street_type invalido: %s", edge.ID, edge.StreetType)
		}

	}
	if !r.connected() {
		return fmt.Errorf("mapa possui componentes desconectados")
	}
	return nil
}

func validRoadType(t domain.RoadType) bool {
	switch t {
	case domain.RoadCycleTrack, domain.RoadBikeLane, domain.RoadResidential, domain.RoadCommonStreet, domain.RoadBusyAvenue, domain.RoadHighway, domain.RoadUnsafe:
		return true
	default:
		return false
	}
}

func (r *Repository) connected() bool {
	if len(r.city.Nodes) == 0 {
		return true
	}
	seen := map[string]bool{}
	queue := []string{r.city.Nodes[0].ID}
	seen[queue[0]] = true
	undirected := map[string][]string{}
	for _, edge := range r.city.Edges {
		undirected[edge.From] = append(undirected[edge.From], edge.To)
		undirected[edge.To] = append(undirected[edge.To], edge.From)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, next := range undirected[current] {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return len(seen) == len(r.city.Nodes)
}
