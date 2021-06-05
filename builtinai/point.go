package builtinai

import (
	"fmt"
	"othello/board"
)

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("<%2d, %2d>", p.x, p.y)
}

func (p point) toBoardPoint() board.Point {
	return board.NewPoint(p.x, p.y)
}
