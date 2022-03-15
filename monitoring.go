package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type MonitorSystem interface {
	handleAlarm(string, string, float64)
	generateEventID(string, string) int
}

// MoMonitorSystem receive the alarms of the servers and generate messages
type PrinterMonitorSystem struct {
	sync.Mutex
	eventid map[string]int
	usedIDs map[int]bool
}

func (m *PrinterMonitorSystem) generateEventID(server string, alarm string) int {
	// Initialize eventid map if it is nil
	if m.eventid == nil {
		m.eventid = make(map[string]int)
	}

	// Initialize usedIDs map if it is nil
	if m.usedIDs == nil {
		m.usedIDs = make(map[int]bool)
	}

	// Check if server+alarm is already in the map
	if _, ok := m.eventid[server+alarm]; ok {
		return m.eventid[server+alarm]
	}

	// Generate a random integer
	for {
		id := rand.Intn(1000000)
		if !m.usedIDs[id] {
			m.usedIDs[id] = true
			m.eventid[server+alarm] = id
			return id
		}
	}
}

func (m *PrinterMonitorSystem) handleAlarm(server string, alarm string, time float64) {
	m.Lock()
	defer m.Unlock()

	// Get time in unix epoch format
	e := fmt.Sprintf("%.0f,%s,%s,%v\n", time, server, alarm, m.generateEventID(server, alarm))

	// TODO guardar en un fichero
	fmt.Printf(e)
}
