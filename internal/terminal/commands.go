package terminal

import (
	"bufio"
	"io"
	"strings"
)

type Command struct {
	Raw string
}

// ReadCommands sends nonempty, trimmed lines from input to commands.
// It blocks until input ends or an error occurs and leaves commands open.
// Run it in a goroutine to read terminal input alongside the window loop.
func ReadCommands(input io.Reader, commands chan<- Command) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		commands <- Command{Raw: line}
	}
	return scanner.Err()
}
