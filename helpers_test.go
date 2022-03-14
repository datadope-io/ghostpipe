package main

import (
	"fmt"
	"sync"
)

type fakeMonSys struct {
	sync.RWMutex
	Alarms []string
}

func (m *fakeMonSys) handleAlarm(server string, alarm string) {
	m.Lock()
	defer m.Unlock()
	// Get time in unix epoch format
	m.Alarms = append(m.Alarms, fmt.Sprintf("%s,%s", server, alarm))
}
