package builtinai

import (
	"fmt"
)

const (
	PHASE1DEPTH6 = 16 // 6x6
	PHASE2DEPTH6 = 20 // 6x6
)

type AI6 struct {
	color    color
	opponent color

	totalValue int

	// table map[bboard6]int

	// board size
	size int

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

func NewAI6(cl color, lv Level) *AI6 {
	ai := AI6{
		color: cl,
		// table:    make(map[bboard6]int),
		opponent: cl.reverse(),
	}

	ai.level = int(lv)
	ai.totalValue = 1476
	ai.size = 6

	return &ai
}

func (ai *AI6) Move(input string) (string, error) {
	aibd := newBboard6(input)
	ai.nodes = 0

	ai.setPhase(aibd)
	ai.setDepth()

	var best node
	// for depth := 2; depth <= ai.depth; depth += 2 {
	best = ai.alphaBetaHelper(aibd, ai.depth)
	// }
	ai.printValue(best)

	bestPoint := point{best.loc % ai.size, best.loc / ai.size}
	if !aibd.putAndCheck(ai.color, best.loc) {
		return "", fmt.Errorf("cannot put: %v, builtin ai %v", bestPoint, ai.color)
	}
	return bestPoint.String(), nil
}

func (ai AI6) Close() {}

func (ai *AI6) printValue(best node) {
	if ai.phase == 1 {
		finValue := float64(best.value) / float64(ai.totalValue) * float64(ai.size*ai.size)
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+.2f}\n", ai.reachedDepth, ai.nodes, finValue)
	} else {
		finValue := best.value
		fmt.Printf("built-in AI: {depth: %d, nodes: %d, value: %+d}\n", ai.reachedDepth, ai.nodes, finValue)
	}
}

func (ai *AI6) setPhase(bd bboard6) {
	emptyCount := bd.emptyCount()
	phase2 := PHASE2DEPTH6 + (ai.level-4)*4 // level
	if emptyCount > phase2 {
		ai.phase = 1
	} else {
		ai.phase = 2
	}
}

func (ai *AI6) setDepth() {
	if ai.phase == 1 {
		ai.depth = PHASE1DEPTH6 + (ai.level-4)*4
		if ai.depth <= 0 {
			ai.depth = 1
		}
	} else {
		ai.depth = MAXINT // until end of game
	}
}

func (ai *AI6) heuristic(bd bboard6) int {
	if ai.phase == 1 { // phase 1
		return bd.eval(ai.color)
	} else { // phase 2
		return bd.count(ai.color) - bd.count(ai.opponent)
	}
}

func (ai *AI6) sortedValidNodes(bd bboard6, cl color) (all nodes) {
	// usually possible point wont surpass 16
	all = make(nodes, 0, 16)
	if ai.phase == 1 { // phase 1 sort by eval
		allValid := bd.allValidLoc(cl)
		for loc := 0; loc < ai.size*ai.size; loc++ {
			if (u1<<loc)&allValid != 0 {
				tmp := bd.cpy()
				tmp.put(cl, loc)
				all = append(all, node{loc, tmp.eval(cl)})
			}
		}
		all.sortDesc()
	} else { // phase 2 sort by mobility
		op := cl.reverse()
		allValid := bd.allValidLoc(cl)
		for loc := 0; loc < ai.size*ai.size; loc++ {
			if (u1<<loc)&allValid != 0 {
				tmp := bd.cpy()
				tmp.put(cl, loc)
				v := tmp.mobility(op)
				all = append(all, node{loc, v})
			}
		}
		// the smaller the opponent's mobility is, the better.
		all.sortAsc()
	}
	return
}

// func (ai *AI6) searchTable(bd bboard6) int {

// 	return v
// }

func (ai *AI6) alphaBetaHelper(bd bboard6, depth int) node {
	return ai.alphaBeta(bd, depth, MININT, MAXINT, true)
}

func (ai *AI6) alphaBeta(bd bboard6, depth int, alpha int, beta int, maxLayer bool) node {
	ai.nodes++

	if depth == 0 {
		ai.reachedDepth = ai.depth
		// if v, exi := ai.table[bd]; exi {
		// 	return node{-1, v}
		// }
		v := ai.heuristic(bd)
		// ai.table[bd] = v
		return node{-1, v}
	}
	if bd.isOver() {
		ai.reachedDepth = ai.depth - depth
		// if v, exi := ai.table[bd]; exi {
		// 	return node{-1, v}
		// }
		v := ai.heuristic(bd)
		// ai.table[bd] = v
		return node{-1, v}
	}

	if maxLayer {
		maxValue := MININT
		bestNode := node{-1, maxValue}

		aiValid := ai.sortedValidNodes(bd, ai.color)
		if len(aiValid) == 0 { // 沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, false)
		}

		for _, n := range aiValid {
			tmp := bd.cpy()
			tmp.put(ai.color, n.loc)
			eval := ai.alphaBeta(tmp, depth-1, alpha, beta, false).value

			if eval > maxValue {
				maxValue = eval
				bestNode = n
			}
			alpha = max(alpha, maxValue)
			if beta <= alpha {
				break
			}
		}

		return node{bestNode.loc, maxValue}
	} else {
		minValue := MAXINT
		bestNode := node{-1, minValue}

		opValid := ai.sortedValidNodes(bd, ai.opponent)
		if len(opValid) == 0 { // 對手沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, true)
		}

		for _, n := range opValid {
			tmp := bd.cpy()
			tmp.put(ai.opponent, n.loc)
			eval := ai.alphaBeta(tmp, depth-1, alpha, beta, true).value

			if eval < minValue {
				minValue = eval
				bestNode = n
			}

			beta = min(beta, minValue)
			if beta <= alpha {
				break
			}
		}

		return node{bestNode.loc, minValue}
	}
}
