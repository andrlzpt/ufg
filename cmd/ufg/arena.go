package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

func drawArenaFeatures(
	buffer *ascii.Buffer,
	match *game.Match,
	arena game.Arena,
) {
	if match.Config.FeedEnabled {
		// Plain floor makes collected trails visible.
		for y := arena.MinY; y <= arena.MaxY; y++ {
			ascii.DrawHorizontalLine(
				buffer,
				arena.MinX,
				y,
				arena.MaxX-arena.MinX+1,
				' ',
			)
		}

		for cell := range match.Pellets {
			buffer.Set(cell.X, cell.Y, '•')
		}
	}

	for cell := range match.Walls {
		buffer.Set(cell.X, cell.Y, '#')
	}
}
