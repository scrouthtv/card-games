package doko

import "github.com/scrouthtv/card-games/logic"

type dokoCall struct {
	name   string
	match  func(doko *Doko, player int) bool
	runner func(doko *Doko, player int)
}

var calls []dokoCall = []dokoCall{callReshuffle, callHealthy}

var callReshuffle dokoCall = dokoCall{
	"reshuffle",
	func(doko *Doko, player int) bool {
		var c *logic.Card
		var worstTrump *logic.Card = logic.NewCard(logic.Diamonds, logic.Jack)
		var nines int = 0
		var highestTrumpValue int = 0

		for _, c = range *doko.hands[player] {
			if c.Value() == 9 {
				nines++
			}
			if doko.trumpValue(c) > highestTrumpValue {
				highestTrumpValue = doko.trumpValue(c)
			}
		}

		return nines >= 5 || highestTrumpValue <= doko.trumpValue(worstTrump)
	},
	func(doko *Doko, player int) {
		doko.Reset()
	},
}

var callHealthy dokoCall = dokoCall{
	"healthy",
	func(doko *Doko, player int) bool {
		return true
	},
	func(doko *Doko, player int) {
		if player == 3 {
			doko.active = 0
			doko.playingState = phasePlay
		} else {
			doko.active++
		}
	},
}

func callByName(name string) *dokoCall {
	var c dokoCall
	for _, c = range calls {
		if c.name == name {
			return &c
		}
	}
	return nil
}
