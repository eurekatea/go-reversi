package builtinai

import (
	"othello/board"
	"testing"
)

func (ai *AI) oldvalidPos(bd board.Board, cl board.Color) (all nodes) {
	all = make(nodes, 0, 16)
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			p := board.NewPoint(i, j)
			if bd.IsValidPoint(cl, p) {
				temp := bd.Copy()
				temp.Put(cl, p)
				all = append(all, newNode(i, j, ai.heuristic(temp, cl)))
			}
		}
	}
	return
}

func BenchmarkOrig(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++X++++X++++XOOO+++OOX+++O+++++++++")

	for i := 0; i < b.N; i++ {
		ai.oldvalidPos(bd, ai.color)
	}
}

func BenchmarkNewone(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++X++++X++++XOOO+++OOX+++O+++++++++")

	for i := 0; i < b.N; i++ {
		ai.validPos(bd, ai.color)
	}
}
