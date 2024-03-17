package src

import (
	"reflect"
	"testing"
)

func TestLevelFromStrings(t *testing.T) {
	level := LevelFromStrings([]string{
		"######",
		"#@OX.#",
		"######",
	})

	expectedLevel := Level{
		Width:  6,
		Height: 3,
		Player: Position{X: 1, Y: 1},
		Boxes: []Position{
			{X: 2, Y: 1},
			{X: 3, Y: 1},
		},
		Targets: []Position{
			{X: 3, Y: 1},
			{X: 4, Y: 1},
		},
		Walls: []Position{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0},
			{X: 0, Y: 1}, {X: 5, Y: 1},
			{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2},
		},
	}

	if !reflect.DeepEqual(level, expectedLevel) {
		t.Errorf("Level loaded incorrectly.\nExpected %+v\nGot      %+v", expectedLevel, level)
	}
}
