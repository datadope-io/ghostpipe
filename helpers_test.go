package main

import (
	"fmt"
	"sync"
)

type fakeMonSys struct {
	sync.RWMutex
	Alarms []string
}

func (m *fakeMonSys) generateEventID(server string, alarm string) int {
	eventid := 0

	switch server {
	case "srv1":
		eventid = 100
	case "db1":
		eventid = 200
	case "backend1":
		eventid = 300
	case "frontend1":
		eventid = 400
	default:
		panic("Unknown server, must be initiliazed")
	}

	switch alarm {
	case "CPU":
		eventid += 1
	case "Memory":
		eventid += 2
	case "Disk":
		eventid += 3
	case "Ping":
		eventid += 4
	case "DBEngine":
		eventid += 5
	case "Proc":
		eventid += 6
	case "DBConnection":
		eventid += 7
	case "BackendConnection":
		eventid += 8
	default:
		panic("Unknown alarm, must be initiliazed")
	}

	return eventid
}

func (m *fakeMonSys) handleAlarm(server string, alarm string, time float64) {
	m.Lock()
	defer m.Unlock()
	// Get time in unix epoch format
	m.Alarms = append(m.Alarms, fmt.Sprintf("%.0f,%s,%s", time, server, alarm))
}
