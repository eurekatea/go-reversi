package builtinai

import (
	"fmt"
	"math"
	"othello/board"
)

type Level int

func (l Level) String() string {
	switch l {
	case 0:
		return "beginner"
	case 1:
		return "amateur"
	case 2:
		return "professional"
	case 3:
		return "expert"
	case 4:
		return "master"
	default:
		return "unknown"
	}
}

const (
	BEGINNER Level = iota
	AMATEUR
	PROFESSIONAL
	EXPERT
	MASTER

	DEPTH      = 12
	STEP2DEPTH = 18
	MININT     = math.MinInt32
	MAXINT     = math.MaxInt32
)

func abs(v int) int {
	if v > 0 {
		return v
	} else {
		return -v
	}
}

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

type AI struct {
	color    color
	opponent color

	valueDisk  [][]int
	totalValue int

	step int

	boardSize int

	// currently limit depth
	depth int

	// maximum reached depth
	reachedDepth int
	// traversed nodes count
	nodes int

	// the larger the stronger, level is between 0~4
	level int
}

func New(cl color, boardSize int, lv Level) *AI {
	ai := AI{
		color:     cl,
		opponent:  cl.reverse(),
		boardSize: boardSize,
		depth:     DEPTH,
	}

	ai.level = int(lv)

	if boardSize == 6 {
		ai.valueDisk = VALUE6x6
		ai.totalValue = TOTAL6x6
		ai.depth = 6 + ai.level*2 // highest: 14
	} else {
		ai.valueDisk = VALUE8x8
		ai.totalValue = TOTAL8x8
		ai.depth = 2 + ai.level*2 // highest: 10
	}
	return &ai
}

func (ai *AI) Move(bd board.Board) (board.Point, error) {
	aibd := newBoardFromStr(bd.String())
	ai.nodes = 0

	ai.setStepDepth(aibd)

	best := ai.alphaBetaHelper(aibd, ai.depth)
	ai.printValue(best)

	bestPoint := point{x: best.x, y: best.y}
	if !aibd.putAndCheck(ai.color, bestPoint) {
		return bestPoint.toBoardPoint(), fmt.Errorf("cannot put: %v, i'm %v", bestPoint, ai.color)
	}
	return bestPoint.toBoardPoint(), nil
}

func (ai *AI) printValue(best node) {
	if ai.step == 1 {
		finValue := float64(best.value) / float64(ai.totalValue) * float64(ai.boardSize*ai.boardSize)
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %.2f}\n", ai.reachedDepth, ai.nodes, finValue)
	} else {
		finValue := best.value
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %d}\n", ai.reachedDepth, ai.nodes, finValue)
	}
}

func (ai *AI) setStepDepth(bd aiboard) {
	emptyCount := bd.emptyCount()

	step2Depth := STEP2DEPTH + (ai.level-4)*4 // level
	if emptyCount > step2Depth {
		ai.step = 1
	} else {
		ai.step = 2
		ai.depth = MAXINT // until end of game
	}
}

func (ai *AI) heuristic(bd aiboard) int {
	if ai.step == 1 { // step 1
		return bd.eval(ai.color, ai.opponent, ai.valueDisk)
	} else { // step 2
		return bd.countPieces(ai.color) - bd.countPieces(ai.opponent)
	}
}

func (ai *AI) sortedValidNodes(bd aiboard, cl color) (all nodes) {
	all = make(nodes, 0, 16) // usually possible point wont surpass 16
	if ai.step == 1 {
		for i := 0; i < bd.size(); i++ {
			for j := 0; j < bd.size(); j++ {
				p := point{x: i, y: j}
				if bd.isValidPoint(cl, p) {
					newValue := ai.valueDisk[i][j] // better one
					all = append(all, newNode(i, j, newValue))
				}
			}
		}
		all.sortDesc()
	} else {
		opponent := cl.reverse()
		for i := 0; i < bd.size(); i++ {
			for j := 0; j < bd.size(); j++ {
				p := point{x: i, y: j}
				if bd.isValidPoint(cl, p) {
					hs := bd.put(cl, p)
					v := bd.mobility(opponent)
					all = append(all, newNode(i, j, v))
					bd.revert(hs)
				}
			}
		}
		// the smaller the opponent mobility is, the better.
		all.sortAsc()
	}
	return
}

func (ai *AI) alphaBetaHelper(bd aiboard, depth int) node {
	return ai.alphaBeta(bd, depth, MININT, MAXINT, true)
}

func (ai *AI) alphaBeta(bd aiboard, depth int, alpha int, beta int, maxLayer bool) node {
	ai.nodes++

	if depth == 0 {
		ai.reachedDepth = ai.depth
		return newNode(-1, -1, ai.heuristic(bd))
	}
	if bd.isOver() {
		ai.reachedDepth = ai.depth - depth
		return newNode(-1, -1, ai.heuristic(bd))
	}

	if maxLayer {
		maxValue := MININT
		bestNode := newNode(-1, -1, maxValue)

		aiValid := ai.sortedValidNodes(bd, ai.color)
		if len(aiValid) == 0 { // 沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, false)
		}

		for _, n := range aiValid {
			hs := bd.put(ai.color, point{x: n.x, y: n.y})
			eval := ai.alphaBeta(bd, depth-1, alpha, beta, false).value
			bd.revert(hs)

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
		minValue := MAXINT
		bestNode := newNode(-1, -1, minValue)

		opValid := ai.sortedValidNodes(bd, ai.opponent)
		if len(opValid) == 0 { // 對手沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, true)
		}

		for _, n := range opValid {
			hs := bd.put(ai.opponent, point{x: n.x, y: n.y})
			eval := ai.alphaBeta(bd, depth-1, alpha, beta, true).value
			bd.revert(hs)

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
