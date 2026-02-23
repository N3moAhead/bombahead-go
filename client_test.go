package bombahead

import (
	"encoding/json"
	"testing"
)

func TestDecodeCell_StringAndNumericEncodings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want CellType
	}{
		{name: "string air", raw: `"AIR"`, want: Air},
		{name: "string wall", raw: `"WALL"`, want: Wall},
		{name: "numeric 0", raw: `0`, want: Wall},
		{name: "numeric 1", raw: `1`, want: Air},
		{name: "numeric 2", raw: `2`, want: Box},
		{name: "numeric unknown defaults to air", raw: `99`, want: Air},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := decodeCell(json.RawMessage(tc.raw))
			if err != nil {
				t.Fatalf("decodeCell() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("decodeCell() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDecodeCell_InvalidEncoding(t *testing.T) {
	t.Parallel()

	got, err := decodeCell(json.RawMessage(`true`))
	if err == nil {
		t.Fatal("decodeCell() expected error for invalid encoding, got nil")
	}
	if got != Air {
		t.Fatalf("decodeCell() fallback = %q, want %q", got, Air)
	}
}

func TestParseClassicState_AssignsMeAndOpponents(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"players":[
			{"id":"p1","pos":{"x":0,"y":0},"health":3,"score":1},
			{"id":"p2","pos":{"x":1,"y":0},"health":2,"score":2}
		],
		"field":{"width":3,"height":2,"field":["AIR",0,2,"BOX","WALL",1]},
		"bombs":[{"pos":{"x":1,"y":1},"fuse":2}],
		"explosions":[{"x":0,"y":1}]
	}`)

	state, err := parseClassicState(payload, "p2")
	if err != nil {
		t.Fatalf("parseClassicState() unexpected error: %v", err)
	}
	if state.Me == nil || state.Me.ID != "p2" {
		t.Fatalf("Me not assigned correctly: %+v", state.Me)
	}
	if len(state.Opponents) != 1 || state.Opponents[0].ID != "p1" {
		t.Fatalf("Opponents = %+v, want [p1]", state.Opponents)
	}
	if got := state.Field.CellAt(Position{X: 1, Y: 0}); got != Wall {
		t.Fatalf("CellAt(1,0) = %q, want %q", got, Wall)
	}
	if got := state.Field.CellAt(Position{X: 2, Y: 0}); got != Box {
		t.Fatalf("CellAt(2,0) = %q, want %q", got, Box)
	}
}

func TestParseClassicState_FallbackWithoutWelcome(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"players":[
			{"id":"first","pos":{"x":0,"y":0},"health":3,"score":1},
			{"id":"second","pos":{"x":1,"y":0},"health":2,"score":2}
		],
		"field":{"width":1,"height":1,"field":["AIR"]},
		"bombs":[],
		"explosions":[]
	}`)

	state, err := parseClassicState(payload, "unknown-id")
	if err != nil {
		t.Fatalf("parseClassicState() unexpected error: %v", err)
	}
	if state.Me == nil || state.Me.ID != "first" {
		t.Fatalf("fallback Me = %+v, want first player", state.Me)
	}
	if len(state.Opponents) != 1 || state.Opponents[0].ID != "second" {
		t.Fatalf("fallback Opponents = %+v, want [second]", state.Opponents)
	}
}

func TestParseClassicState_InvalidPayload(t *testing.T) {
	t.Parallel()

	_, err := parseClassicState([]byte(`{"players":[`), "p1")
	if err == nil {
		t.Fatal("parseClassicState() expected error for invalid JSON, got nil")
	}
}
