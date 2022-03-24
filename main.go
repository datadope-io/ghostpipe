//
// Simulate a monitoring system where different kind of servers are interconnected
// and have alarms that can be triggered.
//
// Each server checks each interval if any of its alarms has been triggered and
// if so, sends a message to the monitoring system.
// The alarms state could be modified externally or could be set by the server
// based of the status of other connected servers.
// For example, a backend server has an alarm that triggers when the connection
// to the database is lost. This alarm is based on the status on the db.
package main

import (
	"flag"
	"fmt"
	"os"
)

// AlarmStatus is an enum for the status of an alarm
type AlarmStatus int

const (
	// AlarmEnabled is the alarm ready to be triggered
	AlarmEnabled AlarmStatus = iota
	// AlarmTriggered is the alarm fired
	AlarmTriggered
	// AlarmACK set when the alarm message has been generated and the alarm is still in trigger state
	AlarmACK

	// AlarmCheckInterval is the time interval when alarms are checked
	AlarmCheckInterval = 1

	// IntervalJitter is the max possible jitter for the interval expressed
	// as a percentage of AlarmCheckInterval
	IntervalJitter = 0.2
)

// Create flags to define graph and events output files
var (
	cytoscapeGraphFile = flag.String("cytoscape", "", "File to save the graph in Cytoscape JSON format")
	graphMLFile        = flag.String("graphml", "graph.graphml", "File to save the graph in GraphML format")
	eventsFile         = flag.String("events", "events.csv", "File to save the events in CSV format")
)

func main() {
	flag.Parse()

	// Create the monitoring system
	mon := &PrinterMonitorSystem{}

	// Create the architecture
	a := Architecture{mon: mon}

	// Select which topology to use (uncomment the one you want to use)
	BackendFrontendNoise(&a)

	// Output the graph in different formats
	if *graphMLFile != "" {
		gFile, err := os.Create(*graphMLFile)
		if err != nil {
			panic(err)
		}
		defer gFile.Close()

		gml := a.GraphML()
		fmt.Printf("Writing GraphML graph to %s\n", *graphMLFile)
		err = gml.Encode(gFile, true)
		if err != nil {
			panic(err)
		}
	}

	if *cytoscapeGraphFile != "" {
		gFile, err := os.Create(*cytoscapeGraphFile)
		if err != nil {
			panic(err)
		}
		defer gFile.Close()

		g := a.CytoscapeGraph()
		fmt.Printf("Writing Cytoscape graph to %s\n", *cytoscapeGraphFile)
		_, err = gFile.WriteString(g)
		if err != nil {
			panic(err)
		}
	}

	// Start the simulation.
	fmt.Println("Starting simulator...")

	// Run the simulation for this long
	a.Start(60.0 * 24 * 2)
	fmt.Println("Simulator finished")

	// Write generated events to file
	fmt.Printf("Writing events to file %s\n", *eventsFile)
	mon.WriteEvents(*eventsFile)
}
