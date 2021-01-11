package doko

import "testing"
import "github.com/scrouthtv/card-games/logic"

func (ds *DokoSim) TestHand(t *testing.T, player int, cards *logic.Deck) {
	if t != nil {
		t.Helper()
	}
	if !ds.doko.hands[player].Equal(cards) {
		t.Errorf("Player %d has wrong hand:", player)
		t.Logf("Expected: %s", cards)
		t.Logf("Got: %s", ds.doko.hands[player])
	}
}

func (ds *DokoSim) TestAllHands(t *testing.T, hands map[int]*logic.Deck) {
	if t != nil {
		t.Helper()
	}
	var player int
	var deck *logic.Deck
	for player, deck = range hands {
		ds.TestHand(t, player, deck)
	}
}

func (ds *DokoSim) TestWondeck(t *testing.T, player int, cards *logic.Deck) {
	if t != nil {
		t.Helper()
	}
	if !ds.doko.won[player].Equal(cards) {
		t.Errorf("Player %d has wrong cards won:", player)
		t.Logf("Expected: %s", cards)
		t.Logf("Got: %s", ds.doko.won[player])
	}
}

func (ds *DokoSim) TestAllWondecks(t *testing.T, won map[int]*logic.Deck) {
	if t != nil {
		t.Helper()
	}
	var player int
	var deck *logic.Deck
	for player, deck = range won {
		ds.TestWondeck(t, player, deck)
	}
}

func (ds *DokoSim) TestTable(t *testing.T, table *logic.Deck) {
	if t != nil {
		t.Helper()
	}
	if !ds.doko.table.Equal(table) {
		t.Errorf("Table contents is wrong:")
		t.Logf("Expected: %s", table)
		t.Logf("Got: %s", ds.doko.table)
	} else {
		t.Logf("Table contents is as expected")
	}
}

func (ds *DokoSim) TestProgress(t *testing.T, shouldp int, shouldt int) {
	if t != nil {
		t.Helper()
	}

	var isp, ist int = ds.doko.Progress()
	if shouldp != -1 && isp != shouldp {
		t.Errorf("Wrong progress, should be %d, is %d", shouldp, isp)
	}
	if shouldt != -1 && ist != shouldt {
		t.Errorf("Wrong progress, should be %d, is %d", shouldt, ist)
	}
}

func testBeat(t *testing.T, doko *Doko, def string, atk string, exp bool) {
	var cdef, catk *logic.Card
	var ok bool
	ok, cdef = logic.CardFromShort(def)
	if !ok {
		return
	}
	ok, catk = logic.CardFromShort(atk)
	if !ok {
		return
	}
	var is bool = doko.beats(cdef, catk)
	if exp != is {
		if exp {
			t.Errorf("Attacker didn't win, they should have with %s vs %s",
				def, atk)
		} else {
			t.Errorf("Attacker did win, they shouldn't have with %s vs %s",
				def, atk)
		}
	}
}
