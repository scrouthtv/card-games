// build !release

package doko

import "github.com/scrouthtv/card-games/logic"

// SetHands sets the player's hands
func (d *Doko) SetHands(hands map[int]*logic.Deck) {
	d.hands = hands;
}
