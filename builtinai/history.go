package builtinai

type history struct {
	origColor color
	place     point
	dirs      [][2]int
	flips     []int
}

func newHistory(place point, origColor color) *history {
	return &history{
		origColor: origColor,
		place:     place,
		dirs:      make([][2]int, 0, 4),
		flips:     make([]int, 0, 4),
	}
}
