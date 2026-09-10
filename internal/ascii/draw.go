package ascii

func DrawText(buffer *Buffer, x, y int, text string) {
	column := x

	for _, r := range text {
		buffer.Set(column, y, r)
		column++
	}
}

type TextBorderStyle struct {
	Horizontal  rune
	Vertical    rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
}

var TextBorder = TextBorderStyle{
	Horizontal:  '-',
	Vertical:    '|',
	TopLeft:     '+',
	TopRight:    '+',
	BottomLeft:  '+',
	BottomRight: '+',
}

func DrawHorizontalLine(buffer *Buffer, x, y, length int, r rune) {
	for offset := range length {
		buffer.Set(x+offset, y, r)
	}
}

func DrawVerticalLine(buffer *Buffer, x, y, length int, r rune) {
	for offset := range length {
		buffer.Set(x, y+offset, r)
	}
}

func DrawTextBorder(
	buffer *Buffer,
	x, y, width, height int,
	style TextBorderStyle,
) {
	if width < 2 || height < 2 {
		return
	}

	right := x + width - 1
	bottom := y + height - 1

	DrawHorizontalLine(buffer, x+1, y, width-2, style.Horizontal)
	DrawHorizontalLine(buffer, x+1, bottom, width-2, style.Horizontal)

	DrawVerticalLine(buffer, x, y+1, height-2, style.Vertical)
	DrawVerticalLine(buffer, right, y+1, height-2, style.Vertical)

	buffer.Set(x, y, style.TopLeft)
	buffer.Set(right, y, style.TopRight)
	buffer.Set(x, bottom, style.BottomLeft)
	buffer.Set(right, bottom, style.BottomRight)
}
