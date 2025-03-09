package ui

import (
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	handling "github.com/Leda-Editor/Leda-Text-Editor/pkg/handling"
	"github.com/fyne-io/terminal"
)

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Reference to the Fyne Application.
	App fyne.App
	// Window provides access to the OS window.
	Window fyne.Window

	// Core state.
	// Editor retains raw text in an edit buffer.
	Editor *widget.Entry
	// Markdown retains rich text interactions: clicks, hovers and longpresses.
	Markdown *widget.RichText
	// MenuBar adds a menu to the window.
	MenuBar *fyne.Container
	// Terminal window
	Terminal *terminal.Terminal
	// Theme allows to customize theme, such as font size.
	Theme *Theme
	// CharacterLabel & LineLabel creates labels for the respective counters.
	CharacterLabel *widget.Label
	LineLabel      *widget.Label

	// current file
	CurrentFileLabel *widget.Label

	// Search/Replace Sidebar
	// SearchAreaContainer Holds Search UI
	SearchAreaContainer *fyne.Container
	// SearchTermEntry where you can type text to find.
	SearchTermEntry *widget.Entry
	// ReplaceTermEntry where you can type text to replace matched occurrences.
	ReplaceTermEntry *widget.Entry
	// SearchResults displays number of matches.
	SearchResults *widget.Label
	// Matches hold indices of all occurrences.
	Matches []int
	// CurrentMatchIdx keeps track of current match.
	CurrentMatchIdx int
	// SidebarVisible indicates whether sidebar is currently visible.
	SidebarVisible bool
	// OriginalText stores original text before search markers are added.
	OriginalText string

	// Markdown visibility toggle
	ShowMarkdown bool

	ZoomLabel *widget.Label
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	theme := NewTheme(app)

	ui := &UI{
		App:                 app,
		Window:              win,
		Editor:              widget.NewMultiLineEntry(),
		Markdown:            widget.NewRichTextFromMarkdown(""),
		Theme:               theme,
		CharacterLabel:      widget.NewLabelWithStyle("Characters: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		LineLabel:           widget.NewLabelWithStyle("Lines: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		CurrentFileLabel:    widget.NewLabelWithStyle("Current File: None", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		SearchAreaContainer: container.NewVBox(),
		SearchTermEntry:     widget.NewEntry(),
		ReplaceTermEntry:    widget.NewEntry(),
		SearchResults:       widget.NewLabel("Results: 0"),
		SidebarVisible:      false,
		Matches:             []int{},
		CurrentMatchIdx:     -1,
		OriginalText:        "",
		ShowMarkdown:        true,
		ZoomLabel:           widget.NewLabelWithStyle("ZoomL 100%", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	}

	config, err := handling.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		config = &handling.Config{}
	}

	ui.Theme.SetThemeFromConfig(config)

	ui.MenuBar = ui.CreateMenuBar()
	ui.Terminal = terminal.New()

	go func() {
		// Wait before starting the shell to prevent resizing errors
		time.Sleep(time.Millisecond * 500)
		if err := ui.Terminal.RunLocalShell(); err != nil {
			fmt.Println("Error running local shell:", err)
		}
	}()

	ui.Theme.ApplyTheme()
	ApplyUserTheme(ui)
	ui.Window.Content().Refresh()

	registerShortcuts(win, ui)

	// Update Markdown Preview whenever text changes.
	ui.Editor.OnChanged = func(content string) {
		ui.RenderMarkdown(content)
		ui.UpdateCounts(content)
	}

	// update markdown preview when file changes
	handling.OnFileChanged = func(uri fyne.URI) {
		ui.UpdateFileLabel(uri)
	}

	return ui
}

func registerShortcuts(win fyne.Window, ui *UI) {
	// Toggle Find Sidebar (Ctrl + F).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ShowSearchUI(false, ui)
	})
	// Toggle Find & Replace Sidebar (Ctrl + H).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyH, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ShowSearchUI(true, ui)
	})
	// Find Next (Ctrl + G).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyG, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		if ui.SidebarVisible {
			ui.nextMatch()
		}
	})
	// Find Previous (Ctrl + Shift + G).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyG, Modifier: desktop.ControlModifier | desktop.ShiftModifier}, func(shortcut fyne.Shortcut) {
		if ui.SidebarVisible {
			ui.previousMatch()
		}
	})
	// Open File (Ctrl + O)
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		handling.OpenFile(ui.Window, ui.Editor)
	})
	// Save File (Ctrl + S).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		handling.SaveFile(ui.Window, ui.Editor)
	})
	// Save File As (Ctrl + Shift + S).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: desktop.ControlModifier | desktop.ShiftModifier}, func(shortcut fyne.Shortcut) {
		handling.SaveFileAs(ui.Window, ui.Editor)
	})
	// Zoom In (Ctrl + =).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyEqual, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ui.ZoomIn()
	})
	// Zoom Out (Ctrl + -).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyMinus, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ui.ZoomOut()
	})
	// Reset Zoom (Ctrl + 0).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.Key0, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ui.ResetZoom()
	})
	// Toggle Markdown Preview (Ctrl + M).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyM, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ui.toggleMarkdownPreview()
	})
	// Toggle Dark Mode (Ctrl + D).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		ToggleDarkMode(ui.App, ui)
	})
	// Open Custom Theme Settings (Ctrl + Shift + T).
	win.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: desktop.ControlModifier}, func(shortcut fyne.Shortcut) {
		OpenThemePickerModal(ui.App, ui.Window, ui)
	})
}

// Updates Markdown Preview.
func (ui *UI) RenderMarkdown(input string) {
	ui.Markdown.ParseMarkdown(input)
	ui.Markdown.Refresh()
}

// Zoom In/Out.
func (ui *UI) ZoomIn() {
	ui.Theme.ZoomIn(ui)
	ui.UpdateZoomLabel()
}

func (ui *UI) ZoomOut() {
	ui.Theme.ZoomOut(ui)
	ui.UpdateZoomLabel()
}

// Reset Zoom.
func (ui *UI) ResetZoom() {
	ui.Theme.ZoomPercent = 100
	ui.Theme.ApplyTheme()
	ui.UpdateZoomLabel()
	ui.Window.Content().Refresh()
}

// Toggle visibility of Markdown preview.
func (ui *UI) toggleMarkdownPreview() {
	ui.ShowMarkdown = !ui.ShowMarkdown
	ui.UpdateLayout()
}

// Update character & line counts.
func (ui *UI) UpdateCounts(content string) {
	charCount := len(content)
	lineCount := len(widget.NewTextGridFromString(content).Rows)

	// Update the labels.
	ui.CharacterLabel.SetText(fmt.Sprintf("Characters: %d", charCount))
	ui.LineLabel.SetText(fmt.Sprintf("Lines: %d", lineCount))
}

func (ui *UI) UpdateZoomLabel() {
	ui.ZoomLabel.SetText(fmt.Sprintf("Zoom: %d%%", ui.Theme.ZoomPercent))
	ui.Window.Content().Refresh()
}

// called when current file is changed
func (ui *UI) UpdateFileLabel(uri fyne.URI) {
	if ui.CurrentFileLabel == nil {
		ui.CurrentFileLabel = widget.NewLabel("Current File: None")
	}
	filename := "None"
	if uri != nil {
		filename = filepath.Base(uri.Path())
	}
	ui.CurrentFileLabel.SetText(fmt.Sprintf("Current File: %s", filename))
}
