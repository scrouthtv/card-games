package doko

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

func TestMain(m *testing.M) {
	var exit int = m.Run()

	if exit == 0 && testing.CoverMode() != "" {
		var c float64 = testing.Coverage()
		if c < 0.8 {
			fmt.Printf("Tests passed but too little covered, got %2.0f %%, should be 80 %%\n", c*100)
			exit = -1
		}
	}

	os.Exit(exit)
}

func TestCardTracer(t *testing.T) {
	var doko *Doko = NewDoko(&GameStub{logic.StatePreparing})
	doko.Start()
	var ds DokoSim = DokoSim{doko}

	// clone the start decks first:
	var start []*logic.Deck = []*logic.Deck{
		logic.EmptyDeck(),
		logic.EmptyDeck(),
		logic.EmptyDeck(),
		logic.EmptyDeck()}

	var i int
	var deck *logic.Deck
	for i, deck = range doko.start {
		start[i].AddAll(*deck...)
	}

	// Play some cards
	for i = 0; i < 5; i++ {
		ds.playOnce(t)
		ds.playOnce(t)
		ds.playOnce(t)
		ds.playOnce(t)
		ds.assertPickup(t, ds.doko.active, true)
	}

	var j int
	var owner int
	var c *logic.Card
	// Check 1: Has the start array changed?
	for i, deck = range doko.start {
		if len(*deck) != len(*start[i]) {
			t.Errorf("Start cards for player %d changed length from %d to %d",
				i, len(*start[i]), len(*deck))
		}
		for j, c = range *deck {
			if c != start[i].Get(j) {
				t.Errorf("Card %d on player %d changed", j, i)
			}
			// Check 2: Is the tracing functionality implemented correctly?
			owner = doko.origOwner(c)
			if owner != i {
				t.Errorf("Wrong owner determing for card %d orignally on player %d, got %d instead",
					j, i, owner)
			}
		}

	}

	// Check 3: Do all three functions return -1 on wrong inputs?
	c = logic.NewCard(logic.Spades, logic.Ace)
	owner = doko.origOwner(c)
	if owner != -1 {
		t.Errorf("Determined an owner for a card that shouldn't have one")
	}
	owner = doko.whoWon(c)
	if owner != -1 {
		t.Errorf("Determined a winner for a card that shouldn't have one")
	}
	owner = doko.whenWon(c)
	if owner != -1 {
		t.Errorf("Determined a time of winning for a card that shouldn't have one")
	}

	// Check 4: Does the whoWon() work?
	var winner int
	for i, deck = range doko.won {
		if deck == nil {
			continue
		}
		for j, c = range *deck {
			winner = doko.whoWon(c)
			if winner != i {
				t.Errorf("Wrong winner for card %d, should be %d, is %d", j, i, winner)
			}
			winner = doko.whenWon(c)
			if winner != j {
				t.Errorf("Wrong time of winning for card %d, should be %d, is %d", j, j, winner)
			}
		}
	}

}

func TestFriends(t *testing.T) {
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	//var ds *DokoSim = &DokoSim{doko}

	doko.Start()

	doko.hands[0] = logic.DeserializeDeck("hk, ca, c10, ca, d10, h9, d9, s10, cq, cj, dj, dk")
	doko.hands[1] = logic.DeserializeDeck("c10, sa, dj, h10, sq, ck, ck, h9, dq, hj, sq, sa")
	doko.hands[2] = logic.DeserializeDeck("cq, sk, sj, da, s10, s9, dq, ha, hq, hj, d10, dk")
	doko.hands[3] = logic.DeserializeDeck("sj, h10, sk, d9, hk, ha, hq, s9, da, c9, c9, cj")

	var i int
	for i = 0; i < 4; i++ {
		doko.start[i] = doko.hands[i].Clone()
	}

	if !doko.IsFriend(0, 2) {
		t.Error("0 and 2 should be friends")
	}

	if !doko.IsFriend(1, 1) {
		t.Error("1 and 1 should be friends")
	}

	if doko.IsFriend(0, 1) {
		t.Error("0 and 1 should be friends")
	}

	(&DokoSim{doko}).TestFriends(t, []int{0, 2}, []int{1, 3})
}

func (ds *DokoSim) TestFriends(t *testing.T, re []int, contra []int) {
	sort.SliceStable(re, func(i int, j int) bool {
		return re[i] > re[j]
	})
	sort.SliceStable(contra, func(i int, j int) bool {
		return contra[i] > contra[j]
	})
	var isre, iscontra []int = ds.doko.Teams()
	sort.SliceStable(isre, func(i int, j int) bool {
		return isre[i] > isre[j]
	})
	sort.SliceStable(iscontra, func(i int, j int) bool {
		return iscontra[i] > iscontra[j]
	})

	if len(re) != len(isre) {
		t.Errorf("Re team should have %d players, does have %d players",
			len(re), len(isre))
	}
	if len(contra) != len(iscontra) {
		t.Errorf("Contra team should have %d players, does have %d players",
			len(contra), len(iscontra))
	}

	var i, k int
	for i, k = range re {
		if k != isre[i] {
			t.Errorf("Re team should have player %d, does have player %d",
				k, isre[i])
		}
	}
	for i, k = range contra {
		if k != iscontra[i] {
			t.Errorf("Contra team should have player %d, does have player %d",
				k, iscontra[i])
		}
	}
}

func TestDokoSum(t *testing.T) {
	var deck *logic.Deck = logic.NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice()
	var value int = deck.Value(dokoCardValue)

	if value != 240 {
		t.Errorf("Doko deck does not have a value of 240, it has %d instead", value)
	}
}

func TestTricks(t *testing.T) {
	var doko *Doko = NewDoko(nil)

	var tricks map[string]int = make(map[string]int)
	tricks["sa, sk, sa, s9"] = 0
	tricks["s10, sk, ha, s9"] = 0
	tricks["s10, da, ha, d9"] = 1  // they have the higher trump
	tricks["da, s9, d9, dk"] = 0   // they have the highest trump
	tricks["da, h10, da, h10"] = 3 // they played the second hearts 10*/

	var trick string
	var should, is int
	var deck *logic.Deck
	for trick, should = range tricks {
		deck = logic.DeserializeDeck(trick)
		is = doko.trickWinner(deck)
		if should != is {
			t.Errorf("Trick %s should win #%d, but instead #%d", trick, should, is)
		}
	}
}

func TestSort(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	var deck *logic.Deck = logic.NewDeck([]int{logic.Ace, 9, 10, logic.Jack, logic.Queen, logic.King}).Twice().Shuffle()
	doko.Sort(deck)
	t.Log(deck.Short())
}

func TestBeat(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	testBeat(t, doko, "da", "ha", false)
	testBeat(t, doko, "d9", "da", true)
	testBeat(t, doko, "d9", "d9", false)
	testBeat(t, doko, "h10", "h10", true)
	testBeat(t, doko, "da", "da", false)
	testBeat(t, doko, "da", "h10", true)
	testBeat(t, doko, "sa", "dk", true)
	testBeat(t, doko, "ha", "h10", true)
	testBeat(t, doko, "h10", "d10", false)
}

func TestInfo(t *testing.T) {
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	var info logic.GameInfo = doko.Info()
	if info.Maxplayers != 4 {
		t.Error("Doko game should have max players 4")
	}
}
