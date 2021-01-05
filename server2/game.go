package main

import (
	"bytes"
	"log"
)

// Game contains server-relevant information about a game
type Game struct {
	id      byte
	players map[int]*Client
	name    string
	hub     *Hub
	state   byte

	ruleset Ruleset
}

const (
	// StatePreparing indicates that the game is currently preparing (e. g. waiting for players)
	StatePreparing = iota
	// StatePlaying indicates that the game is currently running
	StatePlaying
	// StateEnded indicates that the game has ended
	StateEnded
)

// Ruleset implements all moves a game (type) should have
type Ruleset interface {
	Reset() bool
	PlayerMove(player int, p *Packet) bool
	WriteBinary(player int, buf *bytes.Buffer)
	Info() GameInfo
	TypeID() byte
}

// GameInfo contains user-relevant information about a game
type GameInfo struct {
	ID         byte   `json:"id"`
	Name       string `json:"name"`
	Game       string `json:"game"`
	Players    int    `json:"players"`
	Maxplayers int    `json:"maxplayers"`
}

// Start starts the game
func (g *Game) Start() {
	g.state = StatePlaying
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

	g.players[i] = nil
	return true
}

func (g *Game) playerJoin(player *Client) bool {
	if g.ruleset.Info().Players >= g.ruleset.Info().Maxplayers {
		log.Printf("Error: too many players joined")
		return false
	}
	if g.playerID(player) != -1 {
		log.Printf("Error: Player already joined")
		return false
	}

	g.players[len(g.players)] = player
	g.hub.logGames()
	g.hub.sendUpdates(g)
	return true
}
