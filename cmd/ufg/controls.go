package main

import (
	"fmt"
	"slices"

	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/window"
)

const keyboardDevice = -1

type ControlConfig struct {
	// Keyboard is -1. Controllers use their GLFW slot IDs.
	Devices [2]int
}

func deviceControls(
	state *window.State,
	device int,
) (window.Controls, bool) {
	if device == keyboardDevice {
		return state.Keyboard, true
	}

	for _, controller := range state.Controllers {
		if controller.ID == device {
			return controller.Controls, true
		}
	}

	return window.Controls{}, false
}

func matchInputs(
	state *window.State,
	app *AppState,
) (game.Input, game.Input, string) {
	left, leftConnected := deviceControls(state, app.Controls.Devices[0])
	right := window.Controls{}
	rightConnected := true

	if app.Mode == ModeTwoPlayers {
		right, rightConnected = deviceControls(state, app.Controls.Devices[1])
	}

	missing := ""
	switch {
	case !leftConnected && !rightConnected:
		missing = "both devices missing. open controls."
	case !leftConnected:
		missing = "player 1 device missing. open controls."
	case !rightConnected:
		missing = "player 2 device missing. open controls."
	}

	return playerInput(left), playerInput(right), missing
}

func controlDeviceLabel(app *AppState, device int) string {
	if device == keyboardDevice {
		return "keyboard"
	}

	for _, controller := range app.Controllers {
		if controller.ID == device {
			name := []rune(controller.Name)
			if len(name) > 22 {
				name = append(name[:19], '.', '.', '.')
			}

			return fmt.Sprintf(
				"controller %d: %s",
				device+1,
				string(name),
			)
		}
	}

	return fmt.Sprintf("controller %d (missing)", device+1)
}

func controlsMenuOptions(app *AppState) []string {
	return []string{
		"player 1: " + controlDeviceLabel(app, app.Controls.Devices[0]),
		"player 2: " + controlDeviceLabel(app, app.Controls.Devices[1]),
		"back",
	}
}

func cycleControlDevice(app *AppState, player int) {
	devices := []int{keyboardDevice}
	for _, controller := range app.Controllers {
		devices = append(devices, controller.ID)
	}

	current := app.Controls.Devices[player]
	index := slices.Index(devices, current)
	next := devices[(index+1)%len(devices)]

	other := 1 - player
	if app.Controls.Devices[other] == next {
		app.Controls.Devices[other] = current
	}

	app.Controls.Devices[player] = next
}

func closeControlsMenu(app *AppState) {
	app.ControlsMenu = false
	if app.Screen == ScreenMain {
		app.MenuItem = 4
	} else {
		app.MenuItem = 2
	}
}
