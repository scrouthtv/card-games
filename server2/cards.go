package main

import "strconv"

type Card struct {
	suit  int
	value int
}

const (
	clubs = iota
	diamonds
	hearts
	spades
)

// Short returns a representation of this card in the form
// [cdhs][a2-10jqk]
func (c *Card) Short() string {
	var out string
	switch c.suit {
	case clubs:
		out = "c"
	case diamonds:
		out = "d"
	case hearts:
		out = "h"
	case spades:
		out = "s"
	}
	switch c.value {
	case 1:
		out += "a"
	case 11:
		out += "j"
	case 12:
		out += "q"
	case 13:
		out += "k"
	default:
		out += strconv.Itoa(c.value)
	}
	return out
}
