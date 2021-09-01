//go:build android
// +build android

package game

import "fyne.io/fyne/v2"

var (
	unitSize6x6 = fyne.NewSize(54.0, 54.0)
	unitSize8x8 = fyne.NewSize(40.0, 40.0)

	winSize6x6 = fyne.NewSize(316, 426)
	winSize8x8 = fyne.NewSize(420, 530)
)
