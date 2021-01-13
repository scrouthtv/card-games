package logic

import "testing"

func TestDeckClone(t *testing.T) {
	var d1 *Deck = NewDeck([]int{Ace, Jack, King})
	var d2 *Deck = d1
	var d3 *Deck = d1.Clone()

	var n int = len(d1.Remove(Card{Diamonds, Ace}, 1))
	if n != 1 {
		t.Errorf("Expected 1 card to be removed, got %d", n)
	}

	var i int
	var c *Card

	if d1.Length() != d2.Length() {
		t.Errorf("d1 and d2 are of different length: %d vs %d",
			d1.Length(), d2.Length())
	}

	for i, c = range *d1 {
		if c != d2.Get(i) {
			t.Errorf("d1(%d) should be equal to d2(%d), it isn't: %s vs %s",
				i, i, d1.Get(i).Short(), d2.Get(i).Short())
		}
	}

	if d1.Length() != d3.Length() - 1 {
		t.Errorf("d1 should have 1 element less than d2, it doesn't: %d vs %d", d1.Length(), d2.Length())
	}

	t.Run("Calling Twice", func(t *testing.T) {
		d1.Twice()

		if d1 != d2 {
			t.Error("instance of d1 shouldn't have changed")
		}

		if d1.Length() != (d3.Length() - 1) * 2 {
			t.Errorf("twiced deck has wrong length, is %d, should be %d (was before %d)",
				d1.Length(), (d3.Length() - 1) * 2, d3.Length() - 1)
		}
	})

	t.Run("Merging d1 and d3", func(t *testing.T){ 
		var pred3 *Deck = d3.Clone()
		d1.Merge(d3)
		if pred3.Short() != d3.Short() {
			t.Errorf("d3 has been changed by merging, should be \n%s is \n%s",
				pred3.Short(), d3.Short())
		}
	})
}

func TestStringSerialize(t *testing.T) {
	var d *Deck = NewDeck([]int{Ace, 7, 8, 9, 10, Jack, Queen, King})

	var short string = d.Short()

	// re-deck
	var red *Deck = DeserializeDeck(short)

	if !d.Equal(red) {
		t.Errorf("pre- and post- serialization decks should be same: \n%s vs \n%s", d.Short(), red.Short())
	}
}

func TestDeckValue(t *testing.T) {
	var skatValueFunc func(c *Card) int = func(c *Card) int {
		switch c.Value() {
		case Ace:
			return 11;
		case 10:
			return 10;
		case King:
			return 4;
		case Queen:
			return 3;
		case Jack:
			return 2;
		default:
			return 0;
		}
	}

	var skatDeck *Deck = NewDeck([]int{Ace, 7, 8, 9, 10, Jack, Queen, King})
	var skatValueIs int = skatDeck.Value(skatValueFunc)
	var skatValueSb int = 120
	if skatValueIs != skatValueSb {
		t.Errorf("Wrong value for a skat deck calculated: %d, should be %d",
			skatValueIs, skatValueSb)
	}
}
