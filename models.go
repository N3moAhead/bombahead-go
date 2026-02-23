package bombahead

import "math"

// Position represents a coordinate on the game board
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// DistanceTo calculates the Manhattan distance to another position
func (p Position) DistanceTo(other Position) int {
	return int(math.Abs(float64(p.X-other.X)) + math.Abs(float64(p.Y-other.Y)))
}

// Player represents a bot in the game
type Player struct {
	ID     string   `json:"id"`
	Pos    Position `json:"pos"`
	Health int      `json:"health"`
	Score  int      `json:"score"`
}

// Bomb represents a bomb placed on the field
type Bomb struct {
	Pos  Position `json:"pos"`
	Fuse int      `json:"fuse"`
}

// Field represents the game board
type Field struct {
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Cells  []CellType `json:"-"`
}

// CellAt returns the cell type at the given position
// If the position is out of bounds, Wall is returned
func (f Field) CellAt(pos Position) CellType {
	if pos.X < 0 || pos.X >= f.Width || pos.Y < 0 || pos.Y >= f.Height {
		return Wall
	}
	idx := pos.Y*f.Width + pos.X
	if idx < 0 || idx >= len(f.Cells) {
		return Wall
	}
	return f.Cells[idx]
}

// GameState contains all information about the current state of the game
type GameState struct {
	CurrentTick int        `json:"currentTick"`
	Me          *Player    `json:"me"`
	Opponents   []Player   `json:"opponents"`
	Players     []Player   `json:"players"`
	Field       Field      `json:"field"`
	Bombs       []Bomb     `json:"bombs"`
	Explosions  []Position `json:"explosions"`
}
