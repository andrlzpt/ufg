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
	Screen   Screen
	Mode     Mode
	MenuItem int
	Left     game.Player
	Right    game.Player
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

	playerY := float64(blocks.Height-ascii.PlayerHeight) / 2

	app := AppState{
		Screen: ScreenMain,
		Left: game.Player{
			X:          playerPadding,
			Y:          playerY,
			Speed:      40,
			Facing:     game.FacingRight,
			ArmControl: game.ArmControlLastPressed,
		},
		Right: game.Player{
			X:          float64(blocks.Width - playerPadding - ascii.PlayerWidth),
			Y:          playerY,
			Speed:      40,
			Facing:     game.FacingLeft,
			ArmControl: game.ArmControlLastPressed,
		},
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
		if app.Screen == ScreenMain {
			handleMenu(state, &app, text)
		}
		blocks.Clear()
		text.Clear()
		if app.Screen == ScreenArena {
			ascii.DrawBackground(blocks, 0, 0, blocks.Width, blocks.Height, ascii.BackgroundWood)
			ascii.DrawArenaBorder(blocks, 0, 0, blocks.Width, blocks.Height, arenaDepth)

			leftInput := playerInput(state.Keyboard)
			rightInput := game.Input{}

			paused := false
			if app.Mode == ModeTwoPlayers {
				paused = !state.Controller.Connected
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

			if !app.Left.MatchOver(&app.Right) {
				if paused {
					drawCenteredText(text, text.Height/2,
						"connect player two controller")
				} else {
					app.Left.Update(leftInput, state.DeltaTime, minY, maxY, arena, &app.Right)
					app.Right.Update(rightInput, state.DeltaTime, minY, maxY, arena, &app.Left)
				}
			}

			drawPlayer(blocks, &app.Left, false)
			drawPlayer(blocks, &app.Right, true)

			drawCenteredText(
				text,
				1,
				fmt.Sprintf("player 1: %d    player 2: %d",
					app.Left.Score, app.Right.Score),
			)

			if app.Left.MatchOver(&app.Right) {
				winner := "player 1"
				if app.Left.KnockedOut {
					winner = "player2"
				}

				drawCenteredText(
					text,
					text.Height/2,
					winner+" wins by knockout",
				)

				drawCenteredText(
					text,
					text.Height/2+2,
					"type 'go main' to return",
				)
			}

		} else {
			ascii.DrawBorder(blocks, 0, 0, blocks.Width, blocks.Height, 3, ascii.LineZigzag)
			drawMenu(text, app.MenuItem)
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
	if len(fields) == 2 && fields[0] == "go" {
		switch fields[1] {
		case "main":
			app.Screen = ScreenMain
			return
		case "arena":
			app.Screen = ScreenArena
			return
		}
	}
	fmt.Println("Use: go arena or go main")
}
