package game

import (
	"fmt"
	"othello/board"
	builtinai "othello/builtinAI"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	AgentNone Agent = iota
	AgentHuman
	AgentBuiltIn
	AgentExternal
)

type Agent int

type Agents struct {
	BlackAgent Agent
	WhiteAgent Agent
	BlackPath  string
	WhitePath  string
}

func NewAgents() Agents {
	return Agents{
		BlackAgent: AgentNone,
		BlackPath:  "",
		WhiteAgent: AgentNone,
		WhitePath:  "",
	}
}

func (agents Agents) Selected() bool {
	return agents.BlackAgent != AgentNone && agents.WhiteAgent != AgentNone
}

type computer interface {
	Move(board.Board) (board.Point, error)
}

type game struct {
	window fyne.Window
	bd     board.Board
	units  [][]*unit
	com1   computer
	com2   computer
	now    board.Color
	over   bool
}

func New(a fyne.App, window fyne.Window, agents Agents, size int) *fyne.Container {
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

	if agents.BlackAgent == AgentBuiltIn {
		g.com1 = builtinai.New(board.BLACK, size)
	} else if agents.BlackAgent == AgentExternal {
		g.com1 = newCom(board.BLACK, agents.BlackPath)
	}
	if agents.WhiteAgent == AgentBuiltIn {
		g.com2 = builtinai.New(board.WHITE, size)
	} else if agents.WhiteAgent == AgentExternal {
		g.com2 = newCom(board.WHITE, agents.WhitePath)
	}

	if g.com1 != nil || g.com2 != nil {
		go g.round()
	}

	g.update(board.NewPoint(-1, -1))

	restart := widget.NewButton("restart", func() {
		dialog.NewConfirm("confirm", "restart?", func(b bool) {
			if b {
				g.over = true
				window.SetContent(New(a, window, agents, size))
			}
		}, window).Show()
	})

	// resize to minimum size
	window.Resize(fyne.NewSize(1, 1))

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
	var p board.Point
	var err error
	for !g.over {
		if g.isBot(g.now) {
			if g.now == board.BLACK {
				start := time.Now()
				p, err = g.com1.Move(g.bd.Copy())
				fmt.Println("black side spent:", time.Since(start))
			} else {
				start := time.Now()
				p, err = g.com2.Move(g.bd.Copy())
				fmt.Println("white side spent:", time.Since(start))
			}
			g.bd.Put(g.now, p)
			if err != nil {
				g.aiError(err)
			}
			g.now = g.now.Opponent()
			g.update(p)
		}
		time.Sleep(time.Millisecond * 30)
	}
}

func (g *game) update(current board.Point) {
	count := g.showValidAndCount(current)
	if count == 0 {
		g.now = g.now.Opponent()
		g.showValidAndCount(current)
	}
	if g.over = g.bd.IsOver(); g.over {
		g.gameOver()
	}
}

func (g *game) gameOver() {
	winner := g.bd.Winner()
	var text string
	if winner == board.NONE {
		text = "draw"
	} else {
		text = winner.String() + " won"
	}
	text += "\n"
	text += fmt.Sprintf("black pieces: %d\n", g.bd.CountPieces(board.BLACK))
	text += fmt.Sprintf("white pieces: %d\n", g.bd.CountPieces(board.WHITE))
	d := dialog.NewInformation("Game Over", text, g.window)
	d.Resize(fyne.NewSize(250, 0))
	d.Show()
}

func (g *game) showValidAndCount(current board.Point) int {
	count := 0
	for i, line := range g.units {
		for j, u := range line {
			cl := g.bd.AtXY(i, j)
			u.setColor(cl)
			if g.bd.IsValidPoint(g.now, board.NewPoint(i, j)) {
				u.SetResource(possible)
				count++
			}
			if current.X == i && current.Y == j {
				u.setColorCurrent(cl)
			}
		}
	}
	return count
}

func (g *game) aiError(err error) {
	d := dialog.NewError(err, g.window)
	d.SetOnClosed(func() { panic(err) })
	d.Show()
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

	u.g.now = u.g.now.Opponent()
	u.g.update(p)
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

func (u *unit) setColorCurrent(cl board.Color) {
	if cl == board.BLACK {
		u.SetResource(blackCurr)
	} else {
		u.SetResource(whiteCurr)
	}
}
