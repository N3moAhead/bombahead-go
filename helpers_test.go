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

func TestGetNextActionTowards(t *testing.T) {
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
	}
	h := NewGameHelpers(state)

	if got := h.GetNextActionTowards(Position{X: 1, Y: 1}, Position{X: 3, Y: 1}); got != MoveRight {
		t.Fatalf("GetNextActionTowards() = %q, want %q", got, MoveRight)
	}
	if got := h.GetNextActionTowards(Position{X: 2, Y: 2}, Position{X: 2, Y: 2}); got != DoNothing {
		t.Fatalf("GetNextActionTowards(start==target) = %q, want %q", got, DoNothing)
	}

	state.Field.Cells = []CellType{
		Air, Air, Air,
		Air, Wall, Wall,
		Air, Wall, Air,
	}
	state.Field.Width = 3
	state.Field.Height = 3
	if got := h.GetNextActionTowards(Position{X: 0, Y: 0}, Position{X: 2, Y: 2}); got != DoNothing {
		t.Fatalf("GetNextActionTowards(unreachable) = %q, want %q", got, DoNothing)
	}
}

func TestIsSafe_WithBlastBlockingAndChainReaction(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Field: Field{
			Width:  5,
			Height: 5,
			Cells: []CellType{
				Air, Air, Air, Air, Air,
				Air, Air, Wall, Air, Air,
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
				Air, Air, Air, Air, Air,
			},
		},
		Bombs: []Bomb{
			{Pos: Position{X: 2, Y: 2}, Fuse: 1},
			{Pos: Position{X: 4, Y: 2}, Fuse: 5},
		},
		Explosions: []Position{{X: 0, Y: 0}},
	}
	h := NewGameHelpers(state)

	if h.IsSafe(Position{X: 2, Y: 3}) {
		t.Fatal("expected cell in direct blast lane to be unsafe")
	}
	if h.IsSafe(Position{X: 4, Y: 4}) {
		t.Fatal("expected chained bomb blast cell to be unsafe")
	}
	if h.IsSafe(Position{X: 0, Y: 0}) {
		t.Fatal("expected active explosion cell to be unsafe")
	}
	if !h.IsSafe(Position{X: 2, Y: 0}) {
		t.Fatal("expected wall-blocked cell to stay safe")
	}
}

func TestGetNearestSafePosition(t *testing.T) {
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
		Bombs: []Bomb{
			{Pos: Position{X: 2, Y: 2}, Fuse: 1},
		},
	}
	h := NewGameHelpers(state)

	if got := h.GetNearestSafePosition(Position{X: 0, Y: 0}); got != (Position{X: 0, Y: 0}) {
		t.Fatalf("GetNearestSafePosition(safe start) = %+v, want start", got)
	}

	if got := h.GetNearestSafePosition(Position{X: 2, Y: 2}); got != (Position{X: 3, Y: 1}) {
		t.Fatalf("GetNearestSafePosition(unsafe start) = %+v, want {3,1}", got)
	}
}

func TestFindNearestBox(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Field: Field{
			Width:  4,
			Height: 3,
			Cells: []CellType{
				Air, Box, Air, Air,
				Air, Wall, Wall, Air,
				Air, Air, Air, Air,
			},
		},
	}
	h := NewGameHelpers(state)

	if pos, ok := h.FindNearestBox(Position{X: 0, Y: 0}); !ok || pos != (Position{X: 1, Y: 0}) {
		t.Fatalf("FindNearestBox() = (%+v, %v), want ({1,0}, true)", pos, ok)
	}

	state.Field.Cells = []CellType{
		Air, Air, Air, Air,
		Air, Wall, Wall, Air,
		Air, Air, Air, Air,
	}
	if _, ok := h.FindNearestBox(Position{X: 0, Y: 0}); ok {
		t.Fatal("expected no box to be found")
	}
}
