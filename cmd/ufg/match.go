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
			Speed:      40,
			Facing:     game.FacingRight,
			ArmControl: game.ArmControlLastPressed,
		},
		game.Player{
			X:          float64(width - playerPadding - ascii.PlayerWidth),
			Y:          playerY,
			Speed:      40,
			Facing:     game.FacingLeft,
			ArmControl: game.ArmControlLastPressed,
		},
	)

	app.Screen = ScreenArena
}

func drawMatchHUD(text *ascii.Buffer, match *game.Match) {
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

	if match.BreakRemaining > 0 && !match.Over() {
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
		return
	}

	if !match.Over() {
		return
	}

	message := "draw"
	if match.Result != game.MatchDraw {
		winner := "player 1"
		if match.Result == game.MatchRightWins {
			winner = "player 2"
		}

		reason := "points"
		if match.Knockout {
			reason = "knockout"
		}
		message = winner + " wins by " + reason
	}

	drawCenteredText(text, text.Height/2, message)
	drawCenteredText(
		text,
		text.Height/2+2,
		"type 'go main' to return",
	)
}
