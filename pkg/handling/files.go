package handling

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/storage"
)

var CurrentFile fyne.URI

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
	}, window)
}

// saves to the currently open file
func SaveFile(window fyne.Window, editor *widget.Entry) {
	if CurrentFile == nil {
		SaveFileAs(window, editor) // No file is open, so use SaveFileAs instead
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
	editor.SetText("")
	CurrentFile = nil 
}
