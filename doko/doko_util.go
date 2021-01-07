package doko

import (
	"sort"

	"github.com/scrouthtv/card-games/logic"
)

// AllowedCards determines which cards the active player is currently
// allowed to play (e. g. if they have to show a color or don't own
// that color)
func (d *Doko) AllowedCards() *logic.Deck {
	if d.table.Length() == 0 {
		return d.hands[d.active]
	}

	var show *logic.Card = d.table.Get(0)
	var allowed *logic.Deck = logic.EmptyDeck()
	var has *logic.Deck = d.hands[d.active]

	var i int
	for i = 0; i < has.Length(); i++ {
		var ownedCard *logic.Card = has.Get(i)
		if d.color(ownedCard) == d.color(show) {
			allowed.AddAll(ownedCard)
		}
	}

	if allowed.Length() == 0 {
		return d.hands[d.active]
	}

	return allowed
}

// teamsKnown counts whether both clubs queens have been played
func (d *Doko) teamsKnown() bool {
	var queens int = 0
	var deck *logic.Deck
	var c *logic.Card

	for _, c = range *d.table {
		if c.Suit() == logic.Clubs && c.Value() == logic.Queen {
			queens++
		}
	}
	for _, deck = range d.won {
		for _, c = range *deck {
			if c.Suit() == logic.Clubs && c.Value() == logic.Queen {
				queens++
			}
		}
	}

	return queens == 2
}

func (d *Doko) origOwner(card *logic.Card) int {
	var player int
	var start *logic.Deck
	var c *logic.Card

	for player, start = range d.start {
		for _, c = range *start {
			if c == card {
				//log.Printf("Card matches: %p == %p", c, card)
				return player
			}
		}
	}
	return -1
}

// beats calculates whether the attacking card atk defeats the defending card def
func (d *Doko) beats(def *logic.Card, atk *logic.Card) bool {
	if d.color(def) == d.color(atk) {
		if d.value(atk) > d.value(def) {
			return true
		} else if d.value(atk) == d.value(def) {
			return def.Suit() == logic.Hearts && def.Value() == 10
		} else {
			return false
		}
	} else if d.color(atk) == -1 {
		// attacker has trump, defender doesn't
		return true
	} else {
		// attacker didn't show def's color
		return false
	}
}

var dokoValueOrder []int = []int{9, logic.Jack, logic.Queen, logic.King, 10, logic.Ace}

func (d *Doko) value(c *logic.Card) int {
	var i, value int

	value = d.trumpValue(c)
	if value != -1 {
		// return trump value instead
		return value
	}

	for i, value = range dokoValueOrder {
		if value == c.Value() {
			return i
		}
	}
	return 0
}

// color returns the color if this card, returning -1 if the card is a trump
func (d *Doko) color(c *logic.Card) int {
	if d.trumpValue(c) == -1 {
		return c.Suit()
	}
	return -1
}

// DokoTrumpOrder specifies the order of trumps in this game
var DokoTrumpOrder []logic.Card = []logic.Card{
	*logic.NewCard(logic.Hearts, 10),

	*logic.NewCard(logic.Clubs, logic.Queen),
	*logic.NewCard(logic.Spades, logic.Queen),
	*logic.NewCard(logic.Hearts, logic.Queen),
	*logic.NewCard(logic.Diamonds, logic.Queen),

	*logic.NewCard(logic.Clubs, logic.Jack),
	*logic.NewCard(logic.Spades, logic.Jack),
	*logic.NewCard(logic.Hearts, logic.Jack),
	*logic.NewCard(logic.Diamonds, logic.Jack),

	*logic.NewCard(logic.Diamonds, logic.Ace),
	*logic.NewCard(logic.Diamonds, 10),
	*logic.NewCard(logic.Diamonds, logic.King),
	*logic.NewCard(logic.Diamonds, 9),
}

// trumpValue returns the trump value for this card
// Hearts 10 returns 13, diamonds 9 returns 1.
// If the card is not a trump, -1 is returned
func (d *Doko) trumpValue(c *logic.Card) int {
	var value int
	var trump logic.Card
	for value, trump = range DokoTrumpOrder {
		if trump == *c {
			return len(DokoTrumpOrder) - value
		}
	}
	return -1
}

var dokoColorSortOrder []int = []int{-1, logic.Clubs, logic.Spades, logic.Hearts, logic.Diamonds}

// Sort sorts the specified deck / hand according to the game's rules
func (d *Doko) Sort(deck *logic.Deck) {
	var arr []*logic.Card = []*logic.Card(*deck)
	sort.SliceStable(arr, func(a int, b int) bool {
		if d.sortColor(d.color(deck.Get(a))) == d.sortColor(d.color(deck.Get(b))) {
			return d.value(deck.Get(a)) > d.value(deck.Get(b))
		}
		return d.sortColor(d.color(deck.Get(a))) < d.sortColor(d.color(deck.Get(b)))
	})
}

func (d *Doko) sortColor(color int) int {
	var i, j int
	for i, j = range dokoColorSortOrder {
		if j == color {
			return i
		}
	}
	return i + 1
}
