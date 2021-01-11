package doko

import "github.com/scrouthtv/card-games/logic"

// specialCard is a generic card that, if caught by the enemy team, gives
// an extra credit at the end
// This implementation marks both of the specified card, as long as the teams
// are not known. Once they are known, only the relevant cards are marked
// At all times, only the relevant cards are giving extra credit
// when using Score()
type specialCard struct {
	card logic.Card
}

const (
	// ReasonFox indicates that the team got a point for catching an enemie's
	// fox (diamonds ace)
	ReasonFox = iota + reasonMaxEyes
	reasonMaxSpecials
)

func newFox() scoring {
	var fox specialCard = specialCard{*logic.NewCard(logic.Diamonds, logic.Ace)}
	return &fox
}

func (f *specialCard) Name() string {
	return "Fuchs"
}

func (f *specialCard) Reason() int {
	return ReasonFox
}

func (f *specialCard) Score(doko *Doko) (int, int) {
	var scores []int = []int{0, 0, 0, 0}

	var special []*logic.Card

	var i int
	var start *logic.Deck
	var c *logic.Card
	for i = 0; i < 4; i++ {
		start = doko.start[i]
		for _, c = range *start {
			if c.Suit() == f.card.Suit() && c.Value() == f.card.Value() {
				special = append(special, c)
			}
		}
	}

	var owner, winner int
	for i, c = range special {
		winner = doko.whoWon(c)
		if winner != -1 {
			owner = doko.origOwner(c)
			if !doko.IsFriend(owner, winner) {
				scores[winner]++
			}
		}
	}

	var repair, contrapair []int = doko.Teams()
	var rescore, contrascore int = 0, 0

	var player int
	for _, player = range repair {
		rescore += scores[player]
	}
	for _, player = range contrapair {
		contrascore += scores[player]
	}

	return rescore, contrascore
}

func (f *specialCard) MarkCards(doko *Doko) []*logic.Card {
	var special []*logic.Card

	var i int
	var start *logic.Deck
	var c *logic.Card
	for i = 0; i < 4; i++ {
		start = doko.start[i]
		for _, c = range *start {
			if c.Suit() == f.card.Suit() && c.Value() == f.card.Value() {
				special = append(special, c)
			}
		}
	}

	var tK bool = doko.teamsKnown()
	var owner, winner int
	for i, c = range special {
		winner = doko.whoWon(c)
		if winner != -1 {
			// somebody already won the card
			owner = doko.origOwner(c)
			if tK && doko.IsFriend(owner, winner) {
				special = append(special[:i], special[i+1:]...)
			} else if owner == winner {
				// special card is always irrelevant if i cought it myself
				special = append(special[:i], special[i+1:]...)
			}
		}
	}

	return special
}
