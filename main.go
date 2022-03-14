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
	AlarmCheckInterval = 60

	// IntervalJitter is the max possible jitter for the interval expressed
	// as a percentage of AlarmCheckInterval
	IntervalJitter = 0.1
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

	fmt.Printf("Writing graph to %s\n", *graphFile)
	g := a.CytoscapeGraph()
	f, err := os.Create(*graphFile)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(g)
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

	// Disconnect db1 each 200' and reconnect it after 60'
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(200))
		for {
			fmt.Println("\nmonkey: disconnect db1")
			db1.PingAlarm = AlarmTriggered

			proc.Wait(proc.Timeout(60))
			//fmt.Println("\nmonkey: reconnect db1")
			db1.PingAlarm = AlarmEnabled

			proc.Wait(proc.Timeout(140))
		}
	})

	// Disconnect backendD each 380' and reconnect it after 60'
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(380))
		for {
			fmt.Println("\nmonkey: disconnect backendD")
			backendD.PingAlarm = AlarmTriggered

			proc.Wait(proc.Timeout(60))
			//fmt.Println("monkey: reconnect backendD")
			backendD.PingAlarm = AlarmEnabled

			proc.Wait(proc.Timeout(320))
		}
	})

	// Start the simulation.
	fmt.Println("Starting simulator...")
	// The parameter is the function able to trigger alarms
	simulation_duration := 60.0 * 10
	a.Start(simulation_duration)
	fmt.Println("Stopping simulator...")

	/*
		time.Sleep(2 * time.Second)
		fmt.Println("db1 down")
		db1.PingAlarm = AlarmTriggered

		time.Sleep(8 * time.Second)

		fmt.Println("Stopping simulator...")
		a.Stop()
	*/
	/*
		time.Sleep(2 * time.Second)

		fmt.Println("db1 down")

		db1.PingAlarm = AlarmTriggered

		// Tras otro periodo, simulamos que se recupera la db
		time.Sleep(4 * time.Second)

		fmt.Println("db1 up")

		db1.PingAlarm = AlarmEnabled

		// Y un rato después se vuelve a caer
		time.Sleep(4 * time.Second)

		fmt.Println("db1 down")

		db1.PingAlarm = AlarmEnabled
		db1.PingAlarm = AlarmTriggered

		_ = backend1 // TODO borrar

		time.Sleep(20 * time.Second)
	*/
}
