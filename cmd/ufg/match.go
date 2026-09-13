package main

import (
	"fmt"
	"math"

	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

// startMatch creates a fresh match using the selected settings.
func startMatch(app *AppState) {
	var opponent game.OpponentConfig

	if app.Mode == ModeOnePlayer {
		var ok bool
		opponent, ok = app.Championship.CurrentOpponent()
		if !ok {
			fmt.Println("No championship opponent available.")
			return
		}
	}

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
		app.Match.Opponent = opponent
		app.Match.Bout = app.Championship.Stage + 1
		app.Match.TotalBouts = app.Championship.TotalBouts()
		app.Match.Machine = game.NewMachine(app.Difficulty, app.Match.Opponent, app.Match.Config.ArmSpeed)
	}

	app.Screen = ScreenArena
	app.Paused = false
	app.MenuItem = 0
	app.DifficultyMenu = false
	app.ArenaMenu = false
	app.ControlsMenu = false

}

func drawMatchHUD(
	text *ascii.Buffer,
	match *game.Match,
	showBreak bool,
) {
	if !showBreak || match.BreakRemaining <= 0 || match.Over() {
		return
	}

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
