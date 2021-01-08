package doko

import (
	"github.com/scrouthtv/card-games/logic"
)

// Teams returns the player teams,
// all re players are in the first array
// all contra players in the second array
// not always do both arrays have 2 ints (e. g. marriage)
func (d *Doko) Teams() ([]int, []int) {
	var repair, contrapair []int
	var i int
	var inv *logic.Deck
	for i, inv = range d.start {
		if inv.Contains(*logic.NewCard(logic.Clubs, logic.Queen)) {
			repair = append(repair, i)
		} else {
			contrapair = append(contrapair, i)
		}
	}
	return repair, contrapair
}

// PlayerTeams returns a slice that maps player# to team
// 0 being a re player, 1 being a contra player
func (d *Doko) PlayerTeams() []int {
	var repair, contrapair []int = d.Teams()
	var playerTeams []int

	var player int
	for _, player = range repair {
		playerTeams[player] = 0
	}
	for _, player = range contrapair {
		playerTeams[player] = 1
	}

	return playerTeams
}

// Scores calculates the value for each player
// The value is the sum of the value of each card they earned
func (d *Doko) Scores() []int {
	var scores []int = make([]int, 4)
	var repair, contrapair []int = d.Teams()
	var recards, contracards *logic.Deck = logic.EmptyDeck(), logic.EmptyDeck()

	var player int
	for _, player = range repair {
		recards.Merge(d.start[player])
	}
	for _, player = range contrapair {
		contracards.Merge(d.start[player])
	}

	var revalue = recards.Value(dokoCardValue)
	var contravalue = contracards.Value(dokoCardValue)

	for _, player = range repair {
		scores[player] = revalue
	}
	for _, player = range contrapair {
		scores[player] = contravalue
	}

	return scores
}
