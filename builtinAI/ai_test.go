package builtinai

import (
	"othello/board"
	"testing"
)

func BenchmarkEval(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++OX+++++OX+++OXOX+++OXO++++O++")
	// fmt.Println(bd.Visualize())

	for i := 0; i < b.N; i++ {
		ai.evalBoard(bd, ai.color)
	}
}

func BenchmarkEvalPointer(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++OX+++++OX+++OXOX+++OXO++++O++")

	for i := 0; i < b.N; i++ {
		ai.evalBoardPointer(bd, ai.color)
	}
}

func (ai AI) validPosNotPointer(bd board.Board, cl board.Color) (all nodes) {
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

func BenchmarkValid(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++OX+++++OX+++OXOX+++OXO++++O++")

	for i := 0; i < b.N; i++ {
		ai.validPosNotPointer(bd, ai.color)
	}
}

func (ai *AI) validPosPointer(bd board.Board, cl board.Color) (all nodes) {
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

func BenchmarkValidPointer(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++OX+++++OX+++OXOX+++OXO++++O++")

	for i := 0; i < b.N; i++ {
		ai.validPosPointer(bd, ai.color)
	}
}
