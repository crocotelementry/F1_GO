//
// To-do:
// 	2.10.19		Clean up repetitive code and create functions to reuse
//
package main

import (
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// All Registered clients
	all_clients map[*Client]bool
	// Registered live live_clients
	live_clients map[*Client]bool
	// Registered dashboard clients.
	dashboard_clients map[*Client]bool
	// Registered history clients.
	history_clients map[*Client]bool
	// Registered time clients
	time_clients map[*Client]bool

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
		broadcast:         make(chan *Udp_data),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		all_clients:       make(map[*Client]bool),
		live_clients:      make(map[*Client]bool),
		dashboard_clients: make(map[*Client]bool),
		history_clients:   make(map[*Client]bool),
		time_clients:      make(map[*Client]bool),
	}

}

// Hub function that handles all connecting and disconnecting clients and the messages that are to be
// broadcasted to them!
func (h *Hub) run() {
	for {
		select {
		// If we have a new dashboard client connecting!
		case client := <-h.register:
			switch client.conn_type {
			case "dashboard":
				h.live_clients[client] = true
				h.dashboard_clients[client] = true
				h.all_clients[client] = true
				log.Println("", client.conn.RemoteAddr(), " ", "dashboard client opened")
			case "time":
				h.live_clients[client] = true
				h.time_clients[client] = true
				h.all_clients[client] = true
				log.Println("", client.conn.RemoteAddr(), " ", "time client opened")
			case "history":
				h.history_clients[client] = true
				h.all_clients[client] = true
				log.Println("", client.conn.RemoteAddr(), " ", "history client opened")
			}

		// If we have a new client that disconected
		case client := <-h.unregister:
			switch client.conn_type {
			case "dashboard":
				if _, ok := h.dashboard_clients[client]; ok {
					log.Println("", client.conn.RemoteAddr(), " ", "0dashboard client closed")
					delete(h.dashboard_clients, client)  // Delete the client
					delete(h.live_clients, client)       // Delete the client
					delete(h.all_clients, client)        // Delete the client
					close(client.Session_packet_send)    // Close all the different packet channels!
					close(client.Lap_packet_send)        // ..
					close(client.Telemetry_packet_send)  // ..
					close(client.Car_status_packet_send) // ..
					close(client.Save_to_database_alert) // ..
				}
			case "time":
				if _, ok := h.time_clients[client]; ok {
					log.Println("", client.conn.RemoteAddr(), " ", "time client closed")
					delete(h.time_clients, client)       // Delete the client
					delete(h.live_clients, client)       // Delete the client
					delete(h.all_clients, client)        // Delete the client
					close(client.Lap_packet_send)        // Close all the different packet channels!
					close(client.Save_to_database_alert) // ..
				}
			case "history":
				if _, ok := h.history_clients[client]; ok {
					log.Println("", client.conn.RemoteAddr(), " ", "history client closed")
					delete(h.history_clients, client)    // Delete the client
					delete(h.all_clients, client)        // Delete the client
					close(client.Save_to_database_alert) // Close the packet channel!
				}
			}

		// If we have a message to broadcast
		case message := <-h.broadcast:
			switch message.Id {
			// case 0:
			// 	client.Motion_packet_send <- message.Motion_packet
			case 1:
				for client := range h.dashboard_clients { // Loop through all our clients
					select {
					case client.Session_packet_send <- message.Session_packet:
						// default:
						// 	log.Println("", client.conn.RemoteAddr(), " ", "1dashboard client closed")
						// 	delete(h.dashboard_clients, client)   // Delete the client
						// 	delete(h.live_clients, client)   // Delete the client
						// 	delete(h.all_clients, client)   // Delete the client
						// 	close(client.Session_packet_send)     // Close all the different packet channels!
						// 	close(client.Lap_packet_send)         // ..
						// 	close(client.Telemetry_packet_send)   // ..
						// 	close(client.Car_status_packet_send)  // ..
						// 	close(client.Save_to_database_alert)  // ..
					}
				}
			case 2:
				for client := range h.dashboard_clients { // Loop through all our dashboard clients
					select {
					case client.Lap_packet_send <- message.Lap_packet:
						// default:
						// 	log.Println("2 problem")
						// 	log.Println("", client.conn.RemoteAddr(), " ", "2dashboard client closed")
						// 	delete(h.dashboard_clients, client)   // Delete the client
						// 	delete(h.live_clients, client)   // Delete the client
						// 	delete(h.all_clients, client)   // Delete the client
						// 	close(client.Session_packet_send)     // Close all the different packet channels!
						// 	close(client.Lap_packet_send)         // ..
						// 	close(client.Telemetry_packet_send)   // ..
						// 	close(client.Car_status_packet_send)  // ..
						// 	close(client.Save_to_database_alert)  // ..
					}
				}
				for client := range h.time_clients { // Loop through all our time clients
					select {
					case client.Lap_packet_send <- message.Lap_packet:
						// default:
						// 	log.Println("", client.conn.RemoteAddr(), " ", "time client closed")
						// 	delete(h.time_clients, client)    // Delete the client
						// 	delete(h.live_clients, client)   // Delete the client
						// 	delete(h.all_clients, client)   // Delete the client
						// 	close(client.Lap_packet_send)         // Close all the different packet channels!
						// 	close(client.Save_to_database_alert)  // ..
					}
				}
			// case 3:
			// 	client.Event_packet_send <- message.Event_packet
			// case 4:
			// 	client.Participant_packet_send <- message.Participant_packet
			// case 5:
			// 	client.Car_setup_packet_send <- message.Car_setup_packet
			case 6:
				for client := range h.dashboard_clients { // Loop through all our clients
					select {
					case client.Telemetry_packet_send <- message.Telemetry_packet:
						// default:
						// 	log.Println("", client.conn.RemoteAddr(), " ", "3dashboard client closed")
						// 	delete(h.dashboard_clients, client)   // Delete the client
						// 	delete(h.live_clients, client)   // Delete the client
						// 	delete(h.all_clients, client)   // Delete the client
						// 	close(client.Session_packet_send)     // Close all the different packet channels!
						// 	close(client.Lap_packet_send)         // ..
						// 	close(client.Telemetry_packet_send)   // ..
						// 	close(client.Car_status_packet_send)  // ..
						// 	close(client.Save_to_database_alert)  // ..
					}
				}
			case 7:
				for client := range h.dashboard_clients { // Loop through all our clients
					select {
					case client.Car_status_packet_send <- message.Car_status_packet:
						// default:
						// 	log.Println("", client.conn.RemoteAddr(), " ", "4dashboard client closed")
						// 	delete(h.dashboard_clients, client)   // Delete the client
						// 	delete(h.live_clients, client)   // Delete the client
						// 	delete(h.all_clients, client)   // Delete the client
						// 	close(client.Session_packet_send)     // Close all the different packet channels!
						// 	close(client.Lap_packet_send)         // ..
						// 	close(client.Telemetry_packet_send)   // ..
						// 	close(client.Car_status_packet_send)  // ..
						// 	close(client.Save_to_database_alert)  // ..
					}
				}
			case 30:
				log.Println("             ", "End of race code received. Sending Save_to_database_alert to clients")
				for client := range h.all_clients { // Loop through all our clients
					client.Save_to_database_alert <- message.Save_to_database_alert
				}
			default:
				// Packet is either not a valid packet number or it is from packets with ids of:
				// 0, 3, 4, 5
				if message.Id != 0 && message.Id != 3 && message.Id != 4 && message.Id != 5 {
					log.Println("Invalid broadcast message, id=", message.Id)
				}
			}

		}
	}
}
