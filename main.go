package main

import (
	"othello/game"
	"othello/othellotheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	titleTextSize = 54
	cardTextSize  = 30
)

var (
	initWinSize      = fyne.NewSize(500, 450)
	selectDialogSize = fyne.NewSize(250, 1)
)

func main() {
	// defer profile.Start(profile.BlockProfile, profile.CPUProfile).Stop()

	a := app.New()
	customTheme := othellotheme.Theme{}
	a.Settings().SetTheme(&customTheme)
	ui := a.NewWindow("othello")
	ui.SetIcon(game.WindowIcon)

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

		aiLevelSelection1 *widget.Select
		aiLevelSelection2 *widget.Select

		playButton *widget.Button

		top  *fyne.Container
		menu *fyne.Container
	)

	aiLevelSelection1 = widget.NewSelect(
		[]string{"beginner", "amateur", "professional", "expert", "master"},

		func(s string) {
			agents.BlackInternalAILevel = aiLevelSelection1.SelectedIndex()
		},
	)

	aiLevelSelection2 = widget.NewSelect(
		[]string{"beginner", "amateur", "professional", "expert", "master"},

		func(s string) {
			agents.WhiteInternalAILevel = aiLevelSelection2.SelectedIndex()
		},
	)

	selection1 = widget.NewSelect(
		[]string{"human", "built-in AI", "external AI"},

		func(s string) {
			if s == "external AI" {
				d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
					if e == nil && uc != nil {
						agents.BlackPath = uc.URI().Path()
						agents.BlackAgent = game.AgentExternal
						if agents.AllSelected() {
							playButton.Enable()
						}
					}
				}, ui)
				d.SetFilter(storage.NewExtensionFileFilter([]string{".exe", ".out", ""}))
				d.Show()
			} else if s == "human" {
				agents.BlackAgent = game.AgentHuman
			} else {
				agents.BlackAgent = game.AgentBuiltIn
				d := dialog.NewCustom("select AI level", "  ok  ", aiLevelSelection1, ui)
				d.Resize(selectDialogSize)
				d.Show()
			}
			if agents.AllSelected() {
				playButton.Enable()
			}
		},
	)

	subTitle1 := canvas.NewText("black side", theme.ForegroundColor())
	subTitle1.TextSize = cardTextSize
	subTitle1.Alignment = fyne.TextAlignCenter
	blackCard = widget.NewCard("", "", container.NewVBox(subTitle1, container.NewCenter(selection1)))

	selection2 = widget.NewSelect(
		[]string{"human", "built-in AI", "external AI"},

		func(s string) {
			if s == "external AI" {
				d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
					if e == nil && uc != nil {
						agents.WhitePath = uc.URI().Path()
						agents.WhiteAgent = game.AgentExternal
						if agents.AllSelected() {
							playButton.Enable()
						}
					}
				}, ui)
				d.SetFilter(storage.NewExtensionFileFilter([]string{".exe", ".out", ""}))
				d.Show()
			} else if s == "human" {
				agents.WhiteAgent = game.AgentHuman
			} else {
				agents.WhiteAgent = game.AgentBuiltIn
				d := dialog.NewCustom("select AI level", "  ok  ", aiLevelSelection2, ui)
				d.Resize(selectDialogSize)
				d.Show()
			}
			if agents.AllSelected() {
				playButton.Enable()
			}
		},
	)

	subtitle2 := canvas.NewText("white side", theme.ForegroundColor())
	subtitle2.TextSize = cardTextSize
	subtitle2.Alignment = fyne.TextAlignCenter
	whiteCard = widget.NewCard("", "", container.NewVBox(subtitle2, container.NewCenter(selection2)))

	selection3 = widget.NewSelect(
		[]string{"6x6", "8x8"},

		func(s string) {
			if s == "8x8" {
				boardSize = 8
			} else {
				boardSize = 6
			}
		},
	)
	selection3.SetSelected("6x6")

	top = container.NewGridWithColumns(2, blackCard, whiteCard)

	subtitle3 := canvas.NewText("board size", theme.ForegroundColor())
	subtitle3.TextSize = cardTextSize
	subtitle3.Alignment = fyne.TextAlignCenter
	center = widget.NewCard("", "", container.NewVBox(subtitle3, container.NewCenter(selection3)))

	playButton = widget.NewButtonWithIcon(
		"start play",
		theme.MediaPlayIcon(),
		func() {
			c := game.New(a, ui, menu, agents, boardSize)
			menu.Hide()
			ui.SetContent(c)
		},
	)
	playButton.Disable()

	title := canvas.NewText("othello", theme.ForegroundColor())
	title.TextSize = titleTextSize
	title.Alignment = fyne.TextAlignCenter

	all = widget.NewCard(
		"",
		"",
		container.NewVBox(
			title,
			container.NewPadded(),
			container.NewMax(top),
			container.NewPadded(),
			container.NewMax(center),
			container.NewPadded(),
			container.NewPadded(),
			container.NewPadded(),
			container.NewPadded(),
			container.NewPadded(),
			container.NewCenter(playButton),
		),
	)

	menu = container.NewMax(all)
	ui.Resize(initWinSize)
	ui.SetFixedSize(true)
	ui.CenterOnScreen()
	ui.SetContent(menu)
	ui.ShowAndRun()
}
