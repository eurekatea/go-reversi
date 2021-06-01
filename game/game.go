package game

import (
	"othello/game/board"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	COOLDOWN = time.Millisecond * 500
)

type game struct {
	window fyne.Window
	bd     board.Board
	units  [][]*unit
	com1   *com
	com2   *com
	now    board.Color
	over   bool
}

func New(a fyne.App, window fyne.Window, comPath [2]string, size int) *fyne.Container {
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
	g.over = false

	if comPath[0] != "human" {
		g.com1 = newCom(bd, board.BLACK, comPath[0])
	}
	if comPath[1] != "human" {
		g.com2 = newCom(bd, board.WHITE, comPath[1])
	}

	if g.com1 != nil || g.com2 != nil {
		go g.round()
	}

	g.update()

	restart := widget.NewButton("restart", func() {
		dialog.NewConfirm("confirm", "restart?", func(b bool) {
			if b {
				g.over = true
				window.SetContent(New(a, window, comPath, size))
			}
		}, window).Show()
	})

	// for 6x6, 8x8 will resize automatically
	window.Resize(fyne.NewSize(315.8, 356.6))

	return container.NewVBox(grid, restart)
}

func (g *game) isBot(cl board.Color) bool {
	if cl == board.BLACK {
		return g.com1 != nil
	} else {
		return g.com2 != nil
	}
}

func (g *game) round() {
	for !g.over {
		if g.isBot(g.now) {
			if g.now == board.BLACK {
				g.com1.move()
			} else {
				g.com2.move()
			}
			g.now = g.now.Opponent()
			g.update()
		}
		time.Sleep(time.Millisecond * 30)
	}
}

func (g *game) update() {
	count := g.showValidAndCount()
	if count == 0 {
		g.now = g.now.Opponent()
		g.showValidAndCount()
	}
	if g.over = g.bd.IsOver(); g.over {
		winner := g.bd.Winner()
		var text string
		if winner == board.NONE {
			text = "draw"
		} else {
			text = winner.String() + " won"
		}
		dialog.NewInformation("Game Over", text, g.window).Show()
	}
}

func (g *game) showValidAndCount() int {
	count := 0
	for i, line := range g.units {
		for j, u := range line {
			u.setColor(g.bd.AtXY(i, j))
			if g.bd.IsValidPoint(g.now, board.NewPoint(i, j)) {
				u.SetResource(possible)
				count++
			}
		}
	}
	return count
}

type unit struct {
	g *game
	widget.Icon
	x, y  int
	color board.Color
}

func newUnit(g *game, cl board.Color, x, y int) *unit {
	u := &unit{g: g, color: cl, x: x, y: y}
	u.setColor(cl)
	u.ExtendBaseWidget(u)
	return u
}

func (u *unit) Tapped(ev *fyne.PointEvent) {
	if u.g.isBot(u.g.now) {
		return
	}
	p := board.NewPoint(u.x, u.y)
	if !u.g.bd.Put(u.g.now, p) {
		return
	}

	temp := u.g.now
	u.g.now = u.g.now.Opponent()
	u.g.update()
	if temp == board.BLACK {
		u.SetResource(blackCurr)
	} else {
		u.SetResource(whiteCurr)
	}
}

func (u *unit) MinSize() fyne.Size {
	return fyne.NewSize(48, 48)
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
