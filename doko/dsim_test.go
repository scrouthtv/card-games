package doko

import "strings"
import "strconv"
import "testing"

import "github.com/scrouthtv/card-games/logic"

// GameStub is a game that is not connected to a hub or clients,
// but instead saves no data besides the current state
type GameStub struct {
	state byte
}

// ID returns 5
func (g *GameStub) ID() byte {
	return 5
}

// Name returns " ## Invalid Test Game ## "
func (g *GameStub) Name() string {
	return " ## Invalid Test Game ## "
}

// State returns the current state
func (g *GameStub) State() byte {
	return g.state
}

// SetState sets the current state
func (g *GameStub) SetState(state byte) {
	g.state = state
}

// PlayerCount returns 4
func (g *GameStub) PlayerCount() int {
	return 4
}

// SendUpdates does nothing
func (g *GameStub) SendUpdates() {

}

type DokoSim struct {
	doko *Doko
}

func (ds *DokoSim) Move(move string) bool {
	var p *logic.Packet = logic.NewPacket(move)

	return ds.doko.PlayerMove(ds.doko.active, p)
}

func (ds *DokoSim) String() string {
	var out strings.Builder
	out.WriteString("Current Player: ")
	out.WriteString(strconv.Itoa(ds.doko.active))

	var i int
	var deck *logic.Deck
	for i = 0; i < len(ds.doko.hands); i++ {
		deck = ds.doko.hands[i]
		out.WriteString("\nHand ")
		out.WriteString(strconv.Itoa(i))
		out.WriteString(": ")
		out.WriteString(deck.Short())
	}

	out.WriteString("\nTable: ")
	out.WriteString(ds.doko.table.Short())

	return out.String()
}

func (ds *DokoSim) addCardByShort(d *logic.Deck, short string) {
	var c *logic.Card
	var ok bool
	ok, c = logic.CardFromShort(short)
	if ok {
		d.AddAll(c)
	}
}

// playTrick plays a full trick. short specifies the cards to be played
// in 4x whitespace-delimited cards
func (ds *DokoSim) playTrick(t *testing.T, short string) {
	t.Helper()
	var cards []string = strings.Split(short, " ")
	if len(cards) != 4 {
		t.Error("Wrong amount of cards specified")
	}
	var move string
	for _, move = range cards {
		ds.assertCardMove(t, move, true)
	}
	ds.assertPickup(t, ds.doko.active, true)
}

func (ds *DokoSim) assertCardMove(t *testing.T, short string, exp bool) bool {
	if t != nil {
		t.Helper()
	}
	var ok bool = ds.Move("card " + short)
	if ok != exp {
		if ok {
			t.Errorf("Card %s did succeed, it shouldn't have", short)
			t.FailNow()
		} else {
			t.Errorf("Card %s didn't succeed, it should have", short)
			t.FailNow()
		}
		return false
	}
	return true
}

func (ds *DokoSim) assertPickup(t *testing.T, player int, exp bool) {
	if t != nil {
		t.Helper()
	}
	var p *logic.Packet = logic.NewPacket("pickup")
	var ok bool = ds.doko.PlayerMove(player, p)
	if ok != exp {
		if ok {
			t.Error("Pickup did succeed, it shouldn't have")
		} else {
			t.Error("Pickup didn't succeed, it should have")
		}
	}
}

// Makes the active player play the first allowed card
func (ds *DokoSim) playOnce(t *testing.T) {
	t.Helper()
	var card *logic.Card = ds.doko.AllowedCards().Get(0)
	ds.assertCardMove(t, card.Short(), true)
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
