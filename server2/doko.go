package main

// Doko is the ruleset for Doppelkopf
type Doko struct {
	g      *Game
	active int
}

// NewDoko generates a new Doppelkopf ruleset hosted by the
// supplied game
func NewDoko(host *Game) *Doko {
	var d Doko = Doko{host, -1}
	return &d
}

// Reset resets this game by clearing everything
// and giving all players a new hand
func (d *Doko) Reset() bool {
	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	var dist [][]*Card = doko.DistributeAll(4)

	var i int
	for i = 0; i < len(dist); i++ {
		d.g.hands[i] = NewInventory(CardsToItems(dist[i]))
	}
	d.g.table = NewInventory([]Item{})

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
		if len(p.Args()) < 1 {
			return false
		}
		var c *Card
		var ok bool
		ok, c = CardFromShort(p.Args()[0])
		if !ok {
			return false
		}
		ok = d.g.hands[d.active].RemoveItem(c)
		if !ok {
			return false
		}
		d.g.table.AddToSlot(0, c)

		if len(d.g.table.Get(0)) == 4 {

		}

		return true
	}

	return false
}

// trickWinner calculates the winner # in this trick
func (d *Doko) trickWinner(trick []*Card) int {
	var winner int = 0
	var wCard = trick[0]

	var i int
	for i = 1; i < len(trick); i++ {
		if d.beats(wCard, trick[i]) {
			winner = i
			wCard = trick[i]
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