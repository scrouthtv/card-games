package doko

import (
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

func TestFoxFinder(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	doko.Reset() // deal cards
	var fox scoring = newFox()

	var foxes []*logic.Card = fox.MarkCards(doko)
	if len(foxes) != 2 {
		t.Errorf("Expected 2 foxes, got %d", len(foxes))
	}

	var i int
	var c *logic.Card
	for i, c = range foxes {
		if c.Suit() != logic.Diamonds || c.Value() != logic.Ace {
			t.Logf("Here, have an i: %d", i)
		}
	}
}
