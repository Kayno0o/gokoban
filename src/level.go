package src

import (
	"strings"
)

type Level struct {
	Width, Height int
	Player        Position
	Boxes         []Position
	Targets       []Position
	Walls         []Position
	Index         int
}

func LevelFromStrings(lines []string) Level {
	level := Level{
		Width:  len(lines[0]),
		Height: len(lines),
	}

	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#':
				level.Walls = append(level.Walls, Position{x, y})
			case '@':
				level.Player = Position{x, y}
			case 'O':
				level.Boxes = append(level.Boxes, Position{x, y})
			case 'X':
				level.Boxes = append(level.Boxes, Position{x, y})
				level.Targets = append(level.Targets, Position{x, y})
			case '.':
				level.Targets = append(level.Targets, Position{x, y})
			}
		}
	}

	return level
}

// Return true if level completed
func (l *Level) MovePlayer(dx, dy int) bool {
	pos := l.Player.Movement(dx, dy)

	// Check if the new position is out of bounds
	if pos.X < 0 || pos.X >= l.Width || pos.Y < 0 || pos.Y >= l.Height {
		return false
	}

	if pos.Some(l.Walls) {
		return false
	}

	for i, box := range l.Boxes {
		if pos.Eq(box) {
			boxPos := box.Movement(dx, dy)
			if boxPos.Some(append(l.Boxes, l.Walls...)) {
				return false
			}
			l.Boxes[i] = boxPos
		}
	}

	l.Player = pos

	for _, box := range l.Boxes {
		if !box.Some(l.Targets) {
			return false
		}
	}

	return true
}

func (l *Level) ToString() []string {
	lines := make([]string, l.Height)
	for i := range lines {
		lines[i] = strings.Repeat(" ", l.Width)
	}

	for _, wall := range l.Walls {
		lines[wall.Y] = replaceAtIndex(lines[wall.Y], wall.X, "#")
	}

	for _, target := range l.Targets {
		lines[target.Y] = replaceAtIndex(lines[target.Y], target.X, ".")
	}

	for _, box := range l.Boxes {
		if box.Some(l.Targets) {
			lines[box.Y] = replaceAtIndex(lines[box.Y], box.X, "X")
		} else {
			lines[box.Y] = replaceAtIndex(lines[box.Y], box.X, "O")
		}
	}

	lines[l.Player.Y] = replaceAtIndex(lines[l.Player.Y], l.Player.X, "@")

	return lines
}

func replaceAtIndex(str string, index int, replacement string) string {
	return str[:index] + replacement + str[index+1:]
}
