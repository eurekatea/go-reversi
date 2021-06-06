package builtinai

import (
	"fmt"
	"math"
	"othello/board"
)

const (
	STEP2DEPTH = 18
	MININT     = math.MinInt32
	MAXINT     = math.MaxInt32
)

type AI struct {
	color    color
	opponent color

	valueDisk  [][]int
	totalValue int

	// step 1 or step 2
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
	}

	ai.level = int(lv)

	if boardSize == 6 {
		ai.valueDisk = VALUE6x6
		ai.totalValue = TOTAL6x6
		ai.depth = 8 + ai.level*2 // highest: 16
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

	bestPoint := point{best.x, best.y}
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
	if ai.boardSize == 6 {
		step2Depth += 2
	}
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
	// usually possible point wont surpass 16
	all = make(nodes, 0, 16)
	if ai.step == 1 { // step 1 sort by eval
		origValue := bd.eval(cl, cl.reverse(), ai.valueDisk)
		for i := 0; i < bd.size(); i++ {
			for j := 0; j < bd.size(); j++ {
				p := point{i, j}
				if bd.isValidPoint(cl, p) {
					newValue := bd.evalAfterPut(origValue, p, cl, ai.valueDisk)
					all = append(all, node{i, j, newValue})
				}
			}
		}
		all.sortDesc()
	} else { // step 2 sort by mobility
		opponent := cl.reverse()
		for i := 0; i < bd.size(); i++ {
			for j := 0; j < bd.size(); j++ {
				p := point{i, j}
				if bd.isValidPoint(cl, p) {
					hs := bd.put(cl, p)
					v := bd.mobility(opponent)
					all = append(all, node{i, j, v})
					bd.revert(hs)
				}
			}
		}
		// the smaller the opponent's mobility is, the better.
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
		return node{-1, -1, ai.heuristic(bd)}
	}
	if bd.isOver() {
		ai.reachedDepth = ai.depth - depth
		return node{-1, -1, ai.heuristic(bd)}
	}

	if maxLayer {
		maxValue := MININT
		bestNode := node{-1, -1, maxValue}

		aiValid := ai.sortedValidNodes(bd, ai.color)
		if len(aiValid) == 0 { // 沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, false)
		}

		for _, n := range aiValid {
			hs := bd.put(ai.color, point{n.x, n.y})
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

		return node{bestNode.x, bestNode.y, maxValue}
	} else {
		minValue := MAXINT
		bestNode := node{-1, -1, minValue}

		opValid := ai.sortedValidNodes(bd, ai.opponent)
		if len(opValid) == 0 { // 對手沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, true)
		}

		for _, n := range opValid {
			hs := bd.put(ai.opponent, point{n.x, n.y})
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

		return node{bestNode.x, bestNode.y, minValue}
	}
}
