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

		/* TODO esta cogiendo únicamente el último valor de la variable tech
		// Simulamos una subida de carga, que dispara la alama de CPU en todos los servidores
		a.AddMonkey(func(proc simgo.Process) {
			proc.Wait(proc.Timeout(60))
			for {
				for _, db := range clusters[tech] {
					db.CPUAlarm = AlarmTriggered
				}

				// Wait for a random time between 100 and 500 seconds
				proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
			}
		})
		*/
	}

	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			for _, db := range clusters["mariadb"] {
				db.CPUAlarm = AlarmTriggered
			}
			proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
		}
	})

	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			for _, db := range clusters["mysql"] {
				db.CPUAlarm = AlarmTriggered
			}

			// Wait for a random time between 100 and 500 seconds
			proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
		}
	})
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			for _, db := range clusters["postgresql"] {
				db.CPUAlarm = AlarmTriggered
			}

			// Wait for a random time between 100 and 500 seconds
			proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
		}
	})
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			for _, db := range clusters["mongodb"] {
				db.CPUAlarm = AlarmTriggered
			}

			// Wait for a random time between 100 and 500 seconds
			proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
		}
	})
	a.AddMonkey(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(60))
		for {
			for _, db := range clusters["elasticsearch"] {
				db.CPUAlarm = AlarmTriggered
			}

			// Wait for a random time between 100 and 500 seconds
			proc.Wait(proc.Timeout(float64(100 + rand.Intn(400))))
		}
	})
}
