package game

// Bounds includes both minimum and maximum cells.
type Bounds struct {
	MinX int
	MinY int
	MaxX int
	MaxY int
}

func (b Bounds) Contains(point Point) bool {
	return point.X >= b.MinX && point.X <= b.MaxX &&
		point.Y >= b.MinY && point.Y <= b.MaxY
}

// RelativeTo converts world bounds to player-local coordinates.
func (b Bounds) RelativeTo(x, y int) Bounds {
	return Bounds{
		MinX: b.MinX - x,
		MinY: b.MinY - y,
		MaxX: b.MaxX - x,
		MaxY: b.MaxY - y,
	}
}
