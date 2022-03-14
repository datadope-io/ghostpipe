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

// CheckAlarms print a message if the database engines is not working
// or the base server has alarms.
func (d *Database) CheckAlarms() {
	if d.DBEngineAlarm == AlarmTriggered {
		d.DBEngineAlarm = AlarmACK
		d.mon.handleAlarm(d.Name, "DBEngine")
	}

	d.Server.CheckAlarms()
}

// Available return true if the db server is considered available, that is,
// if the db engine is available and the server is available.
func (d *Database) Available() bool {
	return d.Server.Available() && d.DBEngineAlarm == AlarmEnabled
}
