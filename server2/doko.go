package main

import (
	"log"
)

// Doko is the ruleset for Doppelkopf
type Doko struct {
	g      IGame
	active int

	// start: maps #player to their initial inventory
	start map[int]*Deck
	// hands: maps #player to inventoray
	hands map[int]*Deck
	// won: maps #player to Deck, each trick they won
	won map[int]*Deck
	// table: inventory
	table *Deck
}

func dokoCardValue(c *Card) int {
	switch c.value {
	case Ace:
		return 11
	case 10:
		return 10
	case King:
		return 4
	case Queen:
		return 3
	case Jack:
		return 2
	}
	return 0
}

// NewDoko generates a new Doppelkopf ruleset hosted by the
// supplied game
func NewDoko(host IGame) *Doko {
	var d Doko = Doko{host, -1, nil, nil, nil, nil}
	d.Reset()
	return &d
}

// Reset resets this game by clearing everything
// and giving all players a new hand
func (d *Doko) Reset() bool {
	d.start = make(map[int]*Deck)
	d.hands = make(map[int]*Deck)
	d.won = make(map[int]*Deck)

	var doko *Deck = NewDeck([]int{Ace, 9, 10, Jack, Queen, King}).Twice().Shuffle()
	var dist []*Deck = doko.DistributeAll(4)

	var i int
	for i = 0; i < len(dist); i++ {
		d.hands[i] = dist[i]
		d.start[i] = dist[i]
	}
	d.table = EmptyDeck()

	return true
}

// Start starts this game
func (d *Doko) Start() {
	d.active = 0
	d.g.SetState(StatePlaying)
}

// Info returns the GameInfo for this Doppelkopf game
func (d *Doko) Info() GameInfo {
	return GameInfo{
		d.g.ID(), d.g.Name(), "Doppelkopf", d.g.PlayerCount(), 4,
	}
}

// TypeID returns a UUID for this ruleset
func (d *Doko) TypeID() byte {
	return dokoGameUUID
}

// PlayerMove applies the move specified by the given packet to this game
// and returns whether the action was successful
func (d *Doko) PlayerMove(player int, p *Packet) bool {
	if player != d.active {
		log.Println("Ignoring because this player is not active, active: ", d.active)
		return false
	}

	switch p.Action() {
	case "card":
		// Check 1: is the move-requesting player currently active
		if d.g.State() != StatePlaying {
			log.Println("Ignoring because we are not playing")
			return false
		}

		// Check 2: is the request complete
		if len(p.Args()) < 1 {
			log.Println("Ignoring because no card was specified")
			return false
		}

		// Check 3: is the requested card a valid card
		var c *Card
		var ok bool
		ok, c = CardFromShort(p.Args()[0])
		if !ok {
			log.Println("Ignoring because invalid card was specified")
			return false
		}

		// Check 4: does the player own this card
		ok = d.hands[d.active].Remove(*c, 1) > 0
		if !ok {
			log.Println("Ignoring because the player does not own this card")
			return false
		}

		// Check 5: is this player allowed to play this card
		if d.table.Length() > 0 {
			test := Card{Hearts, 9}
			if *c == test {
				log.Println("=========================================")
				log.Println("Allowed cards:")
				log.Println("For player", d.active)
				log.Println(d.AllowedCards().Short())
			}
			if !d.AllowedCards().Contains(*c) {
				log.Println("Ignoring because the player is not allowed to play this card")
				return false
			}
		}

		d.table.AddAll(c)
		if len(*d.table) == 4 {
			var winner int = d.trickWinner(d.table)

			// d.active placed the last card, d.active + 1 placed the first card
			// winner is the # in the trick, not the # in the player array
			// suppose the first placer won, then winner is 0, but if 3 placed
			// the first card, they won:
			winner = winner + d.active + 1
			if winner >= 4 {
				winner -= 4
			}
			d.won[winner].Merge(d.table)
			d.table = EmptyDeck()
			d.active = winner
			if len(*d.hands[d.active]) == 0 {
				d.g.SetState(StateEnded)
			}
		} else {
			d.active++
		}

		d.g.SendUpdates()
		return true
	}

	return false
}

// AllowedCards determines which cards the active player is currently
// allowed to play (e. g. if they have to show a color or don't own
// that color)
func (d *Doko) AllowedCards() *Deck {
	log.Println("calculating allowed cards for active player", d.active)
	if d.table.Length() == 0 {
		return d.hands[d.active]
	}

	var show *Card = d.table.Get(0)
	var allowed *Deck = EmptyDeck()
	var has *Deck = d.hands[d.active]

	log.Println("Have to show", show)
	log.Printf("I (%d) own %s", d.active, has.String())

	var i int
	for i = 0; i < has.Length(); i++ {
		var ownedCard *Card = has.Get(i)
		if d.color(ownedCard) == d.color(show) {
			allowed.AddAll(ownedCard)
		} else {
			/*log.Printf("color(owned) != color(show) <=> color(%s) != color(%s) <=> %d != %d",
			ownedCard.String(), show.String(),
			d.color(ownedCard), d.color(show))*/
		}
	}

	if allowed.Length() == 0 {
		log.Println("no cards allowed so far, all allowed")
		return d.hands[d.active]
	}

	log.Println("returning collected")
	return allowed
}

// Scores calculates the value for each player
// The value is the sum of the value of each card they earned
func (d *Doko) Scores() []int {
	var scores []int = make([]int, 4)
	var repair, contrapair []int = d.Teams()
	var recards, contracards *Deck

	var player int
	for _, player = range repair {
		recards.Merge(d.start[player])
	}
	for _, player = range contrapair {
		contracards.Merge(d.start[player])
	}

	var revalue = recards.Value(dokoCardValue)
	var contravalue = contracards.Value(dokoCardValue)

	for _, player = range repair {
		scores[player] = revalue
	}
	for _, player = range contrapair {
		scores[player] = contravalue
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
	var inv *Deck
	for i, inv = range d.start {
		if inv.Contains(Card{Clubs, Queen}) {
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
