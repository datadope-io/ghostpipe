package main

import (
	"fmt"
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
	DNSAlarm    AlarmStatus

	// mon connection to the monitoring system
	mon MonitorSystem
}

type MonitoredServer interface {
	GetName() string
	CheckAlarms(float64)
	// SetAlarm using the string to identify the alarm, set the alarm to the given status
	SetAlarm(string, AlarmStatus)
}

type ArchitectureServer interface {
	GetName() string
	GetType() string
	GetAlarms() []string
}

func NewServer(name string, mon MonitorSystem) *Server {
	return &Server{
		Name: name,
		mon:  mon,
	}
}

func (d *Server) GetName() string {
	return d.Name
}

func (s *Server) GetAlarms() []string {
	return []string{
		"CPU",
		"Memory",
		"Disk",
		"Ping",
		"DNS",
	}
}

func (s *Server) GetType() string {
	return string(ServerNode)
}

// Run check the alarms of each server each interval
func Run(proc simgo.Process, m MonitoredServer) {
	// Desalign the time of checking for each server
	proc.Wait(proc.Timeout(float64(rand.Intn(AlarmCheckInterval))))

	for {
		m.CheckAlarms(proc.Now())
		proc.Wait(proc.Timeout(AlarmCheckInterval))
		// Execution jitter
		proc.Wait(proc.Timeout(AlarmCheckInterval * rand.Float64() * IntervalJitter))
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

	if s.DNSAlarm == AlarmTriggered {
		s.DNSAlarm = AlarmACK
		s.mon.handleAlarm(s.Name, "DNS", t)
	}
}

// Available returns true if the server is considered available
func (s *Server) Available() bool {
	return s.PingAlarm == AlarmEnabled
}

func (s *Server) SetAlarm(alarm string, status AlarmStatus) {
	switch alarm {
	case "CPU":
		s.CPUAlarm = status
	case "Memory":
		s.MemoryAlarm = status
	case "Disk":
		s.DiskAlarm = status
	case "Ping":
		s.PingAlarm = status
	case "DNS":
		s.DNSAlarm = status
	default:
		panic(fmt.Sprintf("Unknown alarm: %s", alarm))
	}
}
