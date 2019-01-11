package main

import (
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from F1 2018 UDP to be broacasted to connected clients
	broadcast chan *Udp_data

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Creates a new hub
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Udp_data),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Hub function that handles all connecting and disconnecting clients and the messages that are to be
// broadcasted to them!
func (h *Hub) run() {
	for {
		select {
		// If we have a new client connecting!
		case client := <-h.register:
			h.clients[client] = true
			log.Println("", client.conn.RemoteAddr(), " ", "client opened")

		// If we have a client that is disconnecting
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("", client.conn.RemoteAddr(), " ", "client closed")
				delete(h.clients, client)								// Delete the client
				close(client.Motion_packet_send)				// Close all the different packet channels!
				close(client.Session_packet_send)				// ..
				close(client.Lap_packet_send)						// ..
				close(client.Event_packet_send)					// ..
				close(client.Participant_packet_send)		// ..
				close(client.Car_setup_packet_send)			// ..
				close(client.Telemetry_packet_send)			// ..
				close(client.Car_status_packet_send)		// ..
			}

		// If we have a message to broadcast
		case message := <-h.broadcast:
			for client := range h.clients {		// Loop through all our clients
				switch message.Id {							// Depending on what packet we have to send, send that packet
				case 0:
					client.Motion_packet_send <- message.Motion_packet
				case 1:
					client.Session_packet_send <- message.Session_packet
				case 2:
					client.Lap_packet_send <- message.Lap_packet
				case 3:
					client.Event_packet_send <- message.Event_packet
				case 4:
					client.Participant_packet_send <- message.Participant_packet
				case 5:
					client.Car_setup_packet_send <- message.Car_setup_packet
				case 6:
					client.Telemetry_packet_send <- message.Telemetry_packet
				case 7:
					client.Car_status_packet_send <- message.Car_status_packet
				}
			}
		}
	}
}
