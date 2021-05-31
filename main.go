package main

import (
	"othello/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowIcon(game.Icon())
	ebiten.SetWindowSize(game.WIN_WIDTH, game.WIN_HEIGHT)
	ebiten.SetWindowTitle("othello")

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
