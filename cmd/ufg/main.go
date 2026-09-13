package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
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
	ArenaMenu      bool
	Controls       ControlConfig
	ControlsMenu   bool
	Controllers    []window.Controller
	MissingDevice  string
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
			ArmSpeed:        game.DefaultArmSpeed,
			RetractSpeedMul: game.DefaultRetractspeedMul,
			MazeEnabled:     true,
			FeedEnabled:     true,
		},
		Championship: game.NewChampionship(),
		Difficulty:   game.DifficultyEasy,
		Controls: ControlConfig{
			Devices: [2]int{keyboardDevice, 0},
		},
	}

	commands := make(chan terminal.Command, 16)
	fmt.Println("Command: speed [number]")
	go func() {
		if err := terminal.ReadCommands(os.Stdin, commands); err != nil {
			fmt.Fprintln(os.Stderr, "reading commands:", err)
		}
	}()

	window.Run(window.Config{
		Width: windowWidth, Height: windowHeight, Title: "untitled fighting game",
	}, func(state *window.State) {
		app.Controllers = state.Controllers
		handleCommands(commands, &app)

		_, _, app.MissingDevice = matchInputs(state, &app)

		if app.Screen == ScreenArena &&
			!app.Match.Over() &&
			app.MissingDevice != "" &&
			!app.Paused {
			app.Paused = true
			app.MenuItem = 2 // Controls.
		}

		menuFrame := handleMenus(state, &app, text)

		// Menu selection may have changed mode or device assignments.
		leftInput, rightInput, missingDevice := matchInputs(state, &app)
		app.MissingDevice = missingDevice

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

			arena := game.Arena{
				Bounds: game.Bounds{
					MinX: arenaDepth + 1,
					MinY: arenaDepth + 1,
					MaxX: blocks.Width - arenaDepth - 2,
					MaxY: blocks.Height - arenaDepth - 2,
				},
			}

			minY := float64(arena.MinY)
			maxY := float64(arena.MaxY - ascii.PlayerHeight + 1)
			match := &app.Match

			if !menuFrame && !app.Paused &&
				!match.Over() && missingDevice == "" {
				match.Update(
					leftInput, rightInput,
					state.DeltaTime, minY, maxY, arena,
				)

				if match.Over() {
					app.MenuItem = 0
				}
			}
			drawArenaFeatures(blocks, match, arena)
			drawItems(blocks, match.Items)
			drawPlayer(blocks, &match.Left, false)
			drawPlayer(blocks, &match.Right, true)

			showMenu := app.Paused || match.Over()
			drawMatchHUD(text, match, !showMenu && missingDevice == "")

			if showMenu {
				drawAppMenu(text, &app)
			}
		}

		renderer.BeginFrame(windowWidth, windowHeight)
		renderer.RenderBlocks(blocks, blockConfig)
		renderer.RenderText(text, textConfig)

		if app.Screen == ScreenArena {
			renderWallHUD(
				&app.Match,
				textConfig,
				!app.Paused && !app.Match.Over(),
				app.Mode == ModeTwoPlayers && app.MissingDevice != "",
			)
		}
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

	if fields[0] == "speed" {
		if len(fields) == 1 {
			fmt.Printf("Selected arm speed: %g cells/s\n", app.Config.ArmSpeed)

			if app.Screen == ScreenArena {
				fmt.Printf(
					"Current match arm speed: %g cells/s\n",
					app.Match.Config.ArmSpeed,
				)
			}
			return
		}

		if len(fields) == 2 {
			speed, err := strconv.ParseFloat(fields[1], 64)
			if err == nil &&
				speed > 0 &&
				!math.IsNaN(speed) &&
				!math.IsInf(speed, 0) {
				app.Config.ArmSpeed = speed
				fmt.Printf(
					"Arm speed %g cells/s. Applies to next match.\n",
					speed,
				)
				return
			}
		}

		fmt.Println("Use: speed <positive number>, e.g. speed 45")
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

	fmt.Println("Use: go arena, go main, or speed [number]")
}
