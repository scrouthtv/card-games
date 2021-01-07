package logic

// Player is a generic player that takes part in a game
type Player interface {
	// Send sends data to the client
	Send(data []byte)
	// Close closes the client
	Close()
}
