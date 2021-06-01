package main

import (
	"othello/game"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	game.New(a, 8)
}
