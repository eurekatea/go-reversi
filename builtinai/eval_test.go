package builtinai

import (
	"testing"
)

func testPartialValueChange(t *testing.T, input string, p point, cl color) {
	ai := New(cl, 6, 0)

	bd := newBoardFromStr(input)

	currentV := bd.eval(ai.color, ai.opponent, ai.valueDisk)
	c := bd.Copy()
	if !c.putAndCheck(cl, p) {
		t.Error(c.visualize())
		c[p.x][p.y] = cl
		t.Error(c.visualize())
		t.Fatal("cannot put")
	}
	newV := c.eval(ai.color, ai.opponent, ai.valueDisk)

	aiV := bd.evalAfterPut(currentV, p, cl, ai.valueDisk)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(c.visualize())
	}
}

func TestPartialValueChange(t *testing.T) {

	testPartialValueChange(t, "+++++++++++++XXX++++OXX+++O+++++++++", point{5, 3}, WHITE)
	testPartialValueChange(t, "++++++++++++XXOOO++XXOO+O+XXO++XXXO+", point{1, 4}, WHITE)
	testPartialValueChange(t, "++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+", point{1, 4}, WHITE)

}
