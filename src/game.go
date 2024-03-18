package src

import (
	"bufio"
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	currentLevel               Level
	levels                     []Level
	State                      string // play | end | next
	MaxWidth, MaxHeight, Scale int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.MaxWidth * g.Scale, g.MaxHeight * g.Scale
}

func (g *Game) Init(assets embed.FS) error {
	file, err := assets.ReadFile("assets/levels.txt")
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
}

func (g *Game) MovePlayer(dx, dy int) {
	isWin := g.currentLevel.MovePlayer(dx, dy)
	if isWin {
		if g.HasNexLevel() && g.State != "next" {
			g.State = "next"
			setTimeout(500*time.Millisecond, func() {
				g.NextLevel()
				g.State = ""
			})
			return
		}

		g.State = "win"
	}
}

func (g *Game) Update() error {
	if repeatingKeyPressed(ebiten.KeyEscape) || repeatingKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyR) {
		g.ResetLevel()
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyArrowUp) {
		g.MovePlayer(0, -1)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyArrowDown) {
		g.MovePlayer(0, 1)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyArrowLeft) {
		g.MovePlayer(-1, 0)
		return nil
	}

	if repeatingKeyPressed(ebiten.KeyArrowRight) {
		g.MovePlayer(1, 0)
		return nil
	}

	return nil
}

func drawRect(img *image.RGBA, clr color.Color) {
	// Draw the rectangle onto the image
	draw.Draw(img, img.Bounds(), &image.Uniform{clr}, image.Point{}, draw.Src)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	level := g.currentLevel.ToString()
	startX := float64(((g.MaxWidth - len(level[0])) * g.Scale) / 2)
	startY := float64(((g.MaxHeight - len(level)) * g.Scale) / 2)

	for y, row := range level {
		for x, char := range row {
			rect := image.Rect(0, 0, g.Scale, g.Scale)
			rectImage := image.NewRGBA(rect)
			if char == '@' {
				drawRect(rectImage, color.RGBA{255, 255, 0, 255})
			} else if char == '#' {
				drawRect(rectImage, color.RGBA{255, 255, 255, 255})
			} else if char == 'X' {
				drawRect(rectImage, color.RGBA{0, 150, 0, 255})
			} else if char == 'O' {
				drawRect(rectImage, color.RGBA{0, 0, 150, 255})
			} else if char == '.' {
				drawRect(rectImage, color.RGBA{252, 132, 3, 255})
			} else if char == ' ' {
				drawRect(rectImage, color.RGBA{0, 0, 0, 255})
			}
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(startX+float64(x*g.Scale), startY+float64(y*g.Scale))

			img := ebiten.NewImageFromImage(rectImage)
			screen.DrawImage(img, options)
		}
	}
}

func (g *Game) DrawLines(screen *ebiten.Image, pos Position, lines []string) {
	ebitenutil.DebugPrintAt(screen, strings.Join(lines, "\n"), pos.X, pos.Y)
}
