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

	MotionData []History_motionData
}

type List_sessionData struct {
	M_header PacketHeader

	SessionData []History_sessionData
}

type List_lapData struct {
	M_header PacketHeader

	LapData []LapData_lap_group
}

type LapData_lap_group struct {
	LapNum       int
	LapData_list []History_lapData
}

type List_telemetryData struct {
	M_header PacketHeader

	TelemetryData []History_telemetryData
}

type List_statusData struct {
	M_header PacketHeader

	StatusData []History_statusData
}

type History_motionData struct {
	Frame_identifier       int
	Suspension_position_rl float32
	Suspension_position_rr float32
	Suspension_position_fl float32
	Suspension_position_fr float32
	M_worldPositionX       float32
	M_worldPositionY       float32
	M_worldPositionZ       float32
}

type History_sessionData struct {
	Mrame_identifier       int
	Suspension_position_rl float32
	Suspension_position_rr float32
	Suspension_position_fl float32
	Suspension_position_fr float32
	M_worldPositionX       float32
	M_worldPositionY       float32
	M_worldPositionZ       float32
	M_totalLaps            int
	M_trackId              int
}

type History_lapData struct {
	Frame_identifier int
	M_lastLapTime    float32
	M_currentLapTime float32
	M_bestLapTime    float32
	M_sector1Time    float32
	M_sector2Time    float32
	M_currentLapNum  int
	M_sector         int
	M_penalties      int
}

type History_telemetryData struct {
	Frame_identifier             int
	M_speed                      int
	M_throttle                   int
	M_brake                      int
	M_gear                       int
	M_engineRPM                  int
	M_brakesTemperature_rl       int
	M_brakesTemperature_rr       int
	M_brakesTemperature_fl       int
	M_brakesTemperature_fr       int
	M_tyresSurfaceTemperature_rl int
	M_tyresSurfaceTemperature_rr int
	M_tyresSurfaceTemperature_fl int
	M_tyresSurfaceTemperature_fr int
	M_tyresPressure_rl           float32
	M_tyresPressure_rr           float32
	M_tyresPressure_fl           float32
	M_tyresPressure_fr           float32
}

type History_statusData struct {
	Frame_identifier int
	M_maxRPM         int
	M_idleRPM        int
	M_maxGears       int
	M_tyresWear_rl   int
	M_tyresWear_rr   int
	M_tyresWear_fl   int
	M_tyresWear_fr   int
	M_tyresDamage_rl int
	M_tyresDamage_rr int
	M_tyresDamage_fl int
	M_tyresDamage_fr int
}
