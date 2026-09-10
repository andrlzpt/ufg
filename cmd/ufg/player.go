package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

func drawPlayer(buffer *ascii.Buffer, player *game.Player, dark bool) {
	x, y := int(player.X), int(player.Y)

	facing := ascii.FacingRight
	forwardX := 1

	upper := &player.LeftArm
	lower := &player.RightArm

	if player.Facing == game.FacingLeft {
		facing = ascii.FacingLeft
		forwardX = -1
		upper, lower = lower, upper
	}

	ascii.DrawPlayer(
		buffer, x, y, facing, dark,
		len(upper.Path) > 0,
		len(lower.Path) > 0,
	)

	drawArmPath(buffer, x, y, player.LeftArm.Path, forwardX, dark)
	drawArmPath(buffer, x, y, player.RightArm.Path, forwardX, dark)

	drawArmGlove(buffer, x, y, &player.LeftArm, dark)
	drawArmGlove(buffer, x, y, &player.RightArm, dark)
}

func drawArmPath(buffer *ascii.Buffer, x, y int, path []game.Point, forwardX int, dark bool) {
	if len(path) == 0 {
		return
	}

	r := rune('#')
	if dark {
		r = '='
	}

	// Connect the resting arm to its glove's starting point.
	start := path[0]
	for offset := 0; offset <= 2; offset++ {
		buffer.Set(
			x+start.X-offset*forwardX,
			y+start.Y,
			r,
		)
	}

	for _, point := range path {
		buffer.Set(x+point.X, y+point.Y, r)
	}
}

func drawArmGlove(buffer *ascii.Buffer, x, y int, arm *game.Arm, dark bool) {
	if len(arm.Path) == 0 {
		return
	}

	tip := arm.Path[len(arm.Path)-1]
	step := arm.TravelStep()

	ascii.DrawGlove(buffer, x+tip.X, y+tip.Y, step.X, step.Y, dark)
}
