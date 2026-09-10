package main

import (
	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/window"
)

func keyboardInput(state *window.State) game.Input {
	input := game.Input{}

	if state.MoveUp {
		input.MoveY--
	}
	if state.MoveDown {
		input.MoveY++
	}

	return input
}
