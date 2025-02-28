package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme defines a theme with user-defined colors.
type CustomTheme struct {
	Base             fyne.Theme
	Background       color.Color
	Foreground       color.Color
	EditorBg         color.Color
	Primary          color.Color
	MenuBg           color.Color
	ButtonBackground color.Color
	FontSize         float32
	DefaultFont      fyne.Resource
	BoldFont         fyne.Resource
	ItalicFont       fyne.Resource
}

// Colour overrides default colors with custom values.
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return t.Background
	case theme.ColorNameForeground:
		return t.Foreground
	case theme.ColorNamePrimary:
		return t.Primary
	case theme.ColorNameInputBackground:
		return t.EditorBg
	case theme.ColorNameMenuBackground:
		return t.MenuBg
	case theme.ColorNameButton:
		return t.ButtonBackground
	case theme.ColorNameSeparator:
		return t.Primary
	case theme.ColorNameOverlayBackground:
		return t.Background
	}
	return theme.DefaultTheme().Color(name, variant) // Fallback to default
}

// Font overrides font based on style.
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return t.BoldFont
	}
	if style.Italic {
		return t.ItalicFont
	}
	return t.DefaultFont
}

// Icon overrides the default icons.
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name) // Keep default icons
}

// Size overrides default sizes.
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) // Keep default sizes
}

// NewCustomTheme initializes a theme with user-defined colors.
func NewCustomTheme(bg, fg, primary, editorBg, menuBg, buttonBg color.Color, defaultFont, boldFont, italicFont fyne.Resource) fyne.Theme {
	return &CustomTheme{
		Background:       bg,
		Foreground:       fg,
		EditorBg:         editorBg,
		Primary:          primary,
		MenuBg:           menuBg,
		ButtonBackground: buttonBg,
		DefaultFont:      defaultFont,
		BoldFont:         boldFont,
		ItalicFont:       italicFont,
	}
}
