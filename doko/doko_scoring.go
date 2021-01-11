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
	reasonMaxEyes
)

func EmptyScore() *DokoScore {
	return &DokoScore{[]int{0, 0, 0, 0}, make([]int, 0), make([]int, 0)}
}

// extraEyes returns points for no 90/60/30 and black
// as well as adding the reasons to the score
func (d *Doko) extraEyes(s *DokoScore, revalue int, contravalue int) (int, int) {
	var lval int
	if revalue > contravalue {
		lval = contravalue
	} else if revalue == contravalue {
		return 0, 0
	} else {
		lval = revalue
	}

	var wscore int = 0
	var wreasons []int

	switch {
	case lval == 0:
		wscore++
		wreasons = append(wreasons, ReasonBlack)
		fallthrough
	case lval < 30:
		wscore++
		wreasons = append(wreasons, ReasonNo30)
		fallthrough
	case lval < 60:
		wscore++
		wreasons = append(wreasons, ReasonNo60)
		fallthrough
	case lval < 90:
		wscore++
		wreasons = append(wreasons, ReasonNo90)
	}

	if revalue > contravalue {
		s.rereasons = append(s.rereasons, wreasons...)
		return wscore, 0
	} else {
		s.contrareasons = append(s.contrareasons, wreasons...)
		return 0, wscore
	}
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
		rescore = 1
		score.rereasons = append(score.rereasons, ReasonWon)
	} else if revalue == contravalue {
		contrascore = 1
		score.contrareasons = append(score.contrareasons, ReasonAgainstTheElders)
	} else {
		contrascore = 2
		score.contrareasons = append(score.contrareasons, ReasonAgainstTheElders, ReasonWon)
	}

	var rtmp, ctmp int
	// Extra score for thresholds:
	rtmp, ctmp = d.extraEyes(score, revalue, contravalue)
	rescore += rtmp
	contrascore += ctmp

	// Extra scoring features:
	var s scoring
	for _, s = range d.features {
		rtmp, ctmp = s.Score(d)
		rescore += rtmp
		contrascore += ctmp
		if rtmp > 0 {
			score.rereasons = append(score.rereasons, s.Reason())
		}
		if ctmp > 0 {
			score.contrareasons = append(score.contrareasons, s.Reason())
		}
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
		recards.Merge(d.won[player])
	}
	for _, player = range contrapair {
		contracards.Merge(d.won[player])
	}

	var revalue = recards.Value(dokoCardValue)
	var contravalue = contracards.Value(dokoCardValue)

	return revalue, contravalue
}
