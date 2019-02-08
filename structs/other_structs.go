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
