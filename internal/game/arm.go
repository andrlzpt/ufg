package game

type Point struct {
	X int
	Y int
}

type Arm struct {
	Path     []Point
	progress float64
}

// Grow saves each visited cell, relative to the player's top-left.
// The first point is the resting glove center.
func (a *Arm) Grow(held bool, deltaTime float64, start Point) {
	const (
		speed     = 300.0 // Cells per second.
		maxLength = 150   // Extension cells.
	)

	if !held || deltaTime <= 0 {
		return
	}

	if len(a.Path) == 0 {
		a.Path = append(a.Path, start)
	}

	a.progress += speed * deltaTime

	for a.progress >= 1 && len(a.Path)-1 < maxLength {
		tip := a.Path[len(a.Path)-1]
		a.Path = append(a.Path, Point{X: tip.X + 1, Y: tip.Y})
		a.progress--
	}

	if len(a.Path)-1 >= maxLength {
		a.progress = 0
	}
}
