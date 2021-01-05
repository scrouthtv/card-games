package main

import "testing"

func TestValues(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	var cards []Card = []Card{
		{Hearts, 10}, {Hearts, Jack}, {Diamonds, 9},
		{Spades, King}, {Hearts, 9},
	}
	var card Card
	for _, card = range cards {
		t.Logf("Card %s has %d trump value", card.String(), doko.trumpValue(&card))
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
	var deck *Deck
	var cards []*Card
	for trick, should = range tricks {
		deck = DeserializeDeck(trick)
		cards = []*Card(*deck)
		is = doko.trickWinner(cards)
		if should != is {
			t.Errorf("Trick %s should win #%d, but instead #%d", trick, should, is)
		}
	}
}

func TestBeat(t *testing.T) {
	var doko *Doko = NewDoko(nil)
	t.Logf("%t", doko.beats(&Card{Diamonds, Ace}, &Card{Hearts, Ace}))   // false
	t.Logf("%t", doko.beats(&Card{Diamonds, 9}, &Card{Diamonds, Ace}))   // true
	t.Logf("%t", doko.beats(&Card{Hearts, 10}, &Card{Hearts, 10}))       // true
	t.Logf("%t", doko.beats(&Card{Diamonds, Ace}, &Card{Diamonds, Ace})) // false
	t.Logf("%t", doko.beats(&Card{Diamonds, Ace}, &Card{Hearts, 10}))    // true
}