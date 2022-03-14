package main

import (
	"encoding/json"
	"fmt"

	"github.com/fschuetz04/simgo"
)

// Architecture store the different servers of our application
type Architecture struct {
	DBs       []*Database
	Backends  []*Backend
	Frontends []*Frontend
	// Monkeys are functions that will "sabotage" the architecture, triggering alarms
	Monkeys []func(simgo.Process)

	// mon connection to the monitoring system
	mon MonitorSystem

	// sim is the simulation object to use a fake time and make the simulation instant
	sim *simgo.Simulation
}

// CytoscapeGraph is the JSON representation of the architecture in the Cytoscape JSON format
type CytoscapeGraph struct {
	Data       []interface{} `json:"data"`
	Directed   bool          `json:"directed"`
	Multigraph bool          `json:"multigraph"`
	Elements   CSElements    `json:"elements"`
}

type CSElements struct {
	Nodes []CSNode `json:"nodes"`
	Edges []CSEdge `json:"edges"`
}

type CSNodeData struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

type CSNode struct {
	Data CSNodeData `json:"data"`
}

type CSEdge struct {
	Data CSEdgeData `json:"data"`
}

type CSEdgeData struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Run start the monitoring of each server
func (a *Architecture) Start(sim_duration float64) {
	a.sim = &simgo.Simulation{}

	for _, db := range a.DBs {
		a.sim.ProcessReflect(Run, db)
	}

	for _, backend := range a.Backends {
		a.sim.ProcessReflect(Run, backend)
	}

	for _, frontend := range a.Frontends {
		a.sim.ProcessReflect(Run, frontend)
	}

	for _, monkey := range a.Monkeys {
		a.sim.Process(monkey)
	}

	a.sim.RunUntil(sim_duration)
}

func (a *Architecture) NewDatabase(name string) *Database {
	d := NewDatabase(name, a.mon)
	a.AddDB(d)
	return d
}

func (a *Architecture) NewBackend(name string, db *Database) *Backend {
	b := NewBackend(name, db, a.mon)
	a.AddBackend(b)
	return b
}

func (a *Architecture) NewFrontend(name string, backend *Backend) *Frontend {
	f := NewFrontend(name, backend, a.mon)
	a.AddFrontend(f)
	return f
}

func (a *Architecture) AddDB(db *Database) {
	a.DBs = append(a.DBs, db)
}

func (a *Architecture) AddBackend(backend *Backend) {
	a.Backends = append(a.Backends, backend)
}

func (a *Architecture) AddFrontend(frontend *Frontend) {
	a.Frontends = append(a.Frontends, frontend)
}

func (a *Architecture) AddMonkey(monkey func(simgo.Process)) {
	a.Monkeys = append(a.Monkeys, monkey)
}

// CytoscapeGraph return a JSON representation of the architecture in the
// cytoscape format
func (a *Architecture) CytoscapeGraph() string {
	nodes := []CSNode{}

	for _, db := range a.DBs {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				ID:    db.Name,
				Value: db.Name,
				Name:  db.Name,
			},
		})
	}

	for _, backend := range a.Backends {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				ID:    backend.Name,
				Value: backend.Name,
				Name:  backend.Name,
			},
		})
	}

	for _, frontend := range a.Frontends {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				ID:    frontend.Name,
				Value: frontend.Name,
				Name:  frontend.Name,
			},
		})
	}

	edges := []CSEdge{}

	for _, backend := range a.Backends {
		edges = append(edges, CSEdge{
			Data: CSEdgeData{
				Source: backend.Name,
				Target: backend.DBEngine.Name,
			},
		})
	}

	for _, frontend := range a.Frontends {
		edges = append(edges, CSEdge{
			Data: CSEdgeData{
				Source: frontend.Name,
				Target: frontend.Backend.Name,
			},
		})
	}

	cs := CytoscapeGraph{
		Data:     []interface{}{},
		Directed: true,
		Elements: CSElements{
			Nodes: nodes,
			Edges: edges,
		},
	}

	g, err := json.Marshal(cs)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s", g)
}
