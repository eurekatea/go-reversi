package main

import (
	"othello/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	ui := a.NewWindow("othello")

	var (
		count        int = 0
		card1, card2 *widget.Card
		sel1, sel2   *widget.Select
		start        *widget.Button
		left, right  *fyne.Container
		pathes       [2]string
	)
	sel1 = widget.NewSelect([]string{"human", "computer"}, func(s string) {
		if s == "computer" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e != nil {
					pathes[0] = uc.URI().String()
				}
			}, ui).Show()
		} else {
			pathes[0] = "human"
		}
		count++
		if count == 2 {
			start.Enable()
		}
	})
	card1 = widget.NewCard("        black", "", sel1)
	left = container.NewVBox(card1, sel1)

	sel2 = widget.NewSelect([]string{"human", "computer"}, func(s string) {
		if s == "computer" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e != nil {
					pathes[1] = uc.URI().String()
				}
			}, ui).Show()
		} else {
			pathes[1] = "human"
		}
		count++
		if count == 2 {
			start.Enable()
		}
	})
	card2 = widget.NewCard("        white", "", sel2)
	right = container.NewVBox(card2, sel2)

	start = widget.NewButton("start", func() {
		c := game.New(a, ui, pathes, 6)
		ui.SetContent(c)
	})
	start.Disable()

	center := container.NewGridWithColumns(2, left, right)
	all := container.NewVBox(
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(center),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(start),
	)

	table := container.New(layout.NewMaxLayout())
	table.Add(all)

	ui.SetContent(table)
	ui.Resize(fyne.NewSize(450, 450))
	ui.SetFixedSize(true)
	ui.CenterOnScreen()
	ui.ShowAndRun()
}
