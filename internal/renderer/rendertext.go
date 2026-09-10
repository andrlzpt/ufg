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
				drawTextGlyph(x, y, r, config)
			}
		}
	}
	gl.End()
}

func drawTextGlyph(cellX, cellY int, r rune, config Config) {
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
			left := pixelToClipX(pixelX, config)
			right := pixelToClipX(pixelX+scaleX, config)
			top := pixelToClipY(pixelY, config)
			bottom := pixelToClipY(pixelY+scaleY, config)

			gl.Vertex2f(left, top)
			gl.Vertex2f(right, top)
			gl.Vertex2f(right, bottom)
			gl.Vertex2f(left, bottom)
		}
	}
}
