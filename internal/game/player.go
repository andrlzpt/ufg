package game

type Player struct {
	X        float64
	Y        float64
	Speed    float64
	LeftArm  Arm
	RightArm Arm
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

// Update applies this frame's input.
// Punching and steering will be implemented in later steps.
func (p *Player) Update(input Input, deltaTime, minY, maxY float64) {
	// Freeze body while this prototype arm is out.
	punching := input.LeftArm.Held || input.RightArm.Held
	extended := len(p.LeftArm.Path) > 0 || len(p.RightArm.Path) > 0

	if !punching && !extended {
		p.Move(input.MoveY, deltaTime, minY, maxY)
	}

	p.LeftArm.Grow(input.LeftArm.Held, deltaTime, Point{X: 11, Y: 2})
	p.RightArm.Grow(input.RightArm.Held, deltaTime, Point{X: 11, Y: 12})
}
