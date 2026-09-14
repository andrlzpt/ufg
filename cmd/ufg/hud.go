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
	// Stay inside the outer and inner wall borders.
	const margin = float32(5)
	dx := float32(arenaDepth*blockCellWidth) - margin
	dy := float32(arenaDepth*blockCellHeight) - margin

	top := renderer.Wall{
		Width: width - 2*margin, Height: 52,
		Corners: [4]renderer.PixelPoint{
			{X: margin, Y: margin},
			{X: width - margin, Y: margin},
			{X: width - dx, Y: dy},
			{X: dx, Y: dy},
		},
	}
	bottom := renderer.Wall{
		Width: width - 2*margin, Height: 52,
		Corners: [4]renderer.PixelPoint{
			{X: dx, Y: height - dy},
			{X: width - dx, Y: height - dy},
			{X: width - margin, Y: height - margin},
			{X: margin, Y: height - margin},
		},
	}

	// Left reads upward. Right reads downward.
	left := renderer.Wall{
		Width: height - 2*margin, Height: 52,
		Corners: [4]renderer.PixelPoint{
			{X: margin, Y: height - margin},
			{X: margin, Y: margin},
			{X: dx, Y: dy},
			{X: dx, Y: height - dy},
		},
	}
	right := renderer.Wall{
		Width: height - 2*margin, Height: 52,
		Corners: [4]renderer.PixelPoint{
			{X: width - margin, Y: margin},
			{X: width - margin, Y: height - margin},
			{X: width - dx, Y: height - dy},
			{X: width - dx, Y: dy},
		},
	}

	top.Text(fmt.Sprintf(
		"round %d/%d | %ds",
		match.Round,
		match.Config.Rounds,
		int(math.Ceil(match.Remaining)),
	), 13, config)

	detail := "two players"
	if match.Machine != nil {
		detail = fmt.Sprintf(
			"bout %d/%d | %s | %s",
			match.Bout,
			match.TotalBouts,
			match.Opponent.Name,
			game.DifficultyName(match.Machine.Difficulty),
		)
	}
	top.Text(detail, 39, config)

	renderPlayerWall(left, "player 1", &match.Left, config)
	renderPlayerWall(right, "player 2", &match.Right, config)

	if !showControls {
		return
	}

	bottom.Text(fmt.Sprintf(
		"S speed | D extra hit | +%d bonus | X swap",
		game.BonusItemPoints,
	), 13, config)

	hint := "esc/start: pause"
	if missingController {
		hint = "device missing. open controls."
	}
	bottom.Text(hint, 39, config)
}

func renderPlayerWall(
	wall renderer.Wall,
	name string,
	player *game.Player,
	config renderer.Config,
) {
	speedActive := player.SpeedRemaining > 0
	left := punchIndicator("L", &player.LeftArm)
	right := punchIndicator("R", &player.RightArm)
	punchActive := left != "" || right != ""

	const speedLabel = "SPEED BOOST!"
	const punchLabel = "DOUBLE PUNCH!"
	const separator = " | "

	label := ""
	if speedActive {
		label = speedLabel
	}
	if punchActive {
		if label != "" {
			label += separator
		}
		label += punchLabel
	}

	if label == "" {
		wall.Text(
			fmt.Sprintf("%s | %d", name, player.Score),
			13,
			config,
		)
		return
	}

	wall.Text(label, 13, config)

	speedCenter := wall.Width / 2
	punchCenter := wall.Width / 2
	cellWidth := float32(config.CellWidth)

	if speedActive && punchActive {
		start := (wall.Width - float32(len(label))*cellWidth) / 2
		speedCenter = start + float32(len(speedLabel))*cellWidth/2
		punchCenter = start +
			float32(len(speedLabel)+len(separator))*cellWidth +
			float32(len(punchLabel))*cellWidth/2
	}

	if speedActive {
		fraction := float32(
			player.SpeedRemaining / game.DoubleSpeedDuration,
		)
		wall.Bar(speedCenter, 35, 140, 8, fraction, config)
	}

	if punchActive {
		counts := ""
		if left != "" {
			counts = fmt.Sprintf("LEFT:x%d", player.LeftArm.HitsRemaining)
		}
		if right != "" {
			if counts != "" {
				counts += " | "
			}
			counts += fmt.Sprintf("RIGHT:x%d", player.RightArm.HitsRemaining)
		}

		wall.TextAt(counts, punchCenter, 39, config)
	}
}

func punchIndicator(label string, arm *game.Arm) string {
	if !arm.DoublePunches ||
		arm.Retracting ||
		len(arm.Path) == 0 ||
		arm.HitsRemaining <= 0 {
		return ""
	}

	return fmt.Sprintf("DOUBLE PUNCH! %s:X%d", label, arm.HitsRemaining)
}
