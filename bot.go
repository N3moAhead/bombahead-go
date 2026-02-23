package bombahead

// Bot defines the interface that developers must implement
type Bot interface {
	GetNextMove(state *GameState, helpers *GameHelpers) Action
}
