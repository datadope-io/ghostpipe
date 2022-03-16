package main

import (
	"fmt"
	"math/rand"
	"os"
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
	events  []string
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
	e := fmt.Sprintf("%.0f,%s,%s,%v\n", time*60, server, alarm, m.generateEventID(server, alarm))

	m.events = append(m.events, e)
}

func (m *PrinterMonitorSystem) WriteEvents(fileName string) {
	// Delete the events file if it exists
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}

	// File where events will be saved
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the CSV header for the events file
	_, err = f.WriteString("time,server,alarm,eventid\n")
	if err != nil {
		panic(err)
	}
	for _, e := range m.events {
		_, err := f.WriteString(e)
		if err != nil {
			panic(err)
		}
	}
}
