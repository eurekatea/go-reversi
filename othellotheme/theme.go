package othellotheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct {
}

func (th Theme) Font(s fyne.TextStyle) fyne.Resource {
	// instead of resourceDotGothic16RegularTtf
	return theme.DefaultTextBoldFont()
}

func (th Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DarkTheme().Color(name, theme.VariantDark)
}

func (th Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (th Theme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case "text":
		return 16
	case "iconInline":
		return 16
	default:
		return theme.DarkTheme().Size(name)
	}
}
