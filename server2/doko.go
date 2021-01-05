package main

import "log"

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
	var dist [][]Card = doko.DistributeAll(4)

	var i int
	for i = 0; i < len(dist); i++ {
		d.g.hands[i] = NewInventory(CardsToItems(dist[i]))
	}
	d.g.table = NewInventory([]Item{})

	return true
}

// PlayerMove applies the move specified by the given packet to this game
// and returns whether the action was successful
func (d *Doko) PlayerMove(p *Packet) bool {
	switch p.Action() {
	case "card":
		if len(p.Args()) < 1 {
			return false
		}
		var i, j int = d.g.hands[d.active].ItemIndex(&Card{1, 1})
		if i == -1 {
			return false
		}
		log.Printf("Using card @ #%d", j)
	}

	return false
}
