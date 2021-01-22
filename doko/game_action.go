package doko

import "bytes"

import "github.com/scrouthtv/card-games/logic"

const (
	playActionUUID = iota
	pickupActionUUID
)

type action interface {
	WriteBinary(buf *bytes.Buffer)
}

type playAction struct {
	player int
	card *logic.Card
}

type pickupAction struct {
	player int
}

func (p *playAction) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte(playActionUUID)
	buf.WriteByte(byte(p.player))
	buf.WriteByte(p.card.ToBinary())
}

func (p *pickupAction) WriteBinary(buf *bytes.Buffer) {
	buf.WriteByte(pickupActionUUID)
	buf.WriteByte(byte(p.player))
}
