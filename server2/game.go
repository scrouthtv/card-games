package main

// Game contains server-relevant information about a game
type Game struct {
	id      byte
	players map[int]*Client
	name    string
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

func (g *Game) playerMove(player *Client, move string) bool {
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
