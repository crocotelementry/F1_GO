package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	// "bytes"
	// "encoding/binary"

	"github.com/crocotelementry/F1_GO/structs"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1341
)

var (
	newline      = []byte{'\n'}
	space        = []byte{' '}
	message_json structs.Save_to_database_websocket_recive
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
	// ReadBufferSize:  1341,
	// WriteBufferSize: 1341,
}

// // Struct that is sent over websocket to alert of new data to be saved or not saved to longterm storage
// type Save_to_database_alerts struct{
// 	date string
// 	length int
// }

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// What type of client is it
	// Live, Time, History
	conn_type string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// send chan *Udp_data
	Motion_packet_send      chan structs.PacketMotionData
	Session_packet_send     chan structs.PacketSessionData
	Lap_packet_send         chan structs.PacketLapData
	Event_packet_send       chan structs.PacketEventData
	Participant_packet_send chan structs.PacketParticipantsData
	Car_setup_packet_send   chan structs.PacketCarSetupData
	Telemetry_packet_send   chan structs.PacketCarTelemetryData
	Car_status_packet_send  chan structs.PacketCarStatusData
	Save_to_database_alert  chan structs.Save_to_database_alerts
	Save_to_database_status chan structs.Save_to_database_status
}

// readFromClients reads messages from the websocket connection to the hub.
//
// The application runs readFromClients in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readFromClients() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		err = json.Unmarshal([]byte(message), &message_json)
		// log.Println("Add to database message recieved message json:", message_json)

		switch message_json.Type {
		case "add":
			log.Println("", c.conn.RemoteAddr(), " ", "Session chosen for long term storage in mysql with UID:", message_json.Uid)
			go getRedisDataForMysql(c.hub, message_json.Uid)
		default:
			log.Println("Incorrect statement recieved from websocket client:", message_json.Type)
		}
	}
}

// writeDashboard pumps messages from the hub to the websocket connection.
//
// A goroutine running writeDashboard is started for each connection to dashboard. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeDashboard() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("", c.conn.RemoteAddr(), " ", "Stopping dashboard clients ticker, unregistering, and closing connection")
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Session_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Session_packet_send")
				return
			}
			// Marshal our message into json so we can send it over the websockets
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Session_packet_send to dashboard websocket")
				return
			}

		case message, ok := <-c.Lap_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Lap_packet_send")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Lap_packet_send to dashboard websocket")
				return
			}

		case message, ok := <-c.Telemetry_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Telemetry_packet_send")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Telemetry_packet_send to dashboard websocket")
				return
			}

		case message, ok := <-c.Car_status_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Car_status_packet_send")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Car_status_packet_send to dashboard websocket")
				return
			}

		case message, ok := <-c.Save_to_database_alert:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_alert")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_alert to dashboard websocket")
				return
			}

		case message, ok := <-c.Save_to_database_status:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_status")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_status to time websocket")
				return
			}

		case <-ticker.C: // If our ticker has reached its time
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline

			// If our client has disconected from the websocket on thier end, close the client and its connection by returning and executing our defer statement
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with ticker pingmessage with dashboard client")
				return
			}
		}
	}
}

// writeHistory pumps messages from the websocket connection to the hub.
//
// A goroutine running writeHistory is started for each connection to history. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeHistory() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("", c.conn.RemoteAddr(), " ", "Stopping notlive clients ticker, unregistering, and closing connection")
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// _, message, err := c.conn.ReadMessage()
		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Printf("error: %v", err)
		// 	}
		// 	break
		// }
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Fetch that data it wants from database and send it over the ws
		select {
		case message, ok := <-c.Save_to_database_alert:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_alert")
				return
			}

			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_alert to notlive websocket")
				return
			}

		case message, ok := <-c.Save_to_database_status:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_status")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_status to time websocket")
				return
			}

		case <-ticker.C: // If our ticker has reached its time
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline

			// If our client has disconected from the websocket on thier end, close the client and its connection by returning and executing our defer statement
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with ticker pingmessage with notlive client")
				return
			}
		}
	}
}

// writeTime pumps messages from the hub to the websocket connection.
//
// A goroutine running writeTime is started for each connection to time. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeTime() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("", c.conn.RemoteAddr(), " ", "Stopping time clients ticker, unregistering, and closing connection")
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Lap_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Lap_packet_send")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Lap_packet_send to time websocket")
				return
			}

		case message, ok := <-c.Save_to_database_alert:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_alert")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_alert to time websocket")
				return
			}

		case message, ok := <-c.Save_to_database_status:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("!ok problem with Save_to_database_status")
				return
			}
			// Marshal our message into json so we can send it over the websocket
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with writing Save_to_database_status to time websocket")
				return
			}

		case <-ticker.C: // If our ticker has reached its time
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Add another 10 seconds to the SetWriteDeadline

			// If our client has disconected from the websocket on thier end, close the client and its connection by returning and executing our defer statement
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("", c.conn.RemoteAddr(), " ", "error with ticker pingmessage with time client")
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serve_ws(conn_type string, hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	switch conn_type {
	case "dashboard":
		client := &Client{
			hub:                     hub,
			conn_type:               conn_type,
			conn:                    conn,
			Session_packet_send:     make(chan structs.PacketSessionData),
			Lap_packet_send:         make(chan structs.PacketLapData),
			Telemetry_packet_send:   make(chan structs.PacketCarTelemetryData),
			Car_status_packet_send:  make(chan structs.PacketCarStatusData),
			Save_to_database_alert:  make(chan structs.Save_to_database_alerts),
			Save_to_database_status: make(chan structs.Save_to_database_status),
		}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writeDashboard()
		go client.readFromClients()

	case "history":
		client := &Client{
			hub:                     hub,
			conn_type:               conn_type,
			conn:                    conn,
			Save_to_database_alert:  make(chan structs.Save_to_database_alerts),
			Save_to_database_status: make(chan structs.Save_to_database_status),
		}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writeHistory()

	case "time":
		client := &Client{
			hub:                     hub,
			conn_type:               conn_type,
			conn:                    conn,
			Lap_packet_send:         make(chan structs.PacketLapData),
			Save_to_database_alert:  make(chan structs.Save_to_database_alerts),
			Save_to_database_status: make(chan structs.Save_to_database_status),
		}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writeTime()

	default:
		log.Println("ws client type is invalid, type:", conn_type)
		return
	}
}
