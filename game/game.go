package game

import (
	"os"
	"othello/game/board"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	COOLDOWN = time.Millisecond * 500
)

type game struct {
	window    fyne.Window
	bd        board.Board
	units     [][]*unit
	player1   player
	player2   player
	lastMove  board.Point
	winner    board.Color
	available []board.Point
	now       board.Color
}

func New(a fyne.App, size int) {
	window := a.NewWindow("othello")
	g := &game{}
	bd := board.NewBoard(size)

	units := make([][]*unit, size)
	for i := range units {
		units[i] = make([]*unit, size)
	}
	grid := container.New(layout.NewGridLayout(size))
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			u := newUnit(g, board.NONE, i, j)
			grid.Add(u)
			units[i][j] = u
		}
	}

	g.window = window
	g.units = units
	g.now = board.BLACK
	g.bd = bd
	g.lastMove = board.NewPoint(-9, -9)
	g.winner = board.NONE
	g.available = bd.AllValidPoint(board.BLACK)

	if _, err := os.Stat(AI1); err == nil {
		g.player1 = newCom(bd, board.BLACK, AI1)
	} else {
		g.player1 = newHuman(bd, board.BLACK)
	}

	if _, err := os.Stat(AI2); err == nil {
		g.player2 = newCom(bd, board.WHITE, AI2)
	} else {
		g.player2 = newHuman(bd, board.WHITE)
	}

	g.updateWindow()

	window.SetContent(grid)
	window.ShowAndRun()
}

func (g *game) updateWindow() {
	for i, line := range g.units {
		for j, u := range line {
			u.setColor(g.bd.AtXY(i, j))
			if g.bd.IsValidPoint(g.now, board.NewPoint(i, j)) {
				u.SetResource(possible)
			}
		}
	}
}

type unit struct {
	g *game
	widget.Icon
	x, y  int
	color board.Color
}

func (u *unit) setColor(cl board.Color) {
	if cl == board.BLACK {
		u.SetResource(blackImg)
	} else if cl == board.WHITE {
		u.SetResource(whiteImg)
	} else {
		u.SetResource(noneImg)
	}
}

func (u *unit) Tapped(ev *fyne.PointEvent) {
	p := board.NewPoint(u.x, u.y)
	if !u.g.bd.Put(u.g.now, p) {
		return
	}

	if u.g.now == board.BLACK {
		u.SetResource(blackImg)
	} else {
		u.SetResource(whiteImg)
	}

	u.g.now = u.g.now.Opponent()

	u.g.updateWindow()
}

func (u *unit) MinSize() fyne.Size {
	return fyne.NewSize(64, 64)
}

func newUnit(g *game, cl board.Color, x, y int) *unit {
	u := &unit{g: g, color: cl, x: x, y: y}
	u.setColor(cl)
	u.ExtendBaseWidget(u)
	return u
}
