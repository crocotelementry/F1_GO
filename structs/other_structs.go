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

// type RaceSpeed_struct_mini struct {
// 	raceSpeed [1000]interface {}
// }

// type CatchUp_data struct {
// 	Data []int
// }

type CatchUp_dashboard_struct struct {
	M_header PacketHeader // Header

	RaceSpeed_data           []int
	EngineRevs_data          []int
	GearChanges_data         []int
	ThrottleApplication_data []int
	BrakeApplication_data    []int
}

type CatchUp_time_struct struct {
	M_header PacketHeader // Header

	Lap_num     []int
	Lap_time    []float64
	Sector1Time []float64
	Sector2Time []float64
	Sector3Time []float64
	PitStatus   []int
}

type History_status struct {
	M_header PacketHeader

	History_status string
}

type List_motionData struct {
	M_header PacketHeader

	motionData []History_motionData
}

type List_sessionData struct {
	M_header PacketHeader

	sessionData []History_sessionData
}

type List_lapData struct {
	M_header PacketHeader

	lapData []History_lapData
}

type List_telemetryData struct {
	M_header PacketHeader

	telemetryData []History_telemetryData
}

type List_statusData struct {
	M_header PacketHeader

	statusData []History_statusData
}

type History_motionData struct {
	frame_identifier       int
	suspension_position_rl float32
	suspension_position_rr float32
	suspension_position_fl float32
	suspension_position_fr float32
	m_worldPositionX       float32
	m_worldPositionY       float32
	m_worldPositionZ       float32
}

type History_sessionData struct {
	frame_identifier       int
	suspension_position_rl float32
	suspension_position_rr float32
	suspension_position_fl float32
	suspension_position_fr float32
	m_worldPositionX       float32
	m_worldPositionY       float32
	m_worldPositionZ       float32
	m_totalLaps            int
	m_trackId              int
}

type History_lapData struct {
	frame_identifier int
	m_lastLapTime    float32
	m_currentLapTime float32
	m_bestLapTime    float32
	m_sector1Time    float32
	m_sector2Time    float32
	m_sector         int
}

type History_telemetryData struct {
	frame_identifier             int
	m_speed                      int
	m_throttle                   int
	m_brake                      int
	m_gear                       int
	m_engineRPM                  int
	m_brakesTemperature_rl       int
	m_brakesTemperature_rr       int
	m_brakesTemperature_fl       int
	m_brakesTemperature_fr       int
	m_tyresSurfaceTemperature_rl int
	m_tyresSurfaceTemperature_rr int
	m_tyresSurfaceTemperature_fl int
	m_tyresSurfaceTemperature_fr int
	m_tyresPressure_rl           float32
	m_tyresPressure_rr           float32
	m_tyresPressure_fl           float32
	m_tyresPressure_fr           float32
}

type History_statusData struct {
	frame_identifier int
	m_maxRPM         int
	m_idleRPM        int
	m_maxGears       int
	m_tyresWear_rl   int
	m_tyresWear_rr   int
	m_tyresWear_fl   int
	m_tyresWear_fr   int
	m_tyresDamage_rl int
	m_tyresDamage_rr int
	m_tyresDamage_fl int
	m_tyresDamage_fr int
}
