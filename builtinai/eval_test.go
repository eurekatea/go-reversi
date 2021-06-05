package builtinai

import (
	"testing"
)

func unit(t *testing.T, input string, p point, cl color) {
	ai := New(cl, 6, 0)

	bd := newBoardFromStr(input)

	currentV := ai.evalBoard(bd)
	c := bd.Copy()
	if !c.putAndCheck(cl, p) {
		t.Error(c.visualize())
		c[p.x][p.y] = cl
		t.Error(c.visualize())
		t.Fatal("cannot put")
	}
	newV := ai.evalBoard(c)

	aiV := ai.evalAfterPut(bd, currentV, p, cl)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(c.visualize())
	}
}

func TestPartialValueChange(t *testing.T) {

	unit(t, "+++++++++++++XXX++++OXX+++O+++++++++", point{x: 5, y: 3}, WHITE)
	unit(t, "++++++++++++XXOOO++XXOO+O+XXO++XXXO+", point{x: 1, y: 4}, WHITE)
	unit(t, "++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+", point{x: 1, y: 4}, WHITE)

}
