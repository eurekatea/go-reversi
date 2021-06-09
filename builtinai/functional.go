package builtinai

import "math"

const (
	MININT = math.MinInt32
	MAXINT = math.MaxInt32
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
