package main

import (
	"fmt"
	"sync"
)

type MonitorSystem interface {
	handleAlarm(string, string, float64)
}

// MoMonitorSystem receive the alarms of the servers and generate messages
type PrinterMonitorSystem struct {
	sync.Mutex
}

func (m *PrinterMonitorSystem) handleAlarm(server string, alarm string, time float64) {
	m.Lock()
	defer m.Unlock()
	// Get time in unix epoch format
	fmt.Printf("%.0f,%s,%s\n", time, server, alarm)
}
