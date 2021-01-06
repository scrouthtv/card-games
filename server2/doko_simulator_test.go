package main

import (
	"strconv"
	"strings"
	"testing"
)

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
	var cm clientMessage = clientMessage{nil, []byte(move)}
	var p *Packet = cm.toPacket()

	return ds.doko.PlayerMove(ds.doko.active, p)
}

func (ds *DokoSim) String() string {
	var out strings.Builder
	out.WriteString("Current Player: ")
	out.WriteString(strconv.Itoa(ds.doko.active))

	var i int
	var deck *Deck
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

func (ds *DokoSim) TestHand(t *testing.T, player int, cards *Deck) {
	if !ds.doko.hands[player].Equal(cards) {
		t.Errorf("Player %d has wrong hand:", player)
		t.Logf("Expected: %s", cards)
		t.Logf("Got: %s", ds.doko.hands[player])
	}
}

func (ds *DokoSim) TestAllHands(t *testing.T, hands map[int]*Deck) {
	var player int
	var deck *Deck
	for player, deck = range hands {
		ds.TestHand(t, player, deck)
	}
}

func (ds *DokoSim) TestTable(t *testing.T, table *Deck) {
	if !ds.doko.table.Equal(table) {
		t.Errorf("Table contents is wrong:")
		t.Logf("Expected: %s", table)
		t.Logf("Got: %s", ds.doko.table)
	} else {
		t.Logf("Table contents is as expected")
	}
}

// Example game:
// hk ca c10 ca da h9 d9 s10 cq cj dj dk
// c10 sa dj h10 sq ck ck h9 dq hj sq sa
// cq sk sj da s10 s9 dq ha hq hj d10 dk
// sj h10 sk d9 hk ha hq s9 d10 c9 c9 cj
// 1. Player 0 tries to play sa -> invalid
// 2. Player 0 plays hk
// 3. Player 1 tries to play sa -> invalid
// 4. Player 1 plays h9
// 5.

func TestStubGame(t *testing.T) {
	var gs *GameStub = &GameStub{StatePreparing}
	var doko *Doko = NewDoko(gs)
	var ds DokoSim = DokoSim{doko}

	ds.doko.Start()

	ds.doko.hands[0] = DeserializeDeck("hk, ca, c10, ca, da, h9, d9, s10, cq, cj, dj, dk")
	ds.doko.hands[1] = DeserializeDeck("c10, sa, dj, h10, sq, ck, ck, h9, dq, hj, sq, sa")
	ds.doko.hands[2] = DeserializeDeck("cq, sk, sj, da, s10, s9, dq, ha, hq, hj, d10, dk")
	ds.doko.hands[3] = DeserializeDeck("sj, h10, sk, d9, hk, ha, hq, s9, d10, c9, c9, cj")

	var expectedHands map[int]*Deck = make(map[int]*Deck)
	var i int
	var hand *Deck
	for i, hand = range doko.hands {
		var copy Deck = *hand
		expectedHands[i] = &copy
	}

	var expectedTable *Deck
	var copy Deck = *doko.table
	expectedTable = &copy

	t.Log(ds.String())

	var card *Card = doko.hands[0].Get(0)
	var expCard Card = Card{Hearts, King}
	var badCard Card = Card{Spades, Ace}
	if *card != expCard {
		t.Errorf("First card if player 0 should be hk, is %s", card.Short())
	}

	// 1. Player 0 tries to play sa -> invalid
	ds.assertCardMove(t, badCard.Short(), false)

	expectedTable.AddAll(card)
	t.Log("Player 0 is going to play", card.String(), card.Short())

	t.Log(ds.doko.table)

	// 2. Player 0 plays hk
	ds.assertCardMove(t, expCard.Short(), true)

	ds.TestTable(t, expectedTable)

	if ds.doko.active != 1 {
		t.Error("Wrong player active")
		t.FailNow()
	}
	t.Log("--------------")

	var allowed *Deck = doko.AllowedCards()
	t.Log("Player 1 can now play", allowed.Short())

	var allowExpected = EmptyDeck()
	allowExpected.AddAll(&Card{Hearts, 9})

	if !allowed.Equal(allowExpected) {
		t.Error("Wrong cards allowed")
	} else {
		t.Log("Correct cards allowed")
	}

	ds.assertCardMove(t, "sa", false) // fail bc not allowed
	ds.assertCardMove(t, "hk", false) // fail bc not owned
	ds.assertCardMove(t, "h9", true)  // sucess
}

func (ds *DokoSim) assertCardMove(t *testing.T, short string, exp bool) {
	var ok bool = ds.Move("card " + short)
	if ok != exp {
		if ok {
			t.Error("Move did succed, it shouldn't have")
			t.FailNow()
		} else {
			t.Error("Move didn't succed, it should have")
			t.FailNow()
		}
	}
}
