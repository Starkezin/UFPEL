package simulation

import (
	"math"

	"rei-julian/internal/domain"
)

type State struct {
	Active     bool    `json:"active"`
	EdgeIndex  int     `json:"edge_index"`
	Offset     float64 `json:"offset"`
	StepMeters float64 `json:"step_meters"`
	Done       bool    `json:"done"`
}

type StepResult struct {
	Position domain.CyclistPosition `json:"position"`
	Node     domain.Node            `json:"node"`
	State    State                  `json:"state"`
}

func NewState(stepMeters float64) State {
	if stepMeters <= 0 {
		stepMeters = 120
	}
	return State{Active: true, StepMeters: stepMeters}
}

func Step(route domain.Route, state State) StepResult {
	if len(route.Nodes) == 0 {
		state.Done = true
		state.Active = false
		return StepResult{State: state}
	}
	if len(route.Edges) == 0 || len(route.Nodes) == 1 {
		state.Done = true
		state.Active = false
		return StepResult{Position: domain.CyclistPosition{X: route.Nodes[0].X, Y: route.Nodes[0].Y}, Node: route.Nodes[0], State: state}
	}
	if !state.Active || state.Done {
		last := route.Nodes[len(route.Nodes)-1]
		state.Done = true
		state.Active = false
		return StepResult{Position: domain.CyclistPosition{X: last.X, Y: last.Y}, Node: last, State: state}
	}

	remaining := state.StepMeters
	for remaining > 0 && state.EdgeIndex < len(route.Edges) {
		edgeDistance := math.Max(route.Edges[state.EdgeIndex].Distance, 1)
		left := edgeDistance - state.Offset
		if remaining < left {
			state.Offset += remaining
			remaining = 0
			break
		}
		remaining -= left
		state.EdgeIndex++
		state.Offset = 0
	}
	if state.EdgeIndex >= len(route.Edges) {
		state.EdgeIndex = len(route.Edges) - 1
		state.Offset = route.Edges[state.EdgeIndex].Distance
		state.Done = true
		state.Active = false
	}

	from := route.Nodes[state.EdgeIndex]
	to := route.Nodes[state.EdgeIndex+1]
	edgeDistance := math.Max(route.Edges[state.EdgeIndex].Distance, 1)
	t := math.Max(0, math.Min(1, state.Offset/edgeDistance))
	pos := domain.CyclistPosition{X: from.X + (to.X-from.X)*t, Y: from.Y + (to.Y-from.Y)*t}
	nearest := from
	if t >= 0.5 {
		nearest = to
	}
	return StepResult{Position: pos, Node: nearest, State: state}
}
