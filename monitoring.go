package main

import (
	"fmt"
	"sync"
	"time"
)

type MonitorSystem interface {
	handleAlarm(string, string)
}

// MoMonitorSystem receive the alarms of the servers and generate messages
type PrinterMonitorSystem struct {
	sync.Mutex
}

func (m *PrinterMonitorSystem) handleAlarm(server string, alarm string) {
	m.Lock()
	defer m.Unlock()
	// Get time in unix epoch format
	now := time.Now().Unix()
	fmt.Printf("%d,%s,%s\n", now, server, alarm)
}
