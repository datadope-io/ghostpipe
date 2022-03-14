package main

import (
	"context"
	"math/rand"
	"time"
)

// Server represents a server with its possible alarms.
// Also serve as base for other especialized servers
type Server struct {
	Name        string
	CPUAlarm    AlarmStatus
	MemoryAlarm AlarmStatus
	DiskAlarm   AlarmStatus
	PingAlarm   AlarmStatus

	// mon connection to the monitoring system
	mon MonitorSystem
}

type MonitoredServer interface {
	GetName() string
	CheckAlarms()
}

func Run(ctx context.Context, m MonitoredServer) {
	//fmt.Printf("Starting %s\n", m.GetName())

	// Timer with AlarmCheckInterval
	t := time.NewTicker(AlarmCheckInterval)

	// Add a random delay before starting the loop, between 0 and AlarmCheckInterval
	time.Sleep(time.Duration(rand.Intn(int(AlarmCheckInterval.Seconds()))) * time.Second)

	for {
		// Each tick run the CheckAlarms function.
		// Finish the loop if the context is done.
		select {
		case <-t.C:
			m.CheckAlarms()
		case <-ctx.Done():
			//fmt.Printf("Stopping %s\n", m.GetName())
			t.Stop()
			return
		}
	}
}

// CheckAlarms if the server has alarms and print a message for each triggered alarm
func (s *Server) CheckAlarms() {
	if s.CPUAlarm == AlarmTriggered {
		s.CPUAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "CPU")
	}

	if s.MemoryAlarm == AlarmTriggered {
		s.MemoryAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Memory")
	}

	if s.DiskAlarm == AlarmTriggered {
		s.DiskAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Disk")
	}

	if s.PingAlarm == AlarmTriggered {
		s.PingAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Ping")
	}
}

// Available returns true if the server is considered available
func (s *Server) Available() bool {
	return s.PingAlarm == AlarmEnabled
}
