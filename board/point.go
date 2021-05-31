package board

import "fmt"

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func CenterPoint() Point {
	return Point{X: BOARD_LEN / 2, Y: BOARD_LEN / 2}
}

func (p Point) String() string {
	return fmt.Sprintf("<%2v,%2v>", p.X+1, p.Y+1)
}
