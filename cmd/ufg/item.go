package main

import (
	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
)

func drawItems(buffer *ascii.Buffer, items []game.Item) {
	for _, item := range items {
		var symbol rune

		switch item.Kind {
		case game.ItemDoubleSpeed:
			symbol = 'S'
		case game.ItemDoublePunches:
			symbol = 'D'
		case game.ItemPoint:
			symbol = '+'
		case game.ItemSwapScores:
			symbol = 'X'
		}

		ascii.DrawItem(
			buffer,
			item.Position.X,
			item.Position.Y,
			symbol,
		)
	}
}
