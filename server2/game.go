package main

import "log"

// Game contains server-relevant information about a game
type Game struct {
	id      byte
	players map[int]*Client
	name    string
	hub     *Hub

	// hands: maps #player to #slot to []item
	hands map[int]*Inventory
	// table: maps #slot to []item
	table *Inventory

	ruleset Ruleset
}

// Ruleset implements all moves a game (type) should have
type Ruleset interface {
	Reset() bool
	PlayerMove(player int, p *Packet) bool
}

// GameInfo contains user-relevant information about a game
type GameInfo struct {
	ID         byte   `json:"id"`
	Name       string `json:"name"`
	Game       string `json:"game"`
	Players    int    `json:"players"`
	Maxplayers int    `json:"maxplayers"`
}

func (g *Game) sendHands() {
	var i int
	var c *Client
	for i, c = range g.players {
		var inv *Inventory = g.hands[i]
		log.Printf("player %s should get %s", c, inv.Send())
	}
}

func (g *Game) info() GameInfo {
	return GameInfo{
		g.id, g.name, "Doppelkopf", len(g.players), 4,
	}
}

func (g *Game) playerMove(player *Client, move *Packet) bool {
	return g.ruleset.PlayerMove(g.playerID(player), move)
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

func (g *Game) playerLeave(player *Client) bool {
	var i int = g.playerID(player)
	if i == -1 {
		return false
	}

	delete(g.players, i)
	return true
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
