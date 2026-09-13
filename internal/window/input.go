package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Controls struct {
	MoveY      float64
	SteerX     float64
	SteerY     float64
	LeftPunch  bool
	RightPunch bool
}

type MenuControls struct {
	Up      bool
	Down    bool
	Confirm bool
	Back    bool
	Pause   bool
	StickY  float32
}

type Controller struct {
	ID       int
	Name     string
	Controls Controls
	Menu     MenuControls
}

func readKeyboard(w *glfw.Window) Controls {
	var controls Controls

	if w.GetKey(glfw.KeyW) == glfw.Press {
		controls.MoveY--
		controls.SteerY--
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		controls.MoveY++
		controls.SteerY++
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		controls.SteerX--
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		controls.SteerX++
	}

	controls.LeftPunch = w.GetKey(glfw.KeyJ) == glfw.Press
	controls.RightPunch = w.GetKey(glfw.KeyK) == glfw.Press
	return controls
}

func readControllers() []Controller {
	var controllers []Controller

	for id := glfw.Joystick1; id <= glfw.JoystickLast; id++ {
		state := id.GetGamepadState()
		if state == nil {
			continue
		}

		controllers = append(controllers, Controller{
			ID:   int(id),
			Name: id.GetGamepadName(),
			Controls: Controls{
				MoveY:      axisValue(state.Axes[glfw.AxisLeftY]),
				SteerX:     axisValue(state.Axes[glfw.AxisRightX]),
				SteerY:     axisValue(state.Axes[glfw.AxisRightY]),
				LeftPunch:  state.Buttons[glfw.ButtonLeftBumper] == glfw.Press,
				RightPunch: state.Buttons[glfw.ButtonRightBumper] == glfw.Press,
			},
			Menu: MenuControls{
				Up:      state.Buttons[glfw.ButtonDpadUp] == glfw.Press,
				Down:    state.Buttons[glfw.ButtonDpadDown] == glfw.Press,
				Confirm: state.Buttons[glfw.ButtonA] == glfw.Press,
				Back:    state.Buttons[glfw.ButtonB] == glfw.Press,
				Pause:   state.Buttons[glfw.ButtonStart] == glfw.Press,
				StickY:  state.Axes[glfw.AxisLeftY],
			},
		})
	}

	return controllers
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
