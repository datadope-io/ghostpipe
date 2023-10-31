package main

import (
	"fmt"
	"math/rand"

	"github.com/fschuetz04/simgo"
)

// BackendFrontendNoise es una topología que simula dos bases de datos donde hay
// conectados backends y frontends a esos backends.
// Luego tenemos varios servidores aislados con sus alarmas.
func RelacionesInesperadas(a *Architecture) {
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
	for i := 0; i < 50; i++ {
		noiseServers = append(noiseServers, a.NewServer("noise"+fmt.Sprintf("%d", i)))
	}

	// DNS server connected to all servers
	dns := a.NewDNS("dnsA")
	for _, s := range a.GetAllServers() {
		// Do not add the DNS aserver as a client to itself
		if s.GetName() == dns.GetName() {
			continue
		}
		dns.AddClient(s)
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
}
