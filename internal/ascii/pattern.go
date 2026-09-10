package ascii

// LinePattern selects a cell pattern. Shading uses the block renderer's rune ramp.
type LinePattern int

const (
	LineSolid LinePattern = iota
	LineDotted
	LineDashed
	LineCheckered
	LineGradient
	LineFadeOut
	LineRaised
	LineSunken
	LineDouble
	LineZigzag
)

// DrawPatternedHorizontalLine measures length and thickness in buffer cells.
func DrawPatternedHorizontalLine(buffer *Buffer, x, y, length, thickness int, pattern LinePattern) {
	drawPatternedLine(buffer, x, y, length, thickness, pattern, false)
}

// DrawPatternedVerticalLine rotates the same pattern along the vertical axis.
func DrawPatternedVerticalLine(buffer *Buffer, x, y, length, thickness int, pattern LinePattern) {
	drawPatternedLine(buffer, x, y, length, thickness, pattern, true)
}

func drawPatternedLine(buffer *Buffer, x, y, length, thickness int, pattern LinePattern, vertical bool) {
	for along := range length {
		for across := 0; across < thickness; across++ {
			r := lineRune(along, across, length, thickness, pattern)
			if vertical {
				buffer.Set(x+across, y+along, r)
			} else {
				buffer.Set(x+along, y+across, r)
			}
		}
	}
}

// Spaces deliberately erase gaps inside the line's rectangle.
func lineRune(along, across, length, thickness int, pattern LinePattern) rune {
	const ramp = " .:-=+*#%@"
	switch pattern {
	case LineDotted:
		if along%3 == 0 {
			return '@'
		}
		return ' '
	case LineDashed:
		if along%8 < 5 {
			return '@'
		}
		return ' '
	case LineCheckered:
		if (along+across)%2 == 0 {
			return '@'
		}
		return ':'
	case LineGradient, LineFadeOut:
		index := len(ramp) - 1
		if length > 1 {
			index = along * (len(ramp) - 1) / (length - 1)
		}
		if pattern == LineFadeOut {
			index = len(ramp) - 1 - index
		}
		return rune(ramp[index])
	case LineRaised, LineSunken:
		if thickness == 1 {
			return '#'
		}
		index := across * (len(ramp) - 1) / (thickness - 1)
		if pattern == LineRaised {
			index = len(ramp) - 1 - index
		}
		return rune(ramp[index])
	case LineDouble:
		if across == 0 || across == thickness-1 {
			return '@'
		}
		return ' '
	case LineZigzag:
		if thickness == 1 {
			return '@'
		}
		row := along % (2 * (thickness - 1))
		if row >= thickness {
			row = 2*(thickness-1) - row
		}
		if across == row {
			return '@'
		}
		return ' '
	default:
		return '@'
	}
}
