package handling

import (
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var (
	autoSaveEnabled bool
	autoSaveTimer   *time.Timer
	autoSaveDelay   time.Duration = 5 * time.Second // default 5 seconds
	CurrentFile     fyne.URI
	OnFileChanged func(fyne.URI)
)

// opens a file dialog and loads the selected file's content into the editor.
func OpenFile(window fyne.Window, editor *widget.Entry) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if reader == nil { //if no file is selected
			return
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		editor.SetText(string(data))
		CurrentFile = reader.URI() //stores current url

		if OnFileChanged != nil {
			OnFileChanged(CurrentFile)
		}
	}, window)
}

// saves to the currently open file
func SaveFile(window fyne.Window, editor *widget.Entry) {
	if CurrentFile == nil {
		SaveFileAs(window, editor) // no file is open, so use SaveFileAs instead
		return
	}
	// open the current file for writing
	writer, err := storage.Writer(CurrentFile)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	defer writer.Close()
	// write the content
	_, err = writer.Write([]byte(editor.Text))
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
}

// opens a file dialog and saves the editor's content to the selected file.
func SaveFileAs(window fyne.Window, editor *widget.Entry) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(editor.Text)) // converts to bytes and writes to file
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		CurrentFile = writer.URI() // change current file to new location
	}, window)
}

// clears the editor's content and resets current file
func ClearEditor(editor *widget.Entry) {
	dialog.ShowConfirm("Clear Editor", "Are you sure you want to clear the text and reset the file?", func(confirmed bool) {
		if confirmed {
			editor.SetText("")
			CurrentFile = nil
			if autoSaveEnabled {
				StopAutoSave()
			}
			fmt.Println("Editor cleared, file reset.")
		}
	}, nil)
}

// changes auto-save interval
func SetAutoSaveDelay(seconds int, editor *widget.Entry) {
	autoSaveDelay = time.Duration(seconds) * time.Second
	fmt.Println(autoSaveDelay)
	// if enabled, restart with new delay
	if autoSaveEnabled {
		StopAutoSave()
		StartAutoSave(fyne.CurrentApp().Driver().AllWindows()[0], editor) //sends nil! ERROR !!!
	}
}

// begins automatic saving at interval
func StartAutoSave(window fyne.Window, editor *widget.Entry) {
	if autoSaveEnabled {
		return // running
	}
	autoSaveEnabled = true
	fmt.Println("auto-save enabled")
	fmt.Println(autoSaveDelay)
	// repeating timer
	var scheduleAutoSave func()
	scheduleAutoSave = func() {
		autoSaveTimer = time.AfterFunc(autoSaveDelay, func() {
			// only save if you have a current file and auto-save enabled
			if CurrentFile != nil {
				SaveFile(window, editor)
				fmt.Println("saved!")
			} else {
				StopAutoSave() // Stop auto-save if there's no file
			}
			// schedule next save
			if autoSaveEnabled {
				scheduleAutoSave()
				fmt.Println("auto-saving in: ", autoSaveDelay)
			}
		})
	}
	scheduleAutoSave()
}

// stops automatic saving
func StopAutoSave() {
	if !autoSaveEnabled {
		return // not running
	}
	autoSaveEnabled = false
	fmt.Println("auto-save disabled")
	if autoSaveTimer != nil {
		autoSaveTimer.Stop()
		autoSaveTimer = nil
	}
}

// turns auto-save on or off
func ToggleAutoSave(window fyne.Window, editor *widget.Entry) bool {
	if autoSaveEnabled {
		StopAutoSave()
		return false
	} else {
		StartAutoSave(window, editor)
		return true
	}
}

// returns whether auto-save is enabled
func IsAutoSaveEnabled() bool {
	return autoSaveEnabled
}

// returns current auto-save delay
func GetAutoSaveDelay() time.Duration {
	return autoSaveDelay
}
