package othellotheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct {
}

func (th Theme) Font(s fyne.TextStyle) fyne.Resource {
	return resourceJfOpenhuninn10Ttf
}

func (th Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DarkTheme().Color(name, variant)
}

func (th Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (th Theme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}
