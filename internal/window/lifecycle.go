package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Loop func(state *State)

type State struct {
	MouseX              float64
	MouseY              float64
	IsMouseInsideWindow bool
	DeltaTime           float64
	ElapsedTime         float64
	FrameCount          uint64
	Keyboard            Controls
	Controllers         []Controller
	MenuUpPressed       bool
	MenuDownPressed     bool
	ConfirmPressed      bool
	MousePressed        bool
	MouseMoved          bool
	PausePressed        bool
}

type Config struct {
	Width  int
	Height int
	Title  string
}

func Run(config Config, loop Loop) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)
	window, err := glfw.CreateWindow(config.Width, config.Height, config.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	state := State{}
	window.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		state.IsMouseInsideWindow = entered
	})

	var previousUp, previousDown, previousConfirm, previousMouse, previousPause bool
	lastTime := glfw.GetTime()

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		state.DeltaTime = currentTime - lastTime
		state.ElapsedTime = currentTime
		state.FrameCount++
		lastTime = currentTime

		x, y := window.GetCursorPos()

		state.MouseMoved = x != state.MouseX || y != state.MouseY

		state.MouseX = x
		state.MouseY = y

		up := window.GetKey(glfw.KeyUp) == glfw.Press
		down := window.GetKey(glfw.KeyDown) == glfw.Press
		confirm := window.GetKey(glfw.KeyEnter) == glfw.Press
		mouse := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
		pause := window.GetKey(glfw.KeyEscape) == glfw.Press

		state.MenuUpPressed = up && !previousUp
		state.MenuDownPressed = down && !previousDown
		state.ConfirmPressed = confirm && !previousConfirm
		state.MousePressed = mouse && !previousMouse
		state.PausePressed = pause && !previousPause

		previousUp = up
		previousDown = down
		previousConfirm = confirm
		previousMouse = mouse
		previousPause = pause

		state.Keyboard = readKeyboard(window)
		state.Controllers = readControllers()
		loop(&state)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
