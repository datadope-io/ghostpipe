package main

// Frontend represents a frontend server with its possible alarms and connected to a backend.
type Frontend struct {
	Server
	// ProcAlarm is True if the frontend process is not running
	ProcAlarm AlarmStatus
	// BackendConnectionAlarm is True if the backend is not working
	BackendConnectionAlarm AlarmStatus
	Backend                *Backend
}

// NewFrontend create a new Frontend server, start it and return the pointer to it
func NewFrontend(name string, connectedBackend *Backend, mon MonitorSystem) *Frontend {
	return &Frontend{
		Server: Server{
			Name: name,
			mon:  mon,
		},
		Backend: connectedBackend,
	}
}

func (b *Frontend) GetName() string {
	return b.Name
}

func (s *Frontend) GetAlarms() []string {
	serverAlarms := s.Server.GetAlarms()
	return append(serverAlarms, []string{"Proc", "BackendConnection"}...)
}

func (s *Frontend) GetType() string {
	return string(FrontendNode)
}

// CheckAlarms print a message if the server has alarms.
// It check alarms specific to the Frontend, plus generic alarms for the server
// and also generate an alarm if the backend is not available.
func (b *Frontend) CheckAlarms(t float64) {
	if b.ProcAlarm == AlarmTriggered {
		b.ProcAlarm = AlarmACK
		b.mon.handleAlarm(b.Name, "Proc", t)
	}

	// Set the local backend connection alarm based on the state of the database.
	// Generate a new alarm if we are moving from enabled to triggered.
	if b.Backend.Available() {
		b.BackendConnectionAlarm = AlarmEnabled
	} else {
		if b.BackendConnectionAlarm == AlarmEnabled {
			b.BackendConnectionAlarm = AlarmACK
			b.mon.handleAlarm(b.Name, "BackendConnection", t)
		}
	}

	b.Server.CheckAlarms(t)
}

// Available returns true if the Frontend server is considered available, that is,
// if the Frontend process is running and the database is available.
func (b *Frontend) Available() bool {
	return b.Server.Available() && b.ProcAlarm == AlarmEnabled
}

func (b *Frontend) SetAlarm(alarm string, status AlarmStatus) {
	switch alarm {
	case "Proc":
		b.ProcAlarm = status
	case "BackendConnection":
		b.BackendConnectionAlarm = status
	default:
		b.Server.SetAlarm(alarm, status)
	}
}
