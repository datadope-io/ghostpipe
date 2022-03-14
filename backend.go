package main

// Backend represents a backend server with its possible alarms and connected to a database.
type Backend struct {
	Server
	// ProcAlarm is True if the backend process is not running
	ProcAlarm AlarmStatus
	// DBConnectionAlarm is True if the database is not working
	DBConnectionAlarm AlarmStatus
	DBEngine          *Database
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

	b.Server.CheckAlarms(t)
}

// Available returns true if the backend server is considered available, that is,
// if the backend process is running and the database is available.
func (b *Backend) Available() bool {
	return b.Server.Available() && b.ProcAlarm == AlarmEnabled
}
