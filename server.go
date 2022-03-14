package main

import (
	"math/rand"

	"github.com/fschuetz04/simgo"
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
	CheckAlarms(float64)
}

func Run(proc simgo.Process, m MonitoredServer) {
	// Add a random delay before starting the loop, between 0 and AlarmCheckInterval
	proc.Wait(proc.Timeout(float64(rand.Intn(AlarmCheckInterval))))

	for {
		m.CheckAlarms(proc.Now())
		proc.Wait(proc.Timeout(AlarmCheckInterval))
	}
}

// CheckAlarms if the server has alarms and print a message for each triggered alarm
func (s *Server) CheckAlarms(t float64) {
	if s.CPUAlarm == AlarmTriggered {
		s.CPUAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "CPU", t)
	}

	if s.MemoryAlarm == AlarmTriggered {
		s.MemoryAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Memory", t)
	}

	if s.DiskAlarm == AlarmTriggered {
		s.DiskAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Disk", t)
	}

	if s.PingAlarm == AlarmTriggered {
		s.PingAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "Ping", t)
	}
}

// Available returns true if the server is considered available
func (s *Server) Available() bool {
	return s.PingAlarm == AlarmEnabled
}
