package main

import "log"

// Doko is the ruleset for Doppelkopf
type Doko struct {
	g      *Game
	active int

	// start: maps #player to their initial inventory
	start map[int]*Inventory
	// hands: maps #player to inventory
	hands map[int]*Inventory
	// won: maps #player to Deck
	won map[int]*Deck
	// table: inventory
	table *Inventory
}

// NewDoko generates a new Doppelkopf ruleset hosted by the
// supplied game
func NewDoko(host *Game) *Doko {
	var d Doko = Doko{host, -1, nil, nil, nil, nil}
	d.Reset()
	return &d
}

// Reset resets this game by clearing everything
// and giving all players a new hand
func (d *Doko) Reset() bool {
	d.start = make(map[int]*Inventory)
	d.hands = make(map[int]*Inventory)
	d.won = make(map[int]*Deck)

	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	var dist [][]*Card = doko.DistributeAll(4)

	var i int
	for i = 0; i < len(dist); i++ {
		d.hands[i] = NewInventory(dist[i])
		d.start[i] = NewInventory(dist[i])
	}
	d.table = NewInventory([]*Card{})

	return true
}

// PlayerMove applies the move specified by the given packet to this game
// and returns whether the action was successful
func (d *Doko) PlayerMove(player int, p *Packet) bool {
	if player != d.active {
		return false
	}

	switch p.Action() {
	case "card":
		if d.g.state != StatePlaying {
			return false
		}
		if len(p.Args()) < 1 {
			return false
		}
		var c *Card
		var ok bool
		ok, c = CardFromShort(p.Args()[0])
		if !ok {
			return false
		}
		ok = d.hands[d.active].Get(0).Remove(*c, 1) > 0
		if !ok {
			return false
		}
		d.table.AddToSlot(0, c)

		if len(*d.table.Get(0)) == 4 {
			log.Println("This trick is finished, calculating the winner:")
			var winner int = d.trickWinner(d.table.Get(0))

			// d.active placed the last card, d.active + 1 placed the first card
			// winner is the # in the trick, not the # in the player array
			// suppose the first placer won, then winner is 0, but if 3 placed
			// the first card, they won:
			winner = winner + d.active + 1
			if winner >= 4 {
				winner -= 4
			}
			d.won[winner].Merge(d.table.Get(0))
			d.table.ClearAll()
			d.active = winner
			if len(*d.hands[d.active]) == 0 {
				d.g.state = StateEnded
			}
		} else {
			d.active++
		}

		return true
	}

	return false
}

// Scores calculates the value for each player
// The value is the sum of the value of each card they earned
func (d *Doko) Scores() []int {
	var scores []int = make([]int, 4)
	var repair, contrapair []int = d.Teams()
	var recards, contracards *Deck

	var player int
	for _, player = range repair {
		recards.Merge(d.start[player].Get(0))
	}
	for _, player = range contrapair {
		contracards.Merge(d.start[player].Get(0))
	}

	return scores
}

// Teams returns the player teams,
// all re players are in the first array
// all contra players in the second array
// not always do both arrays have 2 ints (e. g. marriage)
func (d *Doko) Teams() ([]int, []int) {
	var repair, contrapair []int
	var i int
	var inv *Inventory
	for i, inv = range d.start {
		if inv.Get(0).Contains(Card{Clubs, Queen}) {
			repair = append(repair, i)
		} else {
			contrapair = append(contrapair, i)
		}
	}
	return repair, contrapair
}

func (d *Doko) containsColor(deck *Deck, color int) bool {
	return deck.ContainsAny(func(c *Card) bool {
		return d.color(c) == color
	})
}

// Hands returns a map that maps each player to their inventory
func (d *Doko) Hands() map[int]Inventory {
	var hands map[int]Inventory = make(map[int]Inventory)

	if true {
		panic("not impl")
	}

	return hands
}

// trickWinner calculates the winner # in this trick
func (d *Doko) trickWinner(trick *Deck) int {
	var winner int = 0
	var wCard = (*trick)[0]

	var i int
	for i = 1; i < trick.Length(); i++ {
		if d.beats(wCard, (*trick)[i]) {
			winner = i
			wCard = (*trick)[i]
		}
	}

	return winner
}

// beats calculates whether the attacking card atk defeats the defending card def
func (d *Doko) beats(def *Card, atk *Card) bool {
	if d.color(def) == d.color(atk) {
		if d.value(atk) > d.value(def) {
			return true
		} else if atk.value == def.value {
			return *def == Card{Hearts, 10}
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

var dokoValueOrder []int = []int{9, Jack, Queen, King, 10, Ace}

func (d *Doko) value(c *Card) int {
	var i, value int

	value = d.trumpValue(c)
	if value != -1 {
		// return trump value instead
		return value
	}

	for i, value = range dokoValueOrder {
		if value == c.value {
			return i
		}
	}
	return 0
}

// color returns the color if this card, returning -1 if the card is a trump
func (d *Doko) color(c *Card) int {
	if d.trumpValue(c) == -1 {
		return c.suit
	}
	return -1
}

var dokoTrumpOrder []Card = []Card{
	{Hearts, 10},
	{Clubs, Queen}, {Spades, Queen}, {Hearts, Queen}, {Diamonds, Queen},
	{Clubs, Jack}, {Spades, Jack}, {Hearts, Jack}, {Diamonds, Jack},
	{Diamonds, Ace}, {Diamonds, 10}, {Diamonds, King}, {Diamonds, 9},
}

// trumpValue returns the trump value for this card
// Hearts 10 returns 13, diamonds 9 returns 1.
// If the card is not a trump, -1 is returned
func (d *Doko) trumpValue(c *Card) int {
	var value int
	var trump Card
	for value, trump = range dokoTrumpOrder {
		if trump == *c {
			return len(dokoTrumpOrder) - value
		}
	}
	return -1
}
