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

type Save_to_database_status struct {
	M_header       PacketHeader // Header
	Status         string
	UID            string
	Current_packet int
	Packet_0       int
	Packet_0_total int
	Packet_1       int
	Packet_1_total int
	Packet_2       int
	Packet_2_total int
	Packet_3       int
	Packet_3_total int
	Packet_4       int
	Packet_4_total int
	Packet_5       int
	Packet_5_total int
	Packet_6       int
	Packet_6_total int
	Packet_7       int
	Packet_7_total int
	Total_current  int
	Total_packets  int
}
