package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"rei-julian/internal/domain"
	"rei-julian/internal/mapdata"
	"rei-julian/internal/routing"
	"rei-julian/internal/simulation"
	"rei-julian/internal/tracking"
)

type Server struct {
	repo     *mapdata.Repository
	routing  *routing.Service
	tracking *tracking.Service

	mu       sync.Mutex
	sessions map[string]*sessionState
	nextID   int64
}

type sessionState struct {
	CurrentRoute *domain.Route
	Simulation   simulation.State
}

func NewServer(repo *mapdata.Repository, routingService *routing.Service, trackingService *tracking.Service) *Server {
	return &Server{repo: repo, routing: routingService, tracking: trackingService, sessions: map[string]*sessionState{}}
}

func (s *Server) Routes(staticDir string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /map", s.cityMap)
	mux.HandleFunc("GET /stats", s.stats)
	mux.HandleFunc("GET /nodes", s.nodes)
	mux.HandleFunc("GET /search", s.search)
	mux.HandleFunc("POST /routes", s.createRoute)
	mux.HandleFunc("POST /routes/compare", s.compareRoutes)
	mux.HandleFunc("POST /tracking/update", s.updateTracking)
	mux.HandleFunc("POST /simulation/start", s.startSimulation)
	mux.HandleFunc("POST /simulation/step", s.stepSimulation)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) cityMap(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.repo.City())
}

func (s *Server) nodes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.repo.Nodes())
}

func (s *Server) stats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.repo.Stats())
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	writeJSON(w, http.StatusOK, map[string]any{"results": s.repo.SearchNodes(r.URL.Query().Get("q"), limit)})
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID     string `json:"session_id"`
		OriginID      string `json:"origin_id"`
		DestinationID string `json:"destination_id"`
		Profile       string `json:"profile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}

	route, err := s.routing.Calculate(req.OriginID, req.DestinationID, domain.RouteProfile(req.Profile))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sessionID := s.saveRoute(req.SessionID, route)
	writeJSON(w, http.StatusOK, map[string]any{"session_id": sessionID, "route": route})
}

func (s *Server) compareRoutes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginID      string `json:"origin_id"`
		DestinationID string `json:"destination_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	profiles := []domain.RouteProfile{domain.ProfileFastest, domain.ProfileBalanced, domain.ProfileSafest}
	routes := map[string]domain.Route{}
	for _, profile := range profiles {
		route, err := s.routing.Calculate(req.OriginID, req.DestinationID, profile)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		routes[string(profile)] = route
	}
	writeJSON(w, http.StatusOK, map[string]any{"routes": routes})
}

func (s *Server) updateTracking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string  `json:"session_id"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		NodeHint  string  `json:"node_hint,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	state, ok := s.session(req.SessionID)
	if !ok || state.CurrentRoute == nil {
		writeError(w, http.StatusBadRequest, "calcule uma rota antes de atualizar o tracking")
		return
	}

	result, err := s.tracking.Update(domain.CyclistPosition{X: req.X, Y: req.Y, NodeHint: req.NodeHint}, *state.CurrentRoute)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if result.NewRoute != nil {
		s.mu.Lock()
		state.CurrentRoute = result.NewRoute
		s.mu.Unlock()
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) startSimulation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID  string  `json:"session_id"`
		StepMeters float64 `json:"step_meters"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	state, ok := s.session(req.SessionID)
	if !ok || state.CurrentRoute == nil {
		writeError(w, http.StatusBadRequest, "calcule uma rota antes de iniciar a simulacao")
		return
	}
	s.mu.Lock()
	state.Simulation = simulation.NewState(req.StepMeters)
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"session_id": req.SessionID, "simulation": state.Simulation})
}

func (s *Server) stepSimulation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	state, ok := s.session(req.SessionID)
	if !ok || state.CurrentRoute == nil {
		writeError(w, http.StatusBadRequest, "calcule uma rota antes de simular")
		return
	}
	result := simulation.Step(*state.CurrentRoute, state.Simulation)
	s.mu.Lock()
	state.Simulation = result.State
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) saveRoute(sessionID string, route domain.Route) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sessionID == "" || s.sessions[sessionID] == nil {
		s.nextID++
		sessionID = fmt.Sprintf("s-%d-%d", time.Now().UnixNano(), s.nextID)
		s.sessions[sessionID] = &sessionState{}
	}
	s.sessions[sessionID].CurrentRoute = &route
	s.sessions[sessionID].Simulation = simulation.State{}
	return sessionID
}

func (s *Server) session(sessionID string) (*sessionState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.sessions[sessionID]
	return state, ok
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
