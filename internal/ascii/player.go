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
	"          @@@@ ",
	"         @@@@@ ",
	"    ######@@@@ ",
	"   ###         ",
	"  ++++         ",
	" ++%%%%   ",
	" ++%%%%%  ",
	" ++%%%%%@@",
	" ++%%%%%  ",
	" ++%%%%   ",
	"  ++++         ",
	"   ###         ",
	"    ######@@@@ ",
	"         @@@@@ ",
	"          @@@@ ",
}

func DrawPlayer(buffer *Buffer, x, y int, facing Facing) {
	for row, line := range playerShape {
		for column, r := range line {
			if r == ' ' {
				continue
			}

			offset := column
			if facing == FacingLeft {
				offset = PlayerWidth - 1 - column
			}

			buffer.Set(x+offset, y+row, r)
		}
	}
}
