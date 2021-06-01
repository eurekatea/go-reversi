package board

type Color int8

const (
	NONE   Color = 0
	BLACK  Color = 1
	WHITE  Color = -1
	BORDER Color = 127
)

func (cl Color) Opponent() Color {
	return -1 * cl
}

func (cl Color) String() string {
	if cl == BLACK {
		return "BLACK"
	} else if cl == WHITE {
		return "WHITE"
	} else {
		return "NONE"
	}
}
