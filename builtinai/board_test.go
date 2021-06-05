package builtinai

import "testing"

func testValidPoint(t *testing.T, input string, p point, target int) {

	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	count := 0
	for i := 0; i < 8; i++ {
		count += bd.countFlipPieces(WHITE, point{x: 4, y: 3}, DIRECTION[i])
	}
	if count != 1 {
		t.Error(count, "\n", bd.Visualize())
	}
}

func TestValidPoint(t *testing.T) {

	testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{x: 4, y: 3}, 1)
	testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{x: 0, y: 5}, 3)
	// testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{x: 4, y: 3}, 1)

}

func testFlip(t *testing.T, input string, cl color, p point, targetState string) {
	bd := newBoardFromStr(input)

	if _, b := bd.putAndCheck(cl, p); !b {
		t.Error("cannot put")
		t.Error("\n", bd.Visualize())
		bd.assign(cl, p)
		t.Error("\n", bd.Visualize())
		return
	}

	out := bd.String()

	for i := range out {
		if out[i] != targetState[i] {
			t.Error("failed\n", out, "\n", targetState)
			t.Error("\n", bd.Visualize())
			return
		}
	}
}

func TestFlip(t *testing.T) {

	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{x: 4, y: 0}, "++++O++++OO++OOOO+++OXOO++X+XX++++++")
	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{x: 5, y: 2}, "+++++++++XX++OOOOO++OXOO++X+XX++++++")
	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", BLACK, point{x: 5, y: 2}, "+++++++++XX++OOOXX++OXOX++X+XX++++++")

}

func testRevert(t *testing.T, input string, cl color, p point) {

	bd := newBoardFromStr(input)
	orig := bd.String()
	origBoard := bd.Copy()

	var hs *history
	var b bool
	if hs, b = bd.putAndCheck(cl, p); !b {
		t.Error("cannot put")
		t.Error("\n", bd.Visualize())
		bd.assign(cl, p)
		t.Error("\n", bd.Visualize())
		return
	}

	bd.revert(hs)
	afterRevert := bd.String()

	for i := 0; i < len(orig); i++ {
		if orig[i] != afterRevert[i] {
			t.Error("\n", origBoard.Visualize(), "\n", bd.Visualize())
		}
	}
}

func TestRevert(t *testing.T) {

	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{x: 4, y: 0})
	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{x: 5, y: 2})
	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", BLACK, point{x: 5, y: 2})

}
