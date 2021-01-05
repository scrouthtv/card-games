package main

import "log"

// Game contains server-relevant information about a game
type Game struct {
	id      byte
	players map[int]*Client
	name    string
	hub     *Hub
}

// GameInfo contains user-relevant information about a game
type GameInfo struct {
	ID         byte   `json:"id"`
	Name       string `json:"name"`
	Game       string `json:"game"`
	Players    int    `json:"players"`
	Maxplayers int    `json:"maxplayers"`
}

func (g *Game) info() GameInfo {
	return GameInfo{
		g.id, g.name, "Doppelkopf", len(g.players), 4,
	}
}

func (g *Game) playerMove(player *Client, move *Packet) bool {
	var np int
	var i int
	var c *Client
	for i, c = range g.players {
		if c == player {
			np = i + 1
		}
	}
	if np >= len(g.players) {
		np++
	}

	g.players[np].send <- []byte("Hello there, its ur turn")

	return true
}

func (g *Game) playerID(player *Client) int {
	var i int
	var c *Client
	for i, c = range g.players {
		if c == player {
			return i
		}
	}
	return -1
}

func (g *Game) playerJoin(player *Client) bool {
	if g.info().Players >= g.info().Maxplayers {
		log.Printf("Error: too many players joined")
		return false
	}
	if g.playerID(player) != -1 {
		log.Printf("Error: Player already joined")
		return false
	}

	g.players[len(g.players)] = player
	log.Printf("Player %s joined %d:", player, g.id)
	g.hub.logGames()
	return true
}
