package ascii

import "testing"

func TestBeveledBorderLighting(t *testing.T) {
	for _, pattern := range []LinePattern{LineRaised, LineSunken} {
		b := NewBuffer(10, 10)
		b.Set(5, 5, 'X')
		DrawBorder(b, 1, 1, 8, 8, 2, pattern)
		light, dark := rune('@'), rune(':')
		if pattern == LineSunken {
			light, dark = dark, light
		}
		if b.At(4, 1) != light || b.At(1, 4) != light || b.At(4, 8) != dark || b.At(8, 4) != dark {
			t.Fatalf("incorrect edge lighting for %v", pattern)
		}
		if b.At(5, 5) != 'X' || b.At(0, 0) != ' ' {
			t.Fatal("border changed interior or exterior")
		}
	}
}

func TestDoubleBorderCorners(t *testing.T) {
	b := NewBuffer(10, 10)
	DrawBorder(b, 0, 0, 10, 10, 3, LineDouble)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			layer := min(x, y, 9-x, 9-y)
			want := rune(' ')
			if layer == 0 || layer == 2 {
				want = '@'
			}
			if b.At(x, y) != want {
				t.Fatalf("cell %d,%d = %q, want %q", x, y, b.At(x, y), want)
			}
		}
	}
}

func TestBorderClampingAndInvalidSizes(t *testing.T) {
	b := NewBuffer(4, 4)
	DrawBorder(b, 0, 0, 4, 4, 99, LineSolid)
	for _, cell := range b.Cells {
		if cell.Rune != '@' {
			t.Fatal("thickness clamp left a hole")
		}
	}
	b.Clear()
	DrawBorder(b, 0, 0, 1, 4, 2, LineSolid)
	DrawBorder(b, 0, 0, 4, 4, 0, LineSolid)
	for _, cell := range b.Cells {
		if cell.Rune != 0 {
			t.Fatal("invalid border drew cells")
		}
	}
}

func TestTextBorderStyle(t *testing.T) {
	b := NewBuffer(4, 3)
	DrawTextBorder(b, 0, 0, 4, 3, TextBorder)
	for y, row := range []string{"+--+", "|  |", "+--+"} {
		for x, r := range row {
			if b.At(x, y) != r {
				t.Fatalf("incorrect text border at %d,%d", x, y)
			}
		}
	}
}

func TestArenaBorderSeamsAndFloor(t *testing.T) {
	b := NewBuffer(18, 14)
	b.Set(8, 7, 'X')
	DrawArenaBorder(b, 1, 1, 16, 12, 3)
	for offset := 0; offset <= 3; offset++ {
		for _, point := range [][2]int{{1 + offset, 1 + offset}, {16 - offset, 1 + offset}, {1 + offset, 12 - offset}, {16 - offset, 12 - offset}} {
			if b.At(point[0], point[1]) != '@' {
				t.Fatalf("missing seam at %v", point)
			}
		}
	}
	if b.At(8, 2) != ':' || b.At(8, 11) != '=' {
		t.Fatal("wall shading is missing")
	}
	if b.At(8, 4) != '@' || b.At(8, 9) != '@' {
		t.Fatal("inner frame is missing")
	}
	if b.At(8, 7) != 'X' || b.At(0, 0) != ' ' {
		t.Fatal("floor or exterior changed")
	}
}

func TestArenaBorderDepthClamp(t *testing.T) {
	b := NewBuffer(7, 7)
	b.Set(3, 3, 'X')
	DrawArenaBorder(b, 0, 0, 7, 7, 100)
	if b.At(3, 3) != 'X' || b.At(2, 2) != '@' {
		t.Fatal("depth clamp did not preserve floor")
	}
	b.Clear()
	DrawArenaBorder(b, 0, 0, 4, 4, 1)
	DrawArenaBorder(b, 0, 0, 7, 7, 0)
	for _, cell := range b.Cells {
		if cell.Rune != 0 {
			t.Fatal("invalid arena drew cells")
		}
	}
}
