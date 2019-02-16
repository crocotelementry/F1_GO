package structs

import (
	"time"
)

type RaceEventDirectory struct {
	M_header           PacketHeader // Header
	Session_start_time time.Time
	Session_end_time   time.Time
}

type Session_start struct {
	Session_start_time time.Time
}

type Session_end struct {
	Session_end_time time.Time
}

// Structs for database alerts
type Session struct {
	Session_UID        string
	Session_start_time time.Time
	Session_end_time   time.Time
}

type Save_to_database_alerts struct {
	M_header PacketHeader // Header

	Num_of_sessions int
	Sessions        []Session
}

type Save_to_database_websocket_recive struct {
	Type string

	Uid uint64
}
