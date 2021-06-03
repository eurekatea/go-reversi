package builtinai

import (
	"fmt"
	"math"
	"math/rand"
	"othello/board"
	"sort"
)

const (
	DEPTH      = 13
	STEP2DEPTH = 18
)

var (
	DIRECTION = [8][2]int{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}

	VALUE6x6WEAKER = [][]int{
		{320, 20, 80, 80, 20, 320},
		{20, 0, 80, 80, 0, 20},
		{80, 80, 80, 80, 80, 80},
		{80, 80, 80, 80, 80, 80},
		{20, 0, 80, 80, 0, 20},
		{320, 20, 80, 80, 20, 320},
	}

	VALUE6x6 = [][]int{
		{100, -36, 53, 53, -36, 100},
		{-36, -69, -10, -10, -69, -36},
		{53, -10, -2, -2, -10, 53},
		{53, -10, -2, -2, -10, 53},
		{-36, -69, -10, -10, -69, -36},
		{100, -36, 53, 53, -36, 100},
	}

	VALUE8x8 = [][]int{
		{800, -286, 426, -24, -24, 426, -286, 800},
		{-286, -552, -177, -82, -82, -177, -552, -286},
		{426, -177, 62, 8, 8, 62, -177, 426},
		{-24, -82, 8, -18, -18, 8, -82, -24},
		{-24, -82, 8, -18, -18, 8, -82, -24},
		{426, -177, 62, 8, 8, 62, -177, 426},
		{-286, -552, -177, -82, -82, -177, -552, -286},
		{800, -286, 426, -24, -24, 426, -286, 800},
	}

	direction = [8][2]int{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}
)

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

type node struct {
	x, y  int
	value int
}

func newNode(x, y, value int) node {
	return node{x: x, y: y, value: value}
}

type nodes []node

func (ns nodes) Len() int {
	return len(ns)
}

func (ns nodes) Less(i, j int) bool {
	return ns[i].value > ns[j].value // descending order
}

func (ns nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

// provide randomness
func (ns nodes) shuffle() {
	rand.Shuffle(len(ns), func(i, j int) {
		ns[i], ns[j] = ns[j], ns[i]
	})
}

func (ns nodes) sort() {
	sort.Sort(ns)
}

type AI struct {
	color        board.Color
	opponent     board.Color
	valueNetWork [][]int
	boardSize    int

	// empty point of board
	emptyCount int

	// currently limit depth
	depth int

	// maximum reached depth
	reachedDepth int
	// traversed nodes count
	nodes int

	// the larger the stronger, level is between 0~4
	level int
}

func New(cl board.Color, boardSize int, level int) *AI {
	ai := AI{
		color:     cl,
		opponent:  cl.Opponent(),
		boardSize: boardSize,
		depth:     DEPTH,
		level:     level,
	}
	if boardSize == 6 {
		if ai.level < 3 {
			ai.valueNetWork = VALUE6x6WEAKER
		} else {
			ai.valueNetWork = VALUE6x6
		}
	} else {
		ai.valueNetWork = VALUE8x8
	}
	return &ai
}

func (ai *AI) Move(bd board.Board) (board.Point, error) {
	ai.nodes = 0
	ai.emptyCount = bd.EmptyCount()

	ai.setDepthByLevel()

	best := ai.alphaBetaHelper(bd, ai.depth)
	fmt.Printf("built-in AI: {depth: %v, nodes: %v}\n", ai.reachedDepth, ai.nodes)

	return board.NewPoint(best.x, best.y), nil
}

func (ai *AI) setDepthByLevel() {
	offset := ai.level - 4 // -4~0

	step2Max := STEP2DEPTH + (offset * 5)
	if ai.emptyCount > step2Max {
		ai.depth = DEPTH + (offset * 3) // step 1
	} else {
		ai.depth = step2Max // step 2
	}
	if ai.boardSize == 8 { // 8x8 need to reduce depth
		ai.depth -= 3
	}
}

func (ai *AI) heuristic(bd board.Board, color board.Color) int {
	if ai.emptyCount > ai.depth { // step 1
		return ai.evalBoard(bd, color)
	} else { // step 2
		return bd.CountPieces(ai.color) - bd.CountPieces(ai.opponent)
	}
}

func (ai *AI) heuristicAfterPut(bd board.Board, currentValue int, p board.Point, color board.Color) int {
	if ai.emptyCount > ai.depth { // step 1
		return ai.evalAfterPut(bd, currentValue, p, color)
	} else { // step 2
		return ai.countAfterPut(bd, currentValue, p, ai.color)
	}
}

func (ai *AI) evalBoard(bd board.Board, color board.Color) int {
	point := 0
	opponent := color.Opponent()
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			if bd.AtXY(i, j) == color {
				point += ai.valueNetWork[i][j]
			} else if bd.AtXY(i, j) == opponent {
				point -= ai.valueNetWork[i][j]
			}
		}
	}
	return point
}

func (ai *AI) changedValue(bd board.Board, cl board.Color, p board.Point, dir [2]int) int {
	delta := 0
	x, y := p.X, p.Y
	opponent := cl.Opponent()

	x, y = x+dir[0], y+dir[1]
	if bd.AtXY(x, y) != opponent {
		return 0
	}
	delta += ai.valueNetWork[x][y] * 2 // flip opponent to yours, so double

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.AtXY(x, y)
		if now != opponent {
			if now == cl {
				return delta
			} else {
				return 0
			}
		}
		delta += ai.valueNetWork[x][y] * 2 // same as above
	}
}

// don't need to copy
func (ai *AI) evalAfterPut(bd board.Board, currentValue int, p board.Point, cl board.Color) int {
	for i := 0; i < 8; i++ {
		currentValue += ai.changedValue(bd, cl, p, direction[i])
	}
	currentValue += ai.valueNetWork[p.X][p.Y]
	return currentValue
}

// don't need to copy board
func (ai *AI) countAfterPut(bd board.Board, currentCount int, p board.Point, cl board.Color) int {
	for i := 0; i < 8; i++ {
		currentCount += bd.CountFlipPieces(cl, p, direction[i])
	}
	return currentCount + 1 // include p itself
}

func (ai *AI) validPos(bd board.Board, cl board.Color) (all nodes) {
	all = make(nodes, 0, 16)
	// nowValue := ai.heuristic(bd, cl)
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			p := board.NewPoint(i, j)
			if bd.IsValidPoint(cl, p) {
				newValue := ai.valueNetWork[i][j] // better one
				// newValue := ai.heuristicAfterPut(bd, nowValue, p, cl) // old one, performed not good as this one
				all = append(all, newNode(i, j, newValue))
			}
		}
	}
	return
}

func (ai *AI) sortedValidPos(bd board.Board, cl board.Color) (all nodes) {
	all = ai.validPos(bd, cl)
	all.shuffle()
	all.sort()
	return
}

func (ai *AI) alphaBetaHelper(bd board.Board, depth int) node {
	return ai.alphaBeta(bd, depth, math.MinInt32, math.MaxInt32, true)
}

func (ai *AI) alphaBeta(bd board.Board, depth int, alpha int, beta int, maxLayer bool) node {
	if depth == 0 {
		ai.reachedDepth = ai.depth
		ai.nodes++
		return newNode(0, 0, ai.heuristic(bd, ai.color))
	}

	if maxLayer {
		bestNode := node{}
		maxValue := math.MinInt32
		all := ai.sortedValidPos(bd, ai.color)
		if len(all) == 0 {
			ai.reachedDepth = ai.depth - depth
			ai.nodes++
			return newNode(0, 0, ai.heuristic(bd, ai.color))
		}
		for _, n := range all {
			temp := bd.Copy()
			temp.PutWithoutCheck(ai.color, board.NewPoint(n.x, n.y))
			eval := ai.alphaBeta(temp, depth-1, alpha, beta, false).value

			if eval > maxValue {
				maxValue = eval
				bestNode = n
			}
			alpha = max(alpha, maxValue)
			if beta <= alpha {
				break
			}
		}
		return newNode(bestNode.x, bestNode.y, maxValue)
	} else {
		bestNode := node{}
		minValue := math.MaxInt32
		all := ai.sortedValidPos(bd, ai.opponent)
		if len(all) == 0 {
			ai.reachedDepth = ai.depth - depth
			ai.nodes++
			return newNode(0, 0, ai.heuristic(bd, ai.color))
		}
		for _, n := range all {
			temp := bd.Copy()
			temp.PutWithoutCheck(ai.opponent, board.NewPoint(n.x, n.y))
			eval := ai.alphaBeta(temp, depth-1, alpha, beta, true).value

			if eval < minValue {
				minValue = eval
				bestNode = n
			}

			beta = min(beta, minValue)
			if beta <= alpha {
				break
			}
		}

		return newNode(bestNode.x, bestNode.y, minValue)
	}
}
