package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphML(t *testing.T) {
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

	expectedXML := `  <graphml xmlns="http://graphml.graphdrawing.org/xmlns" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
      <key id="d0" for="node" attr.name="id" attr.type="string"></key>
      <key id="d1" for="node" attr.name="label" attr.type="string"></key>
      <key id="d2" for="node" attr.name="name" attr.type="string"></key>
      <key id="d3" for="node" attr.name="type" attr.type="string"></key>
      <key id="d4" for="edge" attr.name="type" attr.type="string"></key>
      <key id="d5" for="edge" attr.name="weight" attr.type="int"></key>
      <graph id="g0" edgedefault="undirected">
          <desc>ghostpipe-graph</desc>
          <node id="n0">
              <desc>srv1</desc>
              <data key="d0">srv1</data>
              <data key="d1">srv1</data>
              <data key="d2">srv1</data>
              <data key="d3">server</data>
          </node>
          <node id="n1">
              <desc>srv1-CPU</desc>
              <data key="d0">101</data>
              <data key="d1">CPU</data>
              <data key="d2">srv1-CPU</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n2">
              <desc>srv1-Memory</desc>
              <data key="d0">102</data>
              <data key="d1">Memory</data>
              <data key="d2">srv1-Memory</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n3">
              <desc>srv1-Disk</desc>
              <data key="d0">103</data>
              <data key="d1">Disk</data>
              <data key="d2">srv1-Disk</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n4">
              <desc>srv1-Ping</desc>
              <data key="d0">104</data>
              <data key="d1">Ping</data>
              <data key="d2">srv1-Ping</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n5">
              <desc>db1</desc>
              <data key="d0">db1</data>
              <data key="d1">db1</data>
              <data key="d2">db1</data>
              <data key="d3">db</data>
          </node>
          <node id="n6">
              <desc>db1-CPU</desc>
              <data key="d0">201</data>
              <data key="d1">CPU</data>
              <data key="d2">db1-CPU</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n7">
              <desc>db1-Memory</desc>
              <data key="d0">202</data>
              <data key="d1">Memory</data>
              <data key="d2">db1-Memory</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n8">
              <desc>db1-Disk</desc>
              <data key="d0">203</data>
              <data key="d1">Disk</data>
              <data key="d2">db1-Disk</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n9">
              <desc>db1-Ping</desc>
              <data key="d0">204</data>
              <data key="d1">Ping</data>
              <data key="d2">db1-Ping</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n10">
              <desc>db1-DBEngine</desc>
              <data key="d0">205</data>
              <data key="d1">DBEngine</data>
              <data key="d2">db1-DBEngine</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n11">
              <desc>backend1</desc>
              <data key="d0">backend1</data>
              <data key="d1">backend1</data>
              <data key="d2">backend1</data>
              <data key="d3">backend</data>
          </node>
          <node id="n12">
              <desc>backend1-CPU</desc>
              <data key="d0">301</data>
              <data key="d1">CPU</data>
              <data key="d2">backend1-CPU</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n13">
              <desc>backend1-Memory</desc>
              <data key="d0">302</data>
              <data key="d1">Memory</data>
              <data key="d2">backend1-Memory</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n14">
              <desc>backend1-Disk</desc>
              <data key="d0">303</data>
              <data key="d1">Disk</data>
              <data key="d2">backend1-Disk</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n15">
              <desc>backend1-Ping</desc>
              <data key="d0">304</data>
              <data key="d1">Ping</data>
              <data key="d2">backend1-Ping</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n16">
              <desc>backend1-Proc</desc>
              <data key="d0">306</data>
              <data key="d1">Proc</data>
              <data key="d2">backend1-Proc</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n17">
              <desc>backend1-DBConnection</desc>
              <data key="d0">307</data>
              <data key="d1">DBConnection</data>
              <data key="d2">backend1-DBConnection</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n18">
              <desc>frontend1</desc>
              <data key="d0">frontend1</data>
              <data key="d1">frontend1</data>
              <data key="d2">frontend1</data>
              <data key="d3">frontend</data>
          </node>
          <node id="n19">
              <desc>frontend1-CPU</desc>
              <data key="d0">401</data>
              <data key="d1">CPU</data>
              <data key="d2">frontend1-CPU</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n20">
              <desc>frontend1-Memory</desc>
              <data key="d0">402</data>
              <data key="d1">Memory</data>
              <data key="d2">frontend1-Memory</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n21">
              <desc>frontend1-Disk</desc>
              <data key="d0">403</data>
              <data key="d1">Disk</data>
              <data key="d2">frontend1-Disk</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n22">
              <desc>frontend1-Ping</desc>
              <data key="d0">404</data>
              <data key="d1">Ping</data>
              <data key="d2">frontend1-Ping</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n23">
              <desc>frontend1-Proc</desc>
              <data key="d0">406</data>
              <data key="d1">Proc</data>
              <data key="d2">frontend1-Proc</data>
              <data key="d3">alarm</data>
          </node>
          <node id="n24">
              <desc>frontend1-BackendConnection</desc>
              <data key="d0">408</data>
              <data key="d1">BackendConnection</data>
              <data key="d2">frontend1-BackendConnection</data>
              <data key="d3">alarm</data>
          </node>
          <edge id="e0" source="n0" target="n1" directed="false">
              <desc>srv1-CPU</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e1" source="n0" target="n2" directed="false">
              <desc>srv1-Memory</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e2" source="n0" target="n3" directed="false">
              <desc>srv1-Disk</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e3" source="n0" target="n4" directed="false">
              <desc>srv1-Ping</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e4" source="n5" target="n6" directed="false">
              <desc>db1-CPU</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e5" source="n5" target="n7" directed="false">
              <desc>db1-Memory</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e6" source="n5" target="n8" directed="false">
              <desc>db1-Disk</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e7" source="n5" target="n9" directed="false">
              <desc>db1-Ping</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e8" source="n5" target="n10" directed="false">
              <desc>db1-DBEngine</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e9" source="n11" target="n12" directed="false">
              <desc>backend1-CPU</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e10" source="n11" target="n13" directed="false">
              <desc>backend1-Memory</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e11" source="n11" target="n14" directed="false">
              <desc>backend1-Disk</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e12" source="n11" target="n15" directed="false">
              <desc>backend1-Ping</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e13" source="n11" target="n16" directed="false">
              <desc>backend1-Proc</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e14" source="n11" target="n17" directed="false">
              <desc>backend1-DBConnection</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e15" source="n18" target="n19" directed="false">
              <desc>frontend1-CPU</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e16" source="n18" target="n20" directed="false">
              <desc>frontend1-Memory</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e17" source="n18" target="n21" directed="false">
              <desc>frontend1-Disk</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e18" source="n18" target="n22" directed="false">
              <desc>frontend1-Ping</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e19" source="n18" target="n23" directed="false">
              <desc>frontend1-Proc</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e20" source="n18" target="n24" directed="false">
              <desc>frontend1-BackendConnection</desc>
              <data key="d4">trigger</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e21" source="n11" target="n5" directed="false">
              <desc>backend1-db1</desc>
              <data key="d4">connect</data>
              <data key="d5">1</data>
          </edge>
          <edge id="e22" source="n18" target="n11" directed="false">
              <desc>frontend1-backend1</desc>
              <data key="d4">connect</data>
              <data key="d5">1</data>
          </edge>
      </graph>
  </graphml>`

	g := a.GraphML()

	buf := new(bytes.Buffer)
	err := g.Encode(buf, true)
	if err != nil {
		t.Errorf("Error encoding graph: %s", err)
	}

	assert.Equal(t, expectedXML, buf.String())
}

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
			{"data":{"value":"srv1","name":"srv1","type":"server"}},
			{"data":{"value":"101","name":"srv1-CPU","type":"alarm"}},
			{"data":{"value":"102","name":"srv1-Memory","type":"alarm"}},
			{"data":{"value":"103","name":"srv1-Disk","type":"alarm"}},
			{"data":{"value":"104","name":"srv1-Ping","type":"alarm"}},

			{"data":{"value":"db1","name":"db1","type":"db"}},
			{"data":{"value":"201","name":"db1-CPU","type":"alarm"}},
			{"data":{"value":"202","name":"db1-Memory","type":"alarm"}},
			{"data":{"value":"203","name":"db1-Disk","type":"alarm"}},
			{"data":{"value":"204","name":"db1-Ping","type":"alarm"}},
			{"data":{"value":"205","name":"db1-DBEngine","type":"alarm"}},

			{"data":{"value":"backend1","name":"backend1","type":"backend"}},
			{"data":{"value":"301","name":"backend1-CPU","type":"alarm"}},
			{"data":{"value":"302","name":"backend1-Memory","type":"alarm"}},
			{"data":{"value":"303","name":"backend1-Disk","type":"alarm"}},
			{"data":{"value":"304","name":"backend1-Ping","type":"alarm"}},
			{"data":{"value":"306","name":"backend1-Proc","type":"alarm"}},
			{"data":{"value":"307","name":"backend1-DBConnection","type":"alarm"}},

			{"data":{"value":"frontend1","name":"frontend1","type":"frontend"}},
			{"data":{"value":"401","name":"frontend1-CPU","type":"alarm"}},
			{"data":{"value":"402","name":"frontend1-Memory","type":"alarm"}},
			{"data":{"value":"403","name":"frontend1-Disk","type":"alarm"}},
			{"data":{"value":"404","name":"frontend1-Ping","type":"alarm"}},
			{"data":{"value":"406","name":"frontend1-Proc","type":"alarm"}},
			{"data":{"value":"408","name":"frontend1-BackendConnection","type":"alarm"}}
		],
		"edges":[
			{"data":{"source":"srv1","target":"101","type":"trigger","weight":0}},
			{"data":{"source":"srv1","target":"102","type":"trigger","weight":0}},
			{"data":{"source":"srv1","target":"103","type":"trigger","weight":0}},
			{"data":{"source":"srv1","target":"104","type":"trigger","weight":0}},

			{"data":{"source":"db1","target":"201","type":"trigger","weight":0}},
			{"data":{"source":"db1","target":"202","type":"trigger","weight":0}},
			{"data":{"source":"db1","target":"203","type":"trigger","weight":0}},
			{"data":{"source":"db1","target":"204","type":"trigger","weight":0}},
			{"data":{"source":"db1","target":"205","type":"trigger","weight":0}},

			{"data":{"source":"backend1","target":"301","type":"trigger","weight":0}},
			{"data":{"source":"backend1","target":"302","type":"trigger","weight":0}},
			{"data":{"source":"backend1","target":"303","type":"trigger","weight":0}},
			{"data":{"source":"backend1","target":"304","type":"trigger","weight":0}},
			{"data":{"source":"backend1","target":"306","type":"trigger","weight":0}},
			{"data":{"source":"backend1","target":"307","type":"trigger","weight":0}},

			{"data":{"source":"frontend1","target":"401","type":"trigger","weight":0}},
			{"data":{"source":"frontend1","target":"402","type":"trigger","weight":0}},
			{"data":{"source":"frontend1","target":"403","type":"trigger","weight":0}},
			{"data":{"source":"frontend1","target":"404","type":"trigger","weight":0}},
			{"data":{"source":"frontend1","target":"406","type":"trigger","weight":0}},
			{"data":{"source":"frontend1","target":"408","type":"trigger","weight":0}},

			{"data":{"source":"backend1","target":"db1","type":"connect","weight":0}},
			{"data":{"source":"frontend1","target":"backend1","type":"connect","weight":0}}
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
