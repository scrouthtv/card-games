package doko

import (
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

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
		ds.playOnce()
		ds.playOnce()
		ds.playOnce()
		ds.playOnce()
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

	// Check 3: Does the whoWon() work?
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

	//t.Log(ds.String())

	t.Log(doko.Teams())

	if !doko.IsFriend(0, 2) {
		t.Error("0 and 2 should be friends")
	}
}

// Makes the active player play the first allowed card
func (ds *DokoSim) playOnce() {
	var card *logic.Card = ds.doko.AllowedCards().Get(0)
	ds.assertCardMove(nil, card.Short(), true)
}

func TestValues(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	var cards []logic.Card = []logic.Card{
		*logic.NewCard(logic.Hearts, 10),
		*logic.NewCard(logic.Hearts, logic.Jack),
		*logic.NewCard(logic.Diamonds, 9),
		*logic.NewCard(logic.Spades, logic.King),
		*logic.NewCard(logic.Hearts, 9),
	}
	var card logic.Card
	for _, card = range cards {
		t.Logf("Card %s has %d trump value", card.String(), doko.trumpValue(&card))
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

func testBeat(t *testing.T, doko *Doko, def string, atk string, exp bool) {
	var cdef, catk *logic.Card
	var ok bool
	ok, cdef = logic.CardFromShort(def)
	if !ok {
		return
	}
	ok, catk = logic.CardFromShort(atk)
	if !ok {
		return
	}
	var is bool = doko.beats(cdef, catk)
	if exp != is {
		if exp {
			t.Errorf("Attacker didn't win, they should have with %s vs %s",
				def, atk)
		} else {
			t.Errorf("Attacker did win, they shouldn't have with %s vs %s",
				def, atk)
		}
	}
}
