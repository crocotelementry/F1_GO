//
//
// Websocket server written in Golang to serve our pages with udp_data from F1 2018
//
//
// Author: Kristian Nilssen, Seattle

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"F1_GO/structs"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Set the location for the webpages to be served to.
// This case, localhost over port 8080
var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

// Create a struct for our connection
type Conn_struct struct {
	Key        int
	Connection *websocket.Conn
}

// Set up our structs including the header struct so we are able to determine which
// udp packet is incoming and we can deal with it accordingly
var header structs.PacketHeader
var motion_packet structs.PacketMotionData
var session_packet structs.PacketSessionData
var lap_packet structs.PacketLapData
var event_packet structs.PacketEventData
var participant_packet structs.PacketParticipantsData
var car_setup_packet structs.PacketCarSetupData
var telemetry_packet structs.PacketCarTelemetryData
var car_status_packet structs.PacketCarStatusData

// Set up our efficient user oriented structs
var Efficient_header structs.Efficient_header
var User_PacketSessionData structs.User_PacketSessionData
var User_PacketLapData structs.User_PacketLapData
var User_PacketCarTelemetryData structs.User_PacketCarTelemetryData
var User_PacketCarStatusData structs.User_PacketCarStatusData

// Create a socket listing for udp packets on port 5003
// Since this is not on any specific ip address, if closed we can reopen on the same port again.
// To get over this hurdle, we brought the socket out to only declare once, and can be called by
// each of our websocket handlers to grab packets from the game
var addrs, _ = net.ResolveUDPAddr("udp", ":20777")
var sock, err = net.ListenUDP("udp", addrs)

// Create the channels for each of the packet types
var motion_packet_channel = make(chan structs.PacketMotionData)
var session_packet_channel = make(chan structs.PacketSessionData)
var lap_packet_channel = make(chan structs.PacketLapData)
var event_packet_channel = make(chan structs.PacketEventData)
var participant_packet_channel = make(chan structs.PacketParticipantsData)
var car_setup_packet_channel = make(chan structs.PacketCarSetupData)
var telemetry_packet_channel = make(chan structs.PacketCarTelemetryData)
var car_status_packet_channel = make(chan structs.PacketCarStatusData)

// Create two variables, one for the code when a session starts and one for when a session ends
var session_start_code = [4]uint8{83, 69, 78, 68}
var session_end_code = [4]uint8{83, 83, 84, 65}

// Create a channel for conn_checker_sender and conn_delete_reciever
var conn_add_sender = make(chan Conn_struct, 10)
var conn_delete_sender = make(chan int, 10)

// Keep track of the total number of connections
var num_of_conn = 0

func main() {
	// Now run our f1 2018 udp telemetry packets and save them to a temporary database
	// Do this in a Go routine
	go f1_2018_udp_client()

	// Sets the location in which to serve our static files from for our webpage
	var dir string
	flag.StringVar(&dir, "dir", "./web", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	// Creates a new mux
	router := mux.NewRouter()
	// WHen our html page calls its static files from /static/file, this sets the location to grab them from
	// TODO: Is this all this code needed? Couldn't we just set it web/?
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir(dir))))

	// Our handler functions for each page
	// Landing page /aka live telemetry or telemetry_dashboard
	router.HandleFunc("/", liveHandler)
	router.HandleFunc("/ws", live_wsHandler)

	// History page /aka history_dashboard
	router.HandleFunc("/history", historyHandler)
	router.HandleFunc("/history/ws", history_wsHandler)

	// Live time page /aka time_dashboard
	router.HandleFunc("/time", timeHandler)
	router.HandleFunc("/time/ws", time_wsHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

//
func f1_2018_udp_client() {
	// Print to the terminal window showing that the udp client is indeed up and working
	fmt.Println("F1 2018 UDP Client running...\n")

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

	// Create a map (dictionary) that will keep track of our conn's
	conn_map := make(map[int]*websocket.Conn)

	// Create a variable that will tell us if our conn_map has connections open
	conn_in_conn_map := false

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

		// fmt.Println("running")

		// How we know if we have an open conn or not
		// This helps us with not getting locked on our channel outputs to our connections
		select {
		case conn_add_reciever := <-conn_add_sender:
			conn_map[conn_add_reciever.Key] = conn_add_reciever.Connection
			conn_in_conn_map = true
			fmt.Println("Connection opened!")
		case conn_delete_reciever := <-conn_delete_sender:
			delete(conn_map, conn_delete_reciever)
			fmt.Println("Connection closed!")
			if len(conn_map) == 0 {
				conn_in_conn_map = false
			}
		default:
			break
		}

		switch header.M_packetId {
		case 0:
			// If the packet we received is a motion_packet, read its binary into our motion_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &motion_packet); err != nil {
				fmt.Println("binary.Read motion_packet failed:", err)
			}
			// Send the newly found motion packet over our motion packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				motion_packet_channel <- motion_packet
			}

			// Marshal the struct into json so we can save it in our redis database
			json_motion_packet, err := json.Marshal(motion_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_motion_packet); err != nil {
				fmt.Println("Setting packet json_motion_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 1:
			// If the packet we received is the session_packet, read its binary into our session_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &session_packet); err != nil {
				fmt.Println("binary.Read session_packet failed:", err)
			}
			// Send the newly found session packet over our session packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				session_packet_channel <- session_packet
			}

			// Marshal the struct into json so we can save it in our redis database
			json_session_packet, err := json.Marshal(session_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_session_packet); err != nil {
				fmt.Println("Setting packet json_session_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 2:
			// If the packet we received is the lap_packet, read its binary into our lap_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &lap_packet); err != nil {
				fmt.Println("binary.Read lap_packet failed:", err)
			}
			// Send the newly found lap packet over our lap packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				lap_packet_channel <- lap_packet
			}

			json_lap_packet, err := json.Marshal(lap_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_lap_packet); err != nil {
				fmt.Println("Setting packet json_lap_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 3:
			// If the packet we received is the event_packet, read its binary into our event_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &event_packet); err != nil {
				fmt.Println("binary.Read event_packet failed:", err)
			}

			if Equal(event_packet.M_eventStringCode, session_end_code) {
				fmt.Println("Session end code recieved, total number of packets recived and SET in redis database is:", num_redis_set)
				fmt.Println("EVENT PACKET", event_packet, "\n")
			}

			if Equal(event_packet.M_eventStringCode, session_start_code) {
				number_of_sessions_exists_integer_reply, err := redis_conn.Do("EXISTS", "number_of_sessions")

				if err != nil {
					fmt.Println("Checking if number_of_sessions exists failed:", err)
				}
				fmt.Println("number_of_sessions_exists_integer_reply:", number_of_sessions_exists_integer_reply)

				if number_of_sessions_exists_integer_reply == int64(1) {
					// fmt.Println("1")
					// If number_of_sessions exists
					if _, err := redis_conn.Do("INCR", "number_of_sessions"); err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}
				} else {
					// fmt.Println("2")
					// If number_of_sessions doesnt exist
					if _, err := redis_conn.Do("SET", "number_of_sessions", "1"); err != nil {
						fmt.Println("Setting number_of_sessions to 1 failed:", err)
					}
				}

				session_UIDs_exists_integer_reply, err := redis_conn.Do("EXISTS", "session_UIDs")
				if err != nil {
					fmt.Println("Checking if session_UIDs exists failed:", err)
				}
				fmt.Println("session_UIDs_exists_integer_reply:", session_UIDs_exists_integer_reply)
				if session_UIDs_exists_integer_reply == int64(1) {
					// fmt.Println("3")
					// If number_of_sessions exists
					session_UIDs_SADD_integer_reply, err := redis_conn.Do("SADD", "session_UIDs", (event_packet.M_header.M_sessionUID))
					if err != nil {
						fmt.Println("Incrementing number_of_sessions by 1 failed:", err)
					}

					fmt.Println("session_UIDs_SADD_integer_reply:", session_UIDs_SADD_integer_reply)

					if session_UIDs_SADD_integer_reply == int64(0) {
						// fmt.Println("4")
						fmt.Println("Session with the following UID is already added to redis database,\nreceived session start code but session UID did not change from previous session UID:", event_packet.M_header.M_sessionUID)
					} else {
						// fmt.Println("5")
						num_redis_set = 0
					}
				} else {
					// fmt.Println("6")
					// If number_of_sessions doesnt exist
					if _, err := redis_conn.Do("SET", "session_UIDs", (event_packet.M_header.M_sessionUID)); err != nil {
						fmt.Println("Setting number_of_sessions to 1 failed:", err)
					}
				}

				fmt.Println("EVENT PACKET", event_packet)

			}

			// Send the newly found event packet over our event packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				event_packet_channel <- event_packet
			}

			break
		case 4:
			// If the packet we received is the participant_packet, read its binary into our participant_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &participant_packet); err != nil {
				fmt.Println("binary.Read participant_packet failed:", err)
			}
			// Send the newly found participant packet over our participant packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				participant_packet_channel <- participant_packet
			}

			json_participant_packet, err := json.Marshal(participant_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_participant_packet); err != nil {
				fmt.Println("Setting packet json_participant_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 5:
			// If the packet we received is the car_setup_packet, read its binary into our car_setup_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_setup_packet); err != nil {
				fmt.Println("binary.Read car_setup_packet failed:", err)
			}
			// Send the newly found car_setup packet over our car_setup packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				car_setup_packet_channel <- car_setup_packet
			}

			json_car_setup_packet, err := json.Marshal(car_setup_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_car_setup_packet); err != nil {
				fmt.Println("Setting packet json_car_setup_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 6:
			// If the packet we received is the telemetry_packet, read its binary into our telemetry_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &telemetry_packet); err != nil {
				fmt.Println("binary.Read telemetry_packet failed:", err)
			}
			// Send the newly found telemetry packet over our telemetry packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				telemetry_packet_channel <- telemetry_packet
			}

			json_telemetry_packet, err := json.Marshal(telemetry_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_telemetry_packet); err != nil {
				fmt.Println("Setting packet json_telemetry_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		case 7:
			// If the packet we received is the car_status_packet, read its binary into our car_status_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_status_packet); err != nil {
				fmt.Println("binary.Read car_status_packet failed:", err)
			}
			// Send the newly found car_status packet over our car_status packet channel so our websocket handlers can receive it and send it over our websocket
			// If we have no connections made, we need to skip this step since we will be blocking on the channel until one of our cennections can receive it
			if conn_in_conn_map {
				car_status_packet_channel <- car_status_packet
			}

			json_car_status_packet, err := json.Marshal(car_status_packet)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(objectPrefix+strconv.Itoa(int(header.M_packetId)))

			if _, err := redis_conn.Do("SET", objectPrefix+strconv.Itoa(int(header.M_packetId)), json_car_status_packet); err != nil {
				fmt.Println("Setting packet json_car_status_packet failed:", err)
				num_redis_set -= 1
			}

			num_redis_set += 1

			break
		default:
			break
		}
	}
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
	fmt.Println("Successful connection made with Redis database!\n")
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

// liveHandler is called when our browser goes to the page localhost:8080, this serves up our html file along
// with its corresponding javascript and css files
func liveHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/live_dashboard.html")
}

// Called when at the page localhost:8080/history
func historyHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/history_dashboard.html")
}

// Called when at the page localhost:8080/time
func timeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/time_dashboard.html")
}

// Live data Websocket handler, when our javascriot file, which is served along with our html file from our
// root handler, runs it makes a websocket connection to ws://localhost:8080/ws which triggers our hsHandler with the trailing /ws.
// We then start our live_data_udp_client that listens for and formats the incoming UDP data from F1 2018 then writes the packets to our websocket
func live_wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	num_of_conn += 1

	this_conn_data := Conn_struct{Key: num_of_conn, Connection: conn}

	conn_add_sender <- this_conn_data

	// Call our live_data_udp_client function in a gorutine, this will execute concurrently
	go live_data_udp_client(conn, sock, this_conn_data.Key)
}

// History Websocket handler, when our javascriot file, which is served along with our html file from our
// root handler, runs it makes a websocket connection to ws://localhost:8080/history/ws which triggers our hsHandler with the trailing /ws.
// We then start our history_udp_client that listens for and formats the incoming UDP data from F1 2018 then writes the packets to our websocket
func history_wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	// Call our history_udp_client function in a gorutine, this will execute concurrently
	go history_udp_client(conn, sock)
}

// Live time Websocket handler, when our javascriot file, which is served along with our html file from our
// root handler, runs it makes a websocket connection to ws://localhost:8080/time/ws which triggers our hsHandler with the trailing /ws.
// We then start our time_udp_client that listens for and formats the incoming UDP data from F1 2018 then writes the packets to our websocket
func time_wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	num_of_conn += 1

	this_conn_data := Conn_struct{Key: num_of_conn, Connection: conn}

	conn_add_sender <- this_conn_data

	// Call our time_udp_client function in a gorutine, this will execute concurrently
	go time_udp_client(conn, sock, this_conn_data.Key)
}

// Function that is called by our websocket handlers
// This function gets packets from the F1 2018 game and casts the bytes into structs that are sent over our
// websocket to our front end than then displays them.
func live_data_udp_client(conn *websocket.Conn, sock *net.UDPConn, conn_number int) {
	// Defer the closing of our websocket connection. By doing this, when we get an error connecting to the websocket, which is the result of it being closed
	// by our client, we return which in turn will execute the conn.Close() which was defered from executing until the gorutine is over and retuned.
	defer func() {
		conn.Close()
		fmt.Println("conn closed live_dashboard print with conn_number:", conn_number)
		conn_delete_sender <- conn_number
	}()

	for {
		select {
		case motion_packet_channel_msg := <-motion_packet_channel:
			// fmt.Println("motion_packet_channel_msg recieved!")

			json_motion_packet, err := json.Marshal(motion_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_motion_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_motion_packet: %s", err)
				return
			}
			break
		case session_packet_channel_msg := <-session_packet_channel:
			// fmt.Println("session_packet_channel_msg recieved!")

			json_session_packet, err := json.Marshal(session_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_session_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_session_packet: %s", err)
				return
			}
			break
		case lap_packet_channel_msg := <-lap_packet_channel:
			// fmt.Println("lap_packet_channel_msg recieved!")

			json_lap_packet, err := json.Marshal(lap_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_lap_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_lap_packet: %s", err)
				return
			}
			break
		case event_packet_channel_msg := <-event_packet_channel:
			// fmt.Println("event_packet_channel_msg recieved!")

			json_event_packet, err := json.Marshal(event_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_event_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_event_packet: %s", err)
				return
			}
			break
		case participant_packet_channel_msg := <-participant_packet_channel:
			// fmt.Println("participant_packet_channel_msg recieved!")
			// fmt.Println(participant_packet_channel_msg)

			// for _, element := range participant_packet_channel_msg.M_participants {
			//   fmt.Println(string(element.M_name[:]))
			// }

			json_participant_packet, err := json.Marshal(participant_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_participant_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_participant_packet: %s", err)
				return
			}
			break
		case car_setup_packet_channel_msg := <-car_setup_packet_channel:
			// fmt.Println("car_setup_packet_channel_msg recieved!")

			json_car_setup_packet, err := json.Marshal(car_setup_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_car_setup_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_car_setup_packet: %s", err)
				return
			}
			break
		case telemetry_packet_channel_msg := <-telemetry_packet_channel:
			// fmt.Println("telemetry_packet_channel_msg recieved!")

			json_telemetry_packet, err := json.Marshal(telemetry_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_telemetry_packet); err != nil {
				log.Printf("Websocket ", conn_number, " error writing json_telemetry_packet: %s", err)
				return
			}
			break
		case car_status_packet_channel_msg := <-car_status_packet_channel:
			// fmt.Println("car_status_packet_channel_msg recieved!")

			json_car_status_packet, err := json.Marshal(car_status_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_car_status_packet); err != nil {
				log.Printf("Websocket ", conn_number, " error writing json_car_status_packet: %s", err)
				return
			}
			break
		default:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Ping failed on conn:", conn_number)
				return
			}
			break
		}
	}

}

// Function that is called by our websocket handlers
// This function gets packets from the F1 2018 game and casts the bytes into structs that are sent over our
// websocket to our front end than then displays them.
func history_udp_client(conn *websocket.Conn, sock *net.UDPConn) {
	// Defer the closing of our websocket connection. By doing this, when we get an error connecting to the websocket, which is the result of it being closed
	// by our client, we return which in turn will execute the conn.Close() which was defered from executing until the gorutine is over and retuned.
	defer conn.Close()

	for {
		// Create a buffer to read the incoming udp packets
		// Read the udp packets and if we get an error while reading, print out the error
		buf := make([]byte, 1341)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}

		// Set a new reader which we will use to cast into our structs.
		// This reader is for the header, which we determine what packet we have and what index our users car is in.
		// Bytes 3 in the udp packet will be the packet number and byte 20 will be the index of the users car.
		header_bytes_reader := bytes.NewReader([]byte{buf[3], buf[20]})

		// Read the binary of the udp packet header into our struct
		if err := binary.Read(header_bytes_reader, binary.LittleEndian, &Efficient_header); err != nil {
			fmt.Println("binary.Read header failed:", err)
		}

		// Depending on which packet we have, which we find by looking at header.M_packetId
		// We use a switch statement to then read the whole binary udp packet into its associated struct
		switch Efficient_header.M_packetId {
		// case 1:
		//     // If the packet we received is the session_packet, read its binary into our session_packet struct
		//     if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &session_packet); err != nil {
		// 		    fmt.Println("binary.Read session_packet failed:", err)
		//   	}
		// 		// Convert out struct into JSON format
		//     json_session_packet, err := json.Marshal(session_packet)
		//     if err != nil {
		//       fmt.Println(err)
		//     }
		// 		// Write our JSON formatted F1 UDP packet struct to our websocket
		//     if err := conn.WriteMessage(websocket.TextMessage, json_session_packet); err != nil {
		// 			log.Printf("Websocket error writing session_packet: %s", err)
		// 			return
		// 		}
		// 		break
		default:
			break
		}
	}
}

// Function that is called by our websocket handlers
// This function gets packets from the F1 2018 game and casts the bytes into structs that are sent over our
// websocket to our front end than then displays them.
func time_udp_client(conn *websocket.Conn, sock *net.UDPConn, conn_number int) {
	// Defer the closing of our websocket connection. By doing this, when we get an error connecting to the websocket, which is the result of it being closed
	// by our client, we return which in turn will execute the conn.Close() which was defered from executing until the gorutine is over and retuned.
	defer func() {
		conn.Close()
		fmt.Println("conn closed time print with conn_number:", conn_number)
		conn_delete_sender <- conn_number
	}()

	for {
		select {
		case motion_packet_channel_msg := <-motion_packet_channel:
			// fmt.Println("motion_packet_channel_msg recieved!")

			json_motion_packet, err := json.Marshal(motion_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_motion_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_motion_packet: %s", err)
				return
			}
			break
		case session_packet_channel_msg := <-session_packet_channel:
			// fmt.Println("session_packet_channel_msg recieved!")

			json_session_packet, err := json.Marshal(session_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_session_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_session_packet: %s", err)
				return
			}
			break
		case lap_packet_channel_msg := <-lap_packet_channel:
			// fmt.Println("lap_packet_channel_msg recieved!")

			json_lap_packet, err := json.Marshal(lap_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_lap_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_lap_packet: %s", err)
				return
			}
			break
		case event_packet_channel_msg := <-event_packet_channel:
			// fmt.Println("event_packet_channel_msg recieved!")

			json_event_packet, err := json.Marshal(event_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_event_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_event_packet: %s", err)
				return
			}
			break
		case participant_packet_channel_msg := <-participant_packet_channel:
			// fmt.Println("participant_packet_channel_msg recieved!")
			// fmt.Println(participant_packet_channel_msg)

			// for _, element := range participant_packet_channel_msg.M_participants {
			//   fmt.Println(string(element.M_name[:]))
			// }

			json_participant_packet, err := json.Marshal(participant_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_participant_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_participant_packet: %s", err)
				return
			}
			break
		case car_setup_packet_channel_msg := <-car_setup_packet_channel:
			// fmt.Println("car_setup_packet_channel_msg recieved!")

			json_car_setup_packet, err := json.Marshal(car_setup_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_car_setup_packet); err != nil {
				log.Printf("Websocket", conn_number, " error writing json_car_setup_packet: %s", err)
				return
			}
			break
		case telemetry_packet_channel_msg := <-telemetry_packet_channel:
			// fmt.Println("telemetry_packet_channel_msg recieved!")

			json_telemetry_packet, err := json.Marshal(telemetry_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_telemetry_packet); err != nil {
				log.Printf("Websocket ", conn_number, " error writing json_telemetry_packet: %s", err)
				return
			}
			break
		case car_status_packet_channel_msg := <-car_status_packet_channel:
			// fmt.Println("car_status_packet_channel_msg recieved!")

			json_car_status_packet, err := json.Marshal(car_status_packet_channel_msg)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_car_status_packet); err != nil {
				log.Printf("Websocket ", conn_number, " error writing json_car_status_packet: %s", err)
				return
			}
			break
		default:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Ping failed on conn:", conn_number)
				return
			}
			break
		}
	}
}
