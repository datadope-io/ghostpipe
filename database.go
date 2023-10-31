package main

// Database represents a database server with a database engine running (like postgres, mysql, etc)
type Database struct {
	DBEngineAlarm AlarmStatus
	Server
}

// NewDatabase create a new database server, start it and return the pointer to it
func NewDatabase(name string, mon MonitorSystem) *Database {
	return &Database{
		Server: Server{
			Name: name,
			mon:  mon,
		},
	}
}

func (d *Database) GetName() string {
	return d.Name
}

func (s *Database) GetAlarms() []string {
	serverAlarms := s.Server.GetAlarms()
	return append(serverAlarms, []string{"DBEngine"}...)
}

func (s *Database) GetType() string {
	return string(DBNode)
}

// CheckAlarms print a message if the database engines is not working
// or the base server has alarms.
func (d *Database) CheckAlarms(t float64) {
	if d.DBEngineAlarm == AlarmTriggered {
		d.DBEngineAlarm = AlarmACK
		d.mon.handleAlarm(d.Name, "DBEngine", t)
	}

	d.Server.CheckAlarms(t)
}

// Available return true if the db server is considered available, that is,
// if the db engine is available and the server is available.
func (d *Database) Available() bool {
	return d.Server.Available() && d.DBEngineAlarm == AlarmEnabled
}

func (d *Database) SetAlarm(alarm string, status AlarmStatus) {
	switch alarm {
	case "DBEngine":
		d.DBEngineAlarm = AlarmTriggered
	default:
		d.Server.SetAlarm(alarm, status)
	}
}
