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

type DokoScore struct {
	scores        []int
	rereasons     []int
	contrareasons []int
}

const (
	// ReasonWon indicates that the winning team got a point
	// because they won
	ReasonWon = iota
	// ReasonAgainstTheElders indicates that the contra team
	// got a point because they beat the re party
	ReasonAgainstTheElders
	// ResonNo90 indicates that the winning team got a point
	// because the loosing team didn't reach 90 eyes
	ReasonNo90
	// ResonNo60 indicates that the winning team got a point
	// because the loosing team didn't reach 60 eyes
	ReasonNo60
	// ResonNo30 indicates that the winning team got a point
	// because the loosing team didn't reach 30 eyes
	ReasonNo30
	// ReasonBlack indicates that the winning team got a point
	// because the loosing team didn't reach any eyes
	ReasonBlack
)

func EmptyScore() *DokoScore {
	return &DokoScore{[]int{0, 0, 0, 0}, make([]int, 3), make([]int, 3)}
}

// Scores calculates
func (d *Doko) Scores() *DokoScore {
	if d == nil {
		return EmptyScore()
	}

	var score *DokoScore = EmptyScore()

	var revalue, contravalue = d.Values()

	// Score by won cards:
	var rescore, contrascore = 0, 0
	if revalue > contravalue {
		revalue = 1
		score.rereasons = append(score.rereasons, ReasonWon)
	} else if revalue == contravalue {
		contravalue = 1
		score.contrareasons = append(score.contrareasons, ReasonAgainstTheElders)
	} else {
		contravalue = 2
		score.contrareasons = append(score.contrareasons, ReasonAgainstTheElders, ReasonWon)
	}

	// Extra scoring features:
	var rtmp, ctmp int
	var s scoring
	for _, s = range d.features {
		rtmp, ctmp = s.Score(d)
		rescore += rtmp
		contrascore += ctmp
	}

	// Split score for each player:
	var player int
	var repair, contrapair []int = d.Teams()
	for _, player = range repair {
		score.scores[player] = rescore
	}
	for _, player = range contrapair {
		score.scores[player] = contrascore
	}

	return score
}

// Values calculates the hand value for re and contra team
// The value is the sum of the value of each card the team earned
func (d *Doko) Values() (int, int) {
	if d == nil {
		return 0, 0
	}

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

	return revalue, contravalue
}
