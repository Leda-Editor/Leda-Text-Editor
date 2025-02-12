package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	ui "github.com/Leda-Editor/Leda-Text-Editor/pkg/ui"
)

func main() {
	// Initialize Fyne Application.
	app := app.New()

	// Create a new window for the application
	window := app.NewWindow("Leda Text Editor")

	// Initialize UI.
	ledaUI := ui.NewUI(app, window)

	// Set up window layout.
	window.SetContent(ledaUI.Layout())

	// Set window size.
	window.Resize(fyne.NewSize(800, 600))
	// Display the window and start the event loop.
	window.ShowAndRun()
}
