package ascii

import "testing"

func TestLinePatternSamples(t *testing.T) {
	cases := []struct {
		name              string
		pattern           LinePattern
		length, thickness int
		want              []string
	}{
		{"dotted", LineDotted, 7, 1, []string{"@  @  @"}},
		{"dashed", LineDashed, 10, 1, []string{"@@@@@   @@"}},
		{"checkered", LineCheckered, 4, 2, []string{"@:@:", ":@:@"}},
		{"gradient", LineGradient, 10, 1, []string{" .:-=+*#%@"}},
		{"fade", LineFadeOut, 10, 1, []string{"@%#*+=-:. "}},
		{"raised", LineRaised, 2, 3, []string{"@@", "++", "  "}},
		{"sunken", LineSunken, 2, 3, []string{"  ", "==", "@@"}},
		{"double", LineDouble, 3, 3, []string{"@@@", "   ", "@@@"}},
		{"zigzag", LineZigzag, 5, 3, []string{"@   @", " @ @ ", "  @  "}},
		{"thin zigzag", LineZigzag, 3, 1, []string{"@@@"}},
		{"single gradient", LineGradient, 1, 1, []string{"@"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewBuffer(tc.length, tc.thickness)
			v := NewBuffer(tc.thickness, tc.length)
			DrawPatternedHorizontalLine(h, 0, 0, tc.length, tc.thickness, tc.pattern)
			DrawPatternedVerticalLine(v, 0, 0, tc.length, tc.thickness, tc.pattern)
			for y, row := range tc.want {
				for x, want := range row {
					if h.At(x, y) != want || v.At(y, x) != want {
						t.Fatalf("cell (%d,%d): horizontal %q, vertical %q, want %q", x, y, h.At(x, y), v.At(y, x), want)
					}
				}
			}
		})
	}
}

func TestPatternClippingPreservesPhase(t *testing.T) {
	b := NewBuffer(3, 2)
	DrawPatternedHorizontalLine(b, -1, 0, 5, 1, LineDotted)
	if b.At(0, 0) != ' ' || b.At(1, 0) != ' ' || b.At(2, 0) != '@' {
		t.Fatal("clipping changed dotted pattern phase")
	}
	DrawPatternedHorizontalLine(b, 0, 1, 0, 2, LineSolid)
	DrawPatternedVerticalLine(b, 0, 1, 3, -1, LineSolid)
	for x := 0; x < b.Width; x++ {
		if b.At(x, 1) != ' ' {
			t.Fatal("nonpositive dimensions drew cells")
		}
	}
}
