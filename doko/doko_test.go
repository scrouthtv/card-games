package doko

import (
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

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
