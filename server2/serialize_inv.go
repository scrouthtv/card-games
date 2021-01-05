package main

import (
	"bytes"
)

const (
	cardMaxSuit  = 4
	dokoGameUUID = 1
)

var valueOrder []int = []int{
	Ace, 2, 3, 4, 5, 6, 7, 8, 9, 10, Jack, Queen, King,
}

// ToBinary returns a byte representation fo this card
func (c *Card) ToBinary() byte {
	return byte(c.value*cardMaxSuit + c.suit)
}

// WriteBinary appends the deck to a bytes buffer
func (d *Deck) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte(byte(d.Length()))

	var c *Card
	for _, c = range *d {
		buf.WriteByte(c.ToBinary())
	}
}

// WriteBinary appends the inventory to a bytes buffer
func (inv *Inventory) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte((byte(inv.Length())))

	var d *Deck
	for _, d = range *inv {
		d.WriteBinary(buf)
	}
}

// WriteBinary writes the game's information
// relevant to the specified player to the buffer
func (d *Doko) WriteBinary(player int, buf *bytes.Buffer) {
	switch d.g.state {
	case StatePreparing:
		buf.WriteByte(StatePreparing)
	case StatePlaying:
		buf.WriteByte(StatePlaying | (byte(d.active) << 2) | (byte(player) << 5))
		d.hands[player].WriteBinary(buf)
		d.table.WriteBinary(buf)
	case StateEnded:
		buf.WriteByte(StateEnded)
		var scores []int = d.Scores()
		var score int
		for _, score = range scores {
			buf.WriteByte(byte(score))
		}
	default:
		buf.WriteByte(0)
	}

}

// WriteBinary appends the game to a bytes buffer
// It sends these fields:
//
func (g *Game) WriteBinary(player int, buf *bytes.Buffer) {
	buf.WriteByte(g.ruleset.TypeID())
	g.ruleset.WriteBinary(player, buf)
}
