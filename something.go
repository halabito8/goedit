package main

import (
	"os"

	"golang.org/x/term"
)

type rawTerminal struct {
	width  int
	height int
	term   *term.Terminal
}

func (sh *rawTerminal) Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

func (sh *rawTerminal) Write(b []byte) (int, error) {
	return os.Stdout.Write(b[:])
}

func editorRefreshScreen(term *term.Terminal) {
	// Clear screen
	term.Write([]byte("\x1b[2J"))
	// Reposition cursor to top left
	term.Write([]byte("\x1b[H"))
}

func die(term *term.Terminal) {
	// Clear screen
	term.Write([]byte("\x1b[2J"))
	// Reposition cursor to top left
	term.Write([]byte("\x1b[H"))
}

func main() {
	oldstate, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldstate)

	terminal := rawTerminal{}

	terminal.term = term.NewTerminal(&rawTerminal{}, "")

	terminal.width, terminal.height, err = term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		die(terminal.term)
	}

	editorRefreshScreen(terminal.term)

	var b [1]byte
	for {
		os.Stdin.Read(b[:])
		if b[0] == byte(3) { // Ctrl-C
			die(terminal.term)
			break
		}

		terminal.Write(b[:])
	}
}
