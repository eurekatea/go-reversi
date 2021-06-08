package board

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func StrToPoint(s string) Point {
	col, row := int(s[0]-'A'), int(s[1]-'a')
	return Point{row, col}
}

func (p Point) PointToStr() string {
	return string(rune('A'+p.Y)) + string(rune('a'+p.X))
}
