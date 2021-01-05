package main

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *clientMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// All current games
	games []*Game
}

type clientMessage struct {
	c   *Client
	msg []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *clientMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Returns a pointer to first game this player is part of
func (h *Hub) playerMessage(sender *Client, p *Packet) bool {
	log.Printf("Player %s sent %s", sender, p)

	var g *Game
	if p.Action() == "join" {
		if len(p.Args()) < 1 {
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
		for _, g = range h.games {
			if g.id == byte(id) {
				return g.playerJoin(sender)
			}
		}
		log.Print("No game found")
		return false
	}

	var c *Client
	g = h.playerGame(sender)
	if g != nil {
		return g.playerMove(c, p)
	}
	log.Printf("Found no matching game for %s", sender)
	return false
}

func (h *Hub) createGame(name string) *Game {
	var g Game = Game{h.gameUUID(), make(map[int]*Client), name, h, StatePreparing, nil}
	g.ruleset = NewDoko(&g)
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
	var client *Client
	for player, client = range g.players {
		log.Printf("Sending update to #%d: %s", player, client)
		var buf bytes.Buffer
		g.WriteBinary(player, &buf)
		log.Printf("Player #%d got %s", player, buf.String())
		client.send <- buf.Bytes()
	}
}

func (h *Hub) logGames() {
	var g *Game
	log.Print("-----------------------------")
	for _, g = range h.games {
		log.Printf("%5d | %10s | %10s | %3d/%3d | %d",
			g.id, g.ruleset.Info().Name, g.ruleset.Info().Game,
			g.ruleset.Info().Players, g.ruleset.Info().Maxplayers,
			g.state)
	}
	log.Print("-----------------------------")
}

func (h *Hub) gameUUID() byte {
	var uuid byte = byte(time.Now().UnixNano())
	for h.gameByUUID(uuid) != nil {
		uuid++
	}
	return uuid
}

func (h *Hub) playerGame(player *Client) *Game {
	var g *Game
	var c *Client
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
	var c *Client
	var msg *clientMessage
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
				close(c.send)
			}
		case msg = <-h.broadcast:
			h.playerMessage(msg.c, msg.toPacket())
			h.sendUpdates(h.playerGame(msg.c))
		}
	}
}
