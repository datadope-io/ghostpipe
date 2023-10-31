package main

// Frontend represents a frontend server with its possible alarms and connected to a backend.
type DNS struct {
	Server
	// ProcAlarm is True if the DNS process is not running
	ProcAlarm AlarmStatus
	Clients   []MonitoredServer
}

// NewFrontend create a new Frontend server, start it and return the pointer to it
func NewDNS(name string, mon MonitorSystem) *DNS {
	return &DNS{
		Server: Server{
			Name: name,
			mon:  mon,
		},
	}
}

func (b *DNS) AddClient(client MonitoredServer) {
	b.Clients = append(b.Clients, client)
}

func (b *DNS) GetName() string {
	return b.Name
}

func (s *DNS) GetAlarms() []string {
	serverAlarms := s.Server.GetAlarms()
	return append(serverAlarms, "Proc")
}

func (s *DNS) GetType() string {
	return string(DNSNode)
}

// CheckAlarms print a message if the server has alarms.
// Trigger the DNS alarm in all the connectected clients.
func (b *DNS) CheckAlarms(t float64) {
	if b.ProcAlarm == AlarmTriggered {
		b.ProcAlarm = AlarmACK
		b.mon.handleAlarm(b.Name, "Proc", t)

		for _, client := range b.Clients {
			client.SetAlarm("DNS", AlarmTriggered)
		}
	}

	b.Server.CheckAlarms(t)
}

// Available returns true if the Frontend server is considered available, that is,
// if the Frontend process is running and the database is available.
func (b *DNS) Available() bool {
	return b.Server.Available() && b.ProcAlarm == AlarmEnabled
}

func (b *DNS) SetAlarm(alarm string, status AlarmStatus) {
	if alarm == "Proc" {
		b.ProcAlarm = status
	} else {
		b.Server.SetAlarm(alarm, status)
	}
}
