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

		var special []*logic.Card
		var feature scoring
		for _, feature = range d.features {
			special = append(special, feature.MarkCards(d)...)
		}
		var playerspecial map[int]*logic.Deck
		var c *logic.Card
		var winner int
		for _, c = range special {
			winner = d.whoWon(c)
			if winner != -1 {
				playerspecial[winner].AddAll(c)
			}
		}

		var player int
		var d *logic.Deck
		for player, d = range playerspecial {
			buf.WriteByte(byte(player))
			d.WriteBinary(buf)
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
