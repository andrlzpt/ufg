package renderer

import "github.com/go-gl/gl/v3.2-compatibility/gl"

type PixelPoint struct {
	X, Y float32
}

// This maps a rectangular HUD onto four wall corners.
// Corners follow the text's reading orientation:
// top-left, top-right, bottom-right, bottom-left.
type Wall struct {
	Corners [4]PixelPoint
	Width   float32
	Height  float32
}

func (w Wall) point(x, y float32) (float32, float32) {
	u := x / w.Width
	v := y / w.Height

	a, b := w.Corners[0], w.Corners[1]
	c, d := w.Corners[2], w.Corners[3]

	topX := a.X + (b.X-a.X)*u
	topY := a.Y + (b.Y-a.Y)*u
	bottomX := d.X + (c.X-d.X)*u
	bottomY := d.Y + (c.Y-d.Y)*u

	return topX + (bottomX-topX)*v,
		topY + (bottomY-topY)*v
}

// Text centers a line across the wall.
func (w Wall) Text(label string, y float32, config Config) {
	w.TextAt(label, w.Width/2, y, config)
}

// TextAt centers a line at local wall coordinates.
func (w Wall) TextAt(
	label string,
	centerX, centerY float32,
	config Config,
) {
	runes := []rune(label)
	width := float32(len(runes) * config.CellWidth)
	height := float32(config.CellHeight)

	transform := func(x, y float32) (float32, float32) {
		return w.point(
			centerX-width/2+x,
			centerY-height/2+y,
		)
	}

	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	for x, r := range runes {
		if r != ' ' {
			drawTextGlyph(x, 0, r, config, transform)
		}
	}
	gl.End()
}

func (w Wall) rect(
	x, y, width, height, brightness float32,
	config Config,
) {
	vertex := func(x, y float32) {
		px, py := w.point(x, y)
		gl.Vertex2f(pixelToClipX(px, config), pixelToClipY(py, config))
	}

	gl.Color3f(brightness, brightness, brightness)
	gl.Begin(gl.QUADS)
	vertex(x, y)
	vertex(x+width, y)
	vertex(x+width, y+height)
	vertex(x, y+height)
	gl.End()
}

// Bar fills left to right in the wall's reading orientation.
func (w Wall) Bar(
	centerX, y, width, height, fraction float32,
	config Config,
) {
	fraction = max(float32(0), min(float32(1), fraction))
	x := centerX - width/2

	w.rect(x, y, width, height, 0.1, config)
	if fraction > 0 {
		w.rect(x, y, width*fraction, height, 0.95, config)
	}
}
