package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/andrlzpt/ufg/internal/ascii"
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
)

type Screen string

const (
	ScreenMain  Screen = "main"
	ScreenArena Screen = "arena"
)

type AppState struct {
	Screen Screen
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
	app := AppState{Screen: ScreenMain}
	commands := make(chan terminal.Command, 16)
	fmt.Println("Type 'go arena' to enter the arena; 'go main' to return.")
	go func() {
		if err := terminal.ReadCommands(os.Stdin, commands); err != nil {
			fmt.Fprintln(os.Stderr, "reading commands:", err)
		}
	}()

	window.Run(window.Config{
		Width: windowWidth, Height: windowHeight, Title: "untitled fighting game",
	}, func(_ *window.State) {
		handleCommands(commands, &app)
		blocks.Clear()
		text.Clear()
		if app.Screen == ScreenArena {
			ascii.DrawBackground(blocks, 0, 0, blocks.Width, blocks.Height, ascii.BackgroundWood)
			ascii.DrawArenaBorder(blocks, 0, 0, blocks.Width, blocks.Height, 14)
			drawCenteredText(text, text.Height/2, "type: go main to get back")
		} else {
			ascii.DrawBorder(blocks, 0, 0, blocks.Width, blocks.Height, 3, ascii.LineZigzag)
			drawCenteredText(text, text.Height/2-1, "untitled fighting game.")
			drawCenteredText(text, text.Height/2, "type 'go arena' to enter the arena")
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
