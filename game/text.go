package game

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type Text struct {
	cont *canvas.Text
}

func NewText(title string, size float32, alignment fyne.TextAlign) Text {
	t := Text{}
	t.cont = canvas.NewText(title, theme.ForegroundColor())
	t.cont.TextSize = size
	t.cont.Alignment = alignment
	return t
}

func (t Text) Update(title string) {
	t.cont.Text = title
	t.cont.Refresh()
}

func (t Text) CanvasText() *canvas.Text {
	return t.cont
}

func (t Text) SetMaxSize(max float32) {
	for t.cont.MinSize().Width > max {
		t.cont.Text = t.cont.Text[:len(t.cont.Text)-1]
	}
}
