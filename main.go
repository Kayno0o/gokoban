package main

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"gokoban.kaynooo.xyz/src"
)

//go:embed assets/*
var AssetsFS embed.FS

func main() {
	g := src.Game{}
	g.Scale = 64

	err := g.Init(AssetsFS)
	if err != nil {
		log.Println(err)
		return
	}

	ebiten.SetWindowSize(g.MaxWidth*g.Scale, g.MaxHeight*g.Scale)
	ebiten.SetWindowTitle("GOKOBAN")

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}
