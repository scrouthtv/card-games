package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

// Deck is a collection of cards
type Deck []*Card

// EmptyDeck generates an empty deck
func EmptyDeck() *Deck {
	var d Deck
	return &d
}

// NewDeck generates a deck with the cards specified via their values in each suit.
func NewDeck(values []int) *Deck {
	var value int
	var suit int
	var deck Deck
	for _, value = range values {
		for _, suit = range []int{Clubs, Diamonds, Hearts, Spades} {
			deck = append(deck, &Card{suit, value})
		}
	}

	return &deck
}

// DeserializeDeck recreates a deck from its Short() representation
func DeserializeDeck(str string) *Deck {
	var deck Deck
	var cstr string
	var ok bool
	for _, cstr = range strings.Split(str, ", ") {
		var card *Card
		ok, card = CardFromShort(cstr)
		if ok {
			deck = append(deck, card)
		} else {
			log.Printf("Can't read card %s", cstr)
		}
	}

	return &deck
}

// Get returns the card at the idx-th index in this deck
// The function panics if idx is out of range
func (d *Deck) Get(idx int) *Card {
	return (*d)[idx]
}

// Subdeck returns a pointer to the portion of the deck
// if start or end are out of bounds, the function panics
func (d *Deck) Subdeck(start int, end int) *Deck {
	var subdeck Deck = (*d)[start:end]
	return &subdeck
}

func (d *Deck) String() string {
	var out strings.Builder

	var i int
	var c *Card
	for i, c = range *d {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(c.String())
	}

	return out.String()
}

func (d *Deck) Short() string {
	var out strings.Builder

	var i int
	var c *Card
	for i, c = range *d {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(c.Short())
	}

	return out.String()
}

// Value calculates the value of this deck using
// a specified function
func (d *Deck) Value(value func(c *Card) int) int {
	var sum int
	var card *Card
	for _, card = range *d {
		sum += value(card)
	}
	return sum
}

// Twice combines this deck with itself, so that every card appears twice
// Returns the same pointer
func (d *Deck) Twice() *Deck {
	*d = append(*d, *d...)
	return d
}

// Contains checks if this deck contains the specified card at least once
func (d *Deck) Contains(card Card) bool {
	var cid *Card
	for _, cid = range *d {
		if *cid == card {
			return true
		}
	}
	return false
}

// ContainsAny checks if the acceptor returns true for any card in this deck
func (d *Deck) ContainsAny(acceptor func(c *Card) bool) bool {
	var cid *Card
	for _, cid = range *d {
		if acceptor(cid) {
			return true
		}
	}
	return false
}

// AddAll appends all specified at the end of this deck
func (d *Deck) AddAll(cards ...*Card) {
	*d = append(*d, cards...)
}

// Merge merges another deck into this, leaving the other intact
func (d *Deck) Merge(other *Deck) {
	*d = append(*d, *other...)
}

// Length returns the amount of cards in this deck
func (d *Deck) Length() int {
	return len(*d)
}

// Remove removes the specified card at most n times from this deck,
// returning how many cards were actually removed
func (d *Deck) Remove(card Card, n int) int {
	if n == 0 {
		return 0
	}

	var i int
	var deleted int = 0
	for i = 0; i < len(*d); i++ {
		if *(*d)[i] == card {
			*d = append((*d)[:i], (*d)[i+1:]...)
			deleted++
			if deleted == n {
				return deleted
			}
		}
	}
	return deleted
}

// Shuffle shuffles a deck
// Returns the same pointer
func (d *Deck) Shuffle() *Deck {
	var deck Deck = *d
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i int, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return d
}

// Equal compares two decks if they contain the same cards in the same order
// nil is equal to any deck with 0 cards
func (d *Deck) Equal(other *Deck) bool {
	var i int
	var c *Card
	if d == other {
		return true
	} else if d == nil {
		return other.Length() == 0
	} else if other == nil {
		return d.Length() == 0
	} else if len(*d) != len(*other) {
		return false
	}
	for i, c = range *d {
		if *c != *(*other)[i] {
			return false
		}
	}
	return true
}

// Distribute creates slices which each contain cards cards
// discarding the rest
// If hands * decks is more than the amount of cards in the deck,
// the function panics
func (d *Deck) Distribute(hands int, cards int) []*Deck {
	if hands*cards > len(*d) {
		panic("Too few cards for distribution")
	}

	var distribution []*Deck = make([]*Deck, hands)
	var i int
	for i = 0; i < hands; i++ {
		distribution[i] = d.Subdeck(i*cards, (i+1)*cards)
	}

	return distribution
}

// DistributeAll distributes all cards in this deck evenly to hands
// amount players
func (d *Deck) DistributeAll(hands int) []*Deck {
	return d.Distribute(hands, len(*d)/hands)
}
