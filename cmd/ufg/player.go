package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

func drawLeftPlayer(buffer *ascii.Buffer, player *game.Player) {
	x, y := int(player.X), int(player.Y)

	ascii.DrawPlayer(
		buffer, x, y, ascii.FacingRight, false, len(player.LeftArm.Path) > 0, len(player.RightArm.Path) > 0,
	)

	drawArm(buffer, x, y, player.LeftArm.Path)
	drawArm(buffer, x, y, player.RightArm.Path)
}

// drawArm draws a right-facing player's saved arm path.
func drawArm(buffer *ascii.Buffer, x, y int, path []game.Point) {
	if len(path) == 0 {
		return
	}

	// Join the resting arm to its starting point.
	start := path[0]
	ascii.DrawHorizontalLine(
		buffer, x+start.X-2, y+start.Y, 3, '#',
	)

	for _, point := range path {
		buffer.Set(x+point.X, y+point.Y, '#')
	}

	tip := path[len(path)-1]
	ascii.DrawGlove(buffer, x+tip.X, y+tip.Y, false)
}
