package doko

import (
	"strconv"
	"strings"
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

var debug bool = false

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

func (ds *DokoSim) TestHand(t *testing.T, player int, cards *logic.Deck) {
	if !ds.doko.hands[player].Equal(cards) {
		t.Errorf("Player %d has wrong hand:", player)
		t.Logf("Expected: %s", cards)
		t.Logf("Got: %s", ds.doko.hands[player])
	}
}

func (ds *DokoSim) TestAllHands(t *testing.T, hands map[int]*logic.Deck) {
	var player int
	var deck *logic.Deck
	for player, deck = range hands {
		ds.TestHand(t, player, deck)
	}
}

func (ds *DokoSim) TestWondeck(t *testing.T, player int, cards *logic.Deck) {
	if !ds.doko.won[player].Equal(cards) {
		t.Errorf("Player %d has wrong cards won:", player)
		t.Logf("Expected: %s", cards)
		t.Logf("Got: %s", ds.doko.won[player])
	}
}

func (ds *DokoSim) TestAllWondecks(t *testing.T, won map[int]*logic.Deck) {
	var player int
	var deck *logic.Deck
	for player, deck = range won {
		ds.TestWondeck(t, player, deck)
	}
}

func (ds *DokoSim) TestTable(t *testing.T, table *logic.Deck) {
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
func TestStubGame(t *testing.T) {
	// SETUP:
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	var ds DokoSim = DokoSim{doko}

	ds.doko.Start()

	ds.doko.hands[0] = logic.DeserializeDeck("hk, ca, c10, ca, da, h9, d9, s10, cq, cj, dj, dk")
	ds.doko.hands[1] = logic.DeserializeDeck("c10, sa, dj, h10, sq, ck, ck, h9, dq, hj, sq, sa")
	ds.doko.hands[2] = logic.DeserializeDeck("cq, sk, sj, da, s10, s9, dq, ha, hq, hj, d10, dk")
	ds.doko.hands[3] = logic.DeserializeDeck("sj, h10, sk, d9, hk, ha, hq, s9, d10, c9, c9, cj")

	var i int
	for i = 0; i < 4; i++ {
		doko.start[i] = doko.hands[i].Clone()
	}

	var expectedHands map[int]*logic.Deck = make(map[int]*logic.Deck)
	var hand *logic.Deck
	for i, hand = range doko.hands {
		var cpy logic.Deck = *hand
		expectedHands[i] = &cpy
	}

	var expectedWon map[int]*logic.Deck = make(map[int]*logic.Deck)
	for i = 0; i < 4; i++ {
		expectedWon[i] = logic.EmptyDeck()
	}

	var expectedTable *logic.Deck
	var cpy logic.Deck = *doko.table
	expectedTable = &cpy

	t.Log(ds.String())

	var card *logic.Card = doko.hands[0].Get(0)
	var expCard logic.Card = *logic.NewCard(logic.Hearts, logic.King)
	var badCard logic.Card = *logic.NewCard(logic.Spades, logic.Ace)

	t.Run("1. Teams should be unknown", func(t *testing.T) {
		if doko.teamsKnown() {
			t.Errorf("Teams should be unknown to this point")
		}
	})

	t.Run("2. Player 0 fails sa", func(t *testing.T) {
		if *card != expCard {
			t.Errorf("First card if player 0 should be hk, is %s", card.Short())
		}

		ds.assertCardMove(t, badCard.Short(), false)
	})

	t.Run("3. Player 0 plays hk", func(t *testing.T) {
		expectedTable.AddAll(card)

		ds.assertCardMove(t, expCard.Short(), true)

		ds.TestTable(t, expectedTable)
	})

	if ds.doko.active != 1 {
		t.Error("Wrong player active")
		t.FailNow()
	}

	t.Run("4. Testing allowed cards", func(t *testing.T) {
		var allowed *logic.Deck = doko.AllowedCards()

		var allowExpected = logic.EmptyDeck()
		allowExpected.AddAll(logic.NewCard(logic.Hearts, 9))

		if !allowed.Equal(allowExpected) {
			t.Error("Wrong cards allowed")
		}
	})

	t.Run("5. Player 1 fails sa", func(t *testing.T) {
		ds.assertCardMove(t, "sa", false) // fail bc not allowed
	})

	t.Run("6. Player 1 fails hk", func(t *testing.T) {
		ds.assertCardMove(t, "hk", false) // fail bc not owned*/
	})

	t.Run("7. Player 1 plays h9", func(t *testing.T) {
		var c *logic.Card = logic.NewCard(logic.Hearts, 9)
		expectedTable.AddAll(c)

		ds.assertCardMove(t, "h9", true) // sucess

		ds.TestTable(t, expectedTable)
	})

	t.Run("8. Player 2 tests invalid cards", func(t *testing.T) {
		ds.assertCardMove(t, "", false)
		ds.assertCardMove(t, "s8", false)
		ds.assertCardMove(t, "d9", false)
	})

	t.Run("9. Player 1 tries to play", func(t *testing.T) {
		var p *logic.Packet = logic.NewPacket("card c10")

		if ds.doko.PlayerMove(1, p) {
			t.Error("Move did succeed, it shouldn't have")
		}
	})

	t.Run("10. Player 2 plays ha", func(t *testing.T) {
		var c *logic.Card = logic.NewCard(logic.Hearts, logic.Ace)
		expectedTable.AddAll(c)

		ds.assertCardMove(t, "ha", true)

		ds.TestTable(t, expectedTable)
	})

	t.Run("11. Player 3 fails h10", func(t *testing.T) {
		ds.assertCardMove(t, "h10", false)
	})

	t.Run("12. Player 3 plays hk", func(t *testing.T) {
		ds.assertCardMove(t, "hk", true)
	})

	t.Run("13. Test trick", func(t *testing.T) {
		// hk h9 ha hk, player 2 wins (0-based)
		ds.addCardByShort(expectedWon[2], "hk")
		ds.addCardByShort(expectedWon[2], "h9")
		ds.addCardByShort(expectedWon[2], "ha")
		ds.addCardByShort(expectedWon[2], "hk")
		if ds.doko.active != 2 {
			t.Error("Player 2 should be active, they aren't")
		}

		// Table should be empty:
		if ds.doko.table.Length() > 0 {
			t.Error("Table should be empty")
		}

		// Player 2 should be the only one with cards:
		ds.TestAllWondecks(t, expectedWon)
	})

	t.Run("14. Trick with #2's clubs queen", func(t *testing.T) {
		ds.assertCardMove(t, "cq", true)
		ds.playOnce()
		ds.playOnce()
		ds.playOnce()

		if ds.doko.won[2].Length() != 8 {
			t.Errorf("Player 2 has the wrong amount of cards won")
			t.Errorf("Is: %d, should be 8", ds.doko.won[2].Length())
		}

		if ds.doko.teamsKnown() {
			t.Error("Teams should not be known at this point")
		}
	})

	t.Run("15. Trick with the other clubs queen", func(t *testing.T) {
		ds.assertCardMove(t, "dk", true)
		ds.playOnce()
		ds.assertCardMove(t, "cq", true)

		if !ds.doko.teamsKnown() {
			t.Error("Teams should be known")
		}

		ds.playOnce()

		if !ds.doko.teamsKnown() {
			t.Error("Teams should stillbe known")
		}

	})

	//t.Log(ds.String())
}

func (ds *DokoSim) addCardByShort(d *logic.Deck, short string) {
	var c *logic.Card
	var ok bool
	ok, c = logic.CardFromShort(short)
	if ok {
		d.AddAll(c)
	}
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
