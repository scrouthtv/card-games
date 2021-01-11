package doko

import (
	"bytes"

	"github.com/scrouthtv/card-games/logic"
)

const (
	// DokoGameUUID is a unique id for the dok (de-) serializer
	DokoGameUUID = 1
)

// WriteBinary writes the game's information
// relevant to the specified player to the buffer
func (d *Doko) WriteBinary(player int, buf *bytes.Buffer) {
	switch d.g.State() {
	case logic.StatePreparing:
		buf.WriteByte(logic.StatePreparing)
	case logic.StatePlaying:
		var playable byte = 0
		if d.playable {
			playable = 1
		}
		buf.WriteByte(logic.StatePlaying | (byte(d.active) << 2) | (byte(player) << 4) | (playable << 6))
		d.hands[player].WriteBinary(buf)
		d.table.WriteBinary(buf)
		var i int
		for i = 0; i < 4; i++ {
			buf.WriteByte(byte(d.won[i].Length()))
		}

		// Collect all special cards:
		var special []*logic.Card
		var feature scoring
		for _, feature = range d.features {
			special = append(special, feature.MarkCards(d)...)
		}

		// Collect the special cards for every player:
		var playerspecial map[int]*logic.Deck = make(map[int]*logic.Deck, 4)
		for i = 0; i < 4; i++ {
			playerspecial[i] = logic.EmptyDeck()
		}
		var c *logic.Card
		var winner int
		for _, c = range special {
			winner = d.whoWon(c)
			if winner != -1 {
				playerspecial[winner].AddAll(c)
			}
		}

		var player int
		var deck *logic.Deck
		for player, deck = range playerspecial {
			buf.WriteByte(byte(player))
			deck.WriteBinary(buf)
		}

	case logic.StateEnded:
		buf.WriteByte(logic.StateEnded)

		d.Scores().WriteBinary(buf)

		var player int
		var deck *logic.Deck
		for player, deck = range d.won {
			buf.WriteByte(byte(player))
			deck.WriteBinary(buf)
		}
	default:
		buf.WriteByte(255)
	}

}

func (s *DokoScore) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte(byte(len(s.scores)))
	var score int
	for _, score = range s.scores {
		buf.WriteByte(byte(score))
	}
	writeIntArray(s.rereasons, buf)
	writeIntArray(s.contrareasons, buf)
}

func writeIntArray(arr []int, buf *bytes.Buffer) {
	buf.WriteByte(byte(len(arr)))

	var x int
	for _, x = range arr {
		buf.WriteByte(byte(x))
	}
}
