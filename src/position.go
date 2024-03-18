package src

type Position struct {
	X, Y int
}

func (p Position) Eq(target Position) bool {
	return p.X == target.X && p.Y == target.Y
}

func (p Position) Some(targets []Position) bool {
	for _, target := range targets {
		if p.Eq(target) {
			return true
		}
	}
	return false
}

func (p Position) Movement(dx, dy int) Position {
	p.X += dx
	p.Y += dy
	return p
}
