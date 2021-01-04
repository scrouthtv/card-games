package main

import (
	"log"
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
func (h *Hub) playerMove(msg *clientMessage) {
	var g *Game
	var c *Client
	for _, g = range h.games {
		for _, c = range g.players {
			if c == msg.c {
				g.playerMove(c, string(msg.msg))
			}
		}
	}
	log.Printf("Found no matching game for %s", msg.c)
	return
}

func (h *Hub) createGame(name string) *Game {
	var g Game = Game{h.gameUUID(), make(map[int]*Client), name}
	h.games = append(h.games, &g)
	return &g
}

func (h *Hub) gameUUID() byte {
	var uuid byte = byte(time.Now().UnixNano())
	for h.gameByUUID(uuid) != nil {
		uuid++
	}
	return uuid
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
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case msg = <-h.broadcast:
			log.Printf("Got a message: %s", msg)
			h.playerMove(msg)
			/*for client := range h.clients {
				select {
				case client.send <- msg.msg: // send the current message to every client
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}*/
		}
	}
}
