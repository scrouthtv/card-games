package doko

import (
	"bytes"
	"testing"

	"github.com/scrouthtv/card-games/logic"
)

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

	// Check 1: Try to play a card now:
	t.Run("0. Shouldn't be able to do anything before the game has started", func(t *testing.T) {
		ds.assertCardMove(t, "h10", false)
		var p *logic.Packet = logic.NewPacket("pickup")
		var ok bool = ds.doko.PlayerMove(ds.doko.active, p)
		if ok {
			t.Error("Shouldn't be able to pick up now")
		}
		ok = ds.doko.PlayerMove(2, p)
		if ok {
			t.Error("Player 2 shouldn't be able to pick up now either")
		}

		ds.TestProgress(t, 0, 48)
	})

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

	t.Run("5. Trick pickup fails", func(t *testing.T) {
		ds.assertPickup(t, 0, false)
		ds.assertPickup(t, 1, false)
		ds.assertPickup(t, 2, false)
		ds.assertPickup(t, 3, false)
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

		var p *logic.Packet = logic.NewPacket("card")
		var ok bool = ds.doko.PlayerMove(ds.doko.active, p)
		if ok {
			t.Error("Shouldn't be able to play empty card")
		}
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

	t.Run("13. Succeed without pickup", func(t *testing.T) {
		if ds.doko.active != 2 {
			t.Error("Player 2 should be active as they won the trick")
		}
		// No cards won so far:
		ds.TestAllWondecks(t, expectedWon)

		// We can't play:
		ds.assertCardMove(t, "cq", false)

		// Only 2 can pick up:
		ds.assertPickup(t, 0, false)
		ds.assertPickup(t, 1, false)
		ds.assertPickup(t, 3, false)
		ds.assertPickup(t, 2, true)
		if ds.doko.active != 2 {
			t.Error("Player 2 should still be active")
		}
	})

	t.Run("14. Test trick", func(t *testing.T) {
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

	t.Run("15. Trick with #2's clubs queen", func(t *testing.T) {
		ds.assertCardMove(t, "cq", true)
		ds.playOnce(t)
		ds.playOnce(t)
		ds.playOnce(t)
		ds.assertPickup(t, 2, true)

		if ds.doko.won[2].Length() != 8 {
			t.Errorf("Player 2 has the wrong amount of cards won")
			t.Errorf("Is: %d, should be 8", ds.doko.won[2].Length())
		}

		if ds.doko.teamsKnown() {
			t.Error("Teams should not be known at this point")
		}
	})

	t.Run("16. Trick with the other clubs queen", func(t *testing.T) {
		ds.assertCardMove(t, "dk", true)
		ds.playOnce(t)
		ds.assertCardMove(t, "cq", true)

		if !ds.doko.teamsKnown() {
			t.Error("Teams should be known")
		}

		ds.playOnce(t)
		ds.assertPickup(t, ds.doko.active, true)

		if !ds.doko.teamsKnown() {
			t.Error("Teams should still be known")
		}

	})

	t.Run("17. Invalid command should fail", func(t *testing.T) {
		var p *logic.Packet = logic.NewPacket("whoami")

		var ok bool = ds.doko.PlayerMove(ds.doko.active, p)
		if ok {
			t.Error("Accepted invalid command whoami")
		}
	})

	t.Run("18. Test binary serialization", func(t *testing.T) {
		var shouldbin []byte = []byte {
			// state              active       me      playable
			byte(logic.StatePlaying | (1 << 2) | (0 << 4) | (1 << 6)),
			// hand
			9, 4, 40, 4, 38, 37, 43, 44, 45, 53,
			// table
			0,
			// won
			0, 4, 8, 0,
			// specials[0]: playerID, length
			0, 0,
			// specials[1]: playerID, length
			1, 0,
			// specials[2]
			2, 0,
			// specials[3]
			3, 0,
		}

		var buf bytes.Buffer
		doko.WriteBinary(0, &buf)
		var isbin []byte = buf.Bytes()

		if len(isbin) != len(shouldbin) {
			t.Errorf("Wrong length for binary: %d, should be %d",
				len(isbin), len(shouldbin))
		}

		var i int
		var is, should byte
		for i, is = range isbin {
			should = shouldbin[i]
			if is != should {
				t.Errorf("Wrong byte @ %d: %d, should be %d", 
					i, is, should)
			}
		}

		if t.Failed() {
			t.Log(isbin)
			t.Log(shouldbin)
		}
	})

	//t.Log(ds.String())
}

func TestGameEnd(t *testing.T) {
	// SETUP:
	var gs *GameStub = &GameStub{logic.StatePreparing}
	var doko *Doko = NewDoko(gs)
	var ds DokoSim = DokoSim{doko}
	ds.doko.Start()

	ds.doko.hands[0] = logic.DeserializeDeck("h10, cq, hq, hq, dq, hj, da, ck, ck, sa, s9, h9")
	ds.doko.hands[1] = logic.DeserializeDeck("h10, sq, cj, cj, sj, sj, hj, dj, dj, ca, c10, s9")
	ds.doko.hands[2] = logic.DeserializeDeck("cq, da, d10, d10, d9, sa, s10, sk, sk, ha, hk, ha")
	ds.doko.hands[3] = logic.DeserializeDeck("sq, dq, ca, c9, dk, dk, d9, c10, c9, s10, h9, hk")

	var i int
	for i = 0; i < 4; i++ {
		doko.start[i] = doko.hands[i].Clone()
	}

	// 0 & 2 play together

	// Play the whole game
	ds.playTrick(t, "ck c10 da c9")   // 2 wins their own fox
	ds.playTrick(t, "ha hk h9 s9")    // 2 wins
	ds.playTrick(t, "sa s10 sa dj")   // 1 wins
	ds.playTrick(t, "h10 d10 d9 h10") // 0 wins, they played the snd h10
	ds.playTrick(t, "da dj d10 sq")   // 3 wins the enemies' fox

	// 0: cq, hq, hq, dq, hj, ck, s9
	// 1: sq, cj, cj, sj, sj, hj, ca
	// 2: cq, d9, s10, sk, sk, hk, ha
	// 3: dq, ca, dk, dk, c10, c9, h9 <-

	ds.playTrick(t, "h9 cq sj ha")  // 0 wins
	ds.playTrick(t, "ck ca d9 c9")  // 2 wins
	ds.playTrick(t, "s10 dk s9 hj") // 1 wins
	ds.playTrick(t, "sq cq dq dq")  // 2 wins

	// 0: hj, hq, hq
	// 1: cj, cj, sj
	// 2: sk, sk, hk <-
	// 3: ca, dk, c10

	ds.playTrick(t, "hk dk hj sj")  // 1 wins
	ds.playTrick(t, "cj sk c10 hq") // 0 wins
	ds.playTrick(t, "hq cj sk ca")  // 0 wins

	for i = 0; i < 4; i++ {
		if doko.hands[i].Length() != 0 {
			t.Errorf("Player %d still has cards", i)
		}
	}
	if doko.table.Length() != 0 {
		t.Error("Table still has cards")
	}

	t.Run("Test won decks", func(t *testing.T) {
		var won map[int]*logic.Deck = make(map[int]*logic.Deck)
		won[0] = logic.DeserializeDeck("h10, d10, d9, h10, h9, cq, sj, ha, cj, sk, c10, hq, hq, cj, sk, ca")
		won[1] = logic.DeserializeDeck("sa, s10, sa, dj, s10, dk, s9, hj, hk, dk, hj, sj")
		won[2] = logic.DeserializeDeck("ck, c10, da, c9, ha, hk, h9, s9, ck, ca, d9, c9, sq, cq, dq, dq")
		won[3] = logic.DeserializeDeck("da, dj, d10, sq")
		ds.TestAllWondecks(t, won)
	})

	t.Run("Meta tests", func(t *testing.T) {
		// 0 & 2 are re, 1 & 3 contra
		ds.TestFriends(t, []int{0, 2}, []int{1, 3})
	})

	// 0: h10, d10, d9, h10, hk, dk, cq, sj, hq, cj, sk, c10, hq, cj, sk, ca
	// 2: ck, c10, da, c9, ha, hk, h9, s9, ca, d9, c9, ck, sq, cq, dq, dq

	// 1: sa, s10, sa, dj, ha, hj, sj, h9, s10, dk, s9, hj
	// 3: da, dj, d10, sq

	t.Run("Test scores", func(t *testing.T) {
		var revalue, contravalue int = doko.Values()

		var reteam, contrateam []int = doko.Teams()

		var expReV, expCoV int = 152, 88

		if revalue != expReV {
			t.Errorf("Re eyes should be %d, is %d", expReV, revalue)
			t.Logf("Re team is %v", reteam)
		}
		if contravalue != expCoV {
			t.Errorf("Contra eyes should be %d, is %d", expCoV, contravalue)
			t.Logf("Contra team is %v", contrateam)
		}

		if t.Failed() {
			t.FailNow()
		}

		var s *Score = doko.Scores()
		// Re (0, 2) +1 won +1 no90
		// Co (1, 3) +1 fox
		var expScores []int = []int{2, 1, 2, 1}
		var expRe []int = []int{ReasonWon, ReasonNo90}
		var expContra []int = []int{ReasonFox}

		cmpia(t, "Scores", s.scores, expScores)
		cmpia(t, "Re r", s.rereasons, expRe)
		cmpia(t, "Contra r", s.contrareasons, expContra)
	})
}
