package bombahead

// Action represents the possible moves a bot can make
// Values must match the server protocol exactly
type Action string

const (
	MoveUp    Action = "move_up"
	MoveDown  Action = "move_down"
	MoveLeft  Action = "move_left"
	MoveRight Action = "move_right"
	PlaceBomb Action = "place_bomb"
	DoNothing Action = "nothing"
)

// CellType represents the content of a cell on the game board
type CellType string

const (
	Air  CellType = "AIR"
	Wall CellType = "WALL"
	Box  CellType = "BOX"
)
