package main

import "strings"

// Packet is a single message that is sent over a WebSocket.
// It consists of a key (the action) and optional arguments.
// Keys and arguments are alphanumeric without spaces
type Packet []string

// Action returns the action of this packet
func (p *Packet) Action() string {
	return (*p)[0]
}

// Args returns the arguments of this picket, if any
func (p *Packet) Args() []string {
	if len(*p) < 1 {
		return []string{}
	}
	return (*p)[1:]
}

func (cm *clientMessage) toPacket() *Packet {
	var msg string = string(cm.msg)
	var p Packet = Packet(strings.Split(msg, " "))
	return &p
}
