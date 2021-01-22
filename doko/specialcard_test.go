package doko

import (
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

func TestFoxMeta(t *testing.T) {
	var fox scoring = newFox()

	if fox.Name() != "Fuchs" {
		t.Errorf("Fox name should be Fuchs, is %s", fox.Name())
	}

	if fox.Reason() != ReasonFox {
		t.Errorf("Fox reason should be %d, is %d", ReasonFox, fox.Reason())
	}
}

func TestFoxFinder(t *testing.T) {
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	doko.Reset() // deal cards
	var ds *DokoSim = &DokoSim{doko}
	var fox scoring = newFox()

	ds.doko.Start()

	ds.doko.hands[0] = logic.DeserializeDeck("hk, ca, h10, ca, d10, h9, d9, s10, cq, cj, dj, dk")
	ds.doko.hands[1] = logic.DeserializeDeck("c10, sa, dj, c10, sq, ck, ck, h9, dq, hj, sq, sa")
	ds.doko.hands[2] = logic.DeserializeDeck("cq, sk, sj, da, s10, s9, dq, ha, hq, hj, d10, dk")
	ds.doko.hands[3] = logic.DeserializeDeck("sj, h10, sk, d9, hk, ha, hq, s9, da, c9, c9, cj")

	var i int
	for i = 0; i < 4; i++ {
		doko.start[i] = doko.hands[i].Clone()
	}

	for i = 0; i < 4; i++ {
		doko.PlayerMove(doko.active, logic.NewPacket("call healthy"))
	}

	var rescore, contrascore int

	t.Run("1. Expecting 2 foxes", func(t *testing.T) {
		var foxes []*logic.Card = fox.MarkCards(doko)
		if len(foxes) != 2 {
			t.Errorf("Expected 2 foxes, got %d", len(foxes))
		}

		if doko.teamsKnown() {
			t.Error("Teams should not be known for now")
		}

		var c *logic.Card
		for _, c = range foxes {
			if c.Suit() != logic.Diamonds || c.Value() != logic.Ace {
				t.Errorf("This card is not a fox: %s", c.Short())
			}
		}

		rescore, contrascore = fox.Score(doko)

		if rescore != 0 {
			t.Error("Re team should have 0 points for foxes")
		}
		if contrascore != 0 {
			t.Error("Contra team should have 0 points for foxes")
		}
	})

	t.Run("2. First fox caught by friendlies", func(t *testing.T) {
		ds.assertCardMove(t, "cq", true)
		ds.assertCardMove(t, "dj", true)
		ds.assertCardMove(t, "da", true)
		t.Log("2 plays the first fox")
		ds.assertCardMove(t, "d9", true)
		ds.assertPickup(t, 0, true)

		if *doko.won[0] == nil {
			t.Error("Player 0 should have gotten this trick")
		}
		var won int = len(*doko.won[0])
		if won != 4 {
			t.Errorf("Player 0 should have won 4 cards by now, instead %d", won)
		}
		t.Log("0 catches 2's fox")
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
		ds.assertCardMove(t, "d10", true)
		ds.assertCardMove(t, "hj", true)
		ds.assertCardMove(t, "cq", true)
		ds.assertCardMove(t, "sj", true)
		ds.assertPickup(t, ds.doko.active, true)

		t.Logf("Friendlies are known: %t", doko.teamsKnown())

		// By now, only one fox is still relevant
		var foxes int = len(fox.MarkCards(doko))
		if foxes != 1 {
			t.Errorf("Wrong amount of foxes marked, got %d, should be 1", foxes)
		}
	})

	t.Run("5. The other fox is played", func(t *testing.T) {
		ds.assertCardMove(t, "dk", true)
		ds.assertCardMove(t, "da", true)
		ds.assertCardMove(t, "h10", true)
		ds.assertCardMove(t, "dq", true)
		ds.assertPickup(t, ds.doko.active, true)
	})

	t.Run("6. Expecting one fox", func(t *testing.T) {
		var foxes int = len(fox.MarkCards(doko))
		if foxes != 1 {
			t.Errorf("Wrong amount of foxes marked, got %d, should be 1", foxes)
		}
	})

	rescore, contrascore = fox.Score(doko)
	if rescore != 1 {
		t.Errorf("Wrong score for re team: %d", rescore)
	}
	if contrascore != 0 {
		t.Errorf("Wrong score for contra team: %d", contrascore)
	}
}
