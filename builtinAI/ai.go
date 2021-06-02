package builtinai

import (
	"fmt"
	"math"
	"othello/board"
	"runtime/debug"
	"sort"
)

const (
	DEPTH  = 12
	DEPTH2 = 17
)

var (
	DIRECTION = [8][2]int{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}

	VALUE6x6 = [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 320, 20, 80, 80, 20, 320, 0},
		{0, 20, 0, 80, 80, 0, 20, 0},
		{0, 80, 80, 80, 80, 80, 80, 0},
		{0, 80, 80, 80, 80, 80, 80, 0},
		{0, 20, 0, 80, 80, 0, 20, 0},
		{0, 320, 20, 80, 80, 20, 320, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}

	VALUE8x8 = [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 101, -43, 38, 7, 0, 42, -20, 102, 0},
		{0, -27, -74, -16, -14, -13, -25, -65, -39, 0},
		{0, 56, -30, 12, 5, -4, 7, -15, 48, 0},
		{0, 1, -8, 1, -1, -4, -2, -12, 3, 0},
		{0, -10, -8, 1, -1, -3, 2, -4, -20, 0},
		{0, 59, -23, 6, 1, 4, 6, -19, 35, 0},
		{0, -6, -55, -18, -8, -15, -31, -82, -58, 0},
		{0, 96, -42, 67, -2, -3, 81, -51, 101, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
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
	return ns[i].value > ns[j].value
}

func (ns nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

type AI struct {
	color        board.Color
	opponent     board.Color
	valueNetWork [][]int
	boardSize    int
	emptyCount   int
	depth        int
	nodes        int
}

func New(cl board.Color, boardSize int) *AI {
	ai := AI{
		color:     cl,
		opponent:  cl.Opponent(),
		boardSize: boardSize,
		depth:     DEPTH,
	}
	if boardSize == 6 {
		ai.valueNetWork = VALUE6x6
	} else {
		ai.valueNetWork = VALUE8x8
	}
	return &ai
}

func (ai *AI) Move(bd board.Board) (board.Point, error) {
	debug.SetGCPercent(10000)
	defer debug.SetGCPercent(100)

	ai.nodes = 0
	ai.emptyCount = bd.EmptyCount()
	if ai.emptyCount > 16 {
		ai.depth = DEPTH
	} else {
		ai.depth = DEPTH2
	}
	if ai.boardSize == 8 {
		ai.depth -= 2
	}

	best := ai.alphaBetaHelper(bd, ai.depth)
	fmt.Printf("built-in AI: {nodes: %v}\n", ai.nodes)

	return board.NewPoint(best.x, best.y), nil
}

func (ai AI) heuristic(bd board.Board, color board.Color) int {
	if ai.emptyCount > 16 {
		return ai.evalBoard(bd, color)
	} else {
		return bd.CountPieces(ai.color) - bd.CountPieces(ai.opponent)
	}
}

func (ai AI) evalBoard(bd board.Board, color board.Color) int {
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

func (ai AI) validPos(bd board.Board, cl board.Color) (all nodes) {
	all = make(nodes, 0, 16)
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			p := board.NewPoint(i, j)
			if bd.IsValidPoint(cl, p) {
				temp := bd.Copy()
				temp.Put(cl, p)
				all = append(all, newNode(i, j, ai.heuristic(temp, cl)))
			}
		}
	}
	return
}

func (ai AI) sortedValidPos(bd board.Board, cl board.Color) (all nodes) {
	all = ai.validPos(bd, cl)
	sort.Sort(all)
	return
}

func (ai *AI) alphaBetaHelper(bd board.Board, depth int) node {
	return ai.alphaBeta(bd, depth, math.MinInt32, math.MaxInt32, true)
}

func (ai *AI) alphaBeta(bd board.Board, depth int, alpha int, beta int, maxLayer bool) node {
	ai.nodes++
	if depth == 0 {
		return newNode(0, 0, ai.heuristic(bd, ai.color))
	}

	if maxLayer {
		bestNode := node{}
		maxValue := math.MinInt32
		all := ai.sortedValidPos(bd, ai.color)
		if len(all) == 0 {
			return newNode(0, 0, ai.heuristic(bd, ai.color))
		}
		for _, n := range all {
			temp := bd.Copy()
			temp.Put(ai.color, board.NewPoint(n.x, n.y))
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
			return newNode(0, 0, ai.heuristic(bd, ai.color))
		}
		for _, n := range all {
			temp := bd.Copy()
			temp.Put(ai.opponent, board.NewPoint(n.x, n.y))
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
