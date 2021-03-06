package main

import (
	"othello/board"
	"othello/builtinai"
	"othello/game"
	"othello/othellotheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	fDialog "github.com/sqweek/dialog"
)

const (
	titleTextSize = 54
	cardTextSize  = 30
)

var (
	initWinSize      = fyne.NewSize(500, 570)
	selectDialogSize = fyne.NewSize(250, 1)
)

func main() {
	// defer profile.Start(profile.BlockProfile, profile.CPUProfile).Stop()

	a := app.New()
	a.Settings().SetTheme(&othellotheme.Theme{})
	ui := a.NewWindow("othello")
	ui.SetIcon(game.WindowIcon)

	var (
		boardSize int

		params game.Parameter = game.NewParam()

		blackCard *widget.Card
		whiteCard *widget.Card
		all       *widget.Card
		center    *widget.Card
		goesFirst *widget.Card

		selection1 *widget.Select
		selection2 *widget.Select

		sizeSelect *widget.RadioGroup
		order      *widget.RadioGroup

		levelSelect1 *widget.Select
		levelSelect2 *widget.Select

		playButton *widget.Button

		top  *fyne.Container
		menu *fyne.Container
	)

	levels := []string{
		builtinai.LV_ONE.String(),
		builtinai.LV_TWO.String(),
		builtinai.LV_THREE.String(),
		builtinai.LV_FOUR.String(),
		builtinai.LV_FIVE.String(),
	}
	levelSelect1 = widget.NewSelect(
		levels,
		func(s string) {
			params.BlackAILevel = builtinai.Level(levelSelect1.SelectedIndex())
		},
	)
	levelSelect2 = widget.NewSelect(
		levels,
		func(s string) {
			params.WhiteAILevel = builtinai.Level(levelSelect2.SelectedIndex())
		},
	)

	selection1 = widget.NewSelect(
		[]string{"human", "built-in AI", "external AI"},

		func(s string) {
			if s == "external AI" {
				dir, err := fDialog.File().Load()
				if err == nil {
					params.BlackPath = dir
					params.BlackAgent = game.AgentExternal
					if params.AllSelected() {
						playButton.Enable()
					}
				}
			} else if s == "human" {
				params.BlackAgent = game.AgentHuman
			} else {
				params.BlackAgent = game.AgentBuiltIn
				d := dialog.NewCustom("select AI level", "  ok  ", levelSelect1, ui)
				d.Resize(selectDialogSize)
				d.Show()
			}
			if params.AllSelected() {
				playButton.Enable()
			}
		},
	)

	subtitle1 := game.NewText("black side", cardTextSize, fyne.TextAlignCenter)
	blackCard = widget.NewCard(
		"",
		"",
		container.NewVBox(
			subtitle1.CanvasText(),
			container.NewCenter(selection1),
		),
	)

	selection2 = widget.NewSelect(
		[]string{"human", "built-in AI", "external AI"},

		func(s string) {
			if s == "external AI" {
				dir, err := fDialog.File().Load()
				if err == nil {
					params.WhitePath = dir
					params.WhiteAgent = game.AgentExternal
					if params.AllSelected() {
						playButton.Enable()
					}
				}
			} else if s == "human" {
				params.WhiteAgent = game.AgentHuman
			} else {
				params.WhiteAgent = game.AgentBuiltIn
				d := dialog.NewCustom("select AI level", "  ok  ", levelSelect2, ui)
				d.Resize(selectDialogSize)
				d.Show()
			}
			if params.AllSelected() {
				playButton.Enable()
			}
		},
	)

	subtitle2 := game.NewText("white side", cardTextSize, fyne.TextAlignCenter)
	whiteCard = widget.NewCard(
		"",
		"",
		container.NewVBox(
			subtitle2.CanvasText(),
			container.NewCenter(selection2),
		),
	)

	top = container.NewGridWithColumns(2, blackCard, whiteCard)

	sizeSelect = widget.NewRadioGroup(
		[]string{"6x6", "8x8"},

		func(s string) {
			if s == "8x8" {
				boardSize = 8
			} else {
				boardSize = 6
			}
		},
	)
	sizeSelect.SetSelected("6x6")
	sizeSelect.Required = true

	subtitle3 := game.NewText("board size", cardTextSize, fyne.TextAlignCenter)
	center = widget.NewCard(
		"",
		"",
		container.NewVBox(
			subtitle3.CanvasText(),
			container.NewCenter(sizeSelect),
		),
	)

	order = widget.NewRadioGroup(
		[]string{"black first", "white first"},
		func(s string) {
			if s == "black first" {
				params.GoesFirst = board.BLACK
			} else {
				params.GoesFirst = board.WHITE
			}
		},
	)
	order.SetSelected("black first")
	order.Required = true

	goesFirst = widget.NewCard(
		"",
		"",
		container.NewCenter(order),
	)

	ruleButton := widget.NewButtonWithIcon(
		"      rule      ",
		theme.QuestionIcon(),
		func() {
			dialog.NewInformation(
				"rule",
				"rules here",
				ui,
			).Show()
		},
	)

	playButton = widget.NewButtonWithIcon(
		"      play      ",
		theme.MediaPlayIcon(),
		func() {
			c := game.New(a, ui, menu, params, boardSize)
			menu.Hide()
			ui.SetContent(c)
		},
	)
	playButton.Disable()

	title := game.NewText("othello", titleTextSize, fyne.TextAlignCenter)

	all = widget.NewCard(
		"",
		"",
		container.NewVBox(
			title.CanvasText(),
			container.NewMax(top),
			container.NewMax(center),
			container.NewMax(goesFirst),
			container.NewCenter(ruleButton),
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
