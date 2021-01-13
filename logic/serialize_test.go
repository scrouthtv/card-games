package logic

import (
	"bytes"
	"testing"
)

func TestSerializeCards(t *testing.T) {
	var cards map[Card]byte = make(map[Card]byte)
	cards[Card{Diamonds, Ace}] = byte(Ace * CardMaxSuit + Diamonds)
	cards[Card{Diamonds, 9}] = byte(9 * CardMaxSuit + Diamonds)
	cards[Card{Hearts, 4}] = byte(4 * CardMaxSuit + Hearts)
	cards[Card{Hearts, King}] = byte(King * CardMaxSuit + Hearts)
	cards[Card{Clubs, Ace}] = byte(Ace * CardMaxSuit + Clubs)
	cards[Card{Clubs, King}] = byte(King * CardMaxSuit + Clubs)
	cards[Card{Spades, Queen}] = byte(Queen * CardMaxSuit + Spades)
	cards[Card{Spades, 10}] = byte(10 * CardMaxSuit + Spades)

	var c Card
	var b, isb byte
	for c, b = range cards {
		isb = c.ToBinary()

		if b != isb {
			t.Errorf("Wrong binary representation for card %s: %d, should be %d", c.Short(), isb, b)
		}
	}
}

func TestSerializeDeck(t *testing.T) {
	var deck *Deck = DeserializeDeck("hk, da, d9, hk, hk, sq")

	var buf bytes.Buffer
	deck.WriteBinary(&buf)
	var b []byte = buf.Bytes()

	if len(b) != deck.Length() + 1 {
		t.Errorf("Wrong length of binary representation: %d, should be %d",
			len(b), deck.Length() + 1)
	}

	if b[0] != byte(deck.Length()) {
		t.Errorf("First byte should be deck length (%d), is %d",
			deck.Length(), b[0])
	}

	var i int
	var c *Card
	var isb, sbb byte
	for i, c = range *deck {
		isb = b[i + 1]
		sbb = c.ToBinary()
		if isb != sbb {
			t.Errorf("%d: Card %s got wrong binary: %d vs %d",
				i, c.Short(), isb, sbb)
		}
	}
}
