package doko

import (
	"log"

	"github.com/scrouthtv/card-games/logic"
)

// Doko is the ruleset for Doppelkopf
type Doko struct {
	g      logic.IGame
	active int

	// start: maps #player to their initial inventory
	start map[int]*logic.Deck
	// hands: maps #player to inventoray
	hands map[int]*logic.Deck
	// won: maps #player to Deck, each trick they won
	won map[int]*logic.Deck
	// table: inventory
	table *logic.Deck

	// playable: if set to true, cards can be played,
	// if set to false, the trick has to be picked up first
	playable bool

	features []scoring
}

type scoring interface {
	Name() string
	Score(doko *Doko) (int, int)
	MarkCards(doko *Doko) []*logic.Card
	Reason() int
}

func dokoCardValue(c *logic.Card) int {
	switch c.Value() {
	case logic.Ace:
		return 11
	case 10:
		return 10
	case logic.King:
		return 4
	case logic.Queen:
		return 3
	case logic.Jack:
		return 2
	}
	return 0
}

// NewDoko generates a new Doppelkopf ruleset hosted by the
// supplied game
func NewDoko(host logic.IGame) *Doko {
	var d Doko = Doko{host, -1, nil, nil, nil, nil,
		true, []scoring{newFox()}}
	d.Reset()
	return &d
}

// Reset resets this game by clearing everything
// and giving all players a new hand
func (d *Doko) Reset() bool {
	d.start = make(map[int]*logic.Deck)
	d.hands = make(map[int]*logic.Deck)
	d.won = make(map[int]*logic.Deck)

	var doko *logic.Deck = logic.NewDeck([]int{logic.Ace, 9, 10, logic.Jack, logic.Queen, logic.King}).Twice().Shuffle()
	var dist []*logic.Deck = doko.DistributeAll(4)

	var i int
	for i = 0; i < len(dist); i++ {
		d.Sort(dist[i])
		d.hands[i] = dist[i]
		d.start[i] = dist[i].Clone()
	}
	d.table = logic.EmptyDeck()

	d.playable = true

	return true
}

func (d *Doko) Active() int {
	return d.active
}

// Start starts this game
func (d *Doko) Start() {
	d.active = 0
	d.g.SetState(logic.StatePlaying)
	d.playable = true
}

// Info returns the GameInfo for this Doppelkopf game
func (d *Doko) Info() logic.GameInfo {
	return logic.GameInfo{
		ID:         d.g.ID(),
		Name:       d.g.Name(),
		Game:       "Doppelkopf",
		Players:    d.g.PlayerCount(),
		Maxplayers: 4,
	}
}

// TypeID returns a UUID for this ruleset
func (d *Doko) TypeID() byte {
	return DokoGameUUID
}

// PlayerMove applies the move specified by the given packet to this game
// and returns whether the action was successful
func (d *Doko) PlayerMove(player int, p *logic.Packet) bool {
	if player != d.active {
		log.Println("Ignoring because this player is not active, active: ", d.active)
		return false
	}

	switch p.Action() {
	case "card":
		// Check 1: are we currently playing
		if d.g.State() != logic.StatePlaying {
			log.Println("Ignoring because we are not playing")
			return false
		}

		// Check 2: does the current trick have to be picked up first
		if !d.playable {
			log.Println("Ignoring because the trick has to be picked up first")
			return false
		}

		// Check 3: is the request complete
		if len(p.Args()) < 1 {
			log.Println("Ignoring because no card was specified")
			return false
		}

		// Check 4: is the requested card a valid card
		var c *logic.Card
		var ok bool
		ok, c = logic.CardFromShort(p.Args()[0])
		if !ok {
			log.Println("Ignoring because invalid card was specified")
			return false
		}

		var removed []*logic.Card
		// Check 5: is this player allowed to play this card
		if d.table.Length() > 0 {
			// If there are already cards on the table, is this card allowed?
			if !d.AllowedCards().Contains(*c) {
				log.Println("Ignoring because the player is not allowed to play this card")
				return false
			}
			removed = d.hands[d.active].Remove(*c, 1)
			if len(removed) != 1 {
				return false
			}
			c = removed[0]
		} else {
			// If there are no cards on the table, does the player own that card?
			removed = d.hands[d.active].Remove(*c, 1)
			if len(removed) < 1 {
				return false
			}
			c = removed[0]
		}

		d.table.AddAll(c)
		if len(*d.table) == 4 {
			var winner int = d.trickWinner(d.table)

			// d.active placed the last card, d.active + 1 placed the first card
			// winner is the # in the trick, not the # in the player array
			// suppose the first placer won, then winner is 0, but if 3 placed
			// the first card, they won:
			winner = winner + d.active + 1
			if winner >= 4 {
				winner -= 4
			}
			d.active = winner
			d.playable = false
		} else {
			d.active++
			if d.active == 4 {
				d.active = 0
			}
		}

		d.g.SendUpdates()
		return true
	case "pickup":
		// Check 1: are we currently playing
		if d.g.State() != logic.StatePlaying {
			log.Println("Ignoring because we are not playing")
			return false
		}

		// Check 2: does the current trick have to be picked up
		if d.playable {
			log.Println("Ignoring because the trick has to be played first")
			return false
		}

		d.playerWonTrick(player)
		if len(*d.hands[d.active]) == 0 {
			d.g.SetState(logic.StateEnded)
		}
		d.playable = true
		return true
	}

	return false
}

// playerWonTrick makes the selected player pick
// up the current trick on the table
func (d *Doko) playerWonTrick(winner int) {
	log.Printf("Player %d won the trick", winner)
	var ok bool
	_, ok = d.won[winner]
	if !ok {
		d.won[winner] = logic.EmptyDeck()
	}
	d.won[winner].Merge(d.table)
	d.table = logic.EmptyDeck()
}

// trickWinner calculates the winner # in this trick
func (d *Doko) trickWinner(trick *logic.Deck) int {
	var winner int = 0
	var wCard = (*trick)[0]

	var i int
	for i = 1; i < trick.Length(); i++ {
		if d.beats(wCard, (*trick)[i]) {
			winner = i
			wCard = (*trick)[i]
		}
	}

	return winner
}
