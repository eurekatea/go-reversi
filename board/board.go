package board

const (
	BOARD_LEN  = 6
	BOARD_REAL = BOARD_LEN + 2
)

type Board struct {
	Content [BOARD_REAL][BOARD_REAL]Color
}

func NewBoard() *Board {
	b := Board{}
	for i := 0; i < BOARD_REAL; i++ {
		b.Content[i][0] = BORDER
		b.Content[0][i] = BORDER
		b.Content[BOARD_REAL-1][i] = BORDER
		b.Content[i][BOARD_REAL-1] = BORDER
	}
	b.Content[3][3] = WHITE
	b.Content[3][4] = BLACK
	b.Content[4][3] = BLACK
	b.Content[4][4] = WHITE
	return &b
}

func (bd Board) String() (res string) {
	for i := 0; i < BOARD_LEN; i++ {
		for j := 0; j < BOARD_LEN; j++ {
			switch bd.AtXY(j, i) {
			case NONE:
				res += "+"
			case BLACK:
				res += "X"
			case WHITE:
				res += "O"
			default:
				panic("err: " + bd.AtXY(j, i).String())
			}
		}
	}
	return
}

func (bd Board) Visualize() (res string) {
	for i := 0; i < BOARD_LEN; i++ {
		for j := 0; j < BOARD_LEN; j++ {
			switch bd.AtXY(j, i) {
			case NONE:
				res += "+ "
			case BLACK:
				res += "X "
			case WHITE:
				res += "O "
			}
		}
		res += "\n"
	}
	return
}

func (bd Board) AtPoint(p Point) Color {
	return bd.Content[p.X+1][p.Y+1]
}

func (bd Board) AtXY(x, y int) Color {
	return bd.Content[x+1][y+1]
}

func (bd *Board) PutPoint(cl Color, p Point) bool {
	if p.X < 0 || p.X >= BOARD_LEN || p.Y < 0 || p.Y >= BOARD_LEN {
		return false
	}
	if bd.Content[p.X+1][p.Y+1] != NONE {
		return false
	}
	if !bd.isValidPoint(cl, p) {
		return false
	}
	bd.Content[p.X+1][p.Y+1] = cl
	bd.flip(cl, p)
	return true
}

var direction = [8][2]int{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

func (bd Board) isValidPoint(cl Color, p Point) bool {
	if bd.AtPoint(p) != NONE {
		return false
	}
	for i := 0; i < 8; i++ {
		if bd.isValidDir(cl, p, direction[i]) {
			return true
		}
	}
	return false
}

func (bd Board) isValidDir(cl Color, p Point, dir [2]int) bool {
	x, y := p.X, p.Y
	opponent := cl.Opponent()

	x, y = x+dir[0], y+dir[1]
	if bd.AtXY(x, y) != opponent {
		return false
	}

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.AtXY(x, y)
		if now != opponent {
			break
		}
	}

	return bd.AtXY(x, y) == cl
}

func (bd Board) AllValidPoint(cl Color) []Point {
	var all []Point
	for i := 0; i < BOARD_LEN; i++ {
		for j := 0; j < BOARD_LEN; j++ {
			p := Point{X: i, Y: j}
			if bd.isValidPoint(cl, p) {
				all = append(all, p)
			}
		}
	}
	return all
}

func (bd *Board) flipPoint(p Point) {
	cl := bd.AtPoint(p)
	bd.Content[p.X+1][p.Y+1] = cl.Opponent()
}

func (bd *Board) flip(cl Color, p Point) {
	opponent := cl.Opponent()
	for i := 0; i < 8; i++ {
		x, y := p.X, p.Y
		if !bd.isValidDir(cl, p, direction[i]) {
			continue
		}
		for {
			x, y = x+direction[i][0], y+direction[i][1]
			p := Point{X: x, Y: y}
			if bd.AtPoint(p) == opponent {
				bd.flipPoint(p)
			} else {
				break
			}
		}
	}
}

func (bd Board) Winner() Color {
	var bCount, wCount int
	for i := 0; i < BOARD_LEN; i++ {
		for j := 0; j < BOARD_LEN; j++ {
			p := bd.AtXY(i, j)
			if p == BLACK {
				bCount++
			} else if p == WHITE {
				wCount++
			}
		}
	}
	if bCount > wCount {
		return BLACK
	} else if bCount < wCount {
		return WHITE
	} else {
		return NONE
	}
}
