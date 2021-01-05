package main

import (
	"testing"
	"time"
)

func TestCarcConsts(t *testing.T) {
	if Ace != 1 {
		t.Error("Ace != 1")
	}
	if Jack != 11 {
		t.Error("Jack != 11")
	}
	if Queen != 12 {
		t.Error("Queen != 12")
	}
	if King != 13 {
		t.Error("King != 13")
	}
}

func TestDeckSerialization(t *testing.T) {
	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	var c *Card
	var copy *Card
	var str string
	var ok bool
	for _, c = range *doko {
		str = c.Short()
		ok, copy = CardFromShort(str)
		if !ok {
			t.Errorf("Could not parse %s", str)
		}
		if copy != c {
			t.Errorf("Card %s serialized to %s deserialized to %s", c.String(), str, copy.String())
		}
	}
}

func TestDeckDistribution(t *testing.T) {
	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	var dist [][]*Card = doko.DistributeAll(4)
	t.Logf("Distributed %s to ", doko.String())
	t.Log("1: ", dist[0])
	t.Log("2: ", dist[1])
	t.Log("3: ", dist[2])
	t.Log("4: ", dist[3])
}

func TestDeckShuffle(t *testing.T) {
	var unshuf1 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice()
	var unshuf2 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice()
	var shuf1 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	time.Sleep(time.Millisecond) // have to wait some time for the rng to get a new seed
	var shuf2 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	time.Sleep(time.Millisecond)
	var shuf3 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()
	time.Sleep(time.Millisecond)
	var shuf4 *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice().Shuffle()

	if !unshuf1.Equal(unshuf2) {
		t.Error("The unshuffled decks are not equal")
	}

	var i, j int = sliceContainsEqualMembers([]*Deck{unshuf1, shuf1, shuf2, shuf3, shuf4})
	if i != -1 {
		t.Errorf("Decks #%d and #%d are equal, they shouldn't be:", i, j)
		t.Error(shuf1)
		t.Error(shuf2)
	}
}

func sliceContainsEqualMembers(decks []*Deck) (int, int) {
	var i, j int
	for i = 0; i < len(decks)-1; i++ {
		for j = i + 1; j < len(decks); j++ {
			if decks[i].Equal(decks[j]) {
				return i, j
			}
		}
	}
	return -1, -1
}

func TestDokoGenerator(t *testing.T) {
	var doko *Deck = NewDeck([]int{1, 9, 10, 11, 12, 13}).Twice()
	var should map[Card]int = make(map[Card]int)
	should[Card{Clubs, Jack}] = 2
	should[Card{Clubs, 3}] = 0
	should[Card{Diamonds, 9}] = 2
	should[Card{Diamonds, 8}] = 0
	should[Card{Hearts, King}] = 2
	should[Card{Hearts, 5}] = 0
	should[Card{Spades, Ace}] = 2
	should[Card{Spades, 7}] = 0

	var is map[Card]int = make(map[Card]int)
	var c *Card
	for _, c = range *doko {
		is[*c]++
	}

	var cs Card
	var shouldAmount int
	for cs, shouldAmount = range should {
		if is[cs] != shouldAmount {
			t.Errorf("Card %s should've appeared %d times, it did appear %d times",
				c.String(), shouldAmount, is[cs])
		}
	}
}
