package main

import (
	"fmt"
	"math"

	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/renderer"
)

// renderWallHUD runs after the block and text layers.
func renderWallHUD(
	match *game.Match,
	config renderer.Config,
	showControls bool,
	missingController bool,
) {
	width := float32(config.WindowWidth)
	height := float32(config.WindowHeight)
	wallWidth := float32((arenaDepth + 1) * blockCellWidth)
	wallHeight := float32((arenaDepth + 1) * blockCellHeight)

	top := fmt.Sprintf(
		"round %d/%d | %ds",
		match.Round,
		match.Config.Rounds,
		int(math.Ceil(match.Remaining)),
	)
	detail := "two players"
	rightName := "player 2"

	if match.Machine != nil {
		detail = fmt.Sprintf(
			"bout %d/%d | %s | %s",
			match.Bout,
			match.TotalBouts,
			match.Opponent.Name,
			game.DifficultyName(match.Machine.Difficulty),
		)
	}

	renderer.RenderTextLine(
		top, width/2, wallHeight/2-13, 0, config,
	)
	renderer.RenderTextLine(
		detail, width/2, wallHeight/2+13, 0, config,
	)

	// Both labels face inward.
	renderer.RenderTextLine(
		fmt.Sprintf("player 1 | %d", match.Left.Score),
		wallWidth/2, height/2, -1, config,
	)
	renderer.RenderTextLine(
		fmt.Sprintf("%s | %d", rightName, match.Right.Score),
		width-wallWidth/2, height/2, 1, config,
	)

	if !showControls {
		return
	}

	bottom := height - wallHeight/2
	hint := "esc: pause"
	if missingController {
		hint = "connect player two controller"
	}

	renderer.RenderTextLine(
		fmt.Sprintf(
			"S speed | D extra hit | +%d bonus | X swap",
			game.BonusItemPoints,
		),
		width/2, bottom-13, 0, config,
	)
	renderer.RenderTextLine(
		hint, width/2, bottom+13, 0, config,
	)
}
