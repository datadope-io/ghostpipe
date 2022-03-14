# Ghostpipe

Simulate alarmas in different architectures.

With Ghostpipe we could define an architecture of different servers interconnected and then simulate what
happens when some servers have problems, getting the events generated for each server.

It also generate a graph with the architecture (by default, ``graph.cyjs``).

We can use ``show_graph.py`` to show that graph.

Each server runs it's own goroutine, simulating the behaviour of a local monitoring agent
generating events when there is some problem.

The problems could be produced manually in the code and others are generated based on a relationships
between the servers. For example, a backend server will generate a ``DBConnectionAlarm`` if its database
is not available.

Example graph generated with Ghostpipe and drawed with ``show_graph.py``:
![example graph](example_graph.png)

Example output:
```
Writing graph to graph.cyjs
Starting simulator...

34,noise8,Memory
108,noise5,CPU
136,noise3,Memory
200,noise3,CPU
218,backendC,DBConnection
220,backendB,DBConnection
236,backendA,DBConnection
259,db1,Ping
291,noise4,Disk
356,noise4,Ping
386,noise3,Disk
398,frontendD1,BackendConnection
405,backendC,DBConnection
406,backendB,DBConnection
414,backendD,Ping
417,backendA,DBConnection
428,noise1,Ping
447,db1,Ping
488,noise1,Disk
567,noise3,CPU

Stopping simulator...
```

## Architecture

Each server have alarms and it is considered available based on some rule based on those alarms.

The alarms in ``Server`` are inherited for the rest of servers.

### Server (base)

| Alarms | Availability | Notes |
|-----|--|--|
| CPU | | |
| Memory | | |
| Disk | | |
| Ping | X | |

### Database

| Alarms | Availability | Notes |
|-----|--|--|
| DBEngine | X | Availability take into account also Server.Ping |

### Backend

All backends should be connected to one database.

| Alarms | Availability | Notes |
|-----|--|--|
| Proc | X | Availability take into account also Server.Ping |
| DBConnection | | Its triggered if the connected DB is not available |

### Frontend

All frontends should be connected to one backend.

| Alarms | Availability | Notes |
|-----|--|--|
| Proc | X | Availability take into account also Server.Ping |
| BackendConnection | | Its triggered if the connected backend is not available |
