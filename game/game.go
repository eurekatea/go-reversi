package game

import (
	"fmt"
	"os"
	"othello/board"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WIN_WIDTH  = 425
	WIN_HEIGHT = 450
	COOLDOWN   = time.Millisecond * 500
)

type game struct {
	turn      bool
	over      bool
	bd        *board.Board
	lastClick time.Time
	player1   player
	player2   player
	winner    board.Color
	available []board.Point
}

func NewGame() *game {
	var player1, player2 player
	bd := board.NewBoard()

	if _, err := os.Stat("engine/AI1" + programName); err == nil {
		player1 = newCom(bd, board.BLACK, "engine/AI1"+programName)
	} else {
		player1 = newHuman(bd, board.BLACK)
	}

	if _, err := os.Stat("engine/AI2" + programName); err == nil {
		player2 = newCom(bd, board.WHITE, "engine/AI2"+programName)
	} else {
		player2 = newHuman(bd, board.WHITE)
	}

	g := &game{
		turn:      true,
		over:      false,
		bd:        bd,
		player1:   player1,
		player2:   player2,
		winner:    board.NONE,
		available: bd.AllValidPoint(board.BLACK),
	}

	return g
}

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if time.Since(g.lastClick) > COOLDOWN {
			fmt.Println(g.bd)
			g.lastClick = time.Now()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if time.Since(g.lastClick) > COOLDOWN {
			g.restart()
			g.lastClick = time.Now()
		}
	}

	if !g.over {
		g.round()
	}
	if g.hint() {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
	g.drawStones(screen)
	if g.over {
		g.drawEnd(screen)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("fps: %.02f", ebiten.CurrentFPS()), WIN_WIDTH-65, 0)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// available wrong wrong wrong
func (g *game) round() {
	if g.turn {
		g.player1.move()
		if g.player1.isDone() {
			g.check(board.BLACK)
		}
	} else {
		g.player2.move()
		if g.player2.isDone() {
			g.check(board.WHITE)
		}
	}
}

func (g *game) check(cl board.Color) {
	g.available = g.bd.AllValidPoint(cl.Opponent())
	if len(g.available) != 0 { // if is 0 then skip opponent
		g.turn = !g.turn
	} else {
		g.available = g.bd.AllValidPoint(cl)
		if len(g.available) == 0 {
			g.over = true
			g.winner = g.bd.Winner()
		}
	}
}

func (g game) hint() bool {
	x, y := ebiten.CursorPosition()

	x = int(float64(x-MARGIN_X)/SPACE + FIX)
	y = int(float64(y-MARGIN_Y)/SPACE + FIX)

	for _, p := range g.available {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}

func (g *game) restart() {
	*g = *NewGame()
}
