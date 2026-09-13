package window

type controllerMenuMemory struct {
	previous   MenuControls
	stickReady bool
}

type controllerMenuTracker struct {
	devices map[int]controllerMenuMemory
}

// addPresses combines controller edges with this frame's keyboard edges.
func (t *controllerMenuTracker) addPresses(state *State) {
	next := make(map[int]controllerMenuMemory, len(state.Controllers))

	for _, controller := range state.Controllers {
		current := controller.Menu
		memory, known := t.devices[controller.ID]

		neutral := current.StickY >= -0.2 && current.StickY <= 0.2

		// Held buttons on connection do not trigger menu actions.
		if !known {
			next[controller.ID] = controllerMenuMemory{
				previous:   current,
				stickReady: neutral,
			}
			continue
		}

		up := current.Up && !memory.previous.Up
		down := current.Down && !memory.previous.Down

		if neutral {
			memory.stickReady = true
		} else if memory.stickReady {
			switch {
			case current.StickY <= -0.6:
				up = true
				memory.stickReady = false
			case current.StickY >= 0.6:
				down = true
				memory.stickReady = false
			}
		}

		state.MenuUpPressed = state.MenuUpPressed || up
		state.MenuDownPressed = state.MenuDownPressed || down
		state.ConfirmPressed = state.ConfirmPressed ||
			(current.Confirm && !memory.previous.Confirm)
		state.BackPressed = state.BackPressed ||
			(current.Back && !memory.previous.Back)
		state.PausePressed = state.PausePressed ||
			(current.Pause && !memory.previous.Pause)

		memory.previous = current
		next[controller.ID] = memory
	}

	// Disconnected devices lose their previous input state.
	t.devices = next
}
