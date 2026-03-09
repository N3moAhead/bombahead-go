package bombahead

import "testing"

func TestIsWalkable(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Field: Field{
			Width:  3,
			Height: 3,
			Cells: []CellType{
				Air, Wall, Air,
				Air, Box, Air,
				Air, Air, Air,
			},
		},
		Bombs: []Bomb{
			{Pos: Position{X: 2, Y: 2}, Fuse: 3},
		},
	}
	h := NewGameHelpers(state)

	if h.IsWalkable(Position{X: -1, Y: 0}) {
		t.Fatal("expected out-of-bounds position to be non-walkable")
	}
	if h.IsWalkable(Position{X: 1, Y: 0}) {
		t.Fatal("expected wall to be non-walkable")
	}
	if h.IsWalkable(Position{X: 1, Y: 1}) {
		t.Fatal("expected box to be non-walkable")
	}
	if h.IsWalkable(Position{X: 2, Y: 2}) {
		t.Fatal("expected bomb position to be non-walkable")
	}
	if !h.IsWalkable(Position{X: 0, Y: 0}) {
		t.Fatal("expected air cell without bomb to be walkable")
	}
}

func TestGetAdjacentWalkablePositions(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Field: Field{
			Width:  3,
			Height: 3,
			Cells: []CellType{
				Air, Wall, Air,
				Air, Air, Air,
				Air, Air, Air,
			},
		},
		Bombs: []Bomb{
			{Pos: Position{X: 2, Y: 1}, Fuse: 3},
		},
	}
	h := NewGameHelpers(state)

	got := h.GetAdjacentWalkablePositions(Position{X: 1, Y: 1})
	want := []Position{
		{X: 1, Y: 2},
		{X: 0, Y: 1},
	}

	if len(got) != len(want) {
		t.Fatalf("len(adjacent) = %d, want %d; got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("adjacent[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestGetNextActionTowards_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		state     *GameState
		start     Position
		target    Position
		expectOk  bool
		expectAct Action
	}{
		{
			name: "start equals target returns DoNothing",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:     Position{X: 2, Y: 2},
			target:    Position{X: 2, Y: 2},
			expectOk:  true,
			expectAct: DoNothing,
		},
		{
			name: "simple move up",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:     Position{X: 2, Y: 2},
			target:    Position{X: 2, Y: 1},
			expectOk:  true,
			expectAct: MoveUp,
		},
		{
			name: "simple move down",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:     Position{X: 2, Y: 2},
			target:    Position{X: 2, Y: 3},
			expectOk:  true,
			expectAct: MoveDown,
		},
		{
			name: "simple move left",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:     Position{X: 2, Y: 2},
			target:    Position{X: 1, Y: 2},
			expectOk:  true,
			expectAct: MoveLeft,
		},
		{
			name: "simple move right",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:     Position{X: 2, Y: 2},
			target:    Position{X: 3, Y: 2},
			expectOk:  true,
			expectAct: MoveRight,
		},
		{
			name: "path around single wall",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Wall, Air,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 2, Y: 0},
			expectOk:  true,
			expectAct: MoveRight,
		},
		{
			name: "unreachable target surrounded by walls",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Wall, Air,
						Wall, Air, Wall,
						Air, Wall, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 1, Y: 1},
			expectOk:  true,
			expectAct: DoNothing,
		},
		{
			name: "path blocked by bomb",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 1, Y: 1}, Fuse: 3}},
			},
			start:     Position{X: 0, Y: 1},
			target:    Position{X: 2, Y: 1},
			expectOk:  true,
			expectAct: MoveUp,
		},
		{
			name: "diagonal path finds route",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 2, Y: 2},
			expectOk:  true,
			expectAct: MoveRight,
		},
		{
			name: "longer path around obstacle",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 4, Y: 0},
			expectOk:  true,
			expectAct: MoveRight,
		},
		{
			name: "target out of bounds",
			state: &GameState{
				Field: Field{Width: 3, Height: 3, Cells: make([]CellType, 9)},
			},
			start:     Position{X: 1, Y: 1},
			target:    Position{X: 5, Y: 5},
			expectOk:  true,
			expectAct: DoNothing,
		},
		{
			name: "start position blocked by bomb but target reachable",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 1, Y: 1}, Fuse: 3}},
			},
			start:     Position{X: 1, Y: 1},
			target:    Position{X: 0, Y: 0},
			expectOk:  true,
			expectAct: MoveUp,
		},
		{
			name: "narrow corridor",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 4,
					Cells: []CellType{
						Wall, Air, Wall,
						Wall, Air, Wall,
						Wall, Air, Wall,
						Wall, Air, Wall,
					},
				},
			},
			start:     Position{X: 1, Y: 0},
			target:    Position{X: 1, Y: 3},
			expectOk:  true,
			expectAct: MoveDown,
		},
		{
			name: "single cell board same position",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 0, Y: 0},
			expectOk:  true,
			expectAct: DoNothing,
		},
		{
			name: "box blocks path",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Box, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			target:    Position{X: 2, Y: 0},
			expectOk:  true,
			expectAct: MoveDown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			result := h.GetNextActionTowards(tt.start, tt.target)
			if result != tt.expectAct {
				t.Errorf("GetNextActionTowards(%+v, %+v) = %v, want %v", tt.start, tt.target, result, tt.expectAct)
			}
		})
	}
}

func TestIsSafe_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		state    *GameState
		pos      Position
		expected bool
	}{
		{
			name: "position out of bounds negative",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: -1, Y: 0},
			expected: false,
		},
		{
			name: "position out of bounds X",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 5, Y: 0},
			expected: false,
		},
		{
			name: "position out of bounds Y",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 0, Y: 5},
			expected: false,
		},
		{
			name: "bomb at position",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 3}},
			},
			pos:      Position{X: 2, Y: 2},
			expected: false,
		},
		{
			name: "safe position no bombs",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 2, Y: 2},
			expected: true,
		},
		{
			name: "bomb blast range up",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 2, Y: 1},
			expected: false,
		},
		{
			name: "bomb blast range down",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 2, Y: 3},
			expected: false,
		},
		{
			name: "bomb blast range left",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 1, Y: 2},
			expected: false,
		},
		{
			name: "bomb blast range right",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 3, Y: 2},
			expected: false,
		},
		{
			name: "bomb blast blocked by wall up",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 2, Y: 0},
			expected: true,
		},
		{
			name: "bomb blast blocked by box",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Box, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 2, Y: 0},
			expected: false,
		},
		{
			name: "active explosion at position",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 2, Y: 2},
			expected: true,
		},
		{
			name: "multiple bombs overlapping danger",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{
					{Pos: Position{X: 1, Y: 1}, Fuse: 1},
					{Pos: Position{X: 3, Y: 3}, Fuse: 1},
				},
			},
			pos:      Position{X: 2, Y: 1},
			expected: false,
		},
		{
			name: "safe corner when bombs elsewhere",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 0, Y: 0},
			expected: true,
		},
		{
			name: "bomb at edge of range",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 0, Y: 0}, Fuse: 1}},
			},
			pos:      Position{X: 0, Y: 2},
			expected: false,
		},
		{
			name: "bomb with fuse > 1 is still dangerous",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 5}},
			},
			pos:      Position{X: 2, Y: 1},
			expected: false,
		},
		{
			name: "bomb with fuse <= 1 is immediate danger",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			pos:      Position{X: 2, Y: 1},
			expected: false,
		},
		{
			name: "bomb with fuse 0 is immediate danger",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 0}},
			},
			pos:      Position{X: 2, Y: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			result := h.IsSafe(tt.pos)
			if result != tt.expected {
				t.Errorf("IsSafe(%+v) = %v, want %v", tt.pos, result, tt.expected)
			}
		})
	}
}

func TestGetNearestSafePosition_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		state    *GameState
		start    Position
		expected Position
	}{
		{
			name: "already safe returns same position",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			start:    Position{X: 2, Y: 2},
			expected: Position{X: 2, Y: 2},
		},
		{
			name: "returns start when already safe",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells:  make([]CellType, 9),
				},
			},
			start:    Position{X: 0, Y: 0},
			expected: Position{X: 0, Y: 0},
		},
		{
			name: "escapes bomb danger",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
				Bombs: []Bomb{
					{Pos: Position{X: 2, Y: 2}, Fuse: 1},
					{Pos: Position{X: 1, Y: 3}, Fuse: 1},
				},
			},
			start:    Position{X: 2, Y: 2},
			expected: Position{X: 3, Y: 1},
		},
		{
			name: "corner is safe when bombs in middle",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			start:    Position{X: 3, Y: 3},
			expected: Position{X: 3, Y: 3},
		},
		{
			name: "finds safe position in narrow corridor",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 5,
					Cells: []CellType{
						Wall, Air, Wall,
						Wall, Air, Wall,
						Wall, Air, Wall,
						Wall, Air, Wall,
						Wall, Air, Wall,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 1, Y: 2}, Fuse: 1}},
			},
			start:    Position{X: 1, Y: 2},
			expected: Position{X: 1, Y: 2},
		},
		{
			name: "start at corner bomb far away",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells:  make([]CellType, 25),
				},
				Bombs: []Bomb{{Pos: Position{X: 4, Y: 4}, Fuse: 1}},
			},
			start:    Position{X: 0, Y: 0},
			expected: Position{X: 0, Y: 0},
		},
		{
			name: "wall blocks path to safety",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
					},
				},
				Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
			},
			start:    Position{X: 0, Y: 0},
			expected: Position{X: 0, Y: 0},
		},
		{
			name: "single cell board safe",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
			},
			start:    Position{X: 0, Y: 0},
			expected: Position{X: 0, Y: 0},
		},
		{
			name: "single cell board with bomb returns start",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
				Bombs: []Bomb{{Pos: Position{X: 0, Y: 0}, Fuse: 1}},
			},
			start:    Position{X: 0, Y: 0},
			expected: Position{X: 0, Y: 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			result := h.GetNearestSafePosition(tt.start)
			if result != tt.expected {
				t.Errorf("GetNearestSafePosition(%+v) = %+v, want %+v", tt.start, result, tt.expected)
			}
		})
	}
}

func TestFindNearestBox_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		state     *GameState
		start     Position
		expectOk  bool
		expectPos Position
	}{
		{
			name: "box directly adjacent",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Box,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  true,
			expectPos: Position{X: 2, Y: 1},
		},
		{
			name: "box at corner",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Box, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  true,
			expectPos: Position{X: 0, Y: 0},
		},
		{
			name: "no boxes in entire field",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Wall, Air,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  false,
			expectPos: Position{},
		},
		{
			name: "boxes unreachable due to walls",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Box, Wall, Box,
						Wall, Air, Wall,
						Box, Wall, Box,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  false,
			expectPos: Position{},
		},
		{
			name: "multiple boxes finds nearest",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Box, Air, Air, Air, Box,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Box, Air, Air, Air, Box,
					},
				},
			},
			start:     Position{X: 2, Y: 2},
			expectOk:  true,
			expectPos: Position{X: 4, Y: 0},
		},
		{
			name: "box behind wall not reachable",
			state: &GameState{
				Field: Field{
					Width:  4,
					Height: 4,
					Cells: []CellType{
						Air, Wall, Box, Air,
						Air, Wall, Air, Air,
						Air, Wall, Air, Air,
						Air, Wall, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			expectOk:  false,
			expectPos: Position{},
		},
		{
			name: "start on box returns box",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Box, Air,
						Air, Air, Air,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  true,
			expectPos: Position{X: 1, Y: 1},
		},
		{
			name: "long path to box around walls",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Wall, Box, Air,
						Air, Wall, Wall, Air, Air,
						Air, Air, Wall, Air, Air,
						Wall, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
			},
			start:     Position{X: 0, Y: 0},
			expectOk:  true,
			expectPos: Position{X: 3, Y: 0},
		},
		{
			name: "single cell board with box",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Box}},
			},
			start:     Position{X: 0, Y: 0},
			expectOk:  true,
			expectPos: Position{X: 0, Y: 0},
		},
		{
			name: "single cell board without box",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
			},
			start:     Position{X: 0, Y: 0},
			expectOk:  false,
			expectPos: Position{},
		},
		{
			name: "box surrounded by boxes still finds one",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Box, Box, Box,
						Box, Air, Box,
						Box, Box, Box,
					},
				},
			},
			start:     Position{X: 1, Y: 1},
			expectOk:  true,
			expectPos: Position{X: 1, Y: 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			pos, ok := h.FindNearestBox(tt.start)
			if ok != tt.expectOk {
				t.Errorf("FindNearestBox(%+v) ok = %v, want %v", tt.start, ok, tt.expectOk)
			}
			if pos != tt.expectPos {
				t.Errorf("FindNearestBox(%+v) pos = %+v, want %+v", tt.start, pos, tt.expectPos)
			}
		})
	}
}

func TestIsWalkable_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		state    *GameState
		pos      Position
		expected bool
	}{
		{
			name: "negative coordinates",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: -1, Y: 0},
			expected: false,
		},
		{
			name: "negative Y coordinate",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 0, Y: -1},
			expected: false,
		},
		{
			name: "out of bounds X",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 5, Y: 0},
			expected: false,
		},
		{
			name: "out of bounds Y",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 0, Y: 5},
			expected: false,
		},
		{
			name: "out of bounds both",
			state: &GameState{
				Field: Field{Width: 5, Height: 5, Cells: make([]CellType, 25)},
			},
			pos:      Position{X: 10, Y: 10},
			expected: false,
		},
		{
			name: "wall cell",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells:  []CellType{Wall, Air, Air, Air, Air, Air, Air, Air, Air},
				},
			},
			pos:      Position{X: 0, Y: 0},
			expected: false,
		},
		{
			name: "box cell",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells:  []CellType{Air, Box, Air, Air, Air, Air, Air, Air, Air},
				},
			},
			pos:      Position{X: 1, Y: 0},
			expected: false,
		},
		{
			name: "bomb at position",
			state: &GameState{
				Field: Field{Width: 3, Height: 3, Cells: make([]CellType, 9)},
				Bombs: []Bomb{{Pos: Position{X: 1, Y: 1}, Fuse: 3}},
			},
			pos:      Position{X: 1, Y: 1},
			expected: false,
		},
		{
			name: "air cell without bomb",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells:  []CellType{Air, Air, Air, Air, Air, Air, Air, Air, Air},
				},
			},
			pos:      Position{X: 1, Y: 1},
			expected: true,
		},
		{
			name: "single cell board walkable",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
			},
			pos:      Position{X: 0, Y: 0},
			expected: true,
		},
		{
			name: "single cell board with wall",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Wall}},
			},
			pos:      Position{X: 0, Y: 0},
			expected: false,
		},
		{
			name: "multiple bombs same position",
			state: &GameState{
				Field: Field{Width: 3, Height: 3, Cells: make([]CellType, 9)},
				Bombs: []Bomb{
					{Pos: Position{X: 1, Y: 1}, Fuse: 3},
					{Pos: Position{X: 1, Y: 1}, Fuse: 2},
				},
			},
			pos:      Position{X: 1, Y: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			result := h.IsWalkable(tt.pos)
			if result != tt.expected {
				t.Errorf("IsWalkable(%+v) = %v, want %v", tt.pos, result, tt.expected)
			}
		})
	}
}

func TestGetAdjacentWalkablePositions_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		state    *GameState
		pos      Position
		expected int
	}{
		{
			name: "center of open board",
			state: &GameState{
				Field: Field{
					Width:  5,
					Height: 5,
					Cells: []CellType{
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
						Air, Air, Air, Air, Air,
					},
				},
			},
			pos:      Position{X: 2, Y: 2},
			expected: 4,
		},
		{
			name: "corner position",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
			},
			pos:      Position{X: 0, Y: 0},
			expected: 2,
		},
		{
			name: "edge position",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
			},
			pos:      Position{X: 1, Y: 0},
			expected: 3,
		},
		{
			name: "all directions blocked by walls",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Wall, Wall, Wall,
						Wall, Air, Wall,
						Wall, Wall, Wall,
					},
				},
			},
			pos:      Position{X: 1, Y: 1},
			expected: 0,
		},
		{
			name: "blocked by bombs",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Air, Air,
						Air, Air, Air,
						Air, Air, Air,
					},
				},
				Bombs: []Bomb{
					{Pos: Position{X: 2, Y: 1}, Fuse: 3},
					{Pos: Position{X: 1, Y: 2}, Fuse: 3},
				},
			},
			pos:      Position{X: 1, Y: 1},
			expected: 2,
		},
		{
			name: "single cell board",
			state: &GameState{
				Field: Field{Width: 1, Height: 1, Cells: []CellType{Air}},
			},
			pos:      Position{X: 0, Y: 0},
			expected: 0,
		},
		{
			name: "surrounded by boxes",
			state: &GameState{
				Field: Field{
					Width:  3,
					Height: 3,
					Cells: []CellType{
						Air, Box, Air,
						Box, Air, Box,
						Air, Box, Air,
					},
				},
			},
			pos:      Position{X: 1, Y: 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewGameHelpers(tt.state)
			result := h.GetAdjacentWalkablePositions(tt.pos)
			if len(result) != tt.expected {
				t.Errorf("GetAdjacentWalkablePositions(%+v) returned %d positions, want %d", tt.pos, len(result), tt.expected)
			}
		})
	}
}

func TestGameHelpers_Integration(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Field: Field{
			Width:  5,
			Height: 5,
			Cells: []CellType{
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
			},
		},
		Bombs: []Bomb{{Pos: Position{X: 2, Y: 2}, Fuse: 1}},
	}
	h := NewGameHelpers(state)

	if !h.IsWalkable(Position{X: 0, Y: 0}) {
		t.Error("expected start position to be walkable")
	}

	if h.IsWalkable(Position{X: 2, Y: 2}) {
		t.Error("expected bomb position to not be walkable")
	}

	if h.IsSafe(Position{X: 2, Y: 2}) {
		t.Error("expected bomb position to not be safe")
	}

	if !h.IsSafe(Position{X: 0, Y: 0}) {
		t.Error("expected corner to be safe")
	}

	action := h.GetNextActionTowards(Position{X: 0, Y: 0}, Position{X: 2, Y: 2})
	if action == DoNothing {
		t.Error("expected action towards target")
	}
}
