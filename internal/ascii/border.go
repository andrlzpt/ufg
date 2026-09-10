package ascii

// DrawBorder draws a block frame using a line pattern. Dimensions and thickness
// are in cells. Thickness is clamped so opposite edges cannot overlap.
// Gaps erase frame cells; the interior remains untouched.
func DrawBorder(buffer *Buffer, x, y, width, height, thickness int, pattern LinePattern) {
	if width < 2 || height < 2 || thickness <= 0 {
		return
	}
	thickness = min(thickness, width/2, height/2)
	if pattern == LineRaised || pattern == LineSunken {
		drawBeveledBorder(buffer, x, y, width, height, thickness, pattern)
		return
	}

	DrawPatternedHorizontalLine(buffer, x, y, width, thickness, pattern)
	DrawPatternedHorizontalLine(buffer, x, y+height-thickness, width, thickness, pattern)
	DrawPatternedVerticalLine(buffer, x, y+thickness, height-2*thickness, thickness, pattern)
	DrawPatternedVerticalLine(buffer, x+width-thickness, y+thickness, height-2*thickness, thickness, pattern)

	// Double frames use concentric rectangles, giving both rails joined corners.
	if pattern == LineDouble {
		for layer := 0; layer < thickness; layer++ {
			r := rune(' ')
			if layer == 0 || layer == thickness-1 {
				r = '@'
			}
			drawFrameLayer(buffer, x+layer, y+layer, width-2*layer, height-2*layer, r, r)
		}
	}
}

func drawBeveledBorder(buffer *Buffer, x, y, width, height, thickness int, pattern LinePattern) {
	// Light falls from the upper left. Reverse highlights for an inset frame.
	const highlights = "@%#"
	const shadows = ":-="
	for layer := 0; layer < thickness; layer++ {
		index := 0
		if thickness > 1 {
			index = layer * 2 / (thickness - 1)
		}
		light, dark := rune(highlights[index]), rune(shadows[index])
		if pattern == LineSunken {
			light, dark = dark, light
		}
		drawFrameLayer(buffer, x+layer, y+layer, width-2*layer, height-2*layer, light, dark)
	}
}

func drawFrameLayer(buffer *Buffer, x, y, width, height int, topLeft, bottomRight rune) {
	DrawHorizontalLine(buffer, x, y, width, topLeft)
	DrawVerticalLine(buffer, x, y+1, height-2, topLeft)
	DrawHorizontalLine(buffer, x, y+height-1, width, bottomRight)
	DrawVerticalLine(buffer, x+width-1, y+1, height-2, bottomRight)
}

// DrawArenaBorder draws an inset floor with shaded walls and diagonal corner
// seams, like looking into an open box. Depth is the wall width in block cells.
// The floor stays untouched. Depth is clamped to leave at least one floor cell.
func DrawArenaBorder(buffer *Buffer, x, y, width, height, depth int) {
	if width < 5 || height < 5 || depth < 1 {
		return
	}
	depth = min(depth, (min(width, height)-3)/2)
	right := x + width - 1
	bottom := y + height - 1

	// Dark far walls and lighter near walls make the frame look recessed.
	for layer := 1; layer < depth; layer++ {
		drawFrameLayer(buffer, x+layer, y+layer,
			width-2*layer, height-2*layer, ':', '=')
	}
	drawFrameLayer(buffer, x, y, width, height, '#', '#')
	drawFrameLayer(buffer, x+depth, y+depth,
		width-2*depth, height-2*depth, '@', '@')

	// Join corresponding outer and inner corners with 45-degree seams.
	for offset := 0; offset <= depth; offset++ {
		buffer.Set(x+offset, y+offset, '@')
		buffer.Set(right-offset, y+offset, '@')
		buffer.Set(x+offset, bottom-offset, '@')
		buffer.Set(right-offset, bottom-offset, '@')
	}
}
