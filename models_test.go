package bombahead

import "testing"

func TestPositionDistanceTo(t *testing.T) {
	t.Parallel()

	p := Position{X: 1, Y: 2}
	other := Position{X: -2, Y: 6}
	if got := p.DistanceTo(other); got != 7 {
		t.Fatalf("DistanceTo() = %d, want 7", got)
	}
}

func TestFieldCellAt(t *testing.T) {
	t.Parallel()

	field := Field{
		Width:  2,
		Height: 2,
		Cells: []CellType{
			Air, Box,
			Wall, Air,
		},
	}

	if got := field.CellAt(Position{X: 1, Y: 0}); got != Box {
		t.Fatalf("CellAt(1,0) = %q, want %q", got, Box)
	}
	if got := field.CellAt(Position{X: -1, Y: 0}); got != Wall {
		t.Fatalf("CellAt(out-of-bounds) = %q, want %q", got, Wall)
	}

	shortField := Field{
		Width:  3,
		Height: 1,
		Cells:  []CellType{Air},
	}
	if got := shortField.CellAt(Position{X: 2, Y: 0}); got != Wall {
		t.Fatalf("CellAt(missing index) = %q, want %q", got, Wall)
	}
}
