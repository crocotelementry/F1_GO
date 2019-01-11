package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/crocotelementry/F1_GO/structs"
	"github.com/gomodule/redigo/redis"
)

var (
	header             structs.PacketHeader
	Motion_packet      structs.PacketMotionData
	Session_packet     structs.PacketSessionData
	Lap_packet         structs.PacketLapData
	Event_packet       structs.PacketEventData
	Participant_packet structs.PacketParticipantsData
	Car_setup_packet   structs.PacketCarSetupData
	Telemetry_packet   structs.PacketCarTelemetryData
	Car_status_packet  structs.PacketCarStatusData
	num_redis_set      = 0
	session_start_code = [4]uint8{83, 69, 78, 68}
	session_end_code   = [4]uint8{83, 83, 84, 65}
	redis_ping_done    = make(chan bool)
)

// Client is a middleman between the websocket connection and the hub.
type Udp_data struct {
	Id int

	Motion_packet      structs.PacketMotionData
	Session_packet     structs.PacketSessionData
	Lap_packet         structs.PacketLapData
	Event_packet       structs.PacketEventData
	Participant_packet structs.PacketParticipantsData
	Car_setup_packet   structs.PacketCarSetupData
	Telemetry_packet   structs.PacketCarTelemetryData
	Car_status_packet  structs.PacketCarStatusData
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

	fmt.Println("PING Response = ", s)

	if s == "PONG" {
		redis_ping_done <- true
	} else {
		redis_ping_done <- false
	}

	// Output: PONG
	return nil
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

func getGameData(hub *Hub) {
	// newPool returns a pointer to a redis.Pool
	redis_pool := newPool()
	// get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()
	// use defer to close the connection when the function completes
	defer redis_conn.Close()

	// call Redis PING command to test connectivity
	err := ping(redis_conn)
	if err != nil {
		fmt.Println("Problem with connection to Redis database", err)
	}

	// Set number of SETs to redis database to zero
	num_redis_set := 0

	// Redis database format:
	// Session_uid:Frame_identifier:Packet_id

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

		objectPrefix := strconv.Itoa(int(header.M_sessionUID)) + ":" + strconv.Itoa(int(header.M_frameIdentifier)) + ":"

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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_motion_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_session_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_lap_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

					if session_UIDs_SADD_integer_reply == int64(0) {
						fmt.Println("\nSession with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", Event_packet.M_header.M_sessionUID, "\n")
					} else {
						num_redis_set = 0
					}
				} else {
					if _, err := redis_conn.Do("SET", "session_UIDs", (Event_packet.M_header.M_sessionUID)); err != nil {
						fmt.Println("Setting number_of_sessions to 1 failed:", err)
					}
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_participant_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_car_setup_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_telemetry_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
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

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_car_status_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				num_redis_set -= 1
			}
			num_redis_set += 1
		default:
			continue
		}
	}
}
