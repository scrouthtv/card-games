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

const (
	ace   = iota + 1  // 0 + 1
	jack  = iota + 10 // 1 + 10
	queen             // 1 + 11
	king              // 1 + 12
)

func (c *Card) String() string {
	var out string
	switch c.suit {
	case clubs:
		out = "club "
	case diamonds:
		out = "diamond "
	case hearts:
		out = "heart "
	case spades:
		out = "spade "
	}

	switch c.value {
	case 1:
		out += "ace"
	case 11:
		out += "jack"
	case 12:
		out += "queen"
	case 13:
		out += "king"
	default:
		out += strconv.Itoa(c.value)
	}

	return out
}

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

// Deck is a collection of cards
type Deck []Card

// NewDeck generates a deck with the cards specified via their values in each suit.
func NewDeck(values []int) *Deck {
	var value int
	var suit int
	var deck Deck
	for _, value = range values {
		for _, suit = range []int{clubs, diamonds, hearts, spades} {
			deck = append(deck, Card{suit, value})
		}
	}

	return &deck
}

// Twice combines a deck with itself, so that every card appears twice
func (d *Deck) Twice() {
	*d = append(*d, *d...)
}
