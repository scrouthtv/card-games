package main

import "testing"

func TestCarcConsts(t *testing.T) {
	if ace != 1 {
		t.Error("Ace != 1")
	}
	if jack != 11 {
		t.Error("Jack != 11")
	}
	if queen != 12 {
		t.Error("Queen != 12")
	}
	if king != 13 {
		t.Error("King != 13")
	}
}

func TestDokoGenerator(t *testing.T) {
	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13})
	deckContains(t, doko, Card{clubs, jack}, true)
	deckContains(t, doko, Card{clubs, ace}, true)
	deckContains(t, doko, Card{clubs, 3}, false)
	deckContains(t, doko, Card{diamonds, 9}, true)
	deckContains(t, doko, Card{diamonds, queen}, true)
	deckContains(t, doko, Card{diamonds, 8}, false)
	t.Log("All cards sucessfully tested")
}

func deckContains(t *testing.T, deck *Deck, card Card, should bool) {
	var c Card
	for _, c = range *deck {
		if (c == card) == should {
			return
		}
	}

	t.Errorf("Deck should %t contain card %s, it does not", should, card.String())
}
