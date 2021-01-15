package main

import (
	"bytes"
	"log"

	"github.com/scrouthtv/card-games/logic"
	"github.com/scrouthtv/card-games/doko"
)

// Game contains server-relevant information about a game
type Game struct {
	id          byte
	players     map[int]logic.Player
	playerNames map[int]string
	name        string
	hub         *Hub
	state       byte

	ruleset     Ruleset
}

// Ruleset implements all moves a game (type) should have
type Ruleset interface {
	Reset() bool
	PlayerMove(player int, p *logic.Packet) bool
	WriteBinary(player int, buf *bytes.Buffer)
	Info() logic.GameInfo
	TypeID() byte
	Start()
}

// StartIfReady starts the game if enough players joined
func (g *Game) StartIfReady() {
	if g.ruleset.Info().Players == g.ruleset.Info().Maxplayers {
		g.Start()
	}
}

// SendUpdates sends the current game data to all clients
func (g *Game) SendUpdates() {
	g.hub.sendUpdates(g)
}

// ID returns the uuid for this game
func (g *Game) ID() byte {
	return g.id
}

// Name returns the name for this game
func (g *Game) Name() string {
	return g.name
}

// State returns the state of this game, one of
// StatePreparing, StateRunning or StateEnded
func (g *Game) State() byte {
	return g.state
}

// SetState sets the state to one of the possible states
func (g *Game) SetState(state byte) {
	g.state = state
}

// PlayerCount returns the amount of players currently in this game
func (g *Game) PlayerCount() int {
	return len(g.players)
}

// Start starts the game
func (g *Game) Start() {
	if g.state == logic.StatePreparing {
		log.Printf("Starting game %d", g.id)
		g.state = logic.StatePlaying
		g.ruleset.Start()
		g.hub.sendUpdates(g)
	}
}

func (g *Game) playerMove(player logic.Player, move *logic.Packet) bool {
	if g.ruleset.PlayerMove(g.playerID(player), move) {
		g.hub.sendUpdates(g)
		return true
	}
	return false
}

func (g *Game) playerID(player logic.Player) int {
	var i int
	var c logic.Player
	for i, c = range g.players {
		if c == player {
			return i
		}
	}
	return -1
}

func (g *Game) playerLeave(player logic.Player) bool {
	var i int = g.playerID(player)
	if i == -1 {
		return false
	}

	delete(g.players, i)
	return true
}

func (g *Game) playerJoin(player logic.Player, name string) bool {
	var playerID int = g.freeID()
	if playerID == -1 {
		log.Printf("Error: too many players joined")
		return false
	}
	if g.playerID(player) != -1 {
		log.Printf("Error: Player already joined")
		return false
	}

	g.players[playerID] = player
	g.playerNames[playerID] = name

	player.Send(g.joinInfo(playerID, player));

	g.hub.logGames()
	g.StartIfReady()
	g.hub.sendUpdates(g)
	return true
}

func (g *Game) joinInfo(id int, player logic.Player) []byte {
	var buf bytes.Buffer
	buf.WriteByte(doko.DokoGameUUID)
	return buf.Bytes()
}

// First free id, -1 if none are free
func (g *Game) freeID() int {
	var id int
	var ok bool
	for id = 0; id < g.ruleset.Info().Maxplayers; id++ {
		_, ok = g.players[id]
		if !ok {
			return id
		}
	}

	return -1
}
