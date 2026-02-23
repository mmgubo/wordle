# Wordle - Terminal Game in Go

A terminal-based Wordle clone written in Go with ANSI color output and cross-platform support.

## Project Structure

```
wordle/
├── main.go               # Core game logic and embedded word list
├── terminal_windows.go   # Windows: enables ANSI Virtual Terminal Processing
├── terminal_unix.go      # Unix/Linux/Mac: no-op (ANSI supported natively)
├── go.mod                # Go module (Go 1.21, no external dependencies)
└── .gitignore            # Excludes compiled binaries
```

## Build & Run

```bash
# Build
go build -o wordle        # Unix/Linux/Mac
go build -o wordle.exe    # Windows

# Run directly without building
go run .

# Run compiled binary
./wordle        # Unix/Linux/Mac
wordle.exe      # Windows
```

## How the Game Works

- A random 5-letter word is chosen from a built-in list of 515 possible answers
- Players have 6 attempts to guess the word
- After each guess, letters are color-coded:
  - **Green** — correct letter, correct position
  - **Yellow** — correct letter, wrong position
  - **Gray** — letter not in the word
- The screen redraws after each guess showing the board and a QWERTY keyboard with letter states

## Key Code

**`evaluate(guess, target string) [5]tileState`** (`main.go`)
Uses a two-pass algorithm to correctly handle duplicate letters:
1. First pass marks exact matches (green) and counts remaining letters
2. Second pass marks present-but-misplaced letters (yellow) using remaining frequency counts

**`draw(guesses []string, results [][5]tileState)`** (`main.go`)
Renders the full UI on each turn: clears the screen, draws the 6×5 tile grid, and renders a QWERTY keyboard tracking the best-known state for each letter.

**`enableANSI()`** (`terminal_windows.go` / `terminal_unix.go`)
Build-tag-selected function. On Windows, calls `SetConsoleMode` via `syscall` to enable ANSI escape code support. On Unix it is a no-op.

## Architecture Notes

- **No external dependencies** — standard library only (`bufio`, `fmt`, `math/rand`, `os`, `strings`, `time`, `syscall`)
- **Cross-platform builds** via Go build tags: `//go:build windows` and `//go:build !windows`
- **Word list embedded in source** — 515 words defined as a string slice in `main.go`; no external files needed at runtime
- **Input validation** — guesses are checked for length, alphabetic characters, and membership in the valid word list before being accepted
