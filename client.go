package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"fmt"
	"encoding/json"

	"github.com/gorilla/websocket"
	"F1_GO/structs"
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
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1341,
	WriteBufferSize: 1341,
}

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
	Motion_packet_send chan structs.PacketMotionData
	Session_packet_send chan structs.PacketSessionData
	Lap_packet_send chan structs.PacketLapData
	Event_packet_send chan structs.PacketEventData
	Participant_packet_send chan structs.PacketParticipantsData
	Car_setup_packet_send chan structs.PacketCarSetupData
	Telemetry_packet_send chan structs.PacketCarTelemetryData
	Car_status_packet_send chan structs.PacketCarStatusData
}


// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) fetchHistory() {
	defer func() {
		// c.hub.unregister <- c
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
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

    // Fetch that data it wants from database and send it over the ws
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Motion_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Session_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Lap_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Event_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Participant_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Car_setup_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Telemetry_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}

		case message, ok := <-c.Car_status_packet_send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("case ok value:", ok)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("if !ok problem")
				log.Println("message looks like:", message)
				log.Println("message is Motion_packet_send")
				return
			}
			json_message_marshaled, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
			}
			// Write our JSON formatted F1 UDP packet struct to our websocket
			if err := c.conn.WriteMessage(websocket.TextMessage, json_message_marshaled); err != nil {
				// log.Printf("error writing json_message_marshaled: %s", err)
				return
			}


		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("ticker case problem")
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

  if conn_type == "live" || conn_type == "time" {
    client := &Client{
			hub: hub,
			conn_type: conn_type,
			conn: conn,
			Motion_packet_send: make(chan structs.PacketMotionData),
			Session_packet_send: make(chan structs.PacketSessionData),
			Lap_packet_send: make(chan structs.PacketLapData),
			Event_packet_send: make(chan structs.PacketEventData),
			Participant_packet_send: make(chan structs.PacketParticipantsData),
			Car_setup_packet_send: make(chan structs.PacketCarSetupData),
			Telemetry_packet_send: make(chan structs.PacketCarTelemetryData),
			Car_status_packet_send: make(chan structs.PacketCarStatusData),
		}
  	client.hub.register <- client

  	// Allow collection of memory referenced by the caller by doing all work in
  	// new goroutines.
    go client.writePump()
    // go client.readPump() // Not used since these connections dont need to send data, only receive.

  } else if conn_type == "history" {
    client := &Client{hub: hub, conn_type: conn_type, conn: conn}

    // Allow collection of memory referenced by the caller by doing all work in
  	// new goroutines.
    go client.fetchHistory()

  } else {
    log.Println("ws client type is invalid, type:", conn_type)
    return
  }
}
