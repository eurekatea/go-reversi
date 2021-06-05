package builtinai

import (
	"testing"
)

func particleChange(t *testing.T, input string, cl color, p point) {
	ai := New(cl, 6, 0)

	bd := newBoardFromStr(input)

	currentV := bd.countPieces(cl)
	c := bd.Copy()
	if _, b := c.putAndCheck(cl, p); !b {
		t.Error(c.Visualize())
		t.Fatal("cannot put")
	}
	newV := c.countPieces(cl)

	aiV := ai.countAfterPut(bd, currentV, p, cl)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(c.Visualize())
	}
}

func TestPartialCountChange(t *testing.T) {

	particleChange(t, "+++X++++X++++XOOO+++OOX+++O+++++++++", BLACK, point{x: 5, y: 2})
	particleChange(t, "++++++++++++XXOOO++XXOO+O+XXO++XXXO+", WHITE, point{x: 1, y: 4})
	particleChange(t, "++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+", WHITE, point{x: 1, y: 4})

}
