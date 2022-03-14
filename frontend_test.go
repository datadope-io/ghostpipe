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

	time := 0.0

	// Set the backend down
	db1.PingAlarm = AlarmTriggered

	// Check alarms
	db1.CheckAlarms(time)
	backend1.CheckAlarms(time)
	frontend1.CheckAlarms(time)

	assert.Equal(t, mon.Alarms, []string{"0,db1,Ping", "0,backend1,DBConnection"})
}

func TestBackendDownAffectFrontendAlarm(t *testing.T) {
	mon := &fakeMonSys{}

	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}
	frontend1 := &Frontend{Server: Server{Name: "frontend1", mon: mon}, Backend: backend1}

	time := 0.0

	// Set the backend down
	backend1.PingAlarm = AlarmTriggered

	// Check alarms
	db1.CheckAlarms(time)
	backend1.CheckAlarms(time)
	frontend1.CheckAlarms(time)

	assert.Equal(t, mon.Alarms, []string{"0,backend1,Ping", "0,frontend1,BackendConnection"})
}
