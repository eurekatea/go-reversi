package builtinai

// func TestPartialCountChange(t *testing.T) {
// 	ai := New(board.BLACK, 6, 0)

// 	bd := board.NewBoard(6)
// 	bd.AssignBoard("+++X++++X++++XOOO+++OOX+++O+++++++++")

// 	p := board.NewPoint(5, 2)

// 	currentV := bd.CountPieces(board.BLACK)
// 	c := bd.Copy()
// 	if !c.Put(board.BLACK, p) {
// 		t.Error(c.Visualize())
// 		t.Fatal("cannot put")
// 	}
// 	newV := c.CountPieces(board.BLACK)

// 	aiV := ai.countAfterPut(bd, currentV, p, board.BLACK)

// 	if newV != aiV {
// 		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
// 		t.Error(c.Visualize())
// 	}

// 	ai = New(board.WHITE, 6, 0)
// 	bd = board.NewBoard(6)
// 	bd.AssignBoard("++++++++++++XXOOO++XXOO+O+XXO++XXXO+")

// 	p = board.NewPoint(1, 4)

// 	currentV = bd.CountPieces(board.WHITE)
// 	c = bd.Copy()
// 	if !c.Put(board.WHITE, p) {
// 		t.Error(c.Visualize())
// 		c.Assign(board.WHITE, p.X, p.Y)
// 		t.Error(c.Visualize())
// 		t.Fatal("cannot put")
// 	}
// 	newV = c.CountPieces(board.WHITE)

// 	aiV = ai.countAfterPut(bd, currentV, p, board.WHITE)

// 	if newV != aiV {
// 		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
// 		t.Error(c.Visualize())
// 	}

// 	ai = New(board.WHITE, 6, 0)
// 	bd = board.NewBoard(6)
// 	bd.AssignBoard("++++++++O+X+XXOOO++XXXXXO+XXO+OOOOO+")

// 	p = board.NewPoint(1, 4)

// 	currentV = bd.CountPieces(board.WHITE)
// 	c = bd.Copy()
// 	if !c.Put(board.WHITE, p) {
// 		t.Error(c.Visualize())
// 		c.Assign(board.WHITE, p.X, p.Y)
// 		t.Error(c.Visualize())
// 		t.Fatal("cannot put")
// 	}
// 	newV = c.CountPieces(board.WHITE)

// 	aiV = ai.countAfterPut(bd, currentV, p, board.WHITE)

// 	if newV != aiV {
// 		t.Error("error, orig:", currentV, "real:", newV, "but:", aiV)
// 		t.Error(c.Visualize())
// 	}
// }
