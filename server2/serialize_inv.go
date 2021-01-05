package main

import "bytes"

const (
	cardMaxSuit = 4
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

}

// WriteBinary appends the game to a bytes buffer
// It sends these fields:
//
func (g *Game) WriteBinary(player int, buf *bytes.Buffer) {

}
