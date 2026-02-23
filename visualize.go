package bombahead

import (
	"fmt"
	"sort"
	"strings"
)

const (
	tileAir       = "  "
	tileWall      = "🧱"
	tileBox       = "📦"
	tileBomb      = "💣"
	tileExplosion = "💥"
	tileMe        = "🤖"
)

var opponentIcons = []string{"👾", "🏃", "🚶", "💃", "🕺", "🦊", "🐼", "🐸"}

// RenderField returns a console-friendly visualization of the current game board.
func RenderField(state *GameState) string {
	if state == nil {
		return "<nil game state>\n"
	}

	w := state.Field.Width
	h := state.Field.Height
	if w <= 0 || h <= 0 {
		return "<empty field>\n"
	}

	grid := make([][]string, h)
	for y := 0; y < h; y++ {
		row := make([]string, w)
		for x := 0; x < w; x++ {
			switch state.Field.CellAt(Position{X: x, Y: y}) {
			case Wall:
				row[x] = tileWall
			case Box:
				row[x] = tileBox
			default:
				row[x] = tileAir
			}
		}
		grid[y] = row
	}

	for _, e := range state.Explosions {
		if inBounds(e, w, h) {
			grid[e.Y][e.X] = tileExplosion
		}
	}

	for _, b := range state.Bombs {
		if inBounds(b.Pos, w, h) {
			grid[b.Pos.Y][b.Pos.X] = tileBomb
		}
	}

	icons := opponentIconMap(state)
	for _, p := range state.Opponents {
		if inBounds(p.Pos, w, h) {
			grid[p.Pos.Y][p.Pos.X] = icons[p.ID]
		}
	}

	if state.Me != nil && inBounds(state.Me.Pos, w, h) {
		grid[state.Me.Pos.Y][state.Me.Pos.X] = tileMe
	}

	var sb strings.Builder
	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("══", w))
	sb.WriteString("╗\n")

	for y := range h {
		sb.WriteString("║")
		sb.WriteString(strings.Join(grid[y], ""))
		sb.WriteString("║\n")
	}

	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("══", w))
	sb.WriteString("╝\n")

	appendPlayersSection(&sb, state, icons)
	appendBombsSection(&sb, state.Bombs)

	sb.WriteString("Legend: [space] AIR  🧱 WALL  📦 BOX  💣 BOMB  💥 EXPLOSION  🤖 ME\n")
	return sb.String()
}

func appendPlayersSection(sb *strings.Builder, state *GameState, icons map[string]string) {
	players := stablePlayers(state)
	if len(players) == 0 {
		return
	}

	sb.WriteString("--- PLAYERS ---\n")
	for _, p := range players {
		icon := icons[p.ID]
		if state.Me != nil && p.ID == state.Me.ID {
			icon = tileMe
		}
		fmt.Fprintf(sb, "%s Player %s | Health: %d, Score: %d | Pos: (%d,%d)\n", icon, shortPlayerID(p.ID), p.Health, p.Score, p.Pos.X, p.Pos.Y)
	}
}

func appendBombsSection(sb *strings.Builder, bombs []Bomb) {
	if len(bombs) == 0 {
		return
	}

	sorted := append([]Bomb(nil), bombs...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Pos.Y != sorted[j].Pos.Y {
			return sorted[i].Pos.Y < sorted[j].Pos.Y
		}
		if sorted[i].Pos.X != sorted[j].Pos.X {
			return sorted[i].Pos.X < sorted[j].Pos.X
		}
		return sorted[i].Fuse < sorted[j].Fuse
	})

	sb.WriteString("--- BOMBS ---\n")
	for _, b := range sorted {
		fmt.Fprintf(sb, "💣 at (%d,%d) | Fuse: %d\n", b.Pos.X, b.Pos.Y, b.Fuse)
	}
}

func stablePlayers(state *GameState) []Player {
	if state == nil {
		return nil
	}

	if len(state.Players) > 0 {
		players := append([]Player(nil), state.Players...)
		sort.Slice(players, func(i, j int) bool {
			return players[i].ID < players[j].ID
		})
		return players
	}

	players := make([]Player, 0, len(state.Opponents)+1)
	if state.Me != nil {
		players = append(players, *state.Me)
	}
	players = append(players, state.Opponents...)
	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})
	return players
}

func opponentIconMap(state *GameState) map[string]string {
	icons := make(map[string]string)
	if state == nil {
		return icons
	}

	taken := map[string]struct{}{}
	if state.Me != nil {
		taken[state.Me.ID] = struct{}{}
	}

	opponentIDs := make([]string, 0, len(state.Opponents))
	for _, p := range state.Opponents {
		if _, skip := taken[p.ID]; skip {
			continue
		}
		taken[p.ID] = struct{}{}
		opponentIDs = append(opponentIDs, p.ID)
	}
	sort.Strings(opponentIDs)

	for i, id := range opponentIDs {
		icons[id] = opponentIcons[i%len(opponentIcons)]
	}

	for _, p := range stablePlayers(state) {
		if _, ok := icons[p.ID]; ok {
			continue
		}
		icons[p.ID] = opponentIcons[len(icons)%len(opponentIcons)]
	}

	return icons
}

func shortPlayerID(id string) string {
	if id == "" {
		return "<unknown>"
	}
	const n = 4
	if len(id) <= n {
		return id
	}
	return "..." + id[len(id)-n:]
}

// PrintField prints the current game board visualization to stdout.
func PrintField(state *GameState) {
	fmt.Print(RenderField(state))
}

func inBounds(pos Position, width, height int) bool {
	return pos.X >= 0 && pos.X < width && pos.Y >= 0 && pos.Y < height
}
