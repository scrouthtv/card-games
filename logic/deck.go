package logic

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

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

// Short returns a string representation of this deck by joinin
// all Shorts() on the cards
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
	if d == nil {
		return 0
	}
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
