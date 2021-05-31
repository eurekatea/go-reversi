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
	bd := board.NewBoard()

	g := &game{
		turn:      true,
		over:      false,
		bd:        bd,
		winner:    board.NONE,
		available: bd.AllValidPoint(board.BLACK),
	}

	if _, err := os.Stat(AI1); err == nil {
		g.player1 = newCom(bd, board.BLACK, AI1)
	} else {
		g.player1 = newHuman(g, bd, board.BLACK)
	}

	if _, err := os.Stat(AI2); err == nil {
		g.player2 = newCom(bd, board.WHITE, AI2)
	} else {
		g.player2 = newHuman(g, bd, board.WHITE)
	}

	return g
}

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if time.Since(g.lastClick) > COOLDOWN {
			fmt.Println(g.bd)
			g.lastClick = time.Now()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyV) {
		if time.Since(g.lastClick) > COOLDOWN {
			fmt.Println(g.bd.Visualize())
			g.lastClick = time.Now()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		if time.Since(g.lastClick) > COOLDOWN {
			g.restart()
			g.lastClick = time.Now()
		}
	}

	if !g.over {
		g.round()
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

func (g *game) restart() {
	*g = *NewGame()
}
