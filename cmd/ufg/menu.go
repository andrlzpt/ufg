package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/window"
)

var titleOptions = []string{"one player", "two players"}
var pauseOptions = []string{"resume", "restart", "main menu"}
var resultOptions = []string{"rematch", "main menu"}

func currentMenu(app *AppState) (string, []string) {
	if app.Screen == ScreenMain {
		return "untitled fighting game.", titleOptions
	}

	if app.Match.Over() {
		message := "draw"

		if app.Match.Result != game.MatchDraw {
			winner := "player 1"
			if app.Match.Result == game.MatchRightWins {
				winner = "player 2"
			}

			reason := "points"
			if app.Match.Knockout {
				reason = "knockout"
			}

			message = winner + " wins by " + reason
		}

		return message, resultOptions
	}

	if app.Paused {
		return "paused", pauseOptions
	}

	return "", nil
}

// menuPosition uses the same coordinates for drawing and mouse hit areas.
func menuPosition(
	text *ascii.Buffer,
	label string,
	item int,
) (int, int) {
	width := len([]rune(label)) + 2
	return (text.Width - width) / 2, text.Height/2 + item*2
}

func drawAppMenu(text *ascii.Buffer, app *AppState) {
	title, options := currentMenu(app)
	if len(options) == 0 {
		return
	}

	drawCenteredText(text, text.Height/2-3, title)

	for i, label := range options {
		prefix := "  "
		if i == app.MenuItem {
			prefix = "> "
		}

		x, y := menuPosition(text, label, i)
		ascii.DrawText(text, x, y, prefix+label)
	}

	drawCenteredText(
		text, text.Height-2,
		"arrows + enter or mouse click",
	)
}

// handleMenus reports whether this frame belongs to menu interaction.
func handleMenus(
	state *window.State,
	app *AppState,
	text *ascii.Buffer,
) bool {
	if state.PausePressed && app.Screen == ScreenArena {
		if app.Match.Over() {
			returnToMain(app)
		} else {
			app.Paused = !app.Paused
			app.MenuItem = 0
		}
		return true
	}

	_, options := currentMenu(app)
	if len(options) == 0 {
		return false
	}

	app.MenuItem = max(0, min(len(options)-1, app.MenuItem))

	if state.MenuUpPressed {
		app.MenuItem = (app.MenuItem + len(options) - 1) % len(options)
	}
	if state.MenuDownPressed {
		app.MenuItem = (app.MenuItem + 1) % len(options)
	}

	hovered := -1
	if state.IsMouseInsideWindow {
		for i, label := range options {
			x, y := menuPosition(text, label, i)
			width := len([]rune(label)) + 2

			left := float64(x * textCellWidth)
			top := float64(y * textCellHeight)
			right := left + float64(width*textCellWidth)
			bottom := top + float64(textCellHeight)

			if state.MouseX >= left && state.MouseX < right &&
				state.MouseY >= top && state.MouseY < bottom {
				hovered = i
				break
			}
		}
	}

	if hovered >= 0 && (state.MouseMoved || state.MousePressed) {
		app.MenuItem = hovered
	}

	if state.ConfirmPressed || (state.MousePressed && hovered >= 0) {
		activateMenuItem(app)
	}

	return true
}

func activateMenuItem(app *AppState) {
	if app.Screen == ScreenMain {
		app.Mode = ModeOnePlayer
		if app.MenuItem == 1 {
			app.Mode = ModeTwoPlayers
		}
		startMatch(app)
		return
	}

	if app.Match.Over() {
		if app.MenuItem == 0 {
			startMatch(app)
		} else {
			returnToMain(app)
		}
		return
	}

	switch app.MenuItem {
	case 0:
		app.Paused = false
	case 1:
		startMatch(app)
	case 2:
		returnToMain(app)
	}
}

func returnToMain(app *AppState) {
	app.Screen = ScreenMain
	app.Paused = false
	app.MenuItem = int(app.Mode)
}
