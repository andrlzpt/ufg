package renderer

import (
	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

type Config struct {
	WindowWidth  int
	WindowHeight int
	CellWidth    int
	CellHeight   int
}

var initialized = false

// BeginFrame initializes OpenGL and clears the screen once before drawing layers.
func BeginFrame(width, height int) {
	if !initialized {
		if err := gl.Init(); err != nil {
			panic(err)
		}
		gl.ClearColor(0, 0, 0, 1)
		initialized = true
	}
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func pixelToClipX(x float32, config Config) float32 {
	return (x/float32(config.WindowWidth))*2 - 1
}

func pixelToClipY(y float32, config Config) float32 {
	return 1 - (y/float32(config.WindowHeight))*2
}
