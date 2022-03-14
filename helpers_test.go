package main

import (
	"fmt"
	"sync"
)

type fakeMonSys struct {
	sync.RWMutex
	Alarms []string
}

func (m *fakeMonSys) handleAlarm(server string, alarm string, time float64) {
	m.Lock()
	defer m.Unlock()
	// Get time in unix epoch format
	m.Alarms = append(m.Alarms, fmt.Sprintf("%.0f,%s,%s", time, server, alarm))
}
