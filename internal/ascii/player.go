package ascii

type Facing int

const (
	FacingRight Facing = iota
	FacingLeft
)

const (
	PlayerWidth  = 15
	PlayerHeight = 15
)

// Resting player, viewed from above, facing right.
// + body, # arms, % head, @ gloves and nose.
var playerShape = [PlayerHeight]string{
	"               ",
	"          @@@@ ",
	"    #####@@@@@ ",
	"   ###    @@@@ ",
	"  ++++         ",
	" ++%%%%   ",
	" ++%%%%%  ",
	" ++%%%%%@@",
	" ++%%%%%  ",
	" ++%%%%   ",
	"  ++++         ",
	"   ###    @@@@ ",
	"    #####@@@@@ ",
	"          @@@@ ",
	"               ",
}

func DrawPlayer(buffer *Buffer, x, y int, facing Facing, dark, upperGloveExtended, lowerGloveExtended bool) {
	for row, line := range playerShape {
		for column, r := range line {
			if r == ' ' {
				continue
			}

			if r == '@' {
				if upperGloveExtended && row <= 3 {
					continue
				}
				if lowerGloveExtended && row >= 11 {
					continue
				}
			}

			if dark {
				switch r {
				case '+':
					r = '-'
				case '#':
					r = '='
				case '%':
					r = '+'
				case '@':
					r = '*'
				}
			}

			offset := column
			if facing == FacingLeft {
				offset = PlayerWidth - 1 - column
			}

			buffer.Set(x+offset, y+row, r)
		}
	}
}

// DrawGlove rotates the resting glove to face dx,dy.
// Direction must be a cardinal unit vector.
func DrawGlove(buffer *Buffer, x, y, dx, dy int, dark bool) {
	if !((dx == 0 && (dy == -1 || dy == 1)) ||
		(dy == 0 && (dx == -1 || dx == 1))) {
		return
	}

	r := rune('@')
	if dark {
		r = '*'
	}

	for across := -1; across <= 1; across++ {
		back := -1
		if across == 0 {
			back = -2
		}

		for forward := back; forward <= 2; forward++ {
			column := x + forward*dx - across*dy
			row := y + forward*dy + across*dx
			buffer.Set(column, row, r)
		}
	}
}
