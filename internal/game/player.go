package game

type Player struct {
	X     float64
	Y     float64
	Speed float64
}

// Move keeps the player's top edge between minY and maxY.
// direction: -1 up, 0 still, 1 down. Speed is cells per second.
func (p *Player) Move(direction, deltaTime, minY, maxY float64) {
	if deltaTime <= 0 || maxY < minY {
		return
	}

	direction = max(-1, min(1, direction))
	p.Y += direction * p.Speed * deltaTime
	p.Y = max(minY, min(maxY, p.Y))
}
