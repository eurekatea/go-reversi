package main

import (
	"othello/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	ui := a.NewWindow("othello")

	var (
		count               int = 0
		size                int = 0
		card1, card2        *widget.Card
		sel1, sel2, selSize *widget.Select
		start               *widget.Button
		pathes              [2]string
	)

	sel1 = widget.NewSelect([]string{"human", "external AI"}, func(s string) {
		if s == "external AI" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e == nil && uc != nil {
					pathes[0] = uc.URI().Path()
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
	card1 = widget.NewCard("            black", "", container.NewCenter(sel1))

	sel2 = widget.NewSelect([]string{"human", "external AI"}, func(s string) {
		if s == "external AI" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e == nil && uc != nil {
					pathes[1] = uc.URI().Path()
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
	card2 = widget.NewCard("           white", "", container.NewCenter(sel2))

	selSize = widget.NewSelect([]string{"6x6", "8x8"}, func(s string) {
		if s == "8x8" {
			size = 8
		} else {
			size = 6
		}
	})
	selSize.SetSelected("6x6")

	top := container.NewGridWithColumns(2, card1, card2)
	center := widget.NewCard("                                size", "", container.NewCenter(selSize))
	start = widget.NewButton("           play           ", func() {
		c := game.New(a, ui, pathes, size)
		ui.SetContent(c)
	})
	start.Disable()

	body := container.NewVBox(
		container.NewPadded(),
		container.NewPadded(top),
		container.NewPadded(center),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewCenter(start),
	)

	title := widget.NewCard("                             othello", "", body)

	ui.SetContent(title)
	ui.Resize(fyne.NewSize(500, 450))
	ui.SetFixedSize(true)
	ui.CenterOnScreen()
	ui.ShowAndRun()
}
