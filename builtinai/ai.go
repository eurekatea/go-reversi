package builtinai

import (
	"fmt"
	"math"
)

const (
	MININT = math.MinInt32
	MAXINT = math.MaxInt32

	PHASE1DEPTH = 16 // 6x6
	PHASE2DEPTH = 20 // 6x6
)

type AI struct {
	color    color
	opponent color

	valueDisk  [][]int
	totalValue int

	// phase 1 or phase 2
	phase int

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
		color:    cl,
		opponent: cl.reverse(),
	}

	ai.level = int(lv)

	if boardSize == 6 {
		ai.valueDisk = VALUE6x6
		ai.totalValue = TOTAL6x6
	} else {
		ai.valueDisk = VALUE8x8
		ai.totalValue = TOTAL8x8
	}
	return &ai
}

func (ai *AI) Move(input string) (string, error) {
	aibd := newBoardFromStr(input)
	boardSize := aibd.size()
	ai.nodes = 0

	ai.setPhase(aibd, boardSize)
	ai.setDepth(boardSize)

	best := ai.alphaBetaHelper(aibd, ai.depth)
	ai.printValue(best, boardSize)

	bestPoint := point{best.x, best.y}
	if !aibd.putAndCheck(ai.color, bestPoint) {
		return "", fmt.Errorf("cannot put: %v, builtin ai %v", bestPoint, ai.color)
	}
	return bestPoint.String(), nil
}

func (ai AI) Close() {}

func (ai *AI) printValue(best node, boardSize int) {
	if ai.phase == 1 {
		finValue := float64(best.value) / float64(ai.totalValue) * float64(boardSize*boardSize)
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+.2f}\n", ai.reachedDepth, ai.nodes, finValue)
	} else {
		finValue := best.value
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+d}\n", ai.reachedDepth, ai.nodes, finValue)
	}
}

func (ai *AI) setPhase(bd aiboard, boardSize int) {
	emptyCount := bd.emptyCount()
	phase2 := PHASE2DEPTH + (ai.level-4)*4 // level
	if boardSize == 8 {
		phase2 -= 2
	}
	if emptyCount > phase2 {
		ai.phase = 1
	} else {
		ai.phase = 2
	}
}

func (ai *AI) setDepth(boardSize int) {
	if ai.phase == 1 {
		if boardSize == 8 {
			ai.depth = PHASE1DEPTH + (ai.level-4)*2 - 4
		} else {
			ai.depth = PHASE1DEPTH + (ai.level-4)*4
		}
		if ai.depth <= 0 {
			ai.depth = 1
		}
	} else {
		ai.depth = MAXINT // until end of game
	}
}

func (ai *AI) heuristic(bd aiboard) int {
	if ai.phase == 1 { // phase 1
		return bd.eval(ai.color, ai.opponent, ai.valueDisk)
	} else { // phase 2
		return bd.countPieces(ai.color) - bd.countPieces(ai.opponent)
	}
}

func (ai *AI) sortedValidNodes(bd aiboard, cl color) (all nodes) {
	// usually possible point wont surpass 16
	all = make(nodes, 0, 16)
	if ai.phase == 1 { // phase 1 sort by eval
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
	} else { // phase 2 sort by mobility
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
