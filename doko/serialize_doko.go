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
		buf.WriteByte(byte(player))
	case logic.StatePlaying:
		var playable byte = byte(d.playingState)
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

		var p int
		var deck *logic.Deck
		for p = 0; p < 4; p++ {
			deck = playerspecial[p]
			buf.WriteByte(byte(p))
			deck.WriteBinary(buf)
		}

		buf.WriteByte(d.lastStamps[player])
		d.lastStamps[player] = d.currentStamp()
		buf.WriteByte(d.lastStamps[player])
		var a action
		buf.WriteByte(byte(len(d.actionQueue[player])))
		for _, a = range d.actionQueue[player] {
			a.WriteBinary(buf)
		}
		d.actionQueue[player] = []action{}

	case logic.StateEnded:
		buf.WriteByte(logic.StateEnded)

		buf.WriteByte(byte(player))

		var reteam []int
		reteam, _ = d.Teams()
		buf.WriteByte(byte(len(reteam)))
		var player int
		for _, player = range reteam {
			buf.WriteByte(byte(player))
		}

		d.Scores().WriteBinary(buf)

		var deck *logic.Deck
		for player = 0; player < 4; player++ {
			deck = d.won[player]
			buf.WriteByte(byte(player))
			deck.WriteBinary(buf)
		}

		// Collect all special cards:
		var special []*logic.Card
		var feature scoring
		for _, feature = range d.features {
			special = append(special, feature.MarkCards(d)...)
		}

		// Collect the special cards for every player:
		var playerspecial map[int]*logic.Deck = make(map[int]*logic.Deck, 4)
		var i int
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

		for player = 0; player < 4; player++ {
			deck = playerspecial[player]
			buf.WriteByte(byte(player))
			deck.WriteBinary(buf)
		}

	default:
		buf.WriteByte(255)
	}

}

// WriteBinary writes a score's values and reasons to a buffer
func (s *Score) WriteBinary(buf *bytes.Buffer) {
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

func (d *Doko) currentStamp() byte {
	var progress, total int = d.Progress()
	return byte(progress * 255 / total)
}
