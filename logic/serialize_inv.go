package logic

import (
	"bytes"
)

const (
	// CardMaxSuit sepcifies how many card suits exist
	CardMaxSuit = 4
)

// ValueOrder maps int value (for serialization) to card value
var ValueOrder []int = []int{
	Ace, 2, 3, 4, 5, 6, 7, 8, 9, 10, Jack, Queen, King,
}

// ToBinary returns a byte representation fo this card
func (c *Card) ToBinary() byte {
	return byte(c.value*CardMaxSuit + c.suit)
}

// WriteBinary appends the deck to a bytes buffer
func (d *Deck) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte(byte(d.Length()))

	var c *Card
	for _, c = range *d {
		buf.WriteByte(c.ToBinary())
	}
}
