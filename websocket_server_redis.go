//
//
// Websocket server written in Golang to serve our pages with udp_data from F1 2018
//
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
	// "unicode/utf8"

	"F1_GO/structs"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	// "github.com/gomodule/redigo"
)

// Set the location for the webpages to be served to.
// This case, localhost over port 8080
var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

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

// Listnes for the udp stream from the F1 2018 game and saves the data to a temporary redis database(Not written yet) and then
// sends the data over channels to our websocket handlers 
func f1_2018_udp_client() {
	// Print to the terminal window showing that the udp client is indeed up and working
	fmt.Println("F1 2018 UDP Client running...")
	for {
		buf := make([]byte, 1341)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("readfromudp error::: ", err)
		}

		// Set a new reader which we will use to cast into our structs.
		// This reader is for the header, which we determine what packet we have and what index our users car is in.
		// Bytes 3 in the udp packet will be the packet number and byte 20 will be the index of the users car.
		header_bytes_reader := bytes.NewReader([]byte{buf[3], buf[20]})
		packet_bytes_reader := bytes.NewReader(buf)

		// Read the binary of the udp packet header into our struct
		if err := binary.Read(header_bytes_reader, binary.LittleEndian, &Efficient_header); err != nil {
			fmt.Println("binary.Read header failed:", err)
		}

		// Depending on which packet we have, which we find by looking at header.M_packetId
		// We use a switch statement to then read the whole binary udp packet into its associated struct
		switch Efficient_header.M_packetId {
		case 0:
			// If the packet we received is a motion_packet, read its binary into our motion_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &motion_packet); err != nil {
				fmt.Println("binary.Read motion_packet failed:", err)
			}
			// Send the newly found motion packet over our motion packet channel so our websocket handlers can receive it and send it over our websocket
			motion_packet_channel <- motion_packet
			break
		case 1:
			// If the packet we received is the session_packet, read its binary into our session_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &session_packet); err != nil {
				fmt.Println("binary.Read session_packet failed:", err)
			}
			// Send the newly found session packet over our session packet channel so our websocket handlers can receive it and send it over our websocket
			session_packet_channel <- session_packet
			break
		case 2:
			// If the packet we received is the lap_packet, read its binary into our lap_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &lap_packet); err != nil {
				fmt.Println("binary.Read lap_packet failed:", err)
			}
			// Send the newly found lap packet over our lap packet channel so our websocket handlers can receive it and send it over our websocket
			lap_packet_channel <- lap_packet
			break
		case 3:
			// If the packet we received is the event_packet, read its binary into our event_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &event_packet); err != nil {
				fmt.Println("binary.Read event_packet failed:", err)
			}
			// Send the newly found event packet over our event packet channel so our websocket handlers can receive it and send it over our websocket
			event_packet_channel <- event_packet
			break
		case 4:
			// If the packet we received is the participant_packet, read its binary into our participant_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &participant_packet); err != nil {
				fmt.Println("binary.Read participant_packet failed:", err)
			}
			// Send the newly found participant packet over our participant packet channel so our websocket handlers can receive it and send it over our websocket
			participant_packet_channel <- participant_packet
			break
		case 5:
			// If the packet we received is the car_setup_packet, read its binary into our car_setup_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_setup_packet); err != nil {
				fmt.Println("binary.Read car_setup_packet failed:", err)
			}
			// Send the newly found car_setup packet over our car_setup packet channel so our websocket handlers can receive it and send it over our websocket
			car_setup_packet_channel <- car_setup_packet
			break
		case 6:
			// If the packet we received is the telemetry_packet, read its binary into our telemetry_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &telemetry_packet); err != nil {
				fmt.Println("binary.Read telemetry_packet failed:", err)
			}
			// Send the newly found telemetry packet over our telemetry packet channel so our websocket handlers can receive it and send it over our websocket
			telemetry_packet_channel <- telemetry_packet
			break
		case 7:
			// If the packet we received is the car_status_packet, read its binary into our car_status_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_status_packet); err != nil {
				fmt.Println("binary.Read car_status_packet failed:", err)
			}
			// Send the newly found car_status packet over our car_status packet channel so our websocket handlers can receive it and send it over our websocket
			car_status_packet_channel <- car_status_packet
			break
		default:
			break
		}
	}
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
	// Call our live_data_udp_client function in a gorutine, this will execute concurrently
	go live_data_udp_client(conn, sock)
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
	// Call our time_udp_client function in a gorutine, this will execute concurrently
	go time_udp_client(conn, sock)
}

// Function that is called by our websocket handlers
// This function gets packets from the F1 2018 game and casts the bytes into structs that are sent over our
// websocket to our front end than then displays them.
func live_data_udp_client(conn *websocket.Conn, sock *net.UDPConn) {
	// Defer the closing of our websocket connection. By doing this, when we get an error connecting to the websocket, which is the result of it being closed
	// by our client, we return which in turn will execute the conn.Close() which was defered from executing until the gorutine is over and retuned.
	defer conn.Close()

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
				log.Printf("Websocket error writing json_motion_packet: %s", err)
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
				log.Printf("Websocket error writing json_session_packet: %s", err)
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
				log.Printf("Websocket error writing json_lap_packet: %s", err)
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
				log.Printf("Websocket error writing json_event_packet: %s", err)
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
				log.Printf("Websocket error writing json_participant_packet: %s", err)
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
				log.Printf("Websocket error writing json_car_setup_packet: %s", err)
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
				log.Printf("Websocket error writing json_telemetry_packet: %s", err)
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
				log.Printf("Websocket error writing json_car_status_packet: %s", err)
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
func time_udp_client(conn *websocket.Conn, sock *net.UDPConn) {
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
		case 2:
			// If the packet we received is the lap_packet, read its binary into our lap_packet struct

			// If our users car index is not 0 then we need to grab the correct bytes that belong to our users car data
			player_car_bytes_multiplier := Efficient_header.M_playerCarIndex * 41

			// Since we only need a few of the items from the lap packet. Read the bytes for the items we do need into our user_car_structs
			// efficient struct. Below are the items we need in the format we grab the bytes below.
			// M_header
			// M_lastLapTime, M_currentLapTime M_bestLapTime, M_sector1Time, M_sector2Time
			// M_carPosition
			// M_currentLapNum
			// M_pitStatus
			// M_sector
			if err := binary.Read(bytes.NewReader(append(append(
				buf[0:21],
				buf[player_car_bytes_multiplier+21:player_car_bytes_multiplier+41]...),
				buf[(player_car_bytes_multiplier+53)],
				buf[(player_car_bytes_multiplier+54)],
				buf[(player_car_bytes_multiplier+55)],
				buf[(player_car_bytes_multiplier+56)])), binary.LittleEndian, &User_PacketLapData); err != nil {
				fmt.Println("binary.Read lap_packet failed:", err)
			}

			// Convert out struct into JSON format
			json_lap_packet, err := json.Marshal(User_PacketLapData)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := conn.WriteMessage(websocket.TextMessage, json_lap_packet); err != nil {
				log.Printf("Websocket error writing lap_packet: %s", err)
				return
			}
			break
		}
	}
}
