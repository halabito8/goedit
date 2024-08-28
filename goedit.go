package main

import (
	"os"
	"strconv"

	"golang.org/x/term"
)

type rawTerminal struct {
	width    int
	height   int
	term     *term.Terminal
	oldState *term.State
	x        int
	y        int
}

// Read reads input from the terminal
func (rt *rawTerminal) Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

// Write writes output to the terminal
func (rt *rawTerminal) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

// editorRefreshScreen clears the screen and repositions the cursor
func editorRefreshScreen(rt *rawTerminal) {
	rt.term.Write([]byte("\x1b[2J")) // Clear screen
	rt.term.Write([]byte("\x1b[H"))  // Reposition cursor to top left
}

// die restores the terminal state and clears the screen
func die(rt *rawTerminal) {
	term.Restore(int(os.Stdin.Fd()), rt.oldState) // Restore terminal state
	rt.term.Write([]byte("\x1b[2J"))              // Clear screen
	rt.term.Write([]byte("\x1b[H"))               // Reposition cursor to top left
}

func main() {
	var err error
	terminal := &rawTerminal{
		x: 0,
		y: 0,
	}

	// Enter raw mode
	terminal.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), terminal.oldState)

	// Initialize terminal and get size
	terminal.term = term.NewTerminal(terminal, "")
	terminal.width, terminal.height, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		die(terminal)
		return
	}

	// Refresh screen
	editorRefreshScreen(terminal)

	// Main loop for reading input
	buffer := make([]byte, 10)
	for {
		n, err := terminal.Read(buffer)
		if err != nil {
			die(terminal)
			return
		}

		// Exit on Ctrl-C
		if buffer[0] == byte(3) {
			die(terminal)
			break
		}

		// Handle escape sequences
		if buffer[0] == '\x1b' {
			if buffer[1] == '[' {
				if buffer[2] >= '0' && buffer[2] <= '9' {
					if buffer[3] == '~' {
						switch buffer[2] {
						case '5': // Page Up
							terminal.Write([]byte("\x1b[H"))
						case '6': // Page Down
							terminal.Write([]byte("\x1b[" + strconv.Itoa(terminal.height) + ";1H"))
						}
					}
				} else {
					switch buffer[2] {
					case 'A': // Up
						if terminal.y != 0 {
							terminal.y--
						}
					case 'B': // Down
						if terminal.y != terminal.height-1 {
							terminal.y--
						}
					case 'C': // Right
						if terminal.x != terminal.width-1 {
							terminal.x++
						}
					case 'D': // Left
						if terminal.x != 0 {
							terminal.x--
						}
					case '5': // Page
					}
				}
			}
		}

		// Write buffer to terminal
		terminal.Write(buffer[:n])
	}
}
