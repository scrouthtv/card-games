package logic

import (
	"strconv"
)

// Card is a french playing card of the suite clubs / diamonds / hearts / spades
// and a value from [1 = ace, 2-10, 11 = jack, 12 = queen, 13 = king]
type Card struct {
	suit  int
	value int
}

const (
	// Clubs is the first suit, aka Kreuz
	Clubs = iota
	// Diamonds is the second suit, aka Karo
	Diamonds
	// Hearts is the third suit, aka Herz
	Hearts
	// Spades is the fourth suit, aka Pik
	Spades
)

const (
	// Ace is a constant for the ace card and equal to 1
	Ace = iota + 1 // 0 + 1
	// Jack is a constant for the jack card and equal to 11
	Jack = iota + 10 // 1 + 10
	// Queen is a constant for the queen card and equal to 12
	Queen // 1 + 11
	// King is a constant for the king card and equal to 13
	King // 1 + 12
)

// NewCard creates a new card with this suit and value
func NewCard(suit int, value int) *Card {
	// TODO test if this card is valid
	var c Card = Card{suit, value}
	return &c
}

// Suit returns the suit of this card as one of the constants defined by this package
func (c *Card) Suit() int {
	return c.suit
}

// Value returns the value of this card
func (c *Card) Value() int {
	return c.value
}

func (c *Card) String() string {
	var out string
	switch c.suit {
	case Clubs:
		out = "club "
	case Diamonds:
		out = "diamond "
	case Hearts:
		out = "heart "
	case Spades:
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
	case Clubs:
		out = "c"
	case Diamonds:
		out = "d"
	case Hearts:
		out = "h"
	case Spades:
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

// CardFromShort tries to parse the card specified by short
// If parsing was sucessful, the method returns true and a pointer to the card
// If not, it returns false and nil
func CardFromShort(short string) (bool, *Card) {

	if len(short) < 2 {
		return false, nil
	}

	var c Card = Card{0, 0}

	switch short[0:1] {
	case "c":
		c.suit = Clubs
	case "d":
		c.suit = Diamonds
	case "h":
		c.suit = Hearts
	case "s":
		c.suit = Spades
	default:
		return false, nil
	}

	switch short[1:] {
	case "a":
		c.value = Ace
	case "j":
		c.value = Jack
	case "q":
		c.value = Queen
	case "k":
		c.value = King
	default:
		var value int
		var err error
		value, err = strconv.Atoi(short[1:])
		if err != nil {
			return false, nil
		}
		c.value = value
	}

	return true, &c
}
