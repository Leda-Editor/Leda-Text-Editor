package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	handling "github.com/Leda-Editor/Leda-Text-Editor/pkg/handling"
)

// displays dialog for configuring auto-save settings
func ShowAutoSaveSettings(ui *UI) {
	// create slider for setting auto-save interval
	intervalSlider := widget.NewSlider(5, 300)
	intervalSlider.SetValue(float64(handling.GetAutoSaveDelay().Seconds()))

	// create label to show current value
	intervalLabel := widget.NewLabel(fmt.Sprintf("Auto-save interval: %d seconds", int(intervalSlider.Value)))

	// update label when the slider changes
	intervalSlider.OnChanged = func(value float64) {
		intervalLabel.SetText(fmt.Sprintf("Auto-save interval: %d seconds", int(value)))
	}

	// create checkbox for enabling/disabling auto-save
	autoSaveCheck := widget.NewCheck("Enable auto-save", func(checked bool) {
		if checked {
			handling.StartAutoSave(ui.Window, ui.Editor)
		} else {
			handling.StopAutoSave()
		}
	})
	autoSaveCheck.SetChecked(handling.IsAutoSaveEnabled())

	// create text entry for manual entry of interval
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%d", int(intervalSlider.Value)))

	// update  slider when text entry changes
	intervalEntry.OnSubmitted = func(s string) {
		value, err := strconv.Atoi(s)
		if err != nil || value < 5 || value > 300 {
			dialog.ShowError(fmt.Errorf("please enter a number between 5 and 300"), ui.Window)
			return
		}
		intervalSlider.SetValue(float64(value))
	}

	// create set button for the manual entry
	setButton := widget.NewButton("Set", func() {
		intervalEntry.OnSubmitted(intervalEntry.Text)
	})

	// create container for manual entry and set button
	manualEntry := container.NewBorder(nil, nil, nil, setButton, intervalEntry)

	// create status label
	statusLabel := widget.NewLabel("")
	if handling.IsAutoSaveEnabled() {
		statusLabel.SetText("Auto-save is currently enabled")
	} else {
		statusLabel.SetText("Auto-save is currently disabled")
	}

	// create save location label
	saveLocationLabel := widget.NewLabel("")
	if handling.CurrentFile != nil {
		saveLocationLabel.SetText(fmt.Sprintf("Saving to: %s", handling.CurrentFile.String()))
	} else {
		saveLocationLabel.SetText("No file selected. Auto-save will not function until you save the file.")
	}

	// create apply button
	applyButton := widget.NewButton("Apply", func() {
		// apply settings
		handling.SetAutoSaveDelay(int(intervalSlider.Value), ui.Editor)

		// update status label
		if handling.IsAutoSaveEnabled() {
			statusLabel.SetText("Auto-save is currently enabled")

		} else {
			statusLabel.SetText("Auto-save is currently disabled")
		}

		// Show confirmation
		ShowTemporaryPopup(ui.Window, "Auto-save settings have been updated.")
	})

	// create content for dialog
	content := container.NewVBox(
		autoSaveCheck,
		widget.NewLabel("Drag to set interval:"),
		intervalSlider,
		intervalLabel,
		widget.NewLabel("Or enter manually (5-300 seconds):"),
		manualEntry,
		widget.NewSeparator(),
		statusLabel,
		saveLocationLabel,
		applyButton,
	)

	// create and show dialog
	settingsDialog := dialog.NewCustom("Auto-save Settings", "Close", content, ui.Window)
	settingsDialog.Resize(fyne.NewSize(400, 300))
	settingsDialog.Show()
}

func ShowTemporaryPopup(window fyne.Window, message string) {
	popup := widget.NewPopUp(widget.NewLabel(message), window.Canvas())
	popup.Show()
	go func() { time.Sleep(5 * time.Second); popup.Hide() }()
}
