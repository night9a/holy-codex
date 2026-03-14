package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CodicTheme is a custom Fyne theme evoking an ancient wooden codex:
// warm parchment backgrounds, dark walnut text, aged-ink accents.
type CodicTheme struct{}

var _ fyne.Theme = (*CodicTheme)(nil)

// ─── Colours ──────────────────────────────────────────────────────────────────

var (
	parchment     = color.NRGBA{R: 0xF2, G: 0xE8, B: 0xD5, A: 0xFF} // warm cream
	walnut        = color.NRGBA{R: 0x2C, G: 0x1A, B: 0x0E, A: 0xFF} // deep brown text
	aged_ink      = color.NRGBA{R: 0x6B, G: 0x3D, B: 0x11, A: 0xFF} // primary/action
	burnt_sienna  = color.NRGBA{R: 0x8B, G: 0x45, B: 0x13, A: 0xFF} // hover/focus
	vellum        = color.NRGBA{R: 0xE8, G: 0xDC, B: 0xC4, A: 0xFF} // input bg
	dark_wood     = color.NRGBA{R: 0x1A, G: 0x0D, B: 0x05, A: 0xFF} // overlay bg
	dust          = color.NRGBA{R: 0xAA, G: 0x99, B: 0x80, A: 0xFF} // disabled / placeholder
	scroll_shadow = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x55} // shadow
)

func (t *CodicTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return parchment
	case theme.ColorNameForeground:
		return walnut
	case theme.ColorNamePrimary:
		return aged_ink
	case theme.ColorNameFocus:
		return burnt_sienna
	case theme.ColorNameButton:
		return aged_ink
	case theme.ColorNameDisabled:
		return dust
	case theme.ColorNameDisabledButton:
		return dust
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x8B, G: 0x45, B: 0x13, A: 0x33}
	case theme.ColorNameInputBackground:
		return vellum
	case theme.ColorNamePlaceHolder:
		return dust
	case theme.ColorNamePressed:
		return burnt_sienna
	case theme.ColorNameScrollBar:
		return aged_ink
	case theme.ColorNameShadow:
		return scroll_shadow
	case theme.ColorNameHeaderBackground:
		return dark_wood
	case theme.ColorNameMenuBackground:
		return vellum
	case theme.ColorNameOverlayBackground:
		return dark_wood
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x8B, G: 0x45, B: 0x13, A: 0x55}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x6B, G: 0x3D, B: 0x11, A: 0x88}
	}
	return theme.DefaultTheme().Color(name, variant)
}

// ─── Fonts ────────────────────────────────────────────────────────────────────

func (t *CodicTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Falls back to Fyne default; swap with an embedded serif font (e.g. IM Fell)
	// by loading it via fyne.NewStaticResource and returning here.
	return theme.DefaultTheme().Font(style)
}

// ─── Icons ────────────────────────────────────────────────────────────────────

func (t *CodicTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// ─── Sizes ────────────────────────────────────────────────────────────────────

func (t *CodicTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNameSubHeadingText:
		return 17
	case theme.SizeNameInnerPadding:
		return 6
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameScrollBar:
		return 6
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputBorder:
		return 2
	}
	return theme.DefaultTheme().Size(name)
}