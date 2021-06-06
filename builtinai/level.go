package builtinai

import "fmt"

type Level int

func (l Level) String() string {
	return fmt.Sprintf("%d", l+1)
}

const (
	LV_ONE Level = iota
	LV_TWO
	LV_THREE
	LV_FOUR
	LV_FIVE
)
