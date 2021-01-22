package main

import "testing"
import "time"
import "os"
import "github.com/scrouthtv/card-games/doko"
import "github.com/scrouthtv/card-games/logic"

// TestEndGame provides a fast track to launch a real game
// that is close to the end to test functionality towards
// the end of a game.
// It is skipped by default and only run if $RUN_GAME is set to 1
func TestEndGame(t *testing.T) {
	var doexp string = os.Getenv("RUN_GAME")
	if doexp != "1" {
		t.Skip("Set RUN_GAME to 1 to export serialize-props.mjs")
	}

	go func(t *testing.T) {
		// Wait for the game to be created:
		for len(hub.games) < 1 {
			time.Sleep(time.Millisecond)
		}
		var game *Game = hub.games[0]
		var doko *doko.Doko = game.ruleset.(*doko.Doko)

		// Wait for four players to join:
		for game.PlayerCount() < 4 {
			time.Sleep(time.Second)
		}

		t.Log("4 players joined")

		var i int
		var ok bool
		for i = 0; i < 58; i++ {
			ok = playOnce(doko)
			if !ok {
				t.Errorf("Failed to play #%d", i)
			}
		}

	}(t)

	main()
}

// TestBadGame launches a game where the first player to join
// has 5 nines and the second player has only trumps worse or
// equal to diamonds jack
// This test is skipped if $RUN_GAME is not set
func TestBadGame(t *testing.T) {
	var doexp string = os.Getenv("RUN_GAME")
	if doexp != "1" {
		t.Skip("Set RUN_GAME to 1 to export serialize-props.mjs")
	}

	go func(t *testing.T) {
		// Wait for the game to be created:
		for len(hub.games) < 1 {
			time.Sleep(time.Millisecond)
		}
		var game *Game = hub.games[0]
		var doko *doko.Doko = game.ruleset.(*doko.Doko)

		// Wait for four players to join:
		for game.PlayerCount() < 4 {
			time.Sleep(time.Second)
		}

		t.Log("4 players joined")

		var hands map[int]*logic.Deck = make(map[int]*logic.Deck)
		hands[0] = logic.DeserializeDeck("h9, h9, s9, s9, d9, h10, h10, cq, cq, sq, sq, hq")
		hands[1] = logic.DeserializeDeck("d9, dk, dk, d10, d10, dj, ha, ha, hk, hk, sk, sk")
		hands[2] = logic.DeserializeDeck("sq, sq, sq, sq, sq, sq, sq, sq, sq, sq, sq, sq") // i don't even care anymore
		hands[3] = logic.DeserializeDeck("sa, sa, sa, sa, sa, sa, sa, sa, sa, sa, sa, sa")
		doko.SetHands(hands)
	}(t)

	main()
}

func playOnce(d *doko.Doko) bool {
	var card *logic.Card = d.AllowedCards().Get(0)
	var p *logic.Packet = logic.NewPacket("card " + card.Short())
	var ok bool = d.PlayerMove(d.Active(), p)
	if !ok {
		return d.PlayerMove(d.Active(), logic.NewPacket("pickup"))
	}
	return true
}
