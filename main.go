package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")

	// Create our UDP socket
	addrs, _  = net.ResolveUDPAddr("udp", ":20777")
	sock, err = net.ListenUDP("udp", addrs)
)

// Prints out if we connect to our databases and then the table scheme
func logPrintFormat() {
	ready := <-redis_ping_done
	if ready == true {
		fmt.Println("Successful connection made with Redis database!\n")
	} else {
		fmt.Println("Problem connecting with Redis database!\n")
	}

	fmt.Println("   DATE      TIME     IP    PORT    ADDRESS/EVENT  ")
	fmt.Println("----------|--------|------|------|-----------------")
	return
}

// Called when when a client lands on localhost:8080/
func liveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("", r.RemoteAddr, " ", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./web/live_dashboard.html")
}

// Called when when a client lands on localhost:8080/history
func historyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("", r.RemoteAddr, " ", r.URL)
	if r.URL.Path != "/history" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./web/history_dashboard.html")
}

// Called when when a client lands on localhost:8080/time
func timeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("", r.RemoteAddr, " ", r.URL)
	if r.URL.Path != "/time" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./web/time_dashboard.html")
}

func main() {
	// Sets the location in which to serve our static files from for our webpage
	var dir string
	flag.StringVar(&dir, "dir", "./web", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	go logPrintFormat()

	// Create a hub object to handle our client connections and needs
	live_hub := newHub()
	go live_hub.run()

	// Start getting data from the game
	// Add to redis and send to hub
	go getGameData(live_hub)

	// Creates a new mux
	router := mux.NewRouter()
	// WHen our html page calls its static files from /static/file, this sets the location to grab them from
	// TODO: Is this all this code needed? Couldn't we just set it web/?
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir(dir))))

	// Our handler functions for each page
	// Landing page /aka live telemetry or telemetry_dashboard
	router.HandleFunc("/", liveHandler)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serve_ws("live", live_hub, w, r)
	})

	// History page /aka history_dashboard
	router.HandleFunc("/history", historyHandler)
	router.HandleFunc("/history/ws", func(w http.ResponseWriter, r *http.Request) {
		serve_ws("history", live_hub, w, r)
	})

	// Live time page /aka time_dashboard
	router.HandleFunc("/time", timeHandler)
	router.HandleFunc("/time/ws", func(w http.ResponseWriter, r *http.Request) {
		serve_ws("time", live_hub, w, r)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
