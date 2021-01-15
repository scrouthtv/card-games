// build !release

package doko

import "github.com/scrouthtv/card-games/logic"

// SetHands sets the player's hands
func (doko *Doko) SetHands(hands map[int]*logic.Deck) {
	doko.hands = hands;

	var i int
	var d *logic.Deck
	for i, d = range hands {
		doko.start[i] = d.Clone()
	}
}
