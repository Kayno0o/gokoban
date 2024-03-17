package src

import (
	"bufio"
	"bytes"
	"embed"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

var LevelsFs embed.FS

type Game struct {
	Screen              tcell.Screen
	currentLevel        Level
	levels              []Level
	State               string // play | end
	MaxWidth, MaxHeight int
}

func (g *Game) Init() error {
	file, err := LevelsFs.ReadFile("assets/levels.txt")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(file)

	scanner := bufio.NewScanner(reader)
	var lines []string
	var levels [][]string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(lines) > 0 {
				levels = append(levels, lines)
				lines = nil
			}
			continue
		}

		lines = append(lines, line)
	}

	if len(lines) > 0 {
		levels = append(levels, lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for i, group := range levels {
		level := LevelFromStrings(group)
		level.Index = i
		g.levels = append(g.levels, level)

		g.MaxHeight = int(math.Max(float64(g.MaxHeight), float64(level.Height)))
		g.MaxWidth = int(math.Max(float64(g.MaxWidth), float64(level.Width)))
	}

	g.currentLevel = g.levels[0]
	g.Draw()

	return nil
}

func (g *Game) HasNexLevel() bool {
	return g.currentLevel.Index+1 < len(g.levels)
}

func (g *Game) NextLevel() {
	if !g.HasNexLevel() {
		return
	}
	g.currentLevel = g.levels[g.currentLevel.Index+1]
}

func (g *Game) ResetLevel() {
	g.currentLevel = g.levels[g.currentLevel.Index]
	g.Draw()
}

func (g *Game) MovePlayer(dx, dy int) {
	isWin := g.currentLevel.MovePlayer(dx, dy)
	g.Draw()
	if isWin {
		time.Sleep(500 * time.Millisecond)

		if g.HasNexLevel() {
			g.NextLevel()
			g.Draw()
			return
		}

		g.End()
	}
}

func (g *Game) Draw() {
	g.Screen.Clear()

	level := g.currentLevel.ToString()
	startX := (g.MaxWidth - len(level[0])) / 2
	startY := (g.MaxHeight - len(level)) / 2
	g.DrawLines(Position{X: startX, Y: startY}, level)

	g.DrawLines(Position{g.MaxWidth + 3, 0}, []string{
		"Level " + strconv.Itoa(g.currentLevel.Index+1) + " / " + strconv.Itoa(len(g.levels)),
		"Press R to restart",
		"",
		"@ - Player",
		"# - Wall",
		"O - Box",
		"X - Box on target",
		". - Target",
	})

	g.Screen.Show()
}

func (g *Game) DrawLines(pos Position, lines []string) {
	for y, line := range lines {
		for x, ch := range line {
			g.Screen.SetContent(pos.X+x, pos.Y+y, ch, nil, tcell.StyleDefault)
		}
	}
}

func (g *Game) DrawEverywhereOnBoard(text string) {
	style := tcell.StyleDefault.Foreground(tcell.ColorRed)

	runes := []rune(text)
	for x := 0; x < g.MaxWidth; x++ {
		for y := 0; y < g.MaxHeight; y++ {
			g.Screen.SetContent(x, y, runes[x%len(runes)], nil, style)
		}
	}

	g.Screen.Show()
}

func (g *Game) End() {
	g.DrawEverywhereOnBoard(" SUPER ")
	time.Sleep(1 * time.Second)
	g.DrawEverywhereOnBoard("  HOT  ")
	time.Sleep(1 * time.Second)
	g.Screen.Clear()

	os.Exit(0)
}
