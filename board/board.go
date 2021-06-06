package board

type Board [][]Color

func NewBoard(size int) Board {
	realSize := size + 2
	bd := make(Board, realSize)
	for i := range bd {
		bd[i] = make([]Color, realSize)
	}
	for i := range bd {
		bd[i][0] = BORDER
		bd[0][i] = BORDER
		bd[realSize-1][i] = BORDER
		bd[i][realSize-1] = BORDER
	}

	bd[realSize/2-1][realSize/2-1] = WHITE
	bd[realSize/2][realSize/2] = WHITE
	bd[realSize/2-1][realSize/2] = BLACK
	bd[realSize/2][realSize/2-1] = BLACK

	return bd
}

func (bd Board) Size() int {
	return len(bd) - 2
}

func (bd Board) Copy() Board {
	nbd := make(Board, bd.Size()+2)
	for i := range bd {
		nbd[i] = make([]Color, bd.Size()+2)
		copy(nbd[i], bd[i])
	}
	return nbd
}

func (bd Board) CopyFromBoard(another Board) {
	for i := range bd {
		for j := range bd[i] {
			bd[i][j] = another[i][j]
		}
	}
}

func (bd Board) AssignBoard(bd2 string) {
	indx := 0
	for i := 0; i < bd.Size(); i++ {
		for j := 0; j < bd.Size(); j++ {
			switch bd2[indx] {
			case '+':
				bd.Assign(NONE, j, i)
			case 'X':
				bd.Assign(BLACK, j, i)
			case 'O':
				bd.Assign(WHITE, j, i)
			default:
				panic("err: " + string(bd2[indx]))
			}
			indx++
		}
	}
}

func (bd Board) String() (res string) {
	for i := 0; i < bd.Size(); i++ {
		for j := 0; j < bd.Size(); j++ {
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
	res = "  "
	for i := 0; i < bd.Size(); i++ {
		res += string(rune('a'+i)) + " "
	}
	res += "\n"
	for i := 0; i < bd.Size(); i++ {
		res += string(rune('A'+i)) + " "
		for j := 0; j < bd.Size(); j++ {
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
	return bd[p.X+1][p.Y+1]
}

func (bd Board) AtXY(x, y int) Color {
	return bd[x+1][y+1]
}

func (bd Board) Assign(cl Color, x, y int) {
	bd[x+1][y+1] = cl
}

func (bd Board) Put(cl Color, p Point) bool {
	if p.X < 0 || p.X >= bd.Size() || p.Y < 0 || p.Y >= bd.Size() {
		return false
	}
	if bd.AtPoint(p) != NONE {
		return false
	}
	if !bd.IsValidPoint(cl, p) {
		return false
	}
	bd.Assign(cl, p.X, p.Y)
	bd.flip(cl, p)
	return true
}

func (bd Board) PutWithoutCheck(cl Color, p Point) {
	bd.Assign(cl, p.X, p.Y)
	bd.flip(cl, p)
}

var direction = [8][2]int{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

func (bd Board) IsValidPoint(cl Color, p Point) bool {
	if bd.AtPoint(p) != NONE {
		return false
	}
	for i := 0; i < 8; i++ {
		if bd.CountFlipPieces(cl, p, direction[i]) > 0 {
			return true
		}
	}
	return false
}

func (bd Board) CountFlipPieces(cl Color, p Point, dir [2]int) int {
	count := 0
	x, y := p.X, p.Y
	opponent := cl.Opponent()

	x, y = x+dir[0], y+dir[1]
	if bd.AtXY(x, y) != opponent {
		return 0
	}
	count++

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.AtXY(x, y)
		if now != opponent {
			if now == cl {
				return count
			} else {
				return 0
			}
		}
		count++
	}
}

func (bd Board) flip(cl Color, p Point) {
	for i := 0; i < 8; i++ {
		if count := bd.CountFlipPieces(cl, p, direction[i]); count > 0 {
			for j := 1; j <= count; j++ {
				bd.Assign(cl, p.X+direction[i][0]*j, p.Y+direction[i][1]*j)
			}
		}
	}
}

func (bd Board) AllValidPoint(cl Color) []Point {
	var all []Point
	for i := 0; i < bd.Size(); i++ {
		for j := 0; j < bd.Size(); j++ {
			p := NewPoint(i, j)
			if bd.IsValidPoint(cl, p) {
				all = append(all, p)
			}
		}
	}
	return all
}

func (bd Board) CountPieces(cl Color) int {
	count := 0
	for i := 0; i < bd.Size(); i++ {
		for j := 0; j < bd.Size(); j++ {
			p := bd.AtXY(i, j)
			if p == cl {
				count++
			}
		}
	}
	return count
}

func (bd Board) EmptyCount() int {
	return bd.CountPieces(NONE)
}

func (bd Board) Winner() Color {
	bCount := bd.CountPieces(BLACK)
	wCount := bd.CountPieces(WHITE)
	if bCount > wCount {
		return BLACK
	} else if bCount < wCount {
		return WHITE
	} else {
		return NONE
	}
}

func (bd Board) IsOver() bool {
	for i := 0; i < bd.Size(); i++ {
		for j := 0; j < bd.Size(); j++ {
			p := NewPoint(i, j)
			if bd.IsValidPoint(BLACK, p) || bd.IsValidPoint(WHITE, p) {
				return false
			}
		}
	}
	return true
}
