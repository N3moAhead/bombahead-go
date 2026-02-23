package bombahead

import (
	"strings"
	"testing"
)

func TestRenderField_NilState(t *testing.T) {
	t.Parallel()

	got := RenderField(nil)
	if got != "<nil game state>\n" {
		t.Fatalf("RenderField(nil) = %q, want %q", got, "<nil game state>\n")
	}
}

func TestRenderField_EmptyField(t *testing.T) {
	t.Parallel()

	got := RenderField(&GameState{Field: Field{Width: 0, Height: 0}})
	if got != "<empty field>\n" {
		t.Fatalf("RenderField(empty) = %q, want %q", got, "<empty field>\n")
	}
}

func TestRenderField_BoardAndOverlays(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Me: &Player{ID: "me", Pos: Position{X: 2, Y: 1}, Health: 3, Score: 10},
		Opponents: []Player{
			{ID: "op-z", Pos: Position{X: 3, Y: 0}, Health: 2, Score: 4},
		},
		Players: []Player{
			{ID: "me", Pos: Position{X: 2, Y: 1}, Health: 3, Score: 10},
			{ID: "op-z", Pos: Position{X: 3, Y: 0}, Health: 2, Score: 4},
		},
		Field: Field{
			Width:  4,
			Height: 3,
			Cells: []CellType{
				Air, Wall, Box, Air,
				Air, Air, Air, Air,
				Wall, Air, Air, Box,
			},
		},
		Bombs: []Bomb{
			{Pos: Position{X: 1, Y: 1}, Fuse: 2},
		},
		Explosions: []Position{
			{X: 0, Y: 1},
		},
	}

	got := RenderField(state)

	wantParts := []string{
		"╔════════╗\n",
		"║  🧱📦👾║\n",
		"║💥💣🤖  ║\n",
		"║🧱    📦║\n",
		"╚════════╝\n",
		"--- PLAYERS ---\n",
		"🤖 Player me | Health: 3, Score: 10 | Pos: (2,1)\n",
		"👾 Player op-z | Health: 2, Score: 4 | Pos: (3,0)\n",
		"--- BOMBS ---\n",
		"💣 at (1,1) | Fuse: 2\n",
		"Legend: [space] AIR",
	}

	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Fatalf("RenderField() missing expected output part %q.\nGot:\n%s", part, got)
		}
	}
}

func TestRenderField_ShortPlayerID_NoPanicAndShortening(t *testing.T) {
	t.Parallel()

	state := &GameState{
		Players: []Player{{ID: "xy", Pos: Position{X: 0, Y: 0}, Health: 1, Score: 1}},
		Field:   Field{Width: 1, Height: 1, Cells: []CellType{Air}},
	}

	got := RenderField(state)
	if !strings.Contains(got, "Player xy") {
		t.Fatalf("expected short ID to be rendered without slicing panic, got:\n%s", got)
	}
}
