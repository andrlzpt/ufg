package ascii

// BackgroundPattern describes a low-contrast texture for block rendering.
type BackgroundPattern int

const (
	BackgroundChecker BackgroundPattern = iota
	BackgroundTiles
	BackgroundBricks
	BackgroundWeave
	BackgroundWood
	BackgroundStipple
	BackgroundSolid
)

// DrawBackground fills a rectangle with a deterministic pattern anchored at x,y.
// Each cell is written, so repainting a texture replaces the previous texture.
func DrawBackground(buffer *Buffer, x, y, width, height int, pattern BackgroundPattern) {
	for row := range height {
		for column := range width {
			buffer.Set(x+column, y+row, backgroundRune(column, row, pattern))
		}
	}
}

func backgroundRune(x, y int, pattern BackgroundPattern) rune {
	switch pattern {
	case BackgroundSolid:
		return ':'
	case BackgroundChecker:
		if (x/6+y/6)%2 == 0 {
			return ':'
		}
		return '-'
	case BackgroundTiles:
		if x%10 == 0 || y%10 == 0 {
			return '.'
		}
		if x%10 == 1 || y%10 == 1 {
			return '='
		}
		return ':'
	case BackgroundBricks:
		offset := (y / 5) % 2 * 6
		if y%5 == 0 || (x+offset)%12 == 0 {
			return '.'
		}
		if y%5 == 1 {
			return '='
		}
		return '-'
	case BackgroundWeave:
		if (x/4+y/4)%2 == 0 {
			if y%2 == 0 {
				return '='
			}
		} else if x%2 == 0 {
			return '='
		}
		return ':'
	case BackgroundWood:
		if y%8 == 0 || (x+(y/8)*13)%32 == 0 {
			return '.'
		}
		if (y+(x/7)%3)%4 == 0 {
			return '-'
		}
		return ':'
	case BackgroundStipple:
		// Coordinate hashing avoids flickering random values between frames.
		n := uint32(x)*374761393 + uint32(y)*668265263
		n = (n ^ (n >> 13)) * 1274126177
		return rune("...:::-"[(n^(n>>16))%7])
	default:
		return ' '
	}
}
