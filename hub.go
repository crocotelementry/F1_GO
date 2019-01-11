package main

import (
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered live clients.
	live_clients map[*Client]bool

	// Registered notlive clients.
	notlive_clients map[*Client]bool

	// Inbound messages from F1 2018 UDP to be broacasted to connected clients
	broadcast chan *Udp_data

	// Register requests from the live clients.
	live_register chan *Client

	// Register requests from the notlive clients.
	notlive_register chan *Client

	// Unregister requests from live clients.
	live_unregister chan *Client

	// Unregister requests from notlive clients.
	notlive_unregister chan *Client
}

// Creates a new hub
func newHub() *Hub {
	return &Hub{
		broadcast:          make(chan *Udp_data),
		live_register:      make(chan *Client),
		notlive_register:   make(chan *Client),
		live_unregister:    make(chan *Client),
		notlive_unregister: make(chan *Client),
		live_clients:       make(map[*Client]bool),
		notlive_clients:    make(map[*Client]bool),
	}

}

// Hub function that handles all connecting and disconnecting clients and the messages that are to be
// broadcasted to them!
func (h *Hub) run() {
	for {
		select {
		// If we have a new live client connecting!
		case client := <-h.live_register:
			h.live_clients[client] = true
			log.Println("", client.conn.RemoteAddr(), " ", "live client opened")
		// If we have a new notlive client connecting!
		case client := <-h.notlive_register:
			h.notlive_clients[client] = true
			log.Println("", client.conn.RemoteAddr(), " ", "history client opened")

			// If we have a live client that is disconnecting
		case client := <-h.live_unregister:
			if _, ok := h.live_clients[client]; ok {
				log.Println("", client.conn.RemoteAddr(), " ", "live client closed")
				delete(h.live_clients, client)        // Delete the client
				close(client.Motion_packet_send)      // Close all the different packet channels!
				close(client.Session_packet_send)     // ..
				close(client.Lap_packet_send)         // ..
				close(client.Event_packet_send)       // ..
				close(client.Participant_packet_send) // ..
				close(client.Car_setup_packet_send)   // ..
				close(client.Telemetry_packet_send)   // ..
				close(client.Car_status_packet_send)  // ..
				close(client.save_to_database_alert)  // ..
			}
			// If we have a notlive client that is disconnecting
		case client := <-h.notlive_unregister:
			if _, ok := h.notlive_clients[client]; ok {
				log.Println("", client.conn.RemoteAddr(), " ", "history client closed")
				delete(h.notlive_clients, client)    // Delete the client
				close(client.save_to_database_alert) // Close the packet channel!
			}

		// If we have a message to broadcast
		case message := <-h.broadcast:
			for client := range h.live_clients { // Loop through all our clients
				switch message.Id { // Depending on what packet we have to send, send that packet
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
