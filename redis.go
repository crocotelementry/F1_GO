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
	atm_motion_packet                      = make(chan structs.PacketMotionData)
	atm_session_packet                     = make(chan structs.PacketSessionData)
	atm_lap_packet                         = make(chan structs.PacketLapData)
	atm_event_packet                       = make(chan structs.PacketEventData)
	atm_participant_packet                 = make(chan structs.PacketParticipantsData)
	atm_car_setup_packet                   = make(chan structs.PacketCarSetupData)
	atm_telemetry_packet                   = make(chan structs.PacketCarTelemetryData)
	atm_car_status_packet                  = make(chan structs.PacketCarStatusData)
	atm_race_event_directory               = make(chan structs.RaceEventDirectory)
	redis_done                             = make(chan bool)
)

// Client is a middleman between the websocket connection and the hub.
type Udp_data struct {
	Id int

	Motion_packet          structs.PacketMotionData
	Session_packet         structs.PacketSessionData
	Lap_packet             structs.PacketLapData
	Event_packet           structs.PacketEventData
	Participant_packet     structs.PacketParticipantsData
	Car_setup_packet       structs.PacketCarSetupData
	Telemetry_packet       structs.PacketCarTelemetryData
	Car_status_packet      structs.PacketCarStatusData
	Save_to_database_alert structs.Save_to_database_alerts
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
func getRedisDataForMysql(chosen_session_uid uint64) {
	// get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()

	// Defer the closing of the redis connection until we return at the end of getRedisData
	defer redis_conn.Close()

	// Start up add_to_longterm_storage in another goroutine
	go add_to_longterm_storage()

	log.Println("1")

	session_uid := strconv.FormatUint(chosen_session_uid, 10)

	log.Println("2")

	// Send over the initial data for race_event_directory
	race_event_directory_data, err := (redis_conn.Do("GET", session_uid+":0:0"))
	if err != nil {
		log.Println("Getting initial data for race_event_directory from redis database failed:", err)
	}

	log.Println("race_event_directory_data", race_event_directory_data)

	err = json.Unmarshal(race_event_directory_data.([]byte), &Motion_packet)
	if err != nil {
		log.Println(err)
	}

	log.Println("4")
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

	for packet_type := 0; packet_type < 8; packet_type += 1 {
		max_packet_number, err := redis.Int(redis_conn.Do("GET", session_uid+":"+strconv.Itoa(packet_type)+":Incrementing_packet_number"))
		if err != nil {
			log.Println("Getting max_packet_number for packet id number", packet_type, "from redis database failed:", err)
		}

		for packet_number := 0; packet_number < max_packet_number; packet_number += 1 {
			packet_data, err := redis_conn.Do("GET", session_uid+":"+strconv.Itoa(packet_type)+":"+strconv.Itoa(packet_number))
			if err != nil {
				log.Println("Getting packet_data for packet id number", packet_type, "with packet_number", packet_number, "from redis database failed:", err)
			}

			switch packet_type {
			case 0:
				err := json.Unmarshal(packet_data.([]byte), &Motion_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_motion_packet <- Motion_packet

			case 1:
				err := json.Unmarshal(packet_data.([]byte), &Session_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_session_packet <- Session_packet

			case 2:
				err := json.Unmarshal(packet_data.([]byte), &Lap_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_lap_packet <- Lap_packet

			case 3:
				err := json.Unmarshal(packet_data.([]byte), &Event_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_event_packet <- Event_packet

			case 4:
				err := json.Unmarshal(packet_data.([]byte), &Participant_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_participant_packet <- Participant_packet

			case 5:
				err := json.Unmarshal(packet_data.([]byte), &Car_setup_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_car_setup_packet <- Car_setup_packet

			case 6:
				err := json.Unmarshal(packet_data.([]byte), &Telemetry_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_telemetry_packet <- Telemetry_packet

			case 7:
				err := json.Unmarshal(packet_data.([]byte), &Car_status_packet)
				if err != nil {
					fmt.Println(err)
				}
				atm_car_status_packet <- Car_status_packet

			}
		}
	}

	redis_done <- true
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

	// Redis database format:
	// Session_uid:packet_id:incrementing_packet_number									This is for packets
	// Session_uid:"Incrementing_packet_number"								This is for knowing the max value of incrementing_packet_number
	// session_UIDs																						This is a list for keeping track of what session_UIDs are in the redis database

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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":0:" + strconv.Itoa(incrementing_motion_packet_number)), json_motion_packet); err != nil {
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":1:" + strconv.Itoa(incrementing_session_packet_number)), json_session_packet); err != nil {
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":2:" + strconv.Itoa(incrementing_lap_packet_number)), json_lap_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_lap_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_lap_packet_number += 1
		case 3:
			// If the packet we received is the event_packet, read its binary into our event_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Event_packet); err != nil {
				fmt.Println("binary.Read event_packet failed:", err)
			}

			if Equal(Event_packet.M_eventStringCode, session_start_code) {
				number_of_sessions_exists_integer_reply, err := redis_conn.Do("EXISTS", "number_of_sessions")

				if err != nil {
					fmt.Println("Checking if number_of_sessions exists failed:", err)
				}

				if number_of_sessions_exists_integer_reply == int64(1) {
					if _, err := redis_conn.Do("INCR", "number_of_sessions"); err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}
				} else {
					if _, err := redis_conn.Do("SET", "number_of_sessions", "1"); err != nil {
						fmt.Println("Setting number_of_sessions to 1 failed:", err)
					}
				}

				session_UIDs_exists_integer_reply, err := redis_conn.Do("EXISTS", "session_UIDs")
				if err != nil {
					fmt.Println("Checking if session_UIDs exists failed:", err)
				}

				if session_UIDs_exists_integer_reply == int64(1) {
					session_UIDs_SADD_integer_reply, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID))
					if err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}

					fmt.Println("sadd session uid:", session_UIDs_SADD_integer_reply)

					if session_UIDs_SADD_integer_reply == int64(0) {
						fmt.Println("\nSession with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", Event_packet.M_header.M_sessionUID, "\n")
					}
					// else {
					// 	incrementing_packet_number = 0
					// }
				} else {
					if _, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID)); err != nil {
						fmt.Println("Setting number_of_sessions to 1 failed:", err)
					}
				}

				// Add session start time to redis database
				// Format is as follows:
				// Session_UID:session_start_time
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":session_start_time"), &structs.Session_start{time.Now()}); err != nil {
					fmt.Println("Setting session_start_time to failed:", err)
				}
			}

			// If we receive a session end code, send a data alert over the websocket to ask the user if they want to save the session for long
			// term use in a MYSQL database, discard the session, or hold on to it until our redis database reaches its size limit or its amount
			// of sessions limit.
			if Equal(Event_packet.M_eventStringCode, session_end_code) {

				// fmt.Println(strconv.FormatUint(Event_packet.M_header.M_sessionUID, 10) + ":Incrementing_packet_number")

				sue_return := session_UID_exists(Event_packet.M_header.M_sessionUID)

				// fmt.Println(sue_return)

				if sue_return == false {
					_, err := redis_conn.Do("SADD", "session_UIDs", (Event_packet.M_header.M_sessionUID))
					if err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}
				}

				// Add session end time to redis database
				// Format is as follows:
				// Session_UID:session_end_time
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":session_end_time"), &structs.Session_end{time.Now()}); err != nil {
					fmt.Println("Setting session_end_time to failed:", err)
				}

				// Session_uid:Incrementing_packet_number
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":Incrementing_packet_number"), strconv.Itoa(int(incrementing_packet_number))); err != nil {
					log.Println("             ", "Setting Incrementing_packet_number failed:", err)
				}

				// set incrementing_motion_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":0:Incrementing_packet_number"), incrementing_motion_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_motion_packet_number failed:", err)
				}

				// set incrementing_session_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":1:Incrementing_packet_number"), incrementing_session_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_session_packet_number failed:", err)
				}

				// set incrementing_lap_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":2:Incrementing_packet_number"), incrementing_lap_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_lap_packet_number failed:", err)
				}

				// set incrementing_event_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":3:Incrementing_packet_number"), incrementing_event_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_event_packet_number failed:", err)
				}

				// set incrementing_participant_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":4:Incrementing_packet_number"), incrementing_participant_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_participant_packet_number failed:", err)
				}

				// set incrementing_car_setup_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":5:Incrementing_packet_number"), incrementing_car_setup_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_car_setup_packet_number failed:", err)
				}

				// set incrementing_telemetry_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":6:Incrementing_packet_number"), incrementing_telemetry_packet_number); err != nil {
					log.Println("             ", "Setting incrementing_telemetry_packet_number failed:", err)
				}

				// set incrementing_car_status_packet_number for motion packets and its session_UIDs
				if _, err := redis_conn.Do("SET", (strconv.FormatUint(header.M_sessionUID, 10) + ":7:Incrementing_packet_number"), incrementing_car_status_packet_number); err != nil {
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":4:" + strconv.Itoa(incrementing_participant_packet_number)), json_participant_packet); err != nil {
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":5:" + strconv.Itoa(incrementing_car_setup_packet_number)), json_car_setup_packet); err != nil {
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":6:" + strconv.Itoa(incrementing_telemetry_packet_number)), json_telemetry_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_telemetry_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_telemetry_packet_number += 1
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

			if _, err := redis_conn.Do("SET", (session_uid_Prefix + ":7:" + strconv.Itoa(incrementing_car_status_packet_number)), json_car_status_packet); err != nil {
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
