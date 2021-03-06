package game

import (
	"fmt"
	"othello/board"
	"othello/builtinai"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	counterTextSize = 24
	timerTextSize   = 13
	nameTextSize    = 13
	maxNameLen      = 20
)

var (
	nullPoint = board.NewPoint(-1, -1)
	unitSize  = fyne.NewSize(48, 48)
)

type game struct {
	window fyne.Window
	bd     board.Board
	units  [][]*unit

	counterBlack Text
	counterWhite Text

	passBtn *widget.Button
	com1    computer
	com2    computer
	now     board.Color

	blackSpent time.Duration
	whiteSpent time.Duration

	haveHuman bool
	over      bool
}

func newNameText(winSize fyne.Size, params Parameter) *fyne.Container {
	var name string

	if params.BlackAgent == AgentHuman {
		name = "human"
	} else if params.BlackAgent == AgentBuiltIn {
		name = "AI: " + params.BlackAILevel.String()
	} else {
		path := strings.Split(params.BlackPath, "/")
		if len(path) != 0 {
			name = "AI: " + path[len(path)-1]
		}
	}
	left := NewText(name, nameTextSize, fyne.TextAlignLeading)
	left.SetMaxSize(winSize.Width / 2)

	if params.WhiteAgent == AgentHuman {
		name = "human"
	} else if params.WhiteAgent == AgentBuiltIn {
		name = "AI: " + params.WhiteAILevel.String()
	} else {
		path := strings.Split(params.WhitePath, "/")
		if len(path) != 0 {
			name = "AI: " + path[len(path)-1]
		}
	}
	right := NewText(name, nameTextSize, fyne.TextAlignTrailing)
	right.SetMaxSize(winSize.Width / 2)

	return container.NewGridWithColumns(2, left.CanvasText(), right.CanvasText())
}

func newCounterText() (Text, Text) {
	counter1 := NewText("", counterTextSize, fyne.TextAlignLeading)
	counter2 := NewText("", counterTextSize, fyne.TextAlignTrailing)

	return counter1, counter2
}

func New(a fyne.App, window fyne.Window, menu *fyne.Container, params Parameter, size int) *fyne.Container {
	g := &game{}

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

	if params.BlackAgent == AgentBuiltIn {
		if size == 6 {
			g.com1 = builtinai.NewAI6(builtinai.BLACK, params.BlackAILevel)
		} else {
			g.com1 = builtinai.NewAI8(builtinai.BLACK, params.BlackAILevel)
		}
	} else if params.BlackAgent == AgentExternal {
		g.com1 = newCom(board.BLACK, params.BlackPath)
	}
	if params.WhiteAgent == AgentBuiltIn {
		if size == 6 {
			g.com2 = builtinai.NewAI6(builtinai.WHITE, params.WhiteAILevel)
		} else {
			g.com2 = builtinai.NewAI8(builtinai.WHITE, params.WhiteAILevel)
		}
	} else if params.WhiteAgent == AgentExternal {
		g.com2 = newCom(board.WHITE, params.WhitePath)
	}

	g.window = window
	g.units = units
	g.now = params.GoesFirst
	g.bd = board.NewBoard(size)
	g.over = false
	g.haveHuman = g.com1 == nil || g.com2 == nil
	g.counterBlack, g.counterWhite = newCounterText()

	counterTile := container.NewGridWithColumns(2, g.counterBlack.CanvasText(), g.counterWhite.CanvasText())
	nameText := newNameText(window.Canvas().Size(), params)

	g.passBtn = widget.NewButtonWithIcon(
		"pass",
		theme.ContentRedoIcon(),
		func() {
			g.passBtn.Disable()
			g.now = g.now.Opponent()
			g.update(nullPoint)
		},
	)
	g.passBtn.Disable()

	restart := widget.NewButtonWithIcon(
		"restart",
		theme.MediaReplayIcon(),
		func() {
			dialog.NewConfirm("confirm", "restart?", func(b bool) {
				if b {
					g.cleanAndExit()
					newGame := New(a, window, menu, params, size)
					window.SetContent(newGame)
				}
			}, window).Show()
		},
	)

	editBtn := widget.NewButtonWithIcon(
		"edit",
		theme.DocumentCreateIcon(),
		func() {

		},
	)

	mainMenu := widget.NewButtonWithIcon(
		"menu",
		theme.HomeIcon(),
		func() {
			dialog.NewConfirm("confirm", "return to menu?", func(b bool) {
				if b {
					g.cleanAndExit()
					menu.Show()
					window.SetContent(menu)
				}
			}, window).Show()
		},
	)

	if g.com1 != nil || g.com2 != nil {
		go g.round()
	}
	g.update(nullPoint)

	return container.NewVBox(
		counterTile,
		nameText,
		container.NewCenter(grid),
		container.NewGridWithColumns(2, g.passBtn, restart),
		container.NewGridWithColumns(2, editBtn, mainMenu),
	)
}

func (g *game) isBot(cl board.Color) bool {
	if cl == board.BLACK {
		return g.com1 != nil
	} else {
		return g.com2 != nil
	}
}

func (g *game) round() {
	var out string
	var err error
	defer g.cleanAndExit()
	for !g.over {
		if g.isBot(g.now) {
			start := time.Now()
			if g.now == board.BLACK {
				out, err = g.com1.Move(g.bd.String())
			} else {
				out, err = g.com2.Move(g.bd.String())
			}
			spent := time.Since(start)
			fmt.Println(g.now, "side spent:", spent)
			if g.now == board.BLACK {
				g.blackSpent += spent
			} else {
				g.whiteSpent += spent
			}
			if err != nil {
				g.aiError(err)
				break
			}
			g.bd.PutStr(g.now, out)
			g.now = g.now.Opponent()
			g.update(board.StrToPoint(out))
		} else {
			time.Sleep(time.Millisecond * 30)
		}
	}
}

func (g *game) update(current board.Point) {
	g.over = g.bd.IsOver()
	count := g.showValidAndCount(current)
	if count == 0 && !g.over {
		if g.haveHuman {
			// current side is human
			if (g.now == board.BLACK && g.com1 == nil) || (g.now == board.WHITE && g.com2 == nil) {
				dialog.NewInformation("info", "you have to pass", g.window).Show()
				g.passBtn.Enable()
			} else { // current is computer
				dialog.NewInformation("info", "computer have to pass\nit's your turn", g.window).Show()
				g.now = g.now.Opponent()
				g.update(nullPoint)
			}
		} else {
			g.now = g.now.Opponent()
			g.showValidAndCount(current)
		}
	}
	g.refreshCounter()
	if g.over {
		g.gameOver()
	}
	fmt.Println(g.bd.String())
}

func (g *game) refreshCounter() {
	blacks := g.bd.CountPieces(board.BLACK)
	whites := g.bd.CountPieces(board.WHITE)
	g.counterBlack.Update(fmt.Sprintf("black: %2d", blacks))
	g.counterWhite.Update(fmt.Sprintf("white: %2d", whites))
}

func (g *game) gameOver() {
	var text string
	winner := g.bd.Winner()
	if winner == board.NONE {
		text = "draw"
	} else {
		text = winner.String() + " won"
	}
	d := dialog.NewInformation("Game Over", text, g.window)
	d.Resize(fyne.NewSize(250, 0))
	d.Show()
	fmt.Println("\ngame over")
	fmt.Println("black total:", g.blackSpent, ", white total:", g.whiteSpent)
}

func (g *game) showValidAndCount(current board.Point) int {
	count := 0
	for i, line := range g.units {
		for j, u := range line {
			cl := g.bd.AtXY(i, j)
			if g.bd.IsValidPoint(g.now, board.NewPoint(i, j)) {
				u.SetResource(possible)
				count++
			} else {
				u.setColor(cl)
			}
			if current.X == i && current.Y == j {
				u.setColorCurrent(cl)
			}
		}
	}
	return count
}

func (g *game) aiError(err error) {
	if !g.over {
		d := dialog.NewError(err, g.window)
		d.SetOnClosed(func() { panic(err) })
		d.Show()
	}
}

func (g *game) cleanAndExit() {
	g.over = true
	if g.com1 != nil {
		g.com1.Close()
	}
	if g.com2 != nil {
		g.com2.Close()
	}
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
	if !u.g.bd.PutPoint(u.g.now, p) {
		return
	}

	u.g.now = u.g.now.Opponent()
	u.g.update(p)
}

func (u *unit) MinSize() fyne.Size {
	return unitSize
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
	} else if cl == board.WHITE {
		u.SetResource(whiteCurr)
	}
}
