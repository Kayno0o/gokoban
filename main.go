package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"gokoban.kaynooo.xyz/src"
)

//go:embed assets/*
var LevelsFs embed.FS

func main() {
	src.LevelsFs = LevelsFs

	g := src.Game{}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer s.Fini()

	g.Screen = s
	err = g.Init()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyESC, tcell.KeyCtrlC:
				return
			case tcell.KeyRune:
				if ev.Rune() == 'q' {
					return
				}
				if ev.Rune() == 'r' {
					g.ResetLevel()
					break
				}
			case tcell.KeyUp:
				g.MovePlayer(0, -1)
			case tcell.KeyDown:
				g.MovePlayer(0, 1)
			case tcell.KeyLeft:
				g.MovePlayer(-1, 0)
			case tcell.KeyRight:
				g.MovePlayer(1, 0)
			}
		}
	}
}
