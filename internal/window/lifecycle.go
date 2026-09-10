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
	MoveUp              bool
	MoveDown            bool
	MenuUpPressed       bool
	MenuDownPressed     bool
	ConfirmPressed      bool
	MousePressed        bool
	MouseMoved          bool
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

	var previousUp, previousDown, previousConfirm, previousMouse bool
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

		state.MenuUpPressed = up && !previousUp
		state.MenuDownPressed = down && !previousDown
		state.ConfirmPressed = confirm && !previousConfirm
		state.MousePressed = mouse && !previousMouse

		previousUp = up
		previousDown = down
		previousConfirm = confirm
		previousMouse = mouse

		state.MoveUp = window.GetKey(glfw.KeyW) == glfw.Press
		state.MoveDown = window.GetKey(glfw.KeyS) == glfw.Press
		loop(&state)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
