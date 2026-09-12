package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/andrlzpt/ufg/internal/ascii"
	"github.com/andrlzpt/ufg/internal/game"
	"github.com/andrlzpt/ufg/internal/renderer"
	"github.com/andrlzpt/ufg/internal/terminal"
	"github.com/andrlzpt/ufg/internal/window"
)

const (
	windowWidth     = 840
	windowHeight    = 624
	blockCellWidth  = 4
	blockCellHeight = 4
	textCellWidth   = 14
	textCellHeight  = 26
	arenaDepth      = 14
	playerPadding   = 15
)

type Screen string

const (
	ScreenMain  Screen = "main"
	ScreenArena Screen = "arena"
)

type Mode int

const (
	ModeOnePlayer Mode = iota
	ModeTwoPlayers
)

type AppState struct {
	Screen         Screen
	Mode           Mode
	MenuItem       int
	Config         game.MatchConfig
	Match          game.Match
	Championship   game.Championship
	Paused         bool
	Difficulty     int
	DifficultyMenu bool
}

func main() {
	runtime.LockOSThread()

	blockConfig := renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: blockCellWidth, CellHeight: blockCellHeight,
	}
	textConfig := renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: textCellWidth, CellHeight: textCellHeight,
	}
	blocks := ascii.NewBuffer(windowWidth/blockCellWidth, windowHeight/blockCellHeight)
	text := ascii.NewBuffer(windowWidth/textCellWidth, windowHeight/textCellHeight)

	app := AppState{
		Screen: ScreenMain,
		Config: game.MatchConfig{
			Rounds:          3,
			RoundDuration:   60,
			RoundBreak:      5,
			ArmSpeed:        30,
			RetractSpeedMul: 2,
		},
		Championship: game.NewChampionship(),
		Difficulty:   game.DifficultyEasy,
	}

	commands := make(chan terminal.Command, 16)
	fmt.Println("Type 'go arena' to enter the arena; 'go main' to return.")
	go func() {
		if err := terminal.ReadCommands(os.Stdin, commands); err != nil {
			fmt.Fprintln(os.Stderr, "reading commands:", err)
		}
	}()

	window.Run(window.Config{
		Width: windowWidth, Height: windowHeight, Title: "untitled fighting game",
	}, func(state *window.State) {
		handleCommands(commands, &app)
		menuFrame := handleMenus(state, &app, text)

		blocks.Clear()
		text.Clear()

		if app.Screen == ScreenMain {
			ascii.DrawBorder(
				blocks, 0, 0, blocks.Width, blocks.Height,
				3, ascii.LineZigzag,
			)
			drawAppMenu(text, &app)
		} else {
			ascii.DrawBackground(
				blocks, 0, 0, blocks.Width, blocks.Height,
				ascii.BackgroundWood,
			)
			ascii.DrawArenaBorder(
				blocks, 0, 0, blocks.Width, blocks.Height,
				arenaDepth,
			)

			leftInput := playerInput(state.Keyboard)
			rightInput := game.Input{}

			missingController := false
			if app.Mode == ModeTwoPlayers {
				missingController = !state.Controller.Connected
				rightInput = playerInput(state.Controller.Controls)
			}

			arena := game.Bounds{
				MinX: arenaDepth + 1,
				MinY: arenaDepth + 1,
				MaxX: blocks.Width - arenaDepth - 2,
				MaxY: blocks.Height - arenaDepth - 2,
			}

			minY := float64(arena.MinY)
			maxY := float64(arena.MaxY - ascii.PlayerHeight + 1)
			match := &app.Match

			if !menuFrame && !app.Paused &&
				!match.Over() && !missingController {
				match.Update(
					leftInput, rightInput,
					state.DeltaTime, minY, maxY, arena,
				)

				if match.Over() {
					app.MenuItem = 0
				}
			}
			drawItems(blocks, match.Items)
			drawPlayer(blocks, &match.Left, false)
			drawPlayer(blocks, &match.Right, true)

			showMenu := app.Paused || match.Over()
			drawMatchHUD(text, match, !showMenu)
			if !showMenu {
				drawCenteredText(
					text, 7,
					"S speed | D extra hit | + point | X swap",
				)
			}

			if showMenu {
				drawAppMenu(text, &app)
			} else {
				if missingController {
					drawCenteredText(
						text, 5,
						"connect player two controller",
					)
				}
				drawCenteredText(text, text.Height-2, "esc: pause")
			}
		}

		renderer.BeginFrame(windowWidth, windowHeight)
		renderer.RenderBlocks(blocks, blockConfig)
		renderer.RenderText(text, textConfig)
	})
}

func drawCenteredText(buffer *ascii.Buffer, row int, text string) {
	ascii.DrawText(buffer, (buffer.Width-len([]rune(text)))/2, row, text)
}

func handleCommands(commands <-chan terminal.Command, app *AppState) {
	for {
		select {
		case command, ok := <-commands:
			if !ok {
				return
			}
			handleCommand(command, app)
		default:
			return
		}
	}
}

func handleCommand(command terminal.Command, app *AppState) {
	fields := strings.Fields(command.Raw)
	if len(fields) == 0 {
		return
	}

	if fields[0] == "difficulty" {
		if len(fields) == 1 {
			fmt.Println(
				"Selected difficulty:",
				game.DifficultyName(app.Difficulty),
			)

			if app.Match.Machine != nil {
				fmt.Println(
					"Current match difficulty:",
					game.DifficultyName(app.Match.Machine.Difficulty),
				)
			}
			return
		}

		if len(fields) == 2 {
			level := 0

			switch strings.ToLower(fields[1]) {
			case "easy", "1":
				level = game.DifficultyEasy
			case "normal", "2":
				level = game.DifficultyNormal
			case "hard", "3":
				level = game.DifficultyHard
			}

			if level != 0 {
				app.Difficulty = level
				fmt.Printf(
					"Difficulty %s selected. Applies to the next match.\n",
					game.DifficultyName(level),
				)
				return
			}
		}

		fmt.Println("Use: difficulty easy|normal|hard")
		return
	}

	if len(fields) == 2 && fields[0] == "go" {
		switch fields[1] {
		case "main":
			returnToMain(app)
			return
		case "arena":
			startMatch(app)
			return
		}
	}

	fmt.Println("Use: go arena, go main, or difficulty easy|normal|hard")
}
