package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/scrouthtv/card-games/doko"
	"github.com/scrouthtv/card-games/logic"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[logic.Player]bool

	// Inbound messages from the clients.
	broadcast chan *playerMessage

	// Register requests from the clients.
	register chan logic.Player

	// Unregister requests from clients.
	unregister chan logic.Player

	// All current games
	games []*Game
}

type playerMessage struct {
	c   logic.Player
	msg []byte
}

func (cm *playerMessage) toPacket() *logic.Packet {
	var msg string = string(cm.msg)
	var p logic.Packet = logic.Packet(strings.Split(msg, " "))
	return &p
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *playerMessage),
		register:   make(chan logic.Player),
		unregister: make(chan logic.Player),
		clients:    make(map[logic.Player]bool),
	}
}

// Returns a pointer to first game this player is part of
func (h *Hub) playerMessage(sender logic.Player, p *logic.Packet) bool {
	log.Printf("Player %s sent %s", sender, p)

	var g *Game
	if p.Action() == "join" {
		if len(p.Args()) < 2 {
			log.Print("Missing game")
			return false
		}

		var id int
		var err error
		id, err = strconv.Atoi(p.Args()[0])
		if err != nil {
			log.Print("Wrong game: ", err)
			return false
		}

		var name string = p.Args()[1]

		for _, g = range h.games {
			if g.id == byte(id) {
				return g.playerJoin(sender, name)
			}
		}
		log.Print("No game found")
		return false
	}

	g = h.playerGame(sender)
	if g != nil {
		return g.playerMove(sender, p)
	}
	log.Printf("Found no matching game for %s", sender)
	return false
}

func (h *Hub) createGame(name string) *Game {
	var g Game = Game{h.gameUUID(), make(map[int]logic.Player), make(map[int]string), name, h, logic.StatePreparing, nil}
	g.ruleset = doko.NewDoko(&g)
	h.games = append(h.games, &g)
	h.logGames()
	return &g
}

func (h *Hub) createGameWithID(id byte, name string) *Game {
	var g Game = Game{id, make(map[int]logic.Player), make(map[int]string), name, h, logic.StatePreparing, nil}
	g.ruleset = doko.NewDoko(&g)
	h.games = append(h.games, &g)
	h.logGames()
	return &g
}

func (h *Hub) sendUpdates(g *Game) {
	log.Println("sending updates")
	if g == nil {
		return
	}

	var player int
	var client logic.Player
	for player, client = range g.players {
		var buf bytes.Buffer
		g.WriteBinary(player, &buf)
		client.Send(buf.Bytes())
	}
}

func (h *Hub) logGames() {
	var g *Game
	log.Print("----------------------------------------------")
	for _, g = range h.games {
		log.Printf("%5d | %10s | %10s | %3d/%3d | %d",
			g.id, g.ruleset.Info().Name, g.ruleset.Info().Game,
			g.ruleset.Info().Players, g.ruleset.Info().Maxplayers,
			g.state)
	}
	log.Print("----------------------------------------------")
}

func (h *Hub) gameUUID() byte {
	var uuid byte = byte(time.Now().UnixNano())
	for h.gameByUUID(uuid) != nil {
		uuid++
	}
	return uuid
}

func (h *Hub) playerGame(player logic.Player) *Game {
	var g *Game
	var c logic.Player
	for _, g = range h.games {
		for _, c = range g.players {
			if c == player {
				return g
			}
		}
	}
	return nil
}

func (h *Hub) gameByUUID(uuid byte) *Game {
	var game *Game
	for _, game = range h.games {
		if game.id == uuid {
			return game
		}
	}
	return nil
}

func (h *Hub) run() {
	log.Print("Started Hub")
	var c logic.Player
	var msg *playerMessage
	for {
		select {
		case c = <-h.register:
			h.clients[c] = true
		case c = <-h.unregister:
			log.Printf("Unregistering %s", c)
			var g *Game = h.playerGame(c)
			if g != nil {
				g.playerLeave(c)
			}
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				c.Close()
			}
		case msg = <-h.broadcast:
			h.playerMessage(msg.c, msg.toPacket())
			//h.sendUpdates(h.playerGame(msg.c))
		}
	}
}
