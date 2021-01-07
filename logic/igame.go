package logic

const (
	// StatePreparing indicates that the game is currently preparing (e. g. waiting for players)
	StatePreparing = iota
	// StatePlaying indicates that the game is currently running
	StatePlaying
	// StateEnded indicates that the game has ended
	StateEnded
)

// IGame is a game storage that connects each player# to a client
// and stores the current game state.
// It is mainly used for sending updates to the client
type IGame interface {
	ID() byte
	Name() string
	State() byte
	SetState(state byte)
	SendUpdates()
	PlayerCount() int
}

// GameInfo contains user-relevant information about a game
type GameInfo struct {
	ID         byte   `json:"id"`
	Name       string `json:"name"`
	Game       string `json:"game"`
	Players    int    `json:"players"`
	Maxplayers int    `json:"maxplayers"`
}
