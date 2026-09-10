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

// DrawGlove draws the upper glove around its center.
// Shape matches the resting upper glove.
func DrawGlove(buffer *Buffer, x, y int, dark bool) {
	r := rune('@')
	if dark {
		r = '*'
	}

	DrawHorizontalLine(buffer, x-1, y-1, 4, r)
	DrawHorizontalLine(buffer, x-2, y, 5, r)
	DrawHorizontalLine(buffer, x-1, y+1, 4, r)
}
