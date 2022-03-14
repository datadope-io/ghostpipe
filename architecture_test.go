package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCytoscapeGraph(t *testing.T) {
	mon := &fakeMonSys{}

	a := Architecture{}

	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}
	frontend1 := &Frontend{Server: Server{Name: "frontend1", mon: mon}, Backend: backend1}

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
			{"data":{"id":"db1","value":"db1","name":"db1"}},
			{"data":{"id":"backend1","value":"backend1","name":"backend1"}},
			{"data":{"id":"frontend1","value":"frontend1","name":"frontend1"}}
		],
		"edges":[
			{"data":{"source":"db1","target":"backend1"}},
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
