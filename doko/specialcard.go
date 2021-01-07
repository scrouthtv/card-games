package doko

import "github.com/scrouthtv/card-games/logic"

type specialCard struct {
	card logic.Card
}

func newFox() scoring {
	var fox specialCard = specialCard{*logic.NewCard(logic.Diamonds, logic.Ace)}
	return &fox
}

func (f *specialCard) Name() string {
	return "Fuchs"
}

func (f *specialCard) Score(doko *Doko) []int {
	var scores []int = []int{0, 0, 0, 0}

	var rescore, contrascore int = 0, 0

	rescore = rescore + contrascore

	panic("not impl")

	return scores
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

	panic("not impl")

	return special
}
