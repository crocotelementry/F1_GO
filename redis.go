package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/crocotelementry/F1_GO/structs"
	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
)

var (
	header                                 structs.PacketHeader
	Motion_packet                          structs.PacketMotionData
	Session_packet                         structs.PacketSessionData
	Lap_packet                             structs.PacketLapData
	Event_packet                           structs.PacketEventData
	Participant_packet                     structs.PacketParticipantsData
	Car_setup_packet                       structs.PacketCarSetupData
	Telemetry_packet                       structs.PacketCarTelemetryData
	Car_status_packet                      structs.PacketCarStatusData
	race_event_directory_struct            structs.RaceEventDirectory
	Session_start                          structs.Session_start
	Session_end                            structs.Session_end
	Save_to_database_alerts                structs.Save_to_database_alerts
	Session                                structs.Session
	session_start_code                     = [4]uint8{83, 69, 78, 68}
	session_end_code                       = [4]uint8{83, 83, 84, 65}
	redis_ping_done                        = make(chan bool)
	redis_pool                             = newPool() // newPool returns a pointer to a redis.Pool
	incrementing_motion_packet_number      = 0
	incrementing_session_packet_number     = 0
	incrementing_lap_packet_number         = 0
	incrementing_event_packet_number       = 0
	incrementing_participant_packet_number = 0
	incrementing_car_setup_packet_number   = 0
	incrementing_telemetry_packet_number   = 0
	incrementing_car_status_packet_number  = 0
	atm_motion_packet                      = make(chan *structs.PacketMotionData)
	atm_session_packet                     = make(chan *structs.PacketSessionData)
	atm_lap_packet                         = make(chan *structs.PacketLapData)
	atm_event_packet                       = make(chan *structs.PacketEventData)
	atm_participant_packet                 = make(chan *structs.PacketParticipantsData)
	atm_car_setup_packet                   = make(chan *structs.PacketCarSetupData)
	atm_telemetry_packet                   = make(chan *structs.PacketCarTelemetryData)
	atm_car_status_packet                  = make(chan *structs.PacketCarStatusData)
	atm_race_event_directory               = make(chan structs.RaceEventDirectory)
	redis_done                             = make(chan bool)
	current_session_uid                    = uint64(0)
	previous_session_uid                   = uint64(0)
)

// Client is a middleman between the websocket connection and the hub.
type Udp_data struct {
	Id int

	Motion_packet           structs.PacketMotionData
	Session_packet          structs.PacketSessionData
	Lap_packet              structs.PacketLapData
	Event_packet            structs.PacketEventData
	Participant_packet      structs.PacketParticipantsData
	Car_setup_packet        structs.PacketCarSetupData
	Telemetry_packet        structs.PacketCarTelemetryData
	Car_status_packet       structs.PacketCarStatusData
	Save_to_database_alert  structs.Save_to_database_alerts
	Save_to_database_status structs.Save_to_database_status
}

// To establish connectivity in redigo, you need to create a redis.Pool object which is a pool of connections to Redis.
func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	// fmt.Println("PING Response = ", s)
	fmt.Print("Redis connection       ")

	if s == "PONG" {
		redis_ping_done <- true
		color.Green("Success")
	} else {
		redis_ping_done <- false
		color.Red("Error")
	}

	// Output: PONG
	return nil
}

// Function that grabs a specific Session_uid worth of data from the redis database and
// sends it over to mysql to be put into long term mysql
func getRedisDataForMysql(hub *Hub, chosen_session_uid uint64) {
	// get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()

	// Defer the closing of the redis connection until we return at the end of getRedisData
	defer redis_conn.Close()

	// Start up add_to_longterm_storage in another goroutine
	go add_to_longterm_storage()

	// log.Println("1")

	session_uid := strconv.FormatUint(chosen_session_uid, 10)

	save_Motion_packet := new(structs.PacketMotionData)
	save_Session_packet := new(structs.PacketSessionData)
	save_Lap_packet := new(structs.PacketLapData)
	save_Event_packet := new(structs.PacketEventData)
	save_Participant_packet := new(structs.PacketParticipantsData)
	save_Car_setup_packet := new(structs.PacketCarSetupData)
	save_Telemetry_packet := new(structs.PacketCarTelemetryData)
	save_Car_status_packet := new(structs.PacketCarStatusData)

	// log.Println("2")

	// Send over the initial data for race_event_directory
	race_event_directory_data, err := (redis_conn.Do("GET", session_uid+":0:0"))
	if err != nil {
		log.Println("Getting initial data for race_event_directory from redis database failed:", err)
	}

	// log.Println("race_event_directory_data", race_event_directory_data)

	err = json.Unmarshal(race_event_directory_data.([]byte), &Motion_packet)
	if err != nil {
		log.Println(err)
	}

	// log.Println("4")
	// atm_race_event_directory <- Header

	// Send over our session_start_time and session_end_time
	// Send over the initial data for race_event_directory
	session_start_time, err := (redis_conn.Do("GET", session_uid+":session_start_time"))
	if err != nil {
		log.Println("Getting session_start_time from redis database failed:", err)
	}

	session_end_time, err := (redis_conn.Do("GET", session_uid+":session_end_time"))
	if err != nil {
		log.Println("Getting session_end_time from redis database failed:", err)
	}

	err = json.Unmarshal(session_start_time.([]byte), &Session_start)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(session_end_time.([]byte), &Session_end)
	if err != nil {
		fmt.Println(err)
	}

	atm_race_event_directory <- structs.RaceEventDirectory{Motion_packet.M_header, Session_start.Session_start_time, Session_end.Session_end_time}

	progress_status := structs.Save_to_database_status{
		M_header: structs.PacketHeader{
			M_packetId: 31,
		},
		Status:         "initial",
		UID:            session_uid,
		Current_packet: 0,
		Packet_0:       0,
		Packet_0_total: 0,
		Packet_1:       0,
		Packet_1_total: 0,
		Packet_2:       0,
		Packet_2_total: 0,
		Packet_3:       0,
		Packet_3_total: 0,
		Packet_4:       0,
		Packet_4_total: 0,
		Packet_5:       0,
		Packet_5_total: 0,
		Packet_6:       0,
		Packet_6_total: 0,
		Packet_7:       0,
		Packet_7_total: 0,
		Total_current:  0,
		Total_packets:  0,
	}

	for packet_type := 0; packet_type < 8; packet_type += 1 {
		max_packet_number, err := redis.Int(redis_conn.Do("GET", session_uid+":"+strconv.Itoa(packet_type)+":Incrementing_packet_number"))
		if err != nil {
			log.Println("Getting max_packet_number for packet id number", packet_type, "from redis database failed:", err)
		}
		switch packet_type {
		case 0:
			progress_status.Packet_0_total = max_packet_number
		case 1:
			progress_status.Packet_1_total = max_packet_number
		case 2:
			progress_status.Packet_2_total = max_packet_number
		case 3:
			progress_status.Packet_3_total = max_packet_number
		case 4:
			progress_status.Packet_4_total = max_packet_number
		case 5:
			progress_status.Packet_5_total = max_packet_number
		case 6:
			progress_status.Packet_6_total = max_packet_number
		case 7:
			progress_status.Packet_7_total = max_packet_number
			progress_status.Total_packets = progress_status.Packet_0_total + progress_status.Packet_1_total + progress_status.Packet_2_total + progress_status.Packet_3_total + progress_status.Packet_4_total + progress_status.Packet_5_total + progress_status.Packet_6_total + progress_status.Packet_7_total
		}
	}

	hub.broadcast <- &Udp_data{
		Id:                      31,
		Save_to_database_status: progress_status,
	}

	progress_status.Status = "Saving"

	for packet_type := 0; packet_type < 8; packet_type += 1 {
		max_packet_number, err := redis.Int(redis_conn.Do("GET", session_uid+":"+strconv.Itoa(packet_type)+":Incrementing_packet_number"))
		if err != nil {
			log.Println("Getting max_packet_number for packet id number", packet_type, "from redis database failed:", err)
		}

		switch packet_type {
		case 0:
			progress_status.Current_packet = 0
		case 1:
			progress_status.Current_packet = 1
		case 2:
			progress_status.Current_packet = 2
		case 3:
			progress_status.Current_packet = 3
		case 4:
			progress_status.Current_packet = 4
		case 5:
			progress_status.Current_packet = 5
		case 6:
			progress_status.Current_packet = 6
		case 7:
			progress_status.Current_packet = 7
		}

		for packet_number := 0; packet_number < max_packet_number; packet_number += 1 {
			progress_status.Total_current += 1

			packet_data, err := redis_conn.Do("GET", session_uid+":"+strconv.Itoa(packet_type)+":"+strconv.Itoa(packet_number))
			if err != nil {
				log.Println("Getting packet_data for packet id number", packet_type, "with packet_number", packet_number, "from redis database failed:", err)
			}

			switch packet_type {
			case 0:
				err := json.Unmarshal(packet_data.([]byte), &save_Motion_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_motion_packet <- save_Motion_packet

				progress_status.Packet_0 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 1:
				err := json.Unmarshal(packet_data.([]byte), &save_Session_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_session_packet <- save_Session_packet

				progress_status.Packet_1 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 2:
				err := json.Unmarshal(packet_data.([]byte), &save_Lap_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_lap_packet <- save_Lap_packet

				progress_status.Packet_2 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 3:
				err := json.Unmarshal(packet_data.([]byte), &save_Event_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_event_packet <- save_Event_packet

				progress_status.Packet_3 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 4:
				err := json.Unmarshal(packet_data.([]byte), &save_Participant_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_participant_packet <- save_Participant_packet

				progress_status.Packet_4 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 5:
				err := json.Unmarshal(packet_data.([]byte), &save_Car_setup_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_car_setup_packet <- save_Car_setup_packet

				progress_status.Packet_5 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 6:
				err := json.Unmarshal(packet_data.([]byte), &save_Telemetry_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_telemetry_packet <- save_Telemetry_packet

				progress_status.Packet_6 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			case 7:
				err := json.Unmarshal(packet_data.([]byte), &save_Car_status_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_car_status_packet <- save_Car_status_packet

				progress_status.Packet_7 = packet_number
				hub.broadcast <- &Udp_data{
					Id:                      31,
					Save_to_database_status: progress_status,
				}

			}
		}
	}

	progress_status.Status = "done"
	redis_done <- true
	hub.broadcast <- &Udp_data{
		Id:                      31,
		Save_to_database_status: progress_status,
	}
	return
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b [4]uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Checks if session_UID exists
func session_UID_exists(session_uid uint64) bool {
	// get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()

	// Defer the closing of the redis connection until we return at the end of getRedisData
	defer redis_conn.Close()

	session_UID_exists, err := redis_conn.Do("SISMEMBER", "session_UIDs", session_uid)
	if err != nil {
		log.Println("Getting session_UIDs from redis database failed:", err)
	}

	if session_UID_exists == int64(1) {
		return true
	}

	return false
}

func getGameData(hub *Hub) {
	// get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()
	// use defer to close the connection when the function completes

	c := make(chan os.Signal, 2)
	// When we close F1_go by using control-c, we will catch it, flush the redis database empty and then
	// close the connection to the redis database before closing F1_GO
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		redis_conn.Do("FlushAll")
		fmt.Println("\n")
		log.Println("               redis flushed")
		redis_conn.Close()
		os.Exit(1)
	}()

	defer func() {
		redis_conn.Close()
	}()

	// call Redis PING command to test connectivity
	err := ping(redis_conn)
	if err != nil {
		fmt.Println("Problem with connection to Redis database", err)
	}

	// Set number of SETs to redis database to zero
	incrementing_packet_number := 0

	// Set new session start code recieved to false.
	recieved_new_uid := false

	// Create a reference point for our current lap. This is for adding things to catchup_packet number 2 stuff
	current_lap_number := uint8(0)

	// Set up sector times
	sector1_time := float32(0)
	sector2_time := float32(0)

	// Redis database format:
	// Session_uid:packet_id:incrementing_packet_number									This is for packets
	// Session_uid:"Incrementing_packet_number"								This is for knowing the max value of incrementing_packet_number
	// session_UIDs																						This is a list for keeping track of what session_UIDs are in the redis database

	// Set the number of sessions to 0
	if _, err := redis_conn.Do("SET", "number_of_sessions", "0"); err != nil {
		fmt.Println("Initializing number_of_sessions to 0 failed:", err)
	}

	for {
		buf := make([]byte, 1341)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("readfromudp error::: ", err)
		}

		// Set a new reader which we will use to cast into our structs.
		// This reader is for the header, which we determine what packet we have and what index our users car is in.
		// Bytes 3 in the udp packet will be the packet number and byte 20 will be the index of the users car.
		header_bytes_reader := bytes.NewReader(buf[0:21])
		packet_bytes_reader := bytes.NewReader(buf)

		// Read the binary of the udp packet header into our struct
		if err := binary.Read(header_bytes_reader, binary.LittleEndian, &header); err != nil {
			fmt.Println("binary.Read header failed:", err)
		}

		session_uid_Prefix := strconv.FormatUint(header.M_sessionUID, 10)

		// Depending on which packet we have, which we find by looking at header.M_packetId
		// We use a switch statement to then read the whole binary udp packet into its associated struct

		if header.M_packetId != 3 {
			if recieved_new_uid == false {
				// Set new session start code recieved to true
				recieved_new_uid = true

				// Set the current_session_uid to the current_session_uid
				previous_session_uid = current_session_uid
				current_session_uid = header.M_sessionUID

				log.Println("New Session Start! uid:", header.M_sessionUID)

				// fmt.Println("current_session_uid:", current_session_uid)

				session_UIDs_SADD_integer_reply, err := redis_conn.Do("SADD", "session_UIDs", (session_uid_Prefix))
				if err != nil {
					fmt.Println("SADD session_UIDs failed:", err)
				}

				// log.Println("Packets recived for new session without recivicing a session_start_code \nCould be due to starting up during an already ongoing session:")

				if session_UIDs_SADD_integer_reply == int64(0) {
					fmt.Println("\nSession with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", session_uid_Prefix, "\n")
				} else {

					if _, err := redis_conn.Do("INCR", "number_of_sessions"); err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}

					if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":session_start_time"), &structs.Session_start{time.Now()}); err != nil {
						fmt.Println("Setting session_start_time to failed:", err)
					}

				}
			}
		}

		// fmt.Println("header.M_sessionUID:", header.M_sessionUID)

		switch header.M_packetId {
		case 0:
			// If the packet we received is a motion_packet, read its binary into our motion_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Motion_packet); err != nil {
				fmt.Println("binary.Read motion_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:            0,
				Motion_packet: Motion_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_motion_packet, err := json.Marshal(Motion_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Motion_packet.M_header.M_sessionUID, 10) + ":0:" + strconv.Itoa(incrementing_motion_packet_number)), json_motion_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_motion_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_motion_packet_number += 1
		case 1:
			// If the packet we received is the session_packet, read its binary into our session_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Session_packet); err != nil {
				fmt.Println("binary.Read session_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:             1,
				Session_packet: Session_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_session_packet, err := json.Marshal(Session_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Session_packet.M_header.M_sessionUID, 10) + ":1:" + strconv.Itoa(incrementing_session_packet_number)), json_session_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_session_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_session_packet_number += 1
		case 2:
			// If the packet we received is the lap_packet, read its binary into our lap_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Lap_packet); err != nil {
				fmt.Println("binary.Read lap_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:         2,
				Lap_packet: Lap_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_lap_packet, err := json.Marshal(Lap_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Lap_packet.M_header.M_sessionUID, 10) + ":2:" + strconv.Itoa(incrementing_lap_packet_number)), json_lap_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_lap_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_lap_packet_number += 1

			// Add data from lap_packet to catchUp stuff
			//
			//

			if sector1_time == 0 && Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector1Time != 0 {
				sector1_time = Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector1Time

			} else if sector2_time == 0 && Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector2Time != 0 {
				sector2_time = Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector2Time

			}

			// If our Lap_packet lap number is not the same as our current_lap_number
			if Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_currentLapNum != current_lap_number {
				// if our current_lap_number is not the default, we started a new lap! add the last laps data to catchup data.
				// If its not, set it to the current lap
				log.Println("lap finished, current lap num:", current_lap_number, "current packet lap num:", Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_currentLapNum)

				if current_lap_number != 0 {

					log.Println("Lap is not 0")

					// sector1_time := Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector1Time
					// sector2_time := Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_sector2Time
					last_lap_time := Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_lastLapTime

					if _, err := redis_conn.Do("LPUSH", "catchup_lap_num", current_lap_number); err != nil {
						fmt.Println("Lpush catchup_lap_num to reddis catchUp failed:", err)
					}

					if _, err := redis_conn.Do("LPUSH", "catchup_lap_time", last_lap_time); err != nil {
						fmt.Println("Lpush catchup_lap_time to reddis catchUp failed:", err)
					}

					if _, err := redis_conn.Do("LPUSH", "catchup_sector1Time", sector1_time); err != nil {
						fmt.Println("Lpush catchup_sector1Time to reddis catchUp failed:", err)
					}

					if _, err := redis_conn.Do("LPUSH", "catchup_sector2Time", sector2_time); err != nil {
						fmt.Println("Lpush catchup_sector2Time to reddis catchUp failed:", err)
					}

					if _, err := redis_conn.Do("LPUSH", "catchup_sector3Time", last_lap_time-sector2_time-sector1_time); err != nil {
						fmt.Println("Lpush catchup_sector3Time to reddis catchUp failed:", err)
					}

					if _, err := redis_conn.Do("LPUSH", "catchup_pitStatus", Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_pitStatus); err != nil {
						fmt.Println("Lpush catchup_pitStatus to reddis catchUp failed:", err)
					}

					current_lap_number = Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_currentLapNum
					sector1_time = float32(0)
					sector2_time = float32(0)

				} else {
					log.Println("lap is 0")
					current_lap_number = Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex].M_currentLapNum
				}
			}

		case 3:
			// If the packet we received is the event_packet, read its binary into our event_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Event_packet); err != nil {
				fmt.Println("binary.Read event_packet failed:", err)
			}

			// fmt.Println("New EVENT PACKET:", Event_packet)

			if Equal(Event_packet.M_eventStringCode, session_start_code) {
				if _, err := redis_conn.Do("DEL", "catchup_lap_num catchup_lap_time catchup_sector1Time catchup_sector2Time catchup_sector3Time catchup_pitStatus"); err != nil {
					fmt.Println("Del Lap times from redis Catchup failed:", err)
				}

				// Create a reference point for our current lap. This is for adding things to catchup_packet number 2 stuff
				current_lap_number = uint8(0)
			}
			// Set new session start code recieved to true
			// recieved_new_uid = true

			// number_of_sessions_exists_integer_reply, err := redis_conn.Do("EXISTS", "number_of_sessions")

			// Set the previous UID to previous and then
			// Set the current_session_uid to the current_session_uid
			// previous_session_uid = current_session_uid
			// current_session_uid = Event_packet.M_header.M_sessionUID

			// fmt.Println("session start code recieved! uid:", header.M_sessionUID)
			// fmt.Println("Current session uid", current_session_uid)

			// if err != nil {
			// 	fmt.Println("Checking if number_of_sessions exists failed:", err)
			// }

			// if _, err := redis_conn.Do("INCR", "number_of_sessions"); err != nil {
			// 	fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
			// }

			// if number_of_sessions_exists_integer_reply == int64(1) {
			// 	if _, err := redis_conn.Do("INCR", "number_of_sessions"); err != nil {
			// 		fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
			// 	}
			// } else {
			// 	if _, err := redis_conn.Do("SET", "number_of_sessions", "1"); err != nil {
			// 		fmt.Println("Setting number_of_sessions to 1 failed:", err)
			// 	}
			// }

			// session_UIDs_exists_integer_reply, err := redis_conn.Do("EXISTS", "session_UIDs")
			// if err != nil {
			// 	fmt.Println("Checking if session_UIDs exists failed:", err)
			// }

			// session_UIDs_SADD_integer_reply, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID))
			// if err != nil {
			// 	fmt.Println("SADD session_UIDs failed:", err)
			// }
			//
			// if session_UIDs_SADD_integer_reply == int64(0) {
			// 	fmt.Println("\nSession with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", Event_packet.M_header.M_sessionUID, "\n")
			// }

			// if session_UIDs_exists_integer_reply == int64(1) {
			// 	session_UIDs_SADD_integer_reply, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID))
			// 	if err != nil {
			// 		fmt.Println("SADD session_UIDs failed:", err)
			// 	}
			//
			// 	fmt.Println("sadd session uid:", session_UIDs_SADD_integer_reply)
			//
			// 	if session_UIDs_SADD_integer_reply == int64(0) {
			// 		fmt.Println("\nSession with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", Event_packet.M_header.M_sessionUID, "\n")
			// 	}
			// 	// else {
			// 	// 	incrementing_packet_number = 0
			// 	// }
			// } else {
			// 	if _, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID)); err != nil {
			// 		fmt.Println("SADD session_UIDs failed:", err)
			// 	}
			// }

			// Add session start time to redis database
			// Format is as follows:
			// Session_UID:session_start_time
			// if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":session_start_time"), &structs.Session_start{time.Now()}); err != nil {
			// 	fmt.Println("Setting session_start_time to failed:", err)
			// }
			// }

			// If we receive a session end code, send a data alert over the websocket to ask the user if they want to save the session for long
			// term use in a MYSQL database, discard the session, or hold on to it until our redis database reaches its size limit or its amount
			// of sessions limit.
			if Equal(Event_packet.M_eventStringCode, session_end_code) {

				// Set new session start code recieved to false
				recieved_new_uid = false

				log.Println("session end code recieved! uid:", header.M_sessionUID)

				// fmt.Println(strconv.FormatUint(Event_packet.M_header.M_sessionUID, 10) + ":Incrementing_packet_number")

				// sue_return := session_UID_exists(Event_packet.M_header.M_sessionUID)
				//
				// // fmt.Println(sue_return)
				//
				// if sue_return == false {
				// 	_, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID))
				// 	if err != nil {
				// 		fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
				// 	}
				// }

				// Add session end time to redis database
				// Format is as follows:
				// Session_UID:session_end_time

				temp_uid := strconv.FormatUint(Event_packet.M_header.M_sessionUID, 10)

				if _, err := redis_conn.Do("SET", (temp_uid + ":session_end_time"), &structs.Session_end{time.Now()}); err != nil {
					fmt.Println("Setting session_end_time to failed:", err)
				}

				// Session_uid:Incrementing_packet_number
				if _, err := redis_conn.Do("SET", (temp_uid + ":Incrementing_packet_number"), strconv.Itoa(int(incrementing_packet_number))); err != nil {
					log.Println("             ", "Setting Incrementing_packet_number failed:", err)
				}

				// set incrementing_motion_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":0:Incrementing_packet_number"), incrementing_motion_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_motion_packet_number failed:", err)
				}

				// set incrementing_session_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":1:Incrementing_packet_number"), incrementing_session_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_session_packet_number failed:", err)
				}

				// set incrementing_lap_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":2:Incrementing_packet_number"), incrementing_lap_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_lap_packet_number failed:", err)
				}

				// set incrementing_event_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":3:Incrementing_packet_number"), incrementing_event_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_event_packet_number failed:", err)
				}

				// set incrementing_participant_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":4:Incrementing_packet_number"), incrementing_participant_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_participant_packet_number failed:", err)
				}

				// set incrementing_car_setup_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":5:Incrementing_packet_number"), incrementing_car_setup_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_car_setup_packet_number failed:", err)
				}

				// set incrementing_telemetry_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":6:Incrementing_packet_number"), incrementing_telemetry_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_telemetry_packet_number failed:", err)
				}

				// set incrementing_car_status_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (temp_uid + ":7:Incrementing_packet_number"), incrementing_car_status_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_car_status_packet_number failed:", err)
				}

				// Recieved a session end code, set all packet counting variables to zero
				incrementing_motion_packet_number = 0
				incrementing_session_packet_number = 0
				incrementing_lap_packet_number = 0
				incrementing_event_packet_number = 0
				incrementing_participant_packet_number = 0
				incrementing_car_setup_packet_number = 0
				incrementing_telemetry_packet_number = 0
				incrementing_car_status_packet_number = 0
				incrementing_packet_number = 0

				// Uncomment for easy command line testing
				// go getRedisData()

				// Get the Session_UIDs of the sessions we have stored in our redis database
				session_UIDs, err := redis_conn.Do("SMEMBERS", "session_UIDs")
				if err != nil {
					log.Println("Getting session_UIDs from redis database failed:", err)
				}

				// Get the session UIDS and add them to our Save_to_database_alert struct to be sent over the websocket
				// fmt.Print("session_UIDs:")
				redis_sessions := []structs.Session{}
				session_uids := reflect.ValueOf(session_UIDs)
				for i := 0; i < session_uids.Len(); i += 1 {
					uid, _ := redis.Uint64(session_uids.Index(i).Interface(), nil)

					uid_string := strconv.FormatUint(uid, 10)

					session_start_time, err := (redis_conn.Do("GET", uid_string+":session_start_time"))
					if err != nil {
						log.Println("Getting session_start_time from redis database failed:", err)
					}

					session_end_time, err := (redis_conn.Do("GET", uid_string+":session_end_time"))
					if err != nil {
						log.Println("Getting session_end_time from redis database failed:", err)
					}

					err = json.Unmarshal(session_start_time.([]byte), &Session_start)
					if err != nil {
						fmt.Println(err)
					}

					err = json.Unmarshal(session_end_time.([]byte), &Session_end)
					if err != nil {
						fmt.Println(err)
					}

					redis_sessions = append(redis_sessions, structs.Session{Session_UID: uid_string, Session_start_time: Session_start.Session_start_time, Session_end_time: Session_end.Session_end_time})
				}

				// fmt.Println("redis_sessions", redis_sessions)

				// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
				// be broadcasted to all connected clients
				hub.broadcast <- &Udp_data{
					Id: 30,

					Save_to_database_alert: structs.Save_to_database_alerts{
						M_header: structs.PacketHeader{
							M_packetId: 30,
						},
						Num_of_sessions: session_uids.Len(),
						Sessions:        redis_sessions,
					},
				}
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:           3,
				Event_packet: Event_packet,
			}

		case 4:
			// If the packet we received is the participant_packet, read its binary into our participant_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Participant_packet); err != nil {
				fmt.Println("binary.Read participant_packet failed:", err)
			}

			hub.broadcast <- &Udp_data{
				Id:                 4,
				Participant_packet: Participant_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_participant_packet, err := json.Marshal(Participant_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Participant_packet.M_header.M_sessionUID, 10) + ":4:" + strconv.Itoa(incrementing_participant_packet_number)), json_participant_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_participant_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_participant_packet_number += 1
		case 5:
			// If the packet we received is the car_setup_packet, read its binary into our car_setup_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Car_setup_packet); err != nil {
				fmt.Println("binary.Read car_setup_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:               5,
				Car_setup_packet: Car_setup_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_car_setup_packet, err := json.Marshal(Car_setup_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Car_setup_packet.M_header.M_sessionUID, 10) + ":5:" + strconv.Itoa(incrementing_car_setup_packet_number)), json_car_setup_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_car_setup_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_car_setup_packet_number += 1
		case 6:
			// If the packet we received is the telemetry_packet, read its binary into our telemetry_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Telemetry_packet); err != nil {
				fmt.Println("binary.Read telemetry_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:               6,
				Telemetry_packet: Telemetry_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_telemetry_packet, err := json.Marshal(Telemetry_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Telemetry_packet.M_header.M_sessionUID, 10) + ":6:" + strconv.Itoa(incrementing_telemetry_packet_number)), json_telemetry_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_telemetry_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_telemetry_packet_number += 1

			// Add data from telemetry to catchUp stuff
			//
			// speed
			if _, err := redis_conn.Do("LPUSH", "raceSpeed", Telemetry_packet.M_carTelemetryData[Telemetry_packet.M_header.M_playerCarIndex].M_speed); err != nil {
				fmt.Println("Lpush M_speed to reddis catchUp failed:", err)
			}
			if _, err := redis_conn.Do("LTRIM", "raceSpeed", 0, 1499); err != nil {
				fmt.Println("LTRIM M_speed to reddis catchUp failed:", err)
			}

			// engine rpms
			if _, err := redis_conn.Do("LPUSH", "engineRevs", Telemetry_packet.M_carTelemetryData[Telemetry_packet.M_header.M_playerCarIndex].M_engineRPM); err != nil {
				fmt.Println("Lpush M_engineRPM to reddis catchUp failed:", err)
			}
			if _, err := redis_conn.Do("LTRIM", "engineRevs", 0, 1499); err != nil {
				fmt.Println("LTRIM M_engineRPM to reddis catchUp failed:", err)
			}

			// gear
			if _, err := redis_conn.Do("LPUSH", "gearChanges", Telemetry_packet.M_carTelemetryData[Telemetry_packet.M_header.M_playerCarIndex].M_gear); err != nil {
				fmt.Println("Lpush M_gear to reddis catchUp failed:", err)
			}
			if _, err := redis_conn.Do("LTRIM", "gearChanges", 0, 1499); err != nil {
				fmt.Println("LTRIM M_gear to reddis catchUp failed:", err)
			}

			// Lpush to catchUp for both throttle and brake applications as separate
			// throttle
			if _, err := redis_conn.Do("LPUSH", "throttleApplication", Telemetry_packet.M_carTelemetryData[Telemetry_packet.M_header.M_playerCarIndex].M_throttle); err != nil {
				fmt.Println("Lpush M_throttle to reddis catchUp failed:", err)
			}
			if _, err := redis_conn.Do("LTRIM", "throttleApplication", 0, 1499); err != nil {
				fmt.Println("LTRIM M_throttle to reddis catchUp failed:", err)
			}

			// brake
			if _, err := redis_conn.Do("LPUSH", "brakeApplication", Telemetry_packet.M_carTelemetryData[Telemetry_packet.M_header.M_playerCarIndex].M_brake); err != nil {
				fmt.Println("Lpush M_brake to reddis catchUp failed:", err)
			}
			if _, err := redis_conn.Do("LTRIM", "brakeApplication", 0, 1499); err != nil {
				fmt.Println("LTRIM M_brake to reddis catchUp failed:", err)
			}

		case 7:
			// If the packet we received is the car_status_packet, read its binary into our car_status_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Car_status_packet); err != nil {
				fmt.Println("binary.Read car_status_packet failed:", err)
			}

			// Send the Udp_data struct containing the packet_id and the packet itself over the hub.bradcast channel to
			// be broadcasted to all connected clients
			hub.broadcast <- &Udp_data{
				Id:                7,
				Car_status_packet: Car_status_packet,
			}

			// Marshal the struct into json so we can save it in our redis database
			json_car_status_packet, err := json.Marshal(Car_status_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Car_status_packet.M_header.M_sessionUID, 10) + ":7:" + strconv.Itoa(incrementing_car_status_packet_number)), json_car_status_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_car_status_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_car_status_packet_number += 1
		default:
			continue
		}
	}
}
