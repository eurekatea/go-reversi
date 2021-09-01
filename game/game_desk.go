//go:build !android
// +build !android

package game

import "fyne.io/fyne/v2"

var (
	unitSize6x6 = fyne.NewSize(48, 48)
	unitSize8x8 = fyne.NewSize(48, 48)

	winSize6x6 = fyne.NewSize(316, 479)
	winSize8x8 = fyne.NewSize(420, 583)
)
