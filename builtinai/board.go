package builtinai

type aiboard [][]color

var (
	DIRECTION = [8][2]int{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}
)

func newBoardFromStr(str string) aiboard {
	var size int
	if len(str) == 36 {
		size = 6
	} else {
		size = 8
	}

	indx := 0
	bd := make(aiboard, size+2)
	for i := range bd {
		bd[i] = make([]color, size+2)
	}

	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			p := point{j, i}
			switch str[indx] {
			case '+':
				bd.assign(NONE, p)
			case 'X':
				bd.assign(BLACK, p)
			case 'O':
				bd.assign(WHITE, p)
			default:
				panic("err: " + string(str[indx]))
			}
			indx++
		}
	}

	for i := 0; i < size+2; i++ {
		bd[i][0] = BORDER
		bd[0][i] = BORDER
		bd[size+1][0] = BORDER
		bd[0][size+1] = BORDER
	}

	return bd
}

func (bd aiboard) size() int {
	return len(bd) - 2
}

func (bd aiboard) Copy() aiboard {
	nbd := make(aiboard, bd.size()+2)
	for i := range bd {
		nbd[i] = make([]color, bd.size()+2)
		copy(nbd[i], bd[i])
	}
	return nbd
}

func (bd aiboard) String() (res string) {
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			switch bd.at(point{j, i}) {
			case NONE:
				res += "+"
			case BLACK:
				res += "X"
			case WHITE:
				res += "O"
			default:
				panic("err: " + bd.at(point{j, i}).String())
			}
		}
	}
	return
}

func (bd aiboard) visualize() (res string) {
	res = "  "
	for i := 0; i < bd.size(); i++ {
		res += string(rune('a'+i)) + " "
	}
	res += "\n"
	for i := 0; i < bd.size(); i++ {
		res += string(rune('A'+i)) + " "
		for j := 0; j < bd.size(); j++ {
			switch bd.at(point{j, i}) {
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

func (bd aiboard) put(cl color, p point) history {
	bd.assign(cl, p)
	return bd.flip(cl, p)
}

func (bd aiboard) putAndCheck(cl color, p point) bool {
	if p.x < 0 || p.x >= bd.size() || p.y < 0 || p.y >= bd.size() {
		return false
	}
	if bd.at(p) != NONE {
		return false
	}
	if !bd.isValidPoint(cl, p) {
		return false
	}
	bd.assign(cl, p)
	bd.flip(cl, p)
	return true
}

func (bd aiboard) assign(cl color, p point) {
	bd[p.x+1][p.y+1] = cl
}

func (bd aiboard) at(p point) color {
	return bd[p.x+1][p.y+1]
}

// undo a move
func (bd aiboard) revert(hs history) {
	for i := range hs.dirs {
		x, y := hs.place.x+hs.dirs[i][0], hs.place.y+hs.dirs[i][1]
		for j := 0; j < hs.flips[i]; j++ {
			bd.assign(hs.origColor, point{x, y})
			x, y = x+hs.dirs[i][0], y+hs.dirs[i][1]
		}
	}
	bd.assign(NONE, hs.place)
}

// self loop unrolling lol
func (bd aiboard) isValidPoint(cl color, p point) bool {
	if bd.at(p) != NONE {
		return false
	}
	op := cl.reverse()
	return bd.countFlipPieces(cl, op, p, DIRECTION[0]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[1]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[2]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[3]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[4]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[5]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[6]) > 0 ||
		bd.countFlipPieces(cl, op, p, DIRECTION[7]) > 0
}

func (bd aiboard) countFlipPieces(cl color, opponent color, p point, dir [2]int) int {
	count := 0
	x, y := p.x+dir[0], p.y+dir[1]
	if bd.at(point{x, y}) != opponent {
		return 0
	}
	count++

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.at(point{x, y})
		if now == opponent {
			count++
		} else {
			if now == cl {
				return count
			} else {
				return 0
			}
		}
	}
}

func (bd aiboard) flip(cl color, p point) history {
	hs := newHistory(p, cl.reverse())
	op := cl.reverse()
	for i := 0; i < 8; i++ {
		if count := bd.countFlipPieces(cl, op, p, DIRECTION[i]); count > 0 {
			x, y := p.x+DIRECTION[i][0], p.y+DIRECTION[i][1]
			for j := 0; j < count; j++ {
				bd.assign(cl, point{x, y})
				x, y = x+DIRECTION[i][0], y+DIRECTION[i][1]
			}
			hs.dirs = append(hs.dirs, DIRECTION[i])
			hs.flips = append(hs.flips, count)
		}
	}
	return hs
}

func (bd aiboard) emptyCount() int {
	return bd.countPieces(NONE)
}

func (bd aiboard) isOver() bool {
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			p := point{i, j}
			if bd.isValidPoint(BLACK, p) || bd.isValidPoint(WHITE, p) {
				return false
			}
		}
	}
	return true
}
