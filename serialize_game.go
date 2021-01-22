package main

import "bytes"

// WriteBinary appends the game to a bytes buffer
// It sends these fields:
//
func (g *Game) WriteBinary(player int, buf *bytes.Buffer) {
	buf.WriteByte(g.ruleset.TypeID())
	g.writePlayerInfo(buf)
	g.ruleset.WriteBinary(player, buf)
}

func (g *Game) writePlayerInfo(buf *bytes.Buffer) {
	var players int = len(g.players)
	buf.WriteByte(byte(players))

	var i int
	var name []byte
	var b byte
	for i = 0; i < players; i++ {
		name = []byte(g.playerNames[i])
		buf.WriteByte(byte(len(name)))
		for _, b = range name {
			buf.WriteByte(b)
		}
	}
}
