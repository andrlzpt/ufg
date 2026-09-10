package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Controls struct {
	MoveY      float64
	SteerX     float64
	SteerY     float64
	LeftPunch  bool
	RightPunch bool
}

type Controller struct {
	Connected bool
	Controls  Controls
}

func readKeyboard(w *glfw.Window) Controls {
	var controls Controls

	if w.GetKey(glfw.KeyW) == glfw.Press {
		controls.MoveY--
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		controls.MoveY++
	}
	if w.GetKey(glfw.KeyLeft) == glfw.Press {
		controls.SteerX--
	}
	if w.GetKey(glfw.KeyRight) == glfw.Press {
		controls.SteerX++
	}
	if w.GetKey(glfw.KeyUp) == glfw.Press {
		controls.SteerY--
	}
	if w.GetKey(glfw.KeyDown) == glfw.Press {
		controls.SteerY++
	}

	controls.LeftPunch = w.GetKey(glfw.KeyJ) == glfw.Press
	controls.RightPunch = w.GetKey(glfw.KeyK) == glfw.Press

	return controls
}

func readController() Controller {
	for id := glfw.Joystick1; id <= glfw.JoystickLast; id++ {
		state := id.GetGamepadState()
		if state == nil {
			continue
		}

		return Controller{
			Connected: true,
			Controls: Controls{
				MoveY:      axisValue(state.Axes[glfw.AxisLeftY]),
				SteerX:     axisValue(state.Axes[glfw.AxisRightX]),
				SteerY:     axisValue(state.Axes[glfw.AxisRightY]),
				LeftPunch:  state.Buttons[glfw.ButtonLeftBumper] == glfw.Press,
				RightPunch: state.Buttons[glfw.ButtonRightBumper] == glfw.Press,
			},
		}
	}

	// Fresh zero state clears controls when unplugged.
	return Controller{}
}

// Ignore small stick drift. Rescale remaining travel to reach ±1.
func axisValue(value float32) float64 {
	const deadzone = 0.2
	v := float64(value)

	if v > deadzone {
		return (v - deadzone) / (1 - deadzone)
	}
	if v < -deadzone {
		return (v + deadzone) / (1 - deadzone)
	}
	return 0
}
