package game

type Facing int

const (
	FacingRight Facing = iota
	FacingLeft
)

type ArmControlMode int

const (
	ArmControlLastPressed ArmControlMode = iota
	ArmControlSimultaneous
)

type armID int

const (
	leftArm armID = iota
	rightArm
)

type Player struct {
	X          float64
	Y          float64
	Speed      float64
	LeftArm    Arm
	RightArm   Arm
	ArmControl ArmControlMode
	Facing     Facing

	activeArm    armID
	previousHeld [2]bool
}

// Move keeps the player's top edge between minY and maxY.
// Speed is cells per second.
func (p *Player) Move(direction, deltaTime, minY, maxY float64) {
	if deltaTime <= 0 || maxY < minY {
		return
	}

	direction = max(-1, min(1, direction))
	p.Y += direction * p.Speed * deltaTime
	p.Y = max(minY, min(maxY, p.Y))
}

// Update steers the last activated held arm.
// Each released arm retracts independently.
func (p *Player) Update(input Input, deltaTime, minY, maxY float64, arena Bounds) {
	if deltaTime <= 0 {
		return
	}

	inputs := [2]ArmInput{input.LeftArm, input.RightArm}

	for id := leftArm; id <= rightArm; id++ {
		arm := p.arm(id)
		held := inputs[id].Held
		pressed := held && !p.previousHeld[id]
		p.previousHeld[id] = held

		if !held && len(arm.Path) > 0 {
			arm.BeginRetract()
		}

		// A returning arm must reach home before a fresh punch.
		if pressed && !arm.Retracting {
			if len(arm.Path) == 0 {
				arm.Start(p.armStart(id), p.forward())
			}
			p.activeArm = id
		}
	}

	// When the active arm releases, the other held arm resumes.
	if !p.canControl(p.activeArm, inputs) {
		other := leftArm
		if p.activeArm == leftArm {
			other = rightArm
		}
		if p.canControl(other, inputs) {
			p.activeArm = other
		}
	}

	// Paths are relative to the player's rendered position.
	armBounds := arena.RelativeTo(int(p.X), int(p.Y))

	// Retraction never waits for the other arm.
	for id := leftArm; id <= rightArm; id++ {
		arm := p.arm(id)

		if arm.Retracting {
			arm.Retract(deltaTime)
		} else if p.canControl(id, inputs) {
			if p.ArmControl == ArmControlSimultaneous || id == p.activeArm {
				arm.Grow(inputs[id], deltaTime, armBounds)
			}
		}
	}

	// Body moves only when both arms are home.
	if len(p.LeftArm.Path) == 0 && len(p.RightArm.Path) == 0 {
		p.Move(input.MoveY, deltaTime, minY, maxY)
	}
}

func (p *Player) canControl(id armID, inputs [2]ArmInput) bool {
	arm := p.arm(id)
	return inputs[id].Held && !arm.Retracting && len(arm.Path) > 0
}

func (p *Player) arm(id armID) *Arm {
	if id == leftArm {
		return &p.LeftArm
	}
	return &p.RightArm
}

// Arm origins match the current 15-by-15 player shape.
// Left/right refer to the fighter's own arms.
func (p *Player) armStart(id armID) Point {
	start := Point{X: 11, Y: 2}
	if id == rightArm {
		start.Y = 12
	}

	if p.Facing == FacingLeft {
		start.X = 14 - start.X
		start.Y = 14 - start.Y
	}

	return start
}

func (p *Player) forward() Direction {
	if p.Facing == FacingLeft {
		return DirectionLeft
	}
	return DirectionRight
}
