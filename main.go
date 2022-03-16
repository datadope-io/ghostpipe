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
	"math/rand"
	"os"

	"github.com/fschuetz04/simgo"
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
	graphFile  = flag.String("graph", "graph.cyjs", "File to save the graph in Cytoscape JSON format")
	eventsFile = flag.String("events", "events.csv", "File to save the events in CSV format")
)

func main() {
	flag.Parse()

	// Create the monitoring system
	mon := &PrinterMonitorSystem{}

	// Create the architecture
	a := Architecture{mon: mon}

	// Add servers to the architecture

	// One database serving three different backends.
	// Each backend with one or more frontends.
	db1 := a.NewDatabase("db1")

	backendA := a.NewBackend("backendA", db1)
	frontendA1 := a.NewFrontend("frontendA1", backendA)

	backendB := a.NewBackend("backendB", db1)
	frontendB1 := a.NewFrontend("frontendB1", backendB)
	frontendB2 := a.NewFrontend("frontendB2", backendB)

	backendC := a.NewBackend("backendC", db1)
	frontendC1 := a.NewFrontend("frontendC1", backendC)
	frontendC2 := a.NewFrontend("frontendC2", backendC)
	frontendC3 := a.NewFrontend("frontendC3", backendC)

	// One app with frontend, backend and database
	db2 := a.NewDatabase("db2")
	backendD := a.NewBackend("backendD", db2)
	frontendD1 := a.NewFrontend("frontendD1", backendD)

	// Several servers as noise
	noiseServers := []*Server{}
	for i := 0; i < 500; i++ {
		noiseServers = append(noiseServers, a.NewServer("noise"+fmt.Sprintf("%d", i)))
	}

	g := a.CytoscapeGraph()
	gFile, err := os.Create(*graphFile)
	if err != nil {
		panic(err)
	}
	defer gFile.Close()

	fmt.Printf("Writing graph to %s\n", *graphFile)
	_, err = gFile.WriteString(g)
	if err != nil {
		panic(err)
	}

	_ = frontendA1
	_ = frontendB1
	_ = frontendB2
	_ = frontendC1
	_ = frontendC2
	_ = frontendC3
	_ = frontendD1

	// Disconnect db1 each 60' and reconnect it after 5'
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			// fmt.Println("\nmonkey: disconnect db1")
			db1.PingAlarm = AlarmTriggered

			proc.Wait(proc.Timeout(5))
			// fmt.Println("\nmonkey: reconnect db1")
			db1.PingAlarm = AlarmEnabled

			proc.Wait(proc.Timeout(55))
		}
	})

	// Disconnect backendD each 120' and reconnect it after 60'
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(120))
		for {
			// fmt.Println("\nmonkey: disconnect backendD")
			backendD.PingAlarm = AlarmTriggered

			proc.Wait(proc.Timeout(60))
			// fmt.Println("monkey: reconnect backendD")
			backendD.PingAlarm = AlarmEnabled

			proc.Wait(proc.Timeout(60))
		}
	})

	// Generate alarm noise
	a.AddMonkey(func(proc simgo.Process) {
		for {
			// Get one of the noise servers
			noiseServer := noiseServers[rand.Intn(len(noiseServers))]

			// Trigger one the alarms of the server
			switch rand.Intn(4) {
			case 0:
				noiseServer.CPUAlarm = AlarmTriggered
			case 1:
				noiseServer.MemoryAlarm = AlarmTriggered
			case 2:
				noiseServer.DiskAlarm = AlarmTriggered
			case 3:
				noiseServer.PingAlarm = AlarmTriggered
			}

			proc.Wait(proc.Timeout(1))
		}
	})

	// Start the simulation.
	fmt.Println("Starting simulator...")

	// Run the simulation for this long
	a.Start(60.0 * 24 * 2)
	fmt.Println("Simulator finished")

	// Write generated events to file
	fmt.Printf("Writing events to file %s\n", *eventsFile)
	mon.WriteEvents(*eventsFile)
}
