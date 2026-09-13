package renderer

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

// RenderBlocks draws brightness cells. Call BeginFrame before drawing layers.
func RenderBlocks(buffer *ascii.Buffer, config Config) {
	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			r := buffer.At(x, y)

			if r == '•' {
				drawCell(x, y, 0, config)
				drawPellet(x, y, config)
				continue
			}

			drawCell(x, y, brightnessFor(r), config)
		}

	}
}

func brightnessFor(r rune) float32 {
	brightness, ok := brightnessByRune[r]
	if !ok {
		return 0
	}
	return brightness
}

var brightnessByRune = map[rune]float32{
	' ': 0.0,
	'.': 0.1,
	':': 0.2,
	'-': 0.3,
	'=': 0.4,
	'+': 0.5,
	'*': 0.6,
	'#': 0.75,
	'%': 0.9,
	'@': 1.0,
}

func drawCell(x, y int, brightness float32, config Config) {
	leftPx := float32(x * config.CellWidth)
	rightPx := float32(x+1) * float32(config.CellWidth)
	topPx := float32(y * config.CellHeight)
	bottomPx := float32(y+1) * float32(config.CellHeight)

	left := pixelToClipX(leftPx, config)
	right := pixelToClipX(rightPx, config)
	top := pixelToClipY(topPx, config)
	bottom := pixelToClipY(bottomPx, config)

	drawRect(left, top, right, bottom, brightness)
}

func drawRect(left, top, right, bottom, brightness float32) {
	gl.Color3f(brightness, brightness, brightness)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(left, top)
	gl.Vertex2f(right, top)
	gl.Vertex2f(right, bottom)
	gl.Vertex2f(left, bottom)
	gl.End()
}

// drawPellet draws a centered dot, at most 2-by-2 pixels.
func drawPellet(x, y int, config Config) {
	width := float32(min(2, config.CellWidth))
	height := float32(min(2, config.CellHeight))

	leftPx := float32(x*config.CellWidth) +
		(float32(config.CellWidth)-width)/2
	topPx := float32(y*config.CellHeight) +
		(float32(config.CellHeight)-height)/2

	drawRect(
		pixelToClipX(leftPx, config),
		pixelToClipY(topPx, config),
		pixelToClipX(leftPx+width, config),
		pixelToClipY(topPx+height, config),
		0.9,
	)
}
