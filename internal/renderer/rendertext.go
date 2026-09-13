package renderer

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// RenderText draws white glyphs over existing layers without clearing.
// Call BeginFrame first. Use 7x13 cells or integer multiples such as 14x26.
func RenderText(buffer *ascii.Buffer, config Config) {
	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			r := buffer.At(x, y)
			if r != ' ' {
				drawTextGlyph(x, y, r, config, nil)
			}
		}
	}
	gl.End()
}

// RenderTextLine draws a centered label at pixel coordinates.
// quarterTurns: 0 horizontal, 1 clockwise, -1 counterclockwise.
func RenderTextLine(
	label string,
	centerX, centerY float32,
	quarterTurns int,
	config Config,
) {
	runes := []rune(label)
	width := float32(len(runes) * config.CellWidth)
	height := float32(config.CellHeight)

	transform := func(x, y float32) (float32, float32) {
		x -= width / 2
		y -= height / 2

		switch quarterTurns {
		case 1:
			x, y = -y, x
		case -1:
			x, y = y, -x
		}

		return centerX + x, centerY + y
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

func drawTextGlyph(cellX, cellY int, r rune, config Config, transform func(float32, float32) (float32, float32)) {
	face := basicfont.Face7x13
	baseline := fixed.P(0, face.Ascent)
	bounds, mask, maskPoint, _, ok := face.Glyph(baseline, r)
	if !ok {
		bounds, mask, maskPoint, _, _ = face.Glyph(baseline, '?')
	}

	scaleX := float32(config.CellWidth) / float32(face.Advance)
	scaleY := float32(config.CellHeight) / float32(face.Height)
	originX := float32(cellX * config.CellWidth)
	originY := float32(cellY * config.CellHeight)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			maskX := maskPoint.X + x - bounds.Min.X
			maskY := maskPoint.Y + y - bounds.Min.Y
			_, _, _, alpha := mask.At(maskX, maskY).RGBA()
			if alpha == 0 {
				continue
			}

			pixelX := originX + float32(x)*scaleX
			pixelY := originY + float32(y)*scaleY
			vertex := func(x, y float32) {
				if transform != nil {
					x, y = transform(x, y)
				}
				gl.Vertex2f(pixelToClipX(x, config), pixelToClipY(y, config))
			}

			vertex(pixelX, pixelY)
			vertex(pixelX+scaleX, pixelY)
			vertex(pixelX+scaleX, pixelY+scaleY)
			vertex(pixelX, pixelY+scaleY)
		}
	}
}
