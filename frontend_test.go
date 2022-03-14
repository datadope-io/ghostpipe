package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBDownNotAffectFrontendAlarm(t *testing.T) {
	mon := &fakeMonSys{}

	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}
	frontend1 := &Frontend{Server: Server{Name: "frontend1", mon: mon}, Backend: backend1}

	// Set the backend down
	db1.PingAlarm = AlarmTriggered

	// Check alarms
	db1.CheckAlarms()
	backend1.CheckAlarms()
	frontend1.CheckAlarms()

	assert.Equal(t, mon.Alarms, []string{"db1,Ping", "backend1,DBConnection"})
}

func TestBackendDownAffectFrontendAlarm(t *testing.T) {
	mon := &fakeMonSys{}

	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}
	frontend1 := &Frontend{Server: Server{Name: "frontend1", mon: mon}, Backend: backend1}

	// Set the backend down
	backend1.PingAlarm = AlarmTriggered

	// Check alarms
	db1.CheckAlarms()
	backend1.CheckAlarms()
	frontend1.CheckAlarms()

	assert.Equal(t, mon.Alarms, []string{"backend1,Ping", "frontend1,BackendConnection"})
}
