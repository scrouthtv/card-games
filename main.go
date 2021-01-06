package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path"
	"strings"
)

var addr = flag.String("addr", ":8080", "http service address")
var hub *Hub = newHub()

func servePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "play.html")
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Requested %s, returning %s", r.URL, path.Join("static/", r.URL.Path))
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/home.html")
	} else {
		http.ServeFile(w, r, path.Join("static/", r.URL.Path))
	}
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
			games = append(games, g.ruleset.Info())
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

func serveDeck(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, strings.ToLower(r.URL.Path[1:]))
}

func main() {
	flag.Parse()
	go hub.run()
	http.HandleFunc("/play", servePlayer)
	http.HandleFunc("/api", serveAPI)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/serialize-props.js", serveProps)
	http.HandleFunc("/deck/", serveDeck)
	http.HandleFunc("/", serveStatic)
	hub.createGameWithID(5, "aa")
	//hub.createGame("bb")
	//hub.createGame("cc")
	log.Printf("Started http server")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
