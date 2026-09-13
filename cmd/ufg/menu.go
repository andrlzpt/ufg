package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/window"
)

var titleOptions = []string{"one player", "two players", "difficulty", "arena options", "controls"}
var pauseOptions = []string{"resume", "restart", "controls", "main menu"}
var resultOptions = []string{"rematch", "main menu"}
var difficultyOptions = []string{"easy", "normal", "hard", "back"}
var victoryOptions = []string{"continue", "rematch", "main menu"}
var retryOptions = []string{"retry", "main menu"}
var championshipOptions = []string{"new championship", "main menu"}

func currentMenu(app *AppState) (string, []string) {
	if app.ControlsMenu {
		return "controls", controlsMenuOptions(app)
	}

	if app.Screen == ScreenMain {
		if app.ArenaMenu {
			return "arena options", []string{
				"maze: " + enabledLabel(app.Config.MazeEnabled),
				"feed: " + enabledLabel(app.Config.FeedEnabled),
				"back",
			}
		}
		if app.DifficultyMenu {
			return "difficulty: " + game.DifficultyName(app.Difficulty),
				difficultyOptions
		}
		return "untitled fighting game.", titleOptions
	}

	if app.Match.Over() {
		return matchResultMenu(&app.Match)
	}

	if app.Paused {
		return "paused", pauseOptions
	}

	return "", nil
}

func menuWidth(options []string) int {
	width := 0
	for _, label := range options {
		width = max(width, len([]rune(label))+2)
	}
	return width
}

// Shared layout for drawing and mouse hit areas.
func menuPosition(
	text *ascii.Buffer,
	options []string,
	item int,
) (int, int) {
	return (text.Width - menuWidth(options)) / 2,
		text.Height/2 + item*2
}

func drawAppMenu(text *ascii.Buffer, app *AppState) {
	title, options := currentMenu(app)
	if len(options) == 0 {
		return
	}

	x, _ := menuPosition(text, options, 0)
	ascii.DrawText(text, x+2, text.Height/2-3, title)

	if app.Screen == ScreenArena &&
		!app.Match.Over() &&
		app.MissingDevice != "" {
		drawCenteredText(text, text.Height/2-5, app.MissingDevice)
	}

	for i, label := range options {
		prefix := "  "
		if i == app.MenuItem {
			prefix = "> "
		}

		x, y := menuPosition(text, options, i)
		ascii.DrawText(text, x, y, prefix+label)
	}

	hint := "arrows/stick: move | enter/A: select | B: back"
	if app.ControlsMenu {
		hint = "enter/A/click: change device | esc/B: back"
	}
	drawCenteredText(text, text.Height-2, hint)
}

// handleMenus reports whether this frame belongs to menu interaction.
func handleMenus(
	state *window.State,
	app *AppState,
	text *ascii.Buffer,
) bool {
	if state.BackPressed && menuBack(app) {
		return true
	}

	if state.PausePressed {
		if app.Screen == ScreenMain {
			menuBack(app)
			return true
		}

		if app.Match.Over() {
			returnToMain(app)
			return true
		}

		if app.Paused {
			if app.MissingDevice == "" {
				app.ControlsMenu = false
				app.Paused = false
				app.MenuItem = 0
			}
		} else {
			app.Paused = true
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
		for i := range options {
			x, y := menuPosition(text, options, i)
			width := menuWidth(options)

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

func menuBack(app *AppState) bool {
	switch {
	case app.ControlsMenu:
		closeControlsMenu(app)

	case app.Screen == ScreenMain && app.ArenaMenu:
		app.ArenaMenu = false
		app.MenuItem = 3

	case app.Screen == ScreenMain && app.DifficultyMenu:
		app.DifficultyMenu = false
		app.MenuItem = 2

	case app.Screen == ScreenArena && app.Match.Over():
		returnToMain(app)

	case app.Screen == ScreenArena && app.Paused:
		if app.MissingDevice == "" {
			app.Paused = false
			app.MenuItem = 0
		}

	default:
		return false
	}

	return true
}

func activateMenuItem(app *AppState) {
	if app.ControlsMenu {
		switch app.MenuItem {
		case 0, 1:
			cycleControlDevice(app, app.MenuItem)
		case 2:
			closeControlsMenu(app)
		}
		return
	}

	if app.Screen == ScreenMain {
		if app.ArenaMenu {
			switch app.MenuItem {
			case 0:
				app.Config.MazeEnabled = !app.Config.MazeEnabled
			case 1:
				app.Config.FeedEnabled = !app.Config.FeedEnabled
			case 2:
				app.ArenaMenu = false
				app.MenuItem = 3
			}
			return
		}
		if app.DifficultyMenu {
			if app.MenuItem == 3 {
				app.DifficultyMenu = false
				app.MenuItem = 2
				return
			}

			app.Difficulty = app.MenuItem + 1
			app.DifficultyMenu = false
			app.MenuItem = 2
			return
		}

		switch app.MenuItem {
		case 0:
			app.Mode = ModeOnePlayer
			startMatch(app)
		case 1:
			app.Mode = ModeTwoPlayers
			startMatch(app)
		case 2:
			app.DifficultyMenu = true
			app.MenuItem = app.Difficulty - 1
		case 3:
			app.ArenaMenu = true
			app.MenuItem = 0
		case 4:
			app.ControlsMenu = true
			app.MenuItem = 0
		}
		return
	}

	if app.Match.Over() {
		_, options := matchResultMenu(&app.Match)
		if app.MenuItem < 0 || app.MenuItem >= len(options) {
			return
		}

		switch options[app.MenuItem] {
		case "continue":
			if app.Championship.Advance(app.Match) {
				startMatch(app)
			}

		case "rematch", "retry":
			startMatch(app)

		case "new championship":
			app.Championship = game.NewChampionship()
			app.Mode = ModeOnePlayer
			startMatch(app)

		case "main menu":
			returnToMain(app)
		}
		return
	}

	switch app.MenuItem {
	case 0:
		if app.MissingDevice == "" {
			app.Paused = false
		}
	case 1:
		startMatch(app)
	case 2:
		app.ControlsMenu = true
		app.MenuItem = 0
	case 3:
		returnToMain(app)
	}
}

func returnToMain(app *AppState) {
	app.Screen = ScreenMain
	app.Paused = false
	app.MenuItem = int(app.Mode)
	app.DifficultyMenu = false
	app.ArenaMenu = false
	app.ControlsMenu = false
}

func matchResultMenu(match *game.Match) (string, []string) {
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

	// Local two-player matches have no championship progression.
	if match.Machine == nil {
		return message, resultOptions
	}

	if match.Result != game.MatchLeftWins {
		return message, retryOptions
	}

	if match.TotalBouts > 0 && match.Bout == match.TotalBouts {
		return "championship complete. you win!", championshipOptions
	}

	return message, victoryOptions
}

func enabledLabel(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}
