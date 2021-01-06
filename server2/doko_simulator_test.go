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
	for i, deck = range ds.doko.hands {
		out.WriteString("\nHand ")
		out.WriteString(strconv.Itoa(i))
		out.WriteString(": ")
		out.WriteString(deck.String())
	}

	out.WriteString("\nTable: ")
	out.WriteString(ds.doko.table.String())

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
		t.Errorf("Table contents is has wrong:")
		t.Logf("Expected: %s", table)
		t.Logf("Got: %s", ds.doko.table)
	}
}

func TestStubGame(t *testing.T) {
	var gs *GameStub = &GameStub{StatePreparing}
	var doko *Doko = NewDoko(gs)
	var ds DokoSim = DokoSim{doko}

	ds.doko.Start()

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
	expectedTable.AddAll(card)
	t.Log("Player 0 is going to play", card.String(), card.Short())
	var ok bool = ds.Move("card " + card.Short())
	t.Log("Success:", ok)
	if !ok {
		t.Error("Move did not succeed, it should've")
	}
	ds.TestTable(t, expectedTable)

	t.Log(ds.String())
}
