package main

import (
	"othello/game"
	"othello/othellotheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// defer profile.Start(profile.BlockProfile, profile.CPUProfile).Stop()

	a := app.New()
	a.Settings().SetTheme(&othellotheme.Theme{})
	ui := a.NewWindow("othello")

	var (
		boardSize int         = 0
		agents    game.Agents = game.NewAgents()

		blackCard *widget.Card
		whiteCard *widget.Card
		all       *widget.Card
		center    *widget.Card

		selection1 *widget.Select
		selection2 *widget.Select
		selection3 *widget.Select

		playButton *widget.Button

		top  *fyne.Container
		body *fyne.Container
	)

	selection1 = widget.NewSelect([]string{"human", "built-in AI", "external AI"}, func(s string) {
		if s == "external AI" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e == nil && uc != nil {
					agents.BlackPath = uc.URI().Path()
					agents.BlackAgent = game.AgentExternal
					if agents.Selected() {
						playButton.Enable()
					}
				} else {
					selection1.ClearSelected()
				}
			}, ui).Show()
		} else if s == "human" {
			agents.BlackAgent = game.AgentHuman
		} else {
			agents.BlackAgent = game.AgentBuiltIn
		}
		if agents.Selected() {
			playButton.Enable()
		}
	})
	blackCard = widget.NewCard("           black", "", container.NewCenter(selection1))

	selection2 = widget.NewSelect([]string{"human", "built-in AI", "external AI"}, func(s string) {
		if s == "external AI" {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e == nil && uc != nil {
					agents.WhitePath = uc.URI().Path()
					agents.WhiteAgent = game.AgentExternal
					if agents.Selected() {
						playButton.Enable()
					} else {
						selection1.ClearSelected()
					}
				}
			}, ui).Show()
		} else if s == "human" {
			agents.WhiteAgent = game.AgentHuman
		} else {
			agents.WhiteAgent = game.AgentBuiltIn
		}
		if agents.Selected() {
			playButton.Enable()
		}
	})
	whiteCard = widget.NewCard("           white", "", container.NewCenter(selection2))

	selection3 = widget.NewSelect([]string{"6x6", "8x8"}, func(s string) {
		if s == "8x8" {
			boardSize = 8
		} else {
			boardSize = 6
		}
	})
	selection3.SetSelected("6x6")

	top = container.NewGridWithColumns(2, blackCard, whiteCard)
	center = widget.NewCard("                       board  size", "", container.NewCenter(selection3))
	playButton = widget.NewButton("           play           ", func() {
		c := game.New(a, ui, agents, boardSize)
		ui.SetContent(c)
	})
	playButton.Disable()

	body = container.NewVBox(
		container.NewPadded(),
		container.NewPadded(top),
		container.NewPadded(center),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewPadded(),
		container.NewCenter(playButton),
	)

	all = widget.NewCard("                             othello", "", body)

	ui.SetContent(all)
	ui.Resize(fyne.NewSize(500, 450))
	ui.SetFixedSize(true)
	ui.CenterOnScreen()
	ui.ShowAndRun()
}
