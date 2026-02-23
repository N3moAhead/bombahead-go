package bombahead

// GameHelpers provides utility functions for analyzing the game state
type GameHelpers struct {
	State *GameState
}

const defaultBombRange = 2

// NewGameHelpers creates a new instance of GameHelpers
func NewGameHelpers(state *GameState) *GameHelpers {
	return &GameHelpers{State: state}
}

// IsWalkable checks if a position is within the board and can be traversed
func (h *GameHelpers) IsWalkable(pos Position) bool {
	if pos.X < 0 || pos.X >= h.State.Field.Width || pos.Y < 0 || pos.Y >= h.State.Field.Height {
		return false
	}

	cell := h.State.Field.CellAt(pos)
	if cell == Wall || cell == Box {
		return false
	}

	for _, bomb := range h.State.Bombs {
		if bomb.Pos == pos {
			return false
		}
	}

	return true
}

// GetAdjacentWalkablePositions returns a list of valid adjacent positions
func (h *GameHelpers) GetAdjacentWalkablePositions(pos Position) []Position {
	candidates := []Position{
		{X: pos.X, Y: pos.Y - 1},
		{X: pos.X + 1, Y: pos.Y},
		{X: pos.X, Y: pos.Y + 1},
		{X: pos.X - 1, Y: pos.Y},
	}

	result := make([]Position, 0, 4)
	for _, next := range candidates {
		if h.IsWalkable(next) {
			result = append(result, next)
		}
	}

	return result
}

// GetNextActionTowards returns the next action from start towards target using BFS
func (h *GameHelpers) GetNextActionTowards(start, target Position) Action {
	if start == target {
		return DoNothing
	}

	prev := h.bfs(start, func(pos Position) bool { return pos == target }, false)
	if prev == nil {
		return DoNothing
	}

	path := rebuildPath(start, target, prev)
	if len(path) < 2 {
		return DoNothing
	}

	return actionFromStep(path[0], path[1])
}

// IsSafe checks if a position is currently safe from known explosions and bomb blast lanes
func (h *GameHelpers) IsSafe(pos Position) bool {
	if pos.X < 0 || pos.X >= h.State.Field.Width || pos.Y < 0 || pos.Y >= h.State.Field.Height {
		return false
	}

	for _, b := range h.State.Bombs {
		if b.Pos == pos {
			return false
		}
	}

	danger := h.computeDangerPositions()
	if danger[pos] {
		return false
	}

	return true
}

// GetNearestSafePosition finds the closest safe position from start using BFS
func (h *GameHelpers) GetNearestSafePosition(start Position) Position {
	if h.IsWalkable(start) && h.IsSafe(start) {
		return start
	}

	prev := h.bfs(start, func(pos Position) bool {
		return h.IsWalkable(pos) && h.IsSafe(pos)
	}, true)
	if prev == nil {
		return start
	}

	queue := []Position{start}
	visited := map[Position]bool{start: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if h.IsWalkable(cur) && h.IsSafe(cur) {
			return cur
		}

		for _, next := range h.GetAdjacentWalkablePositions(cur) {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	return start
}

// FindNearestBox locates the closest box position from start
func (h *GameHelpers) FindNearestBox(start Position) (Position, bool) {
	queue := []Position{start}
	visited := map[Position]bool{start: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if h.State.Field.CellAt(cur) == Box {
			return cur, true
		}

		for _, next := range []Position{
			{X: cur.X, Y: cur.Y - 1},
			{X: cur.X + 1, Y: cur.Y},
			{X: cur.X, Y: cur.Y + 1},
			{X: cur.X - 1, Y: cur.Y},
		} {
			if next.X < 0 || next.X >= h.State.Field.Width || next.Y < 0 || next.Y >= h.State.Field.Height {
				continue
			}
			if visited[next] {
				continue
			}

			cell := h.State.Field.CellAt(next)
			if cell == Wall {
				continue
			}

			visited[next] = true
			queue = append(queue, next)
		}
	}

	return Position{}, false
}

func (h *GameHelpers) bfs(start Position, goal func(Position) bool, allowUnsafeStart bool) map[Position]Position {
	if !allowUnsafeStart && !h.IsWalkable(start) {
		return nil
	}

	queue := []Position{start}
	visited := map[Position]bool{start: true}
	prev := make(map[Position]Position)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur != start && goal(cur) {
			return prev
		}

		for _, next := range h.GetAdjacentWalkablePositions(cur) {
			if visited[next] {
				continue
			}
			visited[next] = true
			prev[next] = cur
			queue = append(queue, next)
		}
	}

	return nil
}

func rebuildPath(start, target Position, prev map[Position]Position) []Position {
	path := []Position{target}
	cur := target
	for cur != start {
		p, ok := prev[cur]
		if !ok {
			return nil
		}
		cur = p
		path = append(path, cur)
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func actionFromStep(from, to Position) Action {
	switch {
	case to.X == from.X && to.Y == from.Y-1:
		return MoveUp
	case to.X == from.X+1 && to.Y == from.Y:
		return MoveRight
	case to.X == from.X && to.Y == from.Y+1:
		return MoveDown
	case to.X == from.X-1 && to.Y == from.Y:
		return MoveLeft
	default:
		return DoNothing
	}
}

func (h *GameHelpers) computeDangerPositions() map[Position]bool {
	danger := make(map[Position]bool)
	if h.State == nil {
		return danger
	}

	bombIndex := make(map[Position]int, len(h.State.Bombs))
	for i, b := range h.State.Bombs {
		bombIndex[b.Pos] = i
	}

	triggered := make([]bool, len(h.State.Bombs))
	queue := make([]int, 0, len(h.State.Bombs))

	for _, e := range h.State.Explosions {
		danger[e] = true
		if idx, ok := bombIndex[e]; ok && !triggered[idx] {
			triggered[idx] = true
			queue = append(queue, idx)
		}
	}

	for i, b := range h.State.Bombs {
		if b.Fuse <= 1 && !triggered[i] {
			triggered[i] = true
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]

		blast := h.blastCells(h.State.Bombs[idx].Pos)
		for _, cell := range blast {
			danger[cell] = true
			if hitIdx, ok := bombIndex[cell]; ok && !triggered[hitIdx] {
				triggered[hitIdx] = true
				queue = append(queue, hitIdx)
			}
		}
	}

	return danger
}

func (h *GameHelpers) blastCells(origin Position) []Position {
	cells := []Position{origin}
	directions := []Position{
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
	}

	for _, d := range directions {
		for step := 1; step <= defaultBombRange; step++ {
			pos := Position{
				X: origin.X + d.X*step,
				Y: origin.Y + d.Y*step,
			}
			if pos.X < 0 || pos.X >= h.State.Field.Width || pos.Y < 0 || pos.Y >= h.State.Field.Height {
				break
			}

			cell := h.State.Field.CellAt(pos)
			if cell == Wall {
				break
			}

			cells = append(cells, pos)
			if cell == Box {
				break
			}
		}
	}

	return cells
}
