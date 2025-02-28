package handling

import (
	"encoding/json"
	"os"
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

// SaveConfig writes the updated config values back to config.json
func SaveConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(config)
}
