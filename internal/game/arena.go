package game

type Arena struct {
	Bounds
	Walls map[Point]bool

	origin Point
}

func (a Arena) Blocked(cell Point) bool {
	if !a.Contains(cell) {
		return true
	}

	world := Point{
		X: cell.X + a.origin.X,
		Y: cell.Y + a.origin.Y,
	}
	return a.Walls[world]
}

// RelativeTo shares the wall map instead of copying it.
func (a Arena) RelativeTo(x, y int) Arena {
	return Arena{
		Bounds: a.Bounds.RelativeTo(x, y),
		Walls:  a.Walls,
		origin: Point{
			X: a.origin.X + x,
			Y: a.origin.Y + y,
		},
	}
}

func (m *Match) prepareArena(bounds Bounds) {
	m.Walls = make(map[Point]bool)
	m.Pellets = make(map[Point]bool)

	if m.Config.MazeEnabled {
		m.buildMaze(bounds)
	}

	arena := Arena{Bounds: bounds, Walls: m.Walls}
	m.reachable = m.reachableFloor(arena)

	if m.Config.FeedEnabled {
		const spacing = 4

		for cell := range m.reachable {
			if (cell.X-bounds.MinX-2)%spacing != 0 ||
				(cell.Y-bounds.MinY-2)%spacing != 0 {
				continue
			}

			if !insideStartingPlayer(cell, &m.Left) &&
				!insideStartingPlayer(cell, &m.Right) {
				m.Pellets[cell] = true
			}
		}
	}

	m.arenaReady = true
}

// buildMaze creates four mirrored L-shaped barriers.
// Both player movement lanes and the central crossing stay open.
func (m *Match) buildMaze(bounds Bounds) {
	leftEdge := max(bounds.MinX+8, int(m.Left.X)+23)
	rightEdge := min(bounds.MaxX-8, int(m.Right.X)-8)

	if rightEdge-leftEdge < 48 ||
		bounds.MaxY-bounds.MinY < 48 {
		return
	}

	leftX := leftEdge + (rightEdge-leftEdge)/4
	rightX := leftEdge + rightEdge - leftX
	middleY := (bounds.MinY + bounds.MaxY) / 2

	top := bounds.MinY + 10
	bottom := bounds.MaxY - 10
	upperEnd := middleY - 8
	lowerStart := middleY + 8
	wing := min(12, (rightX-leftX-12)/2)

	for y := top; y <= bottom; y++ {
		if y > upperEnd && y < lowerStart {
			continue
		}

		for thickness := 0; thickness < 2; thickness++ {
			m.Walls[Point{X: leftX + thickness, Y: y}] = true
			m.Walls[Point{X: rightX - thickness, Y: y}] = true
		}
	}

	for offset := 0; offset <= wing; offset++ {
		for _, y := range []int{upperEnd, lowerStart} {
			m.Walls[Point{X: leftX + offset, Y: y}] = true
			m.Walls[Point{X: rightX - offset, Y: y}] = true
		}
	}
}

// reachableFloor finds floor reachable through corridors with enough
// clearance to turn the full glove. It ignores temporary arm paths.
func (m *Match) reachableFloor(arena Arena) map[Point]bool {
	visited := make(map[Point]bool)
	floor := make(map[Point]bool)
	queue := make([]Point, 0)

	clearCenter := func(center Point) bool {
		// Conservative clearance for any glove orientation.
		for y := -2; y <= 2; y++ {
			for x := -2; x <= 2; x++ {
				if arena.Blocked(Point{
					X: center.X + x,
					Y: center.Y + y,
				}) {
					return false
				}
			}
		}
		return true
	}

	add := func(cell Point) {
		if !visited[cell] && clearCenter(cell) {
			visited[cell] = true
			queue = append(queue, cell)
		}
	}

	for _, player := range [2]*Player{&m.Left, &m.Right} {
		for id := leftArm; id <= rightArm; id++ {
			start := player.armStart(id)
			add(Point{
				X: int(player.X) + start.X,
				Y: int(player.Y) + start.Y,
			})
		}
	}

	directions := [...]Direction{
		DirectionUp,
		DirectionDown,
		DirectionLeft,
		DirectionRight,
	}

	for head := 0; head < len(queue); head++ {
		center := queue[head]

		for _, direction := range directions {
			step := directionStep(direction)

			// Mark cells touched by the actual glove shape.
			gloveTouches(center, step, func(cell Point) bool {
				floor[cell] = true
				return false
			})

			add(Point{
				X: center.X + step.X,
				Y: center.Y + step.Y,
			})
		}
	}

	return floor
}

func insideStartingPlayer(cell Point, player *Player) bool {
	x := int(player.X)
	y := int(player.Y)

	return cell.X >= x && cell.X < x+15 &&
		cell.Y >= y && cell.Y < y+15
}

func (m *Match) collectPellets() {
	for _, player := range [2]*Player{&m.Left, &m.Right} {
		for id := leftArm; id <= rightArm; id++ {
			arm := player.arm(id)
			center := player.armStart(id)
			step := directionStep(player.forward())

			if len(arm.Path) > 0 {
				center = arm.Path[len(arm.Path)-1]
				step = arm.TravelStep()
			}

			center.X += int(player.X)
			center.Y += int(player.Y)

			gloveTouches(center, step, func(cell Point) bool {
				if m.Pellets[cell] {
					delete(m.Pellets, cell)
					player.Score += FeedPoints
				}
				return false
			})
		}
	}
}
