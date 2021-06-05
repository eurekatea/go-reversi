package builtinai

import (
	"testing"
)

func heuristicTest(t *testing.T, input string, cl color, p point) {
	ai := New(cl, 6, 0)

	bd := newBoardFromStr(input)

	currentV := ai.heuristic(bd)
	c := bd.Copy()
	if !c.putAndCheck(cl, p) {
		t.Fatal("cannot put")
	}
	newV := ai.heuristic(c)

	aiV := ai.heuristicAfterPut(bd, currentV, p, cl)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(bd.visualize())
		t.Error(c.visualize())
	}

}

func TestPartialHeuristicChange(t *testing.T) {

	heuristicTest(t, "+++++++++++++XXX++++OXX+++O+++++++++", WHITE, point{x: 5, y: 3})
	heuristicTest(t, "++++++++++++XXOOO++XXOO+O+XXO++XXXO+", WHITE, point{x: 1, y: 4})
	heuristicTest(t, "++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+", WHITE, point{x: 1, y: 4})

}
