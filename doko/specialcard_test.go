package doko

import (
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

func TestFoxFinder(t *testing.T) {
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	doko.Reset() // deal cards
	var ds *DokoSim = &DokoSim{doko}
	var fox scoring = newFox()

	ds.doko.Start()

	ds.doko.hands[0] = logic.DeserializeDeck("hk, ca, c10, ca, d10, h9, d9, s10, cq, cj, dj, dk")
	ds.doko.hands[1] = logic.DeserializeDeck("c10, sa, dj, h10, sq, ck, ck, h9, dq, hj, sq, sa")
	ds.doko.hands[2] = logic.DeserializeDeck("cq, sk, sj, da, s10, s9, dq, ha, hq, hj, d10, dk")
	ds.doko.hands[3] = logic.DeserializeDeck("sj, h10, sk, d9, hk, ha, hq, s9, da, c9, c9, cj")

	t.Run("1. Expecting 2 foxes", func(t *testing.T) {
		var foxes []*logic.Card = fox.MarkCards(doko)
		if len(foxes) != 2 {
			t.Errorf("Expected 2 foxes, got %d", len(foxes))
		}

		var c *logic.Card
		for _, c = range foxes {
			if c.Suit() != logic.Diamonds || c.Value() != logic.Ace {
				t.Errorf("This card is not a fox: %s", c.Short())
			}
		}
	})

	t.Run("2. First fox caught by friendlies", func(t *testing.T) {
		ds.assertCardMove(t, "da", true)
		ds.assertCardMove(t, "dj", true)
		ds.assertCardMove(t, "cq", true)
		ds.assertCardMove(t, "d9", true)

		var won int = len(*doko.won[2])
		if won != 4 {
			t.Errorf("Player 2 should have won 4 cards by now, instead %d", won)
		}
	})

	t.Run("3. Two foxes should be relevant", func(t *testing.T) {
		// Both foxes are still relevant because the
		// players don't know each other's friends
		var foxes int = len(fox.MarkCards(doko))
		if foxes != 2 {
			t.Errorf("Wrong amount of foxes marked, got %d, should be 2", foxes)
		}
	})

	t.Run("4. The other clubs queen is played", func(t *testing.T) {
		ds.assertCardMove(t, "dk", true)
		ds.playOnce()
		ds.assertCardMove(t, "cq", true)

		// By now, the one fox should be irrelevant
		var foxes int = len(fox.MarkCards(doko))
		if foxes != 1 {
			t.Errorf("Wrong amount of foxes marked, got %d, should be 1", foxes)
			t.Logf("0 and 2 are friends: %t", doko.IsFriend(0, 2))
		}

		ds.assertCardMove(t, "da", true)

		// But the other one is still relevant
		foxes = len(fox.MarkCards(doko))
		if foxes != 1 {
			t.Errorf("Wrong amount of foxes marked, got %d, should be 1", foxes)
		}
	})
}
