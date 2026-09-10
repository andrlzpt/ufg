package main

import (
	"math"

	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/window"
)

func playerInput(controls window.Controls) game.Input {
	direction := steerDirection(controls.SteerX, controls.SteerY)

	input := game.Input{
		MoveY: controls.MoveY,
		LeftArm: game.ArmInput{
			Held: controls.LeftPunch,
		},
		RightArm: game.ArmInput{
			Held: controls.RightPunch,
		},
	}

	if input.LeftArm.Held {
		input.LeftArm.Steer = direction
	}
	if input.RightArm.Held {
		input.RightArm.Steer = direction
	}

	return input
}

// Snake turns use one axis. Horizontal wins equal diagonals.
func steerDirection(x, y float64) game.Direction {
	if x == 0 && y == 0 {
		return game.DirectionNone
	}

	if math.Abs(x) >= math.Abs(y) {
		if x < 0 {
			return game.DirectionLeft
		}
		return game.DirectionRight
	}

	if y < 0 {
		return game.DirectionUp
	}
	return game.DirectionDown
}
