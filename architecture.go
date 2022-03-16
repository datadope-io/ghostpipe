package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/fschuetz04/simgo"
)

// Architecture store the different servers of our application
type Architecture struct {
	Servers   []*Server
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
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type CSNode struct {
	Data CSNodeData `json:"data"`
}

type CSEdge struct {
	Data CSEdgeData `json:"data"`
}

type CSEdgeData struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Type   string  `json:"type"`
	Epq    float64 `json:"epq"`
}

// Run start the monitoring of each server
func (a *Architecture) Start(sim_duration float64) {
	a.sim = &simgo.Simulation{}

	// Shuffle the servers to start in random order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a.Servers), func(i, j int) { a.Servers[i], a.Servers[j] = a.Servers[j], a.Servers[i] })
	rand.Shuffle(len(a.DBs), func(i, j int) { a.DBs[i], a.DBs[j] = a.DBs[j], a.DBs[i] })
	rand.Shuffle(len(a.Backends), func(i, j int) { a.Backends[i], a.Backends[j] = a.Backends[j], a.Backends[i] })
	rand.Shuffle(len(a.Frontends), func(i, j int) { a.Frontends[i], a.Frontends[j] = a.Frontends[j], a.Frontends[i] })
	rand.Shuffle(len(a.Monkeys), func(i, j int) { a.Monkeys[i], a.Monkeys[j] = a.Monkeys[j], a.Monkeys[i] })

	for _, server := range a.Servers {
		a.sim.ProcessReflect(Run, server)
	}

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

func (a *Architecture) NewServer(name string) *Server {
	s := NewServer(name, a.mon)
	a.AddServer(s)
	return s
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

func (a *Architecture) AddServer(server *Server) {
	a.Servers = append(a.Servers, server)
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
	edges := []CSEdge{}

	for _, server := range a.Servers {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				Value: server.Name,
				Name:  server.Name,
				Type:  "server",
			},
		})

		// Add alarms nodes
		for _, alarm := range []string{"CPU", "Memory", "Disk", "Ping"} {
			id := fmt.Sprintf("%d", a.mon.generateEventID(server.Name, alarm))
			nodes = append(nodes, CSNode{
				Data: CSNodeData{
					Value: id,
					Name:  fmt.Sprintf("%s-%s", server.Name, alarm),
					Type:  "alarm",
				},
			})

			// Add edges from the server to the alarms
			edges = append(edges, CSEdge{
				Data: CSEdgeData{
					Source: server.Name,
					Target: id,
					Type:   "trigger",
				},
			})
		}

	}

	for _, db := range a.DBs {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				Value: db.Name,
				Name:  db.Name,
				Type:  "db",
			},
		})

		// Add alarms nodes
		for _, alarm := range []string{"CPU", "Memory", "Disk", "Ping", "DBEngine"} {
			id := fmt.Sprintf("%d", a.mon.generateEventID(db.Name, alarm))
			nodes = append(nodes, CSNode{
				Data: CSNodeData{
					Value: id,
					Name:  fmt.Sprintf("%s-%s", db.Name, alarm),
					Type:  "alarm",
				},
			})

			// Add edges from the server to the alarms
			edges = append(edges, CSEdge{
				Data: CSEdgeData{
					Source: db.Name,
					Target: id,
					Type:   "trigger",
				},
			})
		}
	}

	for _, backend := range a.Backends {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				Value: backend.Name,
				Name:  backend.Name,
				Type:  "backend",
			},
		})

		// Add alarms nodes
		for _, alarm := range []string{"CPU", "Memory", "Disk", "Ping", "Proc", "DBConnection"} {
			id := fmt.Sprintf("%d", a.mon.generateEventID(backend.Name, alarm))
			nodes = append(nodes, CSNode{
				Data: CSNodeData{
					Value: id,
					Name:  fmt.Sprintf("%s-%s", backend.Name, alarm),
					Type:  "alarm",
				},
			})

			// Add edges from the server to the alarms
			edges = append(edges, CSEdge{
				Data: CSEdgeData{
					Source: backend.Name,
					Target: id,
					Type:   "trigger",
				},
			})
		}
	}

	for _, frontend := range a.Frontends {
		nodes = append(nodes, CSNode{
			Data: CSNodeData{
				Value: frontend.Name,
				Name:  frontend.Name,
				Type:  "frontend",
			},
		})

		// Add alarms nodes
		for _, alarm := range []string{"CPU", "Memory", "Disk", "Ping", "Proc", "BackendConnection"} {
			id := fmt.Sprintf("%d", a.mon.generateEventID(frontend.Name, alarm))
			nodes = append(nodes, CSNode{
				Data: CSNodeData{
					Value: id,
					Name:  fmt.Sprintf("%s-%s", frontend.Name, alarm),
					Type:  "alarm",
				},
			})
			// Add edges from the server to the alarms
			edges = append(edges, CSEdge{
				Data: CSEdgeData{
					Source: frontend.Name,
					Target: id,
					Type:   "trigger",
				},
			})
		}
	}

	// Add links between different servers
	for _, backend := range a.Backends {
		edges = append(edges, CSEdge{
			Data: CSEdgeData{
				Source: backend.Name,
				Target: backend.DBEngine.Name,
				Type:   "connect",
			},
		})
	}

	for _, frontend := range a.Frontends {
		edges = append(edges, CSEdge{
			Data: CSEdgeData{
				Source: frontend.Name,
				Target: frontend.Backend.Name,
				Type:   "connect",
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
