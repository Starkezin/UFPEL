package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	earthRadiusMeters = 6371000.0
	stepMeters        = 10.0
	historyLimit      = 200
)

type AgentConfig struct {
	ID        string
	Name      string
	Color     string
	StartLat  float64
	StartLon  float64
}

type Position struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Step      int       `json:"step"`
}

type AgentState struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Color    string     `json:"color"`
	Position Position   `json:"position"`
	History  []Position `json:"history"`
}

type Snapshot struct {
	Agents    []AgentState `json:"agents"`
	Timestamp time.Time    `json:"timestamp"`
}

type Tracker struct {
	mu      sync.RWMutex
	agents  map[string]AgentState
	order   []string
	clients map[chan Snapshot]struct{}
}

func NewTracker(configs []AgentConfig) *Tracker {
	agents := make(map[string]AgentState, len(configs))
	order := make([]string, 0, len(configs))
	now := time.Now()

	for _, cfg := range configs {
		pos := Position{
			Latitude:  cfg.StartLat,
			Longitude: cfg.StartLon,
			Timestamp: now,
			Step:      0,
		}

		agents[cfg.ID] = AgentState{
			ID:       cfg.ID,
			Name:     cfg.Name,
			Color:    cfg.Color,
			Position: pos,
			History:  []Position{pos},
		}
		order = append(order, cfg.ID)
	}

	return &Tracker{
		agents:  agents,
		order:   order,
		clients: make(map[chan Snapshot]struct{}),
	}
}

func (t *Tracker) Snapshot() Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.snapshotLocked()
}

func (t *Tracker) snapshotLocked() Snapshot {
	agents := make([]AgentState, 0, len(t.order))
	for _, id := range t.order {
		agent := t.agents[id]
		historyCopy := append([]Position(nil), agent.History...)
		agent.History = historyCopy
		agents = append(agents, agent)
	}

	return Snapshot{
		Agents:    agents,
		Timestamp: time.Now(),
	}
}

func (t *Tracker) UpdatePosition(agentID string, p Position) {
	t.mu.Lock()

	agent, ok := t.agents[agentID]
	if !ok {
		t.mu.Unlock()
		return
	}

	agent.Position = p
	agent.History = append(agent.History, p)
	if len(agent.History) > historyLimit {
		agent.History = append([]Position(nil), agent.History[len(agent.History)-historyLimit:]...)
	}
	t.agents[agentID] = agent

	snapshot := t.snapshotLocked()
	clients := make([]chan Snapshot, 0, len(t.clients))
	for ch := range t.clients {
		clients = append(clients, ch)
	}
	t.mu.Unlock()

	for _, ch := range clients {
		select {
		case ch <- snapshot:
		default:
			// cliente lento, ignora este envio
		}
	}
}

func (t *Tracker) Subscribe() chan Snapshot {
	ch := make(chan Snapshot, 8)

	t.mu.Lock()
	t.clients[ch] = struct{}{}
	t.mu.Unlock()

	return ch
}

func (t *Tracker) Unsubscribe(ch chan Snapshot) {
	t.mu.Lock()
	delete(t.clients, ch)
	t.mu.Unlock()
	close(ch)
}

// movePoint desloca um ponto geografico por certa distancia em uma direcao dada.
// Formula de navegacao sobre esfera.
func movePoint(latDeg, lonDeg, distanceMeters, bearingRad float64) (float64, float64) {
	lat1 := degreesToRadians(latDeg)
	lon1 := degreesToRadians(lonDeg)
	angularDistance := distanceMeters / earthRadiusMeters

	lat2 := math.Asin(
		math.Sin(lat1)*math.Cos(angularDistance) +
			math.Cos(lat1)*math.Sin(angularDistance)*math.Cos(bearingRad),
	)

	lon2 := lon1 + math.Atan2(
		math.Sin(bearingRad)*math.Sin(angularDistance)*math.Cos(lat1),
		math.Cos(angularDistance)-math.Sin(lat1)*math.Sin(lat2),
	)

	return radiansToDegrees(lat2), normalizeLongitude(radiansToDegrees(lon2))
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180.0
}

func radiansToDegrees(r float64) float64 {
	return r * 180.0 / math.Pi
}

func normalizeLongitude(lon float64) float64 {
	for lon > 180.0 {
		lon -= 360.0
	}
	for lon < -180.0 {
		lon += 360.0
	}
	return lon
}

func startSimulation(tracker *Tracker, cfg AgentConfig, seed int64) {
	rng := rand.New(rand.NewSource(seed))

	for {
		delay := time.Duration(4+rng.Float64()*2) * time.Second
		time.Sleep(delay)

		snapshot := tracker.Snapshot()

		var current AgentState
		found := false
		for _, agent := range snapshot.Agents {
			if agent.ID == cfg.ID {
				current = agent
				found = true
				break
			}
		}

		if !found {
			continue
		}

		bearing := rng.Float64() * 2 * math.Pi
		newLat, newLon := movePoint(current.Position.Latitude, current.Position.Longitude, stepMeters, bearing)

		next := Position{
			Latitude:  newLat,
			Longitude: newLon,
			Timestamp: time.Now(),
			Step:      current.Position.Step + 1,
		}

		tracker.UpdatePosition(cfg.ID, next)

		log.Printf(
			"agent=%s step=%d lat=%.8f lon=%.8f",
			cfg.ID,
			next.Step,
			next.Latitude,
			next.Longitude,
		)
	}
}

func positionsHandler(tracker *Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")

		snapshot := tracker.Snapshot()
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			http.Error(w, "erro ao serializar posicoes", http.StatusInternalServerError)
			return
		}
	}
}

func sseHandler(tracker *Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming nao suportado", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Accel-Buffering", "no")

		ch := tracker.Subscribe()
		defer tracker.Unsubscribe(ch)

		initial := tracker.Snapshot()
		initialJSON, _ := json.Marshal(initial)
		fmt.Fprintf(w, "data: %s\n\n", initialJSON)
		flusher.Flush()

		ctx := r.Context()

		for {
			select {
			case <-ctx.Done():
				return
			case snapshot := <-ch:
				payload, err := json.Marshal(snapshot)
				if err != nil {
					continue
				}
				fmt.Fprintf(w, "data: %s\n\n", payload)
				flusher.Flush()
			}
		}
	}
}

func main() {
	configs := []AgentConfig{
		{
			ID:       "onibus-1",
			Name:     "Onibus 1",
			Color:    "#d62828",
			StartLat: -31.770687426923516,
			StartLon: -52.34135057529372,
		},
		{
			ID:       "onibus-2",
			Name:     "Onibus 2",
			Color:    "#1d4ed8",
			StartLat: -31.771587426923514,
			StartLon: -52.33995057529372,
		},
	}

	tracker := NewTracker(configs)

	for i, cfg := range configs {
		go startSimulation(tracker, cfg, time.Now().UnixNano()+int64(i*1000))
	}

	http.HandleFunc("/positions", positionsHandler(tracker))
	http.HandleFunc("/stream", sseHandler(tracker))

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Servidor em http://localhost:8080")
	log.Println("GET /positions -> posicoes atuais e historico em JSON")
	log.Println("GET /stream    -> stream SSE com atualizacoes")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
