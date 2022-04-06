package main

import (
	"fmt"
	"math/rand"

	"github.com/fschuetz04/simgo"
)

// BackendFrontendNoise es una topología que simula dos bases de datos donde hay
// conectados backends y frontends a esos backends.
// Luego tenemos varios servidores aislados con sus alarmas.
func MiniBackendFrontendNoise(a *Architecture) {
	// Add servers to the architecture

	// One database serving three different backends.
	// Each backend with one or more frontends.
	db1 := a.NewDatabase("db1")

	backendA := a.NewBackend("backendA", db1)
	frontendA1 := a.NewFrontend("frontendA1", backendA)

	backendB := a.NewBackend("backendB", db1)
	frontendB1 := a.NewFrontend("frontendB1", backendB)

	// Several servers as noise
	noiseServers := []*Server{}
	for i := 0; i < 5; i++ {
		noiseServers = append(noiseServers, a.NewServer("noise"+fmt.Sprintf("%d", i)))
	}

	_ = frontendA1
	_ = frontendB1

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
