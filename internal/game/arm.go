package game

import "slices"

type Point struct {
	X int
	Y int
}

type Arm struct {
	Path       []Point
	Retracting bool
	direction  Direction
	progress   float64
}

const (
	armSpeed = 30.0 // Cells per second.
)

func (a *Arm) Start(start Point, direction Direction) {
	a.Path = append(a.Path[:0], start)
	a.Retracting = false
	a.direction = direction
	a.progress = 0
}

// Grow extends along the saved grid path.
// No steering keeps the current direction. Direct reversal is ignored.
func (a *Arm) Grow(input ArmInput, deltaTime float64, bounds Bounds) {
	if !input.Held || a.Retracting || deltaTime <= 0 ||
		len(a.Path) == 0 {
		return
	}

	a.progress += armSpeed * deltaTime

	for a.progress >= 1 {
		direction := a.direction
		if canTurn(direction, input.Steer) {
			direction = input.Steer
		}

		step := directionStep(direction)
		tip := a.Path[len(a.Path)-1]

		next := Point{
			X: tip.X + step.X,
			Y: tip.Y + step.Y,
		}

		if !a.canGrow(next, step, bounds) {
			a.progress = 0
			return
		}

		a.direction = direction
		a.Path = append(a.Path, next)
		a.progress--
	}
}

// canGrow checks the full glove shape at its proposed position.
// Its dimensions match ascii.DrawGlove.
func (a *Arm) canGrow(next, step Point, bounds Bounds) bool {
	for across := -1; across <= 1; across++ {
		back := -1
		if across == 0 {
			back = -2
		}

		for forward := back; forward <= 2; forward++ {
			cell := Point{
				X: next.X + forward*step.X - across*step.Y,
				Y: next.Y + forward*step.Y + across*step.X,
			}

			if !bounds.Contains(cell) || a.hitsOwnPath(cell) {
				return false
			}
		}
	}

	return true
}

func (a *Arm) hitsOwnPath(cell Point) bool {
	// The last two cells form the attachment behind the glove.
	// They may overlap the glove naturally.
	end := max(0, len(a.Path)-2)

	return slices.Contains(a.Path[:end], cell)
}

func (a *Arm) BeginRetract() {
	if a.Retracting {
		return
	}

	a.Retracting = true
	a.progress = 0
}

// Retract removes points from the tip back toward the body.
func (a *Arm) Retract(deltaTime float64) {
	if !a.Retracting || deltaTime <= 0 {
		return
	}

	a.progress += armSpeed * deltaTime

	for a.progress >= 1 && len(a.Path) > 1 {
		a.Path = a.Path[:len(a.Path)-1]
		a.progress--
	}

	if len(a.Path) <= 1 {
		a.Path = a.Path[:0]
		a.Retracting = false
		a.direction = DirectionRight
		a.progress = 0
	}
}

// TravelStep gives the glove's direction of travel.
// During retraction, it faces the next point toward home.
func (a *Arm) TravelStep() Point {
	if a.Retracting && len(a.Path) > 1 {
		tip := a.Path[len(a.Path)-1]
		next := a.Path[len(a.Path)-2]
		return Point{X: next.X - tip.X, Y: next.Y - tip.Y}
	}

	return directionStep(a.direction)
}

func canTurn(current, next Direction) bool {
	if next == DirectionNone {
		return false
	}

	currentStep := directionStep(current)
	nextStep := directionStep(next)

	if nextStep == (Point{}) {
		return false
	}

	return nextStep.X != -currentStep.X ||
		nextStep.Y != -currentStep.Y
}

func directionStep(direction Direction) Point {
	switch direction {
	case DirectionUp:
		return Point{Y: -1}
	case DirectionDown:
		return Point{Y: 1}
	case DirectionLeft:
		return Point{X: -1}
	case DirectionRight:
		return Point{X: 1}
	default:
		return Point{}
	}
}
