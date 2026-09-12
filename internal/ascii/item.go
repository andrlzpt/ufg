package ascii

// DrawItem draws a 3-by-5 icon centered on x,y.
func DrawItem(buffer *Buffer, x, y int, symbol rune) {
	var rows [5]string

	switch symbol {
	case 'S':
		rows = [5]string{
			"@@@",
			"@  ",
			"@@@",
			"  @",
			"@@@",
		}
	case 'D':
		rows = [5]string{
			"@@ ",
			"@ @",
			"@ @",
			"@ @",
			"@@ ",
		}
	case '+':
		rows = [5]string{
			" @ ",
			" @ ",
			"@@@",
			" @ ",
			" @ ",
		}
	case 'X':
		rows = [5]string{
			"@ @",
			"@ @",
			" @ ",
			"@ @",
			"@ @",
		}
	default:
		return
	}

	for row, line := range rows {
		for column, r := range line {
			if r != ' ' {
				buffer.Set(x-1+column, y-2+row, r)
			}
		}
	}
}
