package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/window"
)

var menuOptions = [...]struct {
	Label string
	Mode  Mode
}{
	{"one player", ModeOnePlayer},
	{"two players", ModeTwoPlayers},
}

// menuPosition keeps drawing and mouse hit areas aligned.
// Coordinates are text buffer cells.
func menuPosition(text *ascii.Buffer, item int) (x, y int) {
	width := len([]rune(menuOptions[item].Label)) + 2
	return (text.Width - width) / 2, text.Height/2 + item*2
}

func drawMenu(text *ascii.Buffer, selected int) {
	drawCenteredText(text, text.Height/2-3, "untitled fighting game.")

	for i, option := range menuOptions {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}

		x, y := menuPosition(text, i)
		ascii.DrawText(text, x, y, prefix+option.Label)
	}
}

func handleMenu(state *window.State, app *AppState, text *ascii.Buffer) {
	if state.MenuUpPressed {
		app.MenuItem = (app.MenuItem + len(menuOptions) - 1) % len(menuOptions)
	}
	if state.MenuDownPressed {
		app.MenuItem = (app.MenuItem + 1) % len(menuOptions)
	}

	hovered := -1
	if state.IsMouseInsideWindow {
		for i, option := range menuOptions {
			x, y := menuPosition(text, i)
			width := len([]rune(option.Label)) + 2

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

	// A resting mouse must not override keyboard selection.
	if hovered >= 0 && (state.MouseMoved || state.MousePressed) {
		app.MenuItem = hovered
	}

	if state.ConfirmPressed || (state.MousePressed && hovered >= 0) {
		app.Mode = menuOptions[app.MenuItem].Mode
		app.Screen = ScreenArena
	}
}
