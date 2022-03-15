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
			{"data":{"id":"srv1","value":"srv1","name":"srv1"}},
			{"data":{"id":"101","value":"CPU","name":"CPU"}},
			{"data":{"id":"102","value":"Memory","name":"Memory"}},
			{"data":{"id":"103","value":"Disk","name":"Disk"}},
			{"data":{"id":"104","value":"Ping","name":"Ping"}},

			{"data":{"id":"db1","value":"db1","name":"db1"}},
			{"data":{"id":"201","value":"CPU","name":"CPU"}},
			{"data":{"id":"202","value":"Memory","name":"Memory"}},
			{"data":{"id":"203","value":"Disk","name":"Disk"}},
			{"data":{"id":"204","value":"Ping","name":"Ping"}},
			{"data":{"id":"205","value":"DBEngine","name":"DBEngine"}},

			{"data":{"id":"backend1","value":"backend1","name":"backend1"}},
			{"data":{"id":"301","value":"CPU","name":"CPU"}},
			{"data":{"id":"302","value":"Memory","name":"Memory"}},
			{"data":{"id":"303","value":"Disk","name":"Disk"}},
			{"data":{"id":"304","value":"Ping","name":"Ping"}},
			{"data":{"id":"306","value":"Proc","name":"Proc"}},
			{"data":{"id":"307","value":"DBConnection","name":"DBConnection"}},

			{"data":{"id":"frontend1","value":"frontend1","name":"frontend1"}},
			{"data":{"id":"401","value":"CPU","name":"CPU"}},
			{"data":{"id":"402","value":"Memory","name":"Memory"}},
			{"data":{"id":"403","value":"Disk","name":"Disk"}},
			{"data":{"id":"404","value":"Ping","name":"Ping"}},
			{"data":{"id":"406","value":"Proc","name":"Proc"}},
			{"data":{"id":"408","value":"BackendConnection","name":"BackendConnection"}}
		],
		"edges":[
			{"data":{"source":"srv1","target":"101"}},
			{"data":{"source":"srv1","target":"102"}},
			{"data":{"source":"srv1","target":"103"}},
			{"data":{"source":"srv1","target":"104"}},

			{"data":{"source":"db1","target":"201"}},
			{"data":{"source":"db1","target":"202"}},
			{"data":{"source":"db1","target":"203"}},
			{"data":{"source":"db1","target":"204"}},
			{"data":{"source":"db1","target":"205"}},

			{"data":{"source":"backend1","target":"301"}},
			{"data":{"source":"backend1","target":"302"}},
			{"data":{"source":"backend1","target":"303"}},
			{"data":{"source":"backend1","target":"304"}},
			{"data":{"source":"backend1","target":"306"}},
			{"data":{"source":"backend1","target":"307"}},

			{"data":{"source":"frontend1","target":"401"}},
			{"data":{"source":"frontend1","target":"402"}},
			{"data":{"source":"frontend1","target":"403"}},
			{"data":{"source":"frontend1","target":"404"}},
			{"data":{"source":"frontend1","target":"406"}},
			{"data":{"source":"frontend1","target":"408"}},

			{"data":{"source":"backend1","target":"db1"}},
			{"data":{"source":"frontend1","target":"backend1"}}
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
