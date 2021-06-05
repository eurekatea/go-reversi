package builtinai

import (
	"fmt"
	"testing"
)

func (ai *AI) oldvalidPos(bd aiboard, cl color) (all nodes) {
	all = make(nodes, 0, 16)
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			p := point{x: i, y: j}
			if bd.isValidPoint(cl, p) {
				temp := bd.Copy()
				temp.put(cl, p)
				all = append(all, newNode(i, j, ai.heuristic(temp)))
			}
		}
	}
	return
}

func BenchmarkOrig(b *testing.B) {
	ai := New(BLACK, 6, 0)

	bd := newBoardFromStr("+++X++++X++++XOOO+++OOX+++O+++++++++")

	for i := 0; i < b.N; i++ {
		ai.oldvalidPos(bd, ai.color)
	}
}

func BenchmarkNewone(b *testing.B) {
	ai := New(BLACK, 6, 0)

	bd := newBoardFromStr("+++X++++X++++XOOO+++OOX+++O+++++++++")

	for i := 0; i < b.N; i++ {
		ai.validNodes(bd, ai.color)
	}
}

func BenchmarkHeuristic(b *testing.B) {
	ai := New(BLACK, 6, 0)

	bd := newBoardFromStr("+++X++++X++++XOOO+++OOX+++O+++++++++")

	for i := 0; i < b.N; i++ {
		ai.heuristic(bd)
	}
}

func BenchmarkPlus(b *testing.B) {
	k := 0
	for i := 0; i < b.N; i++ {
		k++
	}
	fmt.Println(k)
}

// vs revert, avg: 670 ns/op
func BenchmarkCopy(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		cpy := bd.Copy()
		cpy.put(WHITE, point{x: 4, y: 0})
	}
}

// vs copy, avg: 235 ns/op, around 3 times faster (history pass by pointer)
// after history pass by value: 170 ns/op
func BenchmarkRevert(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		hs := bd.put(WHITE, point{x: 4, y: 0})
		bd.revert(hs)
	}
}

// compiler will inline aiboard.at() so it is no problem
func BenchmarkAtP(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	p := point{x: 4, y: 0}
	for i := 0; i < b.N; i++ {
		_ = bd.at(p)
	}
}

func BenchmarkDirect(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	p := point{x: 4, y: 0}
	for i := 0; i < b.N; i++ {
		_ = bd[p.x+1][p.y+1]
	}
}

func BenchmarkAccessTwoDimSlice(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	x, y := 5, 3
	for i := 0; i < b.N; i++ {
		_ = bd[x][y]
	}
}

func BenchmarkAccessOneDimSlice(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")
	oneDim := make([]color, 36)
	cnt := 0
	for i := 0; i < bd.size(); i++ {
		for j := 0; j < bd.size(); j++ {
			oneDim[cnt] = bd[i][j]
			cnt++
		}
	}

	x := 33
	for i := 0; i < b.N; i++ {
		_ = oneDim[x]
	}
}

// if history pass by pointer: avg 255 ns/op
// but pass by value: avg 160 ns/op
func BenchmarkHs(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		hs := bd.put(WHITE, point{x: 4, y: 0})
		bd.revert(hs)
	}
}
