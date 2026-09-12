package game

import "math/rand/v2"

type ItemKind int

const (
	ItemDoubleSpeed ItemKind = iota
	ItemDoublePunches
	ItemPoint
	ItemSwapScores
)

const (
	itemSpawnInterval   = 6.0
	itemLifetime        = 12.0
	doubleSpeedDuration = 5.0
)

// Items occupy a 3-by-5-cell rectangle around Position.
type Item struct {
	Kind      ItemKind
	Position  Point
	Remaining float64
	zone      int
}

func (item Item) Contains(cell Point) bool {
	return cell.X >= item.Position.X-1 &&
		cell.X <= item.Position.X+1 &&
		cell.Y >= item.Position.Y-2 &&
		cell.Y <= item.Position.Y+2
}

func (m *Match) tickItems(deltaTime float64, arena Bounds) {
	m.Left.SpeedRemaining = max(0, m.Left.SpeedRemaining-deltaTime)
	m.Right.SpeedRemaining = max(0, m.Right.SpeedRemaining-deltaTime)

	alive := m.Items[:0]
	for _, item := range m.Items {
		item.Remaining -= deltaTime
		if item.Remaining > 0 {
			alive = append(alive, item)
		}
	}
	m.Items = alive

	m.spawnTimer -= deltaTime
	if m.spawnTimer > 0 {
		return
	}
	m.spawnTimer = itemSpawnInterval

	// Zone 0: between the two 15-cell-wide player sprites.
	m.spawnItem(
		0,
		ItemKind(rand.IntN(2)),
		max(arena.MinX+1, int(m.Left.X)+16),
		min(arena.MaxX-1, int(m.Right.X)-2),
		arena.MinY+2,
		arena.MaxY-2,
	)

	// Zones 1 and 2: reachable body lanes.
	// Use the same score-item kind on both sides.
	scoreKind := ItemPoint
	if rand.IntN(2) == 1 {
		scoreKind = ItemSwapScores
	}

	for i, player := range [2]*Player{&m.Left, &m.Right} {
		bodyX := int(player.X) + 2
		if player.Facing == FacingLeft {
			bodyX = int(player.X) + 12
		}

		m.spawnItem(
			i+1,
			scoreKind,
			bodyX,
			bodyX,
			arena.MinY+4,
			arena.MaxY-4,
		)
	}
}

func (m *Match) spawnItem(
	zone int,
	kind ItemKind,
	minX, maxX, minY, maxY int,
) {
	if maxX < minX || maxY < minY {
		return
	}

	// At most one item per zone.
	for _, item := range m.Items {
		if item.zone == zone {
			return
		}
	}

	for attempt := 0; attempt < 30; attempt++ {
		item := Item{
			Kind: kind,
			Position: Point{
				X: minX + rand.IntN(maxX-minX+1),
				Y: minY + rand.IntN(maxY-minY+1),
			},
			Remaining: itemLifetime,
			zone:      zone,
		}

		if m.itemSpaceFree(item) {
			m.Items = append(m.Items, item)
			return
		}
	}
}

// Avoid spawning underneath an existing player, arm, glove, or item.
func (m *Match) itemSpaceFree(item Item) bool {
	for _, other := range m.Items {
		for y := item.Position.Y - 2; y <= item.Position.Y+2; y++ {
			for x := item.Position.X - 1; x <= item.Position.X+1; x++ {
				if other.Contains(Point{X: x, Y: y}) {
					return false
				}
			}
		}
	}

	for _, player := range [2]*Player{&m.Left, &m.Right} {
		for y := item.Position.Y - 2; y <= item.Position.Y+2; y++ {
			for x := item.Position.X - 1; x <= item.Position.X+1; x++ {
				local := Point{
					X: x - int(player.X),
					Y: y - int(player.Y),
				}

				if player.bodyContains(local) ||
					player.noseContains(local) ||
					player.armContains(leftArm, local) ||
					player.armContains(rightArm, local) {
					return false
				}
			}
		}

		for id := leftArm; id <= rightArm; id++ {
			if itemTouchesGlove(item, player, id) {
				return false
			}
		}
	}

	return true
}

func (m *Match) collectItems() {
	remaining := m.Items[:0]

	for _, item := range m.Items {
		collected := false

		for _, player := range [2]*Player{&m.Left, &m.Right} {
			switch item.Kind {
			case ItemDoubleSpeed, ItemDoublePunches:
				for id := leftArm; id <= rightArm; id++ {
					arm := player.arm(id)

					// Double Punches needs an active, outgoing punch.
					if item.Kind == ItemDoublePunches &&
						(arm.Retracting || len(arm.Path) == 0 ||
							arm.DoublePunches) {
						continue
					}

					if !itemTouchesGlove(item, player, id) {
						continue
					}

					if item.Kind == ItemDoubleSpeed {
						player.SpeedRemaining = doubleSpeedDuration
					} else {
						arm.DoublePunches = true
						arm.HitsRemaining++
					}

					collected = true
					break
				}

			case ItemPoint, ItemSwapScores:
				if itemTouchesBody(item, player) {
					if item.Kind == ItemPoint {
						player.Score++
					} else {
						m.Left.Score, m.Right.Score =
							m.Right.Score, m.Left.Score
					}
					collected = true
				}
			}

			if collected {
				break
			}
		}

		if !collected {
			remaining = append(remaining, item)
		}
	}

	m.Items = remaining
}

func itemTouchesGlove(item Item, player *Player, id armID) bool {
	arm := player.arm(id)
	center := player.armStart(id)
	step := directionStep(player.forward())

	if len(arm.Path) > 0 {
		center = arm.Path[len(arm.Path)-1]
		step = arm.TravelStep()
	}

	center.X += int(player.X)
	center.Y += int(player.Y)

	return gloveTouches(center, step, item.Contains)
}

func itemTouchesBody(item Item, player *Player) bool {
	for y := item.Position.Y - 2; y <= item.Position.Y+2; y++ {
		for x := item.Position.X - 1; x <= item.Position.X+1; x++ {
			if player.bodyContains(Point{
				X: x - int(player.X),
				Y: y - int(player.Y),
			}) {
				return true
			}
		}
	}
	return false
}

// bodyContains matches the current torso and head, excluding arms,
// gloves, and nose. Coordinates are player-local.
func (p *Player) bodyContains(cell Point) bool {
	if p.Facing == FacingLeft {
		cell.X = 14 - cell.X
	}

	switch cell.Y {
	case 4, 10:
		return cell.X >= 2 && cell.X <= 5
	case 5, 9:
		return cell.X >= 1 && cell.X <= 6
	case 6, 7, 8:
		return cell.X >= 1 && cell.X <= 7
	default:
		return false
	}
}
