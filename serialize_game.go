package main

import "bytes"

// WriteBinary appends the game to a bytes buffer
// It sends these fields:
//
func (g *Game) WriteBinary(player int, buf *bytes.Buffer) {
	buf.WriteByte(g.ruleset.TypeID())
	g.ruleset.WriteBinary(player, buf)
}
