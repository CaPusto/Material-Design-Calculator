package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MaterialTheme struct{}

func (m *MaterialTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	dark := variant == theme.VariantDark

	switch name {
	case theme.ColorNameBackground:
		if dark {
			return color.RGBA{R: 24, G: 24, B: 24, A: 255}
		}
		return color.RGBA{R: 250, G: 250, B: 250, A: 255}
	case theme.ColorNameInputBackground:
		if dark {
			return color.RGBA{R: 32, G: 32, B: 32, A: 255}
		}
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}
	case theme.ColorNameButton:
		if dark {
			return color.RGBA{R: 43, G: 43, B: 43, A: 255}
		}
		return color.RGBA{R: 224, G: 224, B: 224, A: 255}
	case theme.ColorNamePrimary:
		if dark {
			return color.RGBA{R: 128, G: 203, B: 196, A: 255}
		}
		return color.RGBA{R: 0, G: 150, B: 136, A: 255}
	case theme.ColorNameForeground:
		if dark {
			return color.RGBA{R: 230, G: 230, B: 230, A: 255}
		}
		return color.RGBA{R: 33, G: 33, B: 33, A: 255}
	case theme.ColorNameHover:
		if dark {
			return color.RGBA{R: 60, G: 60, B: 60, A: 255}
		}
		return color.RGBA{R: 208, G: 208, B: 208, A: 255}
	case theme.ColorNameDisabled:
		if dark {
			return color.RGBA{R: 100, G: 100, B: 100, A: 255}
		}
		return color.RGBA{R: 158, G: 158, B: 158, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m *MaterialTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 18
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameScrollBar:
		return 6
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func (m *MaterialTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *MaterialTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
