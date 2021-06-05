package builtinai

import (
	"fmt"
	"math"
	"math/rand"
	"othello/board"
	"sort"
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

var (
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

	TOTAL6x6 int
	TOTAL8x8 int
)

func init() {
	for i := 0; i < len(VALUE6x6); i++ {
		for j := 0; j < len(VALUE6x6); j++ {
			TOTAL6x6 += abs(VALUE6x6[i][j])
		}
	}
	for i := 0; i < len(VALUE8x8); i++ {
		for j := 0; j < len(VALUE8x8); j++ {
			TOTAL8x8 += abs(VALUE8x8[i][j])
		}
	}
}

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
	color    color
	opponent color

	valueNetWork [][]int
	totalValue   int

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
		ai.valueNetWork = VALUE6x6
		ai.totalValue = TOTAL6x6
		ai.depth = 4 + ai.level*2
	} else {
		ai.valueNetWork = VALUE8x8
		ai.totalValue = TOTAL8x8
		ai.depth = 2 + ai.level*2
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
	if _, ok := aibd.putAndCheck(ai.color, bestPoint); !ok {
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

	if emptyCount > STEP2DEPTH {
		ai.step = 1
	} else {
		ai.step = 2
		ai.depth += 6
	}
}

func (ai *AI) heuristic(bd aiboard) int {
	if ai.step == 1 { // step 1
		return ai.evalBoard(bd)
	} else { // step 2
		return bd.countPieces(ai.color) - bd.countPieces(ai.opponent)
	}
}

func (ai *AI) heuristicAfterPut(bd aiboard, currentValue int, p point, cl color) int {
	if ai.step == 1 { // step 1
		return ai.evalAfterPut(bd, currentValue, p, cl)
	} else { // step 2
		return ai.countAfterPut(bd, currentValue, p, ai.color)
	}
}

func (ai *AI) evalBoard(bd aiboard) int {
	value := 0
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			p := point{x: i, y: j}
			if bd.at(p) == ai.color {
				value += ai.valueNetWork[i][j]
			} else if bd.at(p) == ai.opponent {
				value -= ai.valueNetWork[i][j]
			}
		}
	}
	return value
}

func (ai *AI) changedValue(bd aiboard, cl color, p point, dir [2]int) int {
	delta := 0
	x, y := p.x, p.y
	opponent := cl.reverse()

	x, y = x+dir[0], y+dir[1]
	if bd.at(point{x: x, y: y}) != opponent {
		return 0
	}
	delta += ai.valueNetWork[x][y] * 2 // flip opponent to yours, so double

	for {
		x, y = x+dir[0], y+dir[1]
		now := bd.at(point{x: x, y: y})
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
func (ai *AI) evalAfterPut(bd aiboard, currentValue int, p point, cl color) int {
	for i := 0; i < 8; i++ {
		currentValue += ai.changedValue(bd, cl, p, DIRECTION[i])
	}
	currentValue += ai.valueNetWork[p.x][p.y]
	return currentValue
}

// don't need to copy board
func (ai *AI) countAfterPut(bd aiboard, currentCount int, p point, cl color) int {
	for i := 0; i < 8; i++ {
		currentCount += bd.countFlipPieces(cl, p, DIRECTION[i])
	}
	return currentCount + 1 // include p itself
}

func (ai *AI) validPos(bd aiboard, cl color) (all nodes) {
	all = make(nodes, 0, 16)
	// nowValue := ai.heuristic(bd, cl)
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			p := point{x: i, y: j}
			if bd.isValidPoint(cl, p) {
				newValue := ai.valueNetWork[i][j] // better one
				// 下面這個？ 之前慢的原因？ 因為heuristic只該計算出一種值？（不該分兩種顏色）
				// 這也是為什麼 min layer 之前 sort by asc 失效的原因？待研究
				// newValue := ai.heuristicAfterPut(bd, nowValue, p, cl) // old one, performed not good as this one
				all = append(all, newNode(i, j, newValue))
			}
		}
	}
	return
}

func (ai *AI) sortedValidPos(bd aiboard, cl color) (all nodes) {
	all = ai.validPos(bd, cl)
	all.shuffle()
	all.sort()
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

	aiValid := ai.validPos(bd, ai.color)
	opValid := ai.validPos(bd, ai.opponent)

	if len(aiValid) == 0 && len(opValid) == 0 {
		ai.reachedDepth = ai.depth - depth
		return newNode(-1, -1, ai.heuristic(bd))
	}

	if maxLayer {
		maxValue := MININT
		bestNode := newNode(-1, -1, maxValue)

		if len(aiValid) == 0 { // 沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, false)
		}
		aiValid.sort()

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

		if len(opValid) == 0 { // 對手沒地方下，換邊
			return ai.alphaBeta(bd, depth, alpha, beta, true)
		}
		opValid.sort()

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
