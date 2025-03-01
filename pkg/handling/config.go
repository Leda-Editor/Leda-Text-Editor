package handling

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Config struct matching JSON structure
type Config struct {
	Theme struct {
		Mode             string `json:"mode"`
		BackgroundColour string `json:"backgroundColour"`
		TextColour       string `json:"textColour"`
		EditorColour     string `json:"editorColour"`
		MenuColour       string `json:"menuColour"`
		ButtonColour     string `json:"buttonColour"`
		PrimaryColour    string `json:"primaryColour"`
	} `json:"theme"`
	ZoomPercent int `json:"zoomPercent"`
	Fonts       struct {
		Default string `json:"default"`
		Bold    string `json:"bold"`
		Italic  string `json:"italic"`
	} `json:"fonts"`
}

// LoadConfig reads the config.json file and parses it into a Config struct
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func OpenConfigFile(window fyne.Window, editor *widget.Entry) {
	// Define the config file path (modify as needed)
	execPath, err := os.Executable()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	configPath := filepath.Join(filepath.Dir(execPath), "config.json")

	// Check if the file exists
	file, err := os.Open(configPath)
	if err != nil {
		workdir, err := os.Getwd()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		configPath = filepath.Join(workdir, "config.json")
		file, err = os.Open(configPath)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
	}

	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Set the text in the editor
	editor.SetText(string(data))
}
