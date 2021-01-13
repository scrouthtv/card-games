package logic

import "strings"

// Packet is a single message that is sent over a WebSocket.
// It consists of a key (the action) and optional arguments.
// Keys and arguments are alphanumeric without spaces
type Packet []string

// NewPacket creates a new packet by splitting the command
// string by whitespace characters
func NewPacket(command string) *Packet {
	var p Packet = Packet(strings.Split(command, " "))
	return &p
}

// Action returns the action of this packet
func (p *Packet) Action() string {
	if p == nil {
		return ""
	}
	return (*p)[0]
}

// Args returns the arguments of this picket, if any
func (p *Packet) Args() []string {
	if p == nil {
		return nil
	}
	return (*p)[1:]
}
