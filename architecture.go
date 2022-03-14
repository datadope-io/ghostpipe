package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Architecture store the different servers of our application
type Architecture struct {
	DBs       []*Database
	Backends  []*Backend
	Frontends []*Frontend

	// mon connection to the monitoring system
	mon MonitorSystem

	// ctx context to signal servers to stop
	ctxCancel context.CancelFunc

	wg sync.WaitGroup
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
func (a *Architecture) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	a.ctxCancel = cancel

	for _, db := range a.DBs {
		go func(d *Database) {
			a.wg.Add(1)
			defer a.wg.Done()
			Run(ctx, d)
		}(db)
	}

	for _, backend := range a.Backends {
		go func(b *Backend) {
			a.wg.Add(1)
			defer a.wg.Done()
			Run(ctx, b)
		}(backend)
	}

	for _, frontend := range a.Frontends {
		go func(f *Frontend) {
			a.wg.Add(1)
			defer a.wg.Done()
			Run(ctx, f)
		}(frontend)
	}
}

// Stop signal all servers to stop and wait for them to stop
func (a *Architecture) Stop() {
	a.ctxCancel()
	a.wg.Wait()
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
