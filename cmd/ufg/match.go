package main

import (
	"fmt"
	"math"

	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

// startMatch creates a fresh match using the selected settings.
func startMatch(app *AppState) {
	width := windowWidth / blockCellWidth
	height := windowHeight / blockCellHeight
	playerY := float64(height-ascii.PlayerHeight) / 2

	app.Match = game.NewMatch(
		app.Config,
		game.Player{
			X:          playerPadding,
			Y:          playerY,
			Speed:      80,
			Facing:     game.FacingRight,
			ArmControl: game.ArmControlLastPressed,
		},
		game.Player{
			X:          float64(width - playerPadding - ascii.PlayerWidth),
			Y:          playerY,
			Speed:      80,
			Facing:     game.FacingLeft,
			ArmControl: game.ArmControlLastPressed,
		},
	)

	if app.Mode == ModeOnePlayer {
		app.Match.Machine = game.NewMachine(app.Difficulty, app.Match.Config.ArmSpeed)
		app.Match.Machine.Reset()
	}

	app.Screen = ScreenArena
	app.Paused = false
	app.MenuItem = 0
	app.DifficultyMenu = false

}

func drawMatchHUD(
	text *ascii.Buffer,
	match *game.Match,
	showBreak bool,
) {
	drawCenteredText(text, 1, fmt.Sprintf(
		"player 1: %d    player 2: %d",
		match.Left.Score,
		match.Right.Score,
	))

	drawCenteredText(text, 3, fmt.Sprintf(
		"round %d/%d    %ds",
		match.Round,
		match.Config.Rounds,
		int(math.Ceil(match.Remaining)),
	))

	if match.Machine != nil {
		drawCenteredText(
			text,
			5,
			"machine: "+game.DifficultyName(match.Machine.Level),
		)
	}

	if showBreak && match.BreakRemaining > 0 && !match.Over() {
		drawCenteredText(
			text,
			text.Height/2,
			"round over. get ready.",
		)
		drawCenteredText(
			text,
			text.Height/2+2,
			fmt.Sprintf(
				"round %d starts in %d",
				match.Round,
				int(math.Ceil(match.BreakRemaining)),
			),
		)
	}
}
