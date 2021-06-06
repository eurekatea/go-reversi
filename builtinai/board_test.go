package builtinai

import "testing"

func testValidPoint(t *testing.T, input string, p point, target int) {

	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	count := 0
	for i := 0; i < 8; i++ {
		count += bd.countFlipPieces(WHITE, BLACK, point{4, 3}, DIRECTION[i])
	}
	if count != 1 {
		t.Error(count, "\n", bd.visualize())
	}
}

func TestValidPoint(t *testing.T) {

	testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{4, 3}, 1)
	testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{0, 5}, 3)
	// testValidPoint(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", point{ 4,  3}, 1)

}

func testFlip(t *testing.T, input string, cl color, p point, targetState string) {
	bd := newBoardFromStr(input)

	if !bd.putAndCheck(cl, p) {
		t.Error("cannot put")
		t.Error("\n", bd.visualize())
		bd.assign(cl, p)
		t.Error("\n", bd.visualize())
		return
	}

	out := bd.String()

	for i := range out {
		if out[i] != targetState[i] {
			t.Error("failed\n", out, "\n", targetState)
			t.Error("\n", bd.visualize())
			return
		}
	}
}

func TestFlip(t *testing.T) {

	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{4, 0}, "++++O++++OO++OOOO+++OXOO++X+XX++++++")
	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{5, 2}, "+++++++++XX++OOOOO++OXOO++X+XX++++++")
	testFlip(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", BLACK, point{5, 2}, "+++++++++XX++OOOXX++OXOX++X+XX++++++")

}

func testRevert(t *testing.T, input string, cl color, p point) {

	bd := newBoardFromStr(input)
	orig := bd.String()
	origBoard := bd.Copy()

	hs := bd.put(cl, p)

	bd.revert(hs)
	afterRevert := bd.String()

	for i := 0; i < len(orig); i++ {
		if orig[i] != afterRevert[i] {
			t.Error("\n", origBoard.visualize(), "\n", bd.visualize())
		}
	}
}

func TestRevert(t *testing.T) {

	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{4, 0})
	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", WHITE, point{5, 2})
	testRevert(t, "+++++++++XX++OOOX+++OXOO++X+XX++++++", BLACK, point{5, 2})

}
