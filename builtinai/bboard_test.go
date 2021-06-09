package builtinai

import (
	"fmt"
	"math/rand"
	"othello/board"
	"testing"
	"time"
)

func TestBboard(t *testing.T) {

	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%3 == 0 {
				input += "X"
			} else if r%3 == 1 {
				input += "O"
			} else {
				input += "+"
			}
		}
		bd := newBboard8(input)

		if bd.String() != input {
			t.Error(bd.String())
			t.Error(input)
		}
	}
}

func TestCount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		bCnt := 0
		wCnt := 0
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%3 == 0 {
				input += "X"
				bCnt++
			} else if r%3 == 1 {
				input += "O"
				wCnt++
			} else {
				input += "+"
			}
		}
		bd := newBboard8(input)

		if bCnt != bd.count(BLACK) {
			t.Error(bCnt, bd.count(BLACK), "black")
		}
		if wCnt != bd.count(WHITE) {
			t.Error(wCnt, bd.count(WHITE), "white")
		}
	}
}

func TestEmptyCount(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		cnt := 0
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%3 == 0 {
				input += "X"
			} else if r%3 == 1 {
				input += "O"
			} else {
				input += "+"
				cnt++
			}
		}
		bd := newBboard8(input)

		if cnt != bd.emptyCount() {
			t.Error(cnt, bd.emptyCount())
		}
	}
}

func testPut(t *testing.T, input string, cl color, x, y int) {
	b := board.NewBoard(8)
	b.AssignBoard(input)

	var c board.Color
	if cl == BLACK {
		c = board.BLACK
	} else {
		c = board.WHITE
	}

	if !b.PutPoint(c, board.NewPoint(x, y)) {
		t.Error("cannot put", input, x, y)
		t.Error(b.Visualize())
	}

	bbd := newBboard8(input)
	if !bbd.putAndCheck(cl, y*8+x) {
		t.Error("cannot put", input)
	}

	if b.String() != bbd.String() {
		t.Error(b.String())
		t.Error(bbd.String())
		t.Error(input, x, y, cl)
		t.Error(b.Visualize())
	}
}

func TestPut(t *testing.T) {
	testPut(t, "+++++++++++++++++++++++++++OX++++++XO+++++++++++++++++++++++++++", BLACK, 2, 3)
	testPut(t, "+++++++++++++++++++X++++++XOOO+++X+OOX+++++O++++++++++++++++++++", BLACK, 6, 3)
	testPut(t, "+++++++++++X++++++++X+++++XXXXO+++XXOO++++X+O+++++++++++++++++++", WHITE, 2, 2)
	testPut(t, "+++++++++++X++++++++X+++++XXXXO+++XXOO++++X+O+++++++++++++++++++", WHITE, 3, 2)
	testPut(t, "+++++++++++X++++++++X++++OOOOOO+++XXOO++++X+O+++++++++++++++++++", BLACK, 4, 6)
	testPut(t, "+++++++++++X++++++++X++++OOOOXO+++XXOXX+++X+OO++++++++++++++++++", BLACK, 2, 2)
	testPut(t, "+++++++++++X+++++++XX++++OOXXXO+++XXOXX+++X+OO++++++++++++++++++", WHITE, 1, 4)
	testPut(t, "+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++", BLACK, 2, 2)
	testPut(t, "+++O+++++++OO+++++XOXX+++OXOOXO+++XOXOOO+OXOOXXX+++OOXX++++XO+++", BLACK, 7, 3)
	testPut(t, "+++O+++++++OO+++++XOXX+++OXOOXXX++XOXOXX+OXOOXXX+++OOXX++++XO+++", WHITE, 7, 7)
	testPut(t, "+++O+++++++OO+++X+XOXX+++XXOOXXX++XOOOXX+OXOOOXX+++OOXO++++XO++O", WHITE, 1, 2)
	testPut(t, "+++O+++++O+OO+++X+OOXX+++XOOOXXX++OOOOXX+OOOOOXX++OOXXX++++XXX+O", BLACK, 1, 4)
	testPut(t, "+++O+++++O+OO+++X+OOXX+++XOOOXXX+XXXXXXX+OXOOOXX++OXXXX++++XXX+O", WHITE, 6, 2)
	testPut(t, "+++O+++++O+OO+++X+OOOOO++XOOOOXX+XXXOXXX+OXOOOXX++OXXXX++++XXX+O", BLACK, 5, 0)
	testPut(t, "+++OX+O++O+XXOXOX+XXOXXX+XOOXOXX+XXXXXXX+OXOXOXX++XXXXX+++XXXX+O", BLACK, 0, 5)
}

func testEval(t *testing.T, input string, cl color) {
	b := newBoardFromStr(input)
	bbd := newBboard8(input)

	if b.eval(cl, cl.reverse(), VALUE8x8) != bbd.eval(cl) {
		t.Error()
	}
}

func TestEval(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		cnt := 0
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%3 == 0 {
				input += "X"
			} else if r%3 == 1 {
				input += "O"
			} else {
				input += "+"
				cnt++
			}
		}

		if rand.Int()%2 == 0 {
			testEval(t, input, BLACK)
		} else {
			testEval(t, input, WHITE)
		}
	}
}

func testValidLoc(t *testing.T, input string) bool {
	bd := newBoardFromStr(input)
	bbd := newBboard8(input)
	fail := false

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b1 := bd.isValidPoint(BLACK, point{i, j})
			b2 := bbd.isValidLoc(BLACK, j*8+i)
			if b1 != b2 {
				bd.assign(BORDER, point{i, j})
				t.Error("Black", i, j, b1, b2)
				fail = true
			}
			b1 = bd.isValidPoint(WHITE, point{i, j})
			b2 = bbd.isValidLoc(WHITE, j*8+i)
			if b1 != b2 {
				bd.assign(BORDER, point{i, j})
				t.Error("White", i, j, b1, b2)
				fail = true
			}
		}
	}
	if fail {
		fmt.Println(bd.visualize())
		return true
	}
	return false
}

func TestIsValidLoc(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%3 == 0 {
				input += "X"
			} else if r%3 == 1 {
				input += "O"
			} else {
				input += "+"
			}
		}

		if testValidLoc(t, input) {
			break
		}
	}
}

func TestOver(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		input := ""
		for j := 0; j < 64; j++ {
			r := rand.Uint64()
			if r%11 < 5 {
				input += "X"
			} else if r%5 < 9 {
				input += "O"
			} else {
				input += "+"
			}
		}

		bd := newBoardFromStr(input)
		bbd := newBboard8(input)

		if bd.isOver() != bbd.isOver() {
			t.Error()
		}
	}
}

// avg: 91.04 ns/op
func BenchmarkEvalNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.eval(BLACK, WHITE, VALUE8x8)
	}
}

// avg: 56.4 ns/op
func BenchmarkEvalB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.eval(BLACK)
	}
}

// avg: 82 ns/op
func BenchmarkCountNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		_ = bd.countPieces(BLACK) - bd.countPieces(WHITE)
	}
}

// avg: 8.4 ns/op
func BenchmarkCountB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		_ = bd.count(BLACK) - bd.count(WHITE)
	}
}

// avg: 564 ns/op 496 B/op 13 allocs/op
func BenchmarkCpy(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.put(WHITE, point{4, 0})
		_ = bd.Copy()
	}
}

// avg: 170 ns/op 96 B/op 2 allocs/op
func BenchmarkRevertbd(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		hs := bd.put(WHITE, point{4, 0})
		bd.revert(hs)
	}
}

// avg: 40 ns/op 0 B/op 0 allocs/op
func BenchmarkCpyb(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	for i := 0; i < b.N; i++ {
		bd.put(WHITE, 32)
		_ = bd.cpy()
	}
}

func BenchmarkAccessNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	p := point{4, 3}
	for i := 0; i < b.N; i++ {
		_ = bd.at(p)
	}
}

// almost the same ↑↓ (0.5 ns/op)

func BenchmarkAccessB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	loc := 28
	for i := 0; i < b.N; i++ {
		_ = bd.at(loc)
	}
}

// avg: 0.677 ns/op
func BenchmarkAssignNorm(b *testing.B) {
	bd := newBoardFromStr("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	p := point{4, 3}
	for i := 0; i < b.N; i++ {
		bd.assign(WHITE, p)
	}
}

// avg: 1.372 ns/op
func BenchmarkAssignB(b *testing.B) {
	bd := newBboard8("+++++++++++XO++++++OOX+++OOOOXO+++OOOOOO+OXOOXXX+++OOXX++++XO+++")
	loc := 28
	for i := 0; i < b.N; i++ {
		bd.assign(WHITE, loc)
	}
}
