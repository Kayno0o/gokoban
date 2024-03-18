package src

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func repeatingKeyPressed(key ebiten.Key) bool {
	d := inpututil.KeyPressDuration(key)
	return d == 1
}

func setTimeout(delay time.Duration, callback func()) {
	go func() {
		time.Sleep(delay)
		callback()
	}()
}
