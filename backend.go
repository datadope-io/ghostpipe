package main

// Backend represents a backend server with its possible alarms and connected to a database.
type Backend struct {
	Server
	// ProcAlarm is True if the backend process is not running
	ProcAlarm AlarmStatus
	// DBConnectionAlarm is True if the database is not working
	DBConnectionAlarm AlarmStatus
	DBEngine          *Database

	// dnsAvailable used to store the state of the DNS server the last time CheckAlarms was called.
	dnsAvailable bool
}

// NewBackend create a new backend server, start it and return the pointer to it
func NewBackend(name string, connectedDB *Database, mon MonitorSystem) *Backend {
	return &Backend{
		Server: Server{
			Name: name,
			mon:  mon,
		},
		DBEngine: connectedDB,
	}
}

func (b *Backend) GetName() string {
	return b.Name
}

func (s *Backend) GetAlarms() []string {
	serverAlarms := s.Server.GetAlarms()
	return append(serverAlarms, []string{"Proc", "DBConnection"}...)
}

func (s *Backend) GetType() string {
	return string(BackendNode)
}

// CheckAlarms print a message if the server has alarms.
// It check alarms specific to the backend, plus generic alarms for the server
// and also generate an alarm if the database is not available.
func (b *Backend) CheckAlarms(t float64) {
	if b.ProcAlarm == AlarmTriggered {
		b.ProcAlarm = AlarmACK
		b.mon.handleAlarm(b.Name, "Proc", t)
	}

	// Set the local db connection alarm based on the state of the database.
	// Generate a new alarm if we are moving from enabled to triggered.
	if b.DBEngine.Available() {
		b.DBConnectionAlarm = AlarmEnabled
	} else {
		if b.DBConnectionAlarm == AlarmEnabled {
			b.DBConnectionAlarm = AlarmACK
			b.mon.handleAlarm(b.Name, "DBConnection", t)
		}
	}

	// If DNS server is not available, backend could not communicate with the database,
	// so we also trigger the DBConnection alarm.
	if b.Server.DNSAlarm == AlarmTriggered && b.DBConnectionAlarm == AlarmEnabled {
		b.DBConnectionAlarm = AlarmACK
		b.dnsAvailable = false
		b.mon.handleAlarm(b.Name, "DBConnection", t)
	}
	// If DNS was unavailable and is now available, we can clear the DBConnection alarm.
	if !b.dnsAvailable && b.Server.DNSAlarm == AlarmEnabled && b.DBConnectionAlarm == AlarmACK {
		b.dnsAvailable = true
		b.DBConnectionAlarm = AlarmEnabled
	}

	b.Server.CheckAlarms(t)
}

// Available returns true if the backend server is considered available, that is,
// if the backend process is running and the database is available.
func (b *Backend) Available() bool {
	return b.Server.Available() && b.ProcAlarm == AlarmEnabled
}

func (b *Backend) SetAlarm(alarm string, status AlarmStatus) {
	switch alarm {
	case "Proc":
		b.ProcAlarm = status
	case "DBConnection":
		b.DBConnectionAlarm = status
	default:
		b.Server.SetAlarm(alarm, status)
	}
}
