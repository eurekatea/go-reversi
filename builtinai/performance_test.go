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
		ai.validPos(bd, ai.color)
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

// vs copy, avg: 235 ns/op, around 3 times faster
func BenchmarkRevert(b *testing.B) {
	bd := newBoardFromStr("+++++++++XX++OOOX+++OXOO++X+XX++++++")

	for i := 0; i < b.N; i++ {
		hs := bd.put(WHITE, point{x: 4, y: 0})
		bd.revert(hs)
	}
}
