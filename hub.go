package main

import (
	"log"
	
	"F1_GO/structs"
)

var (
	new_motion_packet      = make(chan structs.PacketMotionData)
	new_session_packet     = make(chan structs.PacketSessionData)
	new_lap_packet         = make(chan structs.PacketLapData)
	new_event_packet       = make(chan structs.PacketEventData)
	new_participant_packet = make(chan structs.PacketParticipantsData)
	new_car_setup_packet   = make(chan structs.PacketCarSetupData)
	new_telemetry_packet   = make(chan structs.PacketCarTelemetryData)
	new_car_status_packet  = make(chan structs.PacketCarStatusData)
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from F1 2018 UDP
	broadcast chan *Udp_data

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Udp_data),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("", client.conn.RemoteAddr(), " ", "client opened")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("", client.conn.RemoteAddr(), " ", "client closed")
				delete(h.clients, client)
				close(client.Motion_packet_send)
				close(client.Session_packet_send)
				close(client.Lap_packet_send)
				close(client.Event_packet_send)
				close(client.Participant_packet_send)
				close(client.Car_setup_packet_send)
				close(client.Telemetry_packet_send)
				close(client.Car_status_packet_send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				switch message.Id {
				case 0:
					// send_message := message.Motion_packet
					client.Motion_packet_send <- message.Motion_packet
				case 1:
					// send_message := message.Session_packet
					client.Session_packet_send <- message.Session_packet
				case 2:
					// send_message := message.Lap_packet
					client.Lap_packet_send <- message.Lap_packet
				case 3:
					// send_message := message.Event_packet
					client.Event_packet_send <- message.Event_packet
				case 4:
					// send_message := message.Participant_packet
					client.Participant_packet_send <- message.Participant_packet
				case 5:
					// send_message := message.Car_setup_packet
					client.Car_setup_packet_send <- message.Car_setup_packet
				case 6:
					// send_message := message.Telemetry_packet
					client.Telemetry_packet_send <- message.Telemetry_packet
				case 7:
					// send_message := message.Car_status_packet
					client.Car_status_packet_send <- message.Car_status_packet
				}
			}
		}
	}
}
