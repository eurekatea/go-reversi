package builtinai

import (
	"othello/board"
	"testing"
)

func TestPartialValueChange(t *testing.T) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++X++++X++++XOOO+++OOX+++O+++++++++")

	currentV := ai.evalBoard(bd, board.BLACK)
	c := bd.Copy()
	c.Put(board.BLACK, board.NewPoint(3, 1))
	newV := ai.evalBoard(c, board.BLACK)

	aiV := ai.evalAfterPut(bd, currentV, board.NewPoint(3, 1), board.BLACK)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
	}

	ai = New(board.WHITE, 6, 0)
	bd = board.NewBoard(6)
	bd.AssignBoard("++++++++++++XXOOO++XXOO+O+XXO++XXXO+")
	t.Error(bd.Visualize())
	currentV = ai.evalBoard(bd, board.WHITE)
	c = bd.Copy()
	c.Put(board.WHITE, board.NewPoint(4, 2))
	newV = ai.evalBoard(c, board.WHITE)

	aiV = ai.evalAfterPut(bd, currentV, board.NewPoint(4, 2), board.WHITE)

	if newV != aiV {
		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
		t.Error(c.Visualize())
	}
}

func (ai AI) evalBoardNoPointer(bd board.Board, color board.Color) int {
	point := 0
	opponent := color.Opponent()
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			if bd.AtXY(i, j) == color {
				point += ai.valueNetWork[i][j]
			} else if bd.AtXY(i, j) == opponent {
				point -= ai.valueNetWork[i][j]
			}
		}
	}
	return point
}

func BenchmarkEval(b *testing.B) {
	ai := New(board.BLACK, 6, 0)

	bd := board.NewBoard(6)
	bd.AssignBoard("+++++++OX+++++OX+++OXOX+++OXO++++O++")
	// fmt.Println(bd.Visualize())

	for i := 0; i < b.N; i++ {
		ai.evalBoardNoPointer(bd, ai.color)
	}
}

func (ai *AI) evalBoardPointer(bd board.Board, color board.Color) int {
	point := 0
	opponent := color.Opponent()
	for i := 0; i < ai.boardSize; i++ {
		for j := 0; j < ai.boardSize; j++ {
			if bd.AtXY(i, j) == color {
				point += ai.valueNetWork[i][j]
			} else if bd.AtXY(i, j) == opponent {
				point -= ai.valueNetWork[i][j]
			}
		}
	}
	return point
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
