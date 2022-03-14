package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBDownAffectBackendAlarm(t *testing.T) {
	mon := &fakeMonSys{}
	db1 := &Database{Server: Server{Name: "db1", mon: mon}}
	backend1 := &Backend{Server: Server{Name: "backend1", mon: mon}, DBEngine: db1}

	// There should be no alarms, since we haven't checked db neither backend
	assert.Len(t, mon.Alarms, 0)

	// Check for alarms in both serves
	db1.CheckAlarms()
	backend1.CheckAlarms()

	// There should be no alarms, since we db and backend are initializated
	// with no problems
	assert.Len(t, mon.Alarms, 0)

	// Set the db down
	db1.PingAlarm = AlarmTriggered

	// Database now should report an alarm
	db1.CheckAlarms()
	assert.Equal(t, mon.Alarms, []string{"db1,Ping"})

	// The backend should also report an alarm regarding the connection to the db
	backend1.CheckAlarms()
	assert.Equal(t, mon.Alarms, []string{"db1,Ping", "backend1,DBConnection"})
}

// TeTestDBDownAffectMultipleBackends tests the case when a db is down and affects multiple backends
func TestDBDownAffectMultipleBackends(t *testing.T) {
	mon := &fakeMonSys{}
	db1 := &Database{Server: Server{Name: "db1", mon: mon}}

	backends := []*Backend{}

	for i := 0; i < 3; i++ {
		backend := &Backend{Server: Server{Name: fmt.Sprintf("backend%d", i), mon: mon}, DBEngine: db1}
		backends = append(backends, backend)
	}

	// Set the db down
	db1.PingAlarm = AlarmTriggered

	// Check alarms in db and all backends
	db1.CheckAlarms()
	for _, backend := range backends {
		backend.CheckAlarms()
	}

	expectedAlarms := []string{"db1,Ping"}
	for _, backend := range backends {
		expectedAlarms = append(expectedAlarms, fmt.Sprintf("%s,DBConnection", backend.Name))
	}

	assert.Equal(t, mon.Alarms, expectedAlarms)
}

func TestDBRetriggerAffectsBackend(t *testing.T) {
	mon := &fakeMonSys{}

	db1 := &Database{
		Server: Server{
			Name: "db1",
			mon:  mon,
		},
	}

	backend1 := &Backend{
		Server: Server{
			Name: "backend1",
			mon:  mon,
		},
		DBEngine: db1,
	}

	// Trigger the ping alarm in the db, recover it and trigger again
	db1.PingAlarm = AlarmTriggered
	db1.CheckAlarms()
	backend1.CheckAlarms()

	db1.PingAlarm = AlarmEnabled
	db1.CheckAlarms()
	backend1.CheckAlarms()

	// Clear the alarm history to have only the last alarms
	mon.Alarms = []string{}

	// A second trigger of the db ping alarm should generate the alarm for the db and the backend connection
	db1.PingAlarm = AlarmTriggered
	db1.CheckAlarms()
	backend1.CheckAlarms()

	assert.Equal(t, mon.Alarms, []string{"db1,Ping", "backend1,DBConnection"})
}
