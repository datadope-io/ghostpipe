package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCytoscapeGraph(t *testing.T) {
	mon := &fakeMonSys{}

	a := Architecture{mon: mon}

	srv1 := &Server{Name: "srv1"}
	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}
	frontend1 := &Frontend{Server: Server{Name: "frontend1", mon: mon}, Backend: backend1}

	a.AddServer(srv1)
	a.AddDB(db1)
	a.AddBackend(backend1)
	a.AddFrontend(frontend1)

	expectedJSON := `
{
	"data":[],
	"directed":true,
	"multigraph":false,
	"elements": {
		"nodes":[
			{"data":{"id":"srv1","name":"srv1","type":"server"}},
			{"data":{"id":"101","name":"CPU","type":"alarm"}},
			{"data":{"id":"102","name":"Memory","type":"alarm"}},
			{"data":{"id":"103","name":"Disk","type":"alarm"}},
			{"data":{"id":"104","name":"Ping","type":"alarm"}},

			{"data":{"id":"db1","name":"db1","type":"db"}},
			{"data":{"id":"201","name":"CPU","type":"alarm"}},
			{"data":{"id":"202","name":"Memory","type":"alarm"}},
			{"data":{"id":"203","name":"Disk","type":"alarm"}},
			{"data":{"id":"204","name":"Ping","type":"alarm"}},
			{"data":{"id":"205","name":"DBEngine","type":"alarm"}},

			{"data":{"id":"backend1","name":"backend1","type":"backend"}},
			{"data":{"id":"301","name":"CPU","type":"alarm"}},
			{"data":{"id":"302","name":"Memory","type":"alarm"}},
			{"data":{"id":"303","name":"Disk","type":"alarm"}},
			{"data":{"id":"304","name":"Ping","type":"alarm"}},
			{"data":{"id":"306","name":"Proc","type":"alarm"}},
			{"data":{"id":"307","name":"DBConnection","type":"alarm"}},

			{"data":{"id":"frontend1","name":"frontend1","type":"frontend"}},
			{"data":{"id":"401","name":"CPU","type":"alarm"}},
			{"data":{"id":"402","name":"Memory","type":"alarm"}},
			{"data":{"id":"403","name":"Disk","type":"alarm"}},
			{"data":{"id":"404","name":"Ping","type":"alarm"}},
			{"data":{"id":"406","name":"Proc","type":"alarm"}},
			{"data":{"id":"408","name":"BackendConnection","type":"alarm"}}
		],
		"edges":[
			{"data":{"source":"srv1","target":"101","epq":0}},
			{"data":{"source":"srv1","target":"102","epq":0}},
			{"data":{"source":"srv1","target":"103","epq":0}},
			{"data":{"source":"srv1","target":"104","epq":0}},

			{"data":{"source":"db1","target":"201","epq":0}},
			{"data":{"source":"db1","target":"202","epq":0}},
			{"data":{"source":"db1","target":"203","epq":0}},
			{"data":{"source":"db1","target":"204","epq":0}},
			{"data":{"source":"db1","target":"205","epq":0}},

			{"data":{"source":"backend1","target":"301","epq":0}},
			{"data":{"source":"backend1","target":"302","epq":0}},
			{"data":{"source":"backend1","target":"303","epq":0}},
			{"data":{"source":"backend1","target":"304","epq":0}},
			{"data":{"source":"backend1","target":"306","epq":0}},
			{"data":{"source":"backend1","target":"307","epq":0}},

			{"data":{"source":"frontend1","target":"401","epq":0}},
			{"data":{"source":"frontend1","target":"402","epq":0}},
			{"data":{"source":"frontend1","target":"403","epq":0}},
			{"data":{"source":"frontend1","target":"404","epq":0}},
			{"data":{"source":"frontend1","target":"406","epq":0}},
			{"data":{"source":"frontend1","target":"408","epq":0}},

			{"data":{"source":"backend1","target":"db1","epq":0}},
			{"data":{"source":"frontend1","target":"backend1","epq":0}}
		]
	}
}
`
	// Convert to one line format
	expectedJSON = strings.Replace(expectedJSON, "\n", "", -1)
	expectedJSON = strings.Replace(expectedJSON, "\t", "", -1)
	expectedJSON = strings.Replace(expectedJSON, " ", "", -1)

	assert.Equal(t, strings.Replace(expectedJSON, "\n", "", -1), a.CytoscapeGraph())
}
