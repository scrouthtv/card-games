package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var hub *Hub = newHub()

func servePlayer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Requested invalid URL")
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "play.html")
}

func serveLanding(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func serveAPI(w http.ResponseWriter, r *http.Request) {
	//log.Println(r)
	var req map[string][]string = r.URL.Query()
	var jw *json.Encoder = json.NewEncoder(w)
	var ok bool
	if _, ok = req["games"]; ok {
		var games []GameInfo
		var g *Game
		for _, g = range hub.games {
			games = append(games, g.info())
		}
		jw.Encode(games)
		return
	} else if _, ok = req["create"]; ok {
		var name []string
		name, ok = req["name"]
		if ok && len(name) > 0 {
			var g *Game = hub.createGame(name[0])
			log.Printf("New game with id %d created\n", g.id)
			jw.Encode(CreateResponse{true, "", g.id})
		} else {
			jw.Encode(CreateResponse{false, "Missing key: name", 0})
		}
	} else {
		w.Write([]byte("Unknown request"))
	}
}

func main() {
	flag.Parse()
	go hub.run()
	http.HandleFunc("/play", servePlayer)
	//http.HandleFunc("/join", serveLanding)
	http.HandleFunc("/api", serveAPI)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", serveLanding)
	hub.createGame("aa")
	hub.createGame("bb")
	hub.createGame("cc")
	log.Printf("Started http server")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
