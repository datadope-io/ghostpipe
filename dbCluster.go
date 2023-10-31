package main

import (
	"fmt"
	"math/rand"

	"github.com/fschuetz04/simgo"
)

// DBCluster simula varios clusters de servidores donde cuando se
// cae la alarma de disco se cae para todas. Igual para la CPU y memoria
func DBCluster(a *Architecture) {
	// for _, tech := range []string{"mariadb", "mysql", "postgresql", "mongodb", "redis", "memcached", "elasticsearch", "cassandra"} {
	clusterSize := map[string]int{
		"mariadb":       3,
		"mysql":         5,
		"postgresql":    3,
		"mongodb":       5,
		"elasticsearch": 6,
	}
	clusters := make(map[string][]*Database, len(clusterSize))

	for tech, size := range clusterSize {
		clusters[tech] = make([]*Database, size)
		for i := 0; i < size; i++ {
			clusters[tech][i] = a.NewDatabase(fmt.Sprintf("%s-%d", tech, i))
		}

		// Esta llamada busca crear los edges entre los servidores
		a.NewClusterDB(clusters[tech])

		t := tech // Copia de la variable para que no se vea modificada por el loop al crear el func literal

		// Simulamos una subida de carga, que dispara la alama de CPU en todos los servidores
		a.AddMonkey(func(proc simgo.Process) {
			for {
				proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
				for _, db := range clusters[t] {
					db.CPUAlarm = AlarmTriggered
				}
			}
		})

		// Simulamos una subida de carga, que dispara la alama de memoria en todos los servidores
		a.AddMonkey(func(proc simgo.Process) {
			for {
				proc.Wait(proc.Timeout(float64(100 + rand.Intn(700))))
				for _, db := range clusters[t] {
					db.MemoryAlarm = AlarmTriggered
				}
			}
		})

		// Simulamos un llenado de disco
		a.AddMonkey(func(proc simgo.Process) {
			for {
				proc.Wait(proc.Timeout(float64(100 + rand.Intn(1200))))
				for _, db := range clusters[t] {
					db.DiskAlarm = AlarmTriggered
				}
			}
		})
	}

	// Several servers as noise
	noiseServers := []*Server{}
	for i := 0; i < 50; i++ {
		noiseServers = append(noiseServers, a.NewServer("noise"+fmt.Sprintf("%d", i)))
	}

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
