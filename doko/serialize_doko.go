package doko

import (
	"bytes"

	"github.com/scrouthtv/card-games/logic"
)

const (
	// DokoGameUUID is a unique id for the dok (de-) serializer
	DokoGameUUID = 1
)

// WriteBinary writes the game's information
// relevant to the specified player to the buffer
func (d *Doko) WriteBinary(player int, buf *bytes.Buffer) {
	switch d.g.State() {
	case logic.StatePreparing:
		buf.WriteByte(logic.StatePreparing)
	case logic.StatePlaying:
		buf.WriteByte(logic.StatePlaying | (byte(d.active) << 2) | (byte(player) << 5))
		d.hands[player].WriteBinary(buf)
		d.table.WriteBinary(buf)
		var i int
		for i = 0; i < 4; i++ {
			buf.WriteByte(byte(d.won[i].Length()))
		}

	case logic.StateEnded:
		buf.WriteByte(logic.StateEnded)
		var scores []int = d.Scores()
		var score int
		for _, score = range scores {
			buf.WriteByte(byte(score))
		}
	default:
		buf.WriteByte(0)
	}

}
