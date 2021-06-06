package builtinai

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
)
