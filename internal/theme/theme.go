package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Theme represents a color theme
type Theme struct {
	Name   string                 `json:"name"`
	Colors map[string]interface{} `json:"colors"`
}

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return &Theme{
		Name: "default",
		Colors: map[string]interface{}{
			"accent": map[string]string{
				"dark":  "93A9EC",
				"light": "93A9EC",
			},
			"primary": map[string]string{
				"dark":  "FFFFFF",
				"light": "000000",
			},
			"secondary": map[string]string{
				"dark":  "CCCCCC",
				"light": "666666",
			},
		},
	}
}

// LoadTheme loads a theme from a file
func LoadTheme(themePath string) (*Theme, error) {
	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme file: %w", err)
	}

	return &theme, nil
}

// GetThemePath returns the path to a theme file
func GetThemePath(themeName string) (string, error) {
	// Try multiple locations for theme files
	possiblePaths := []string{
		filepath.Join("/usr/local/share/gk/themes", themeName+".json"),
		filepath.Join(os.Getenv("HOME"), ".config/gk/themes", themeName+".json"),
		filepath.Join(os.Getenv("HOME"), ".gk/themes", themeName+".json"),
		filepath.Join(".", "themes", themeName+".json"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("theme '%s' not found", themeName)
}

// ListThemes lists all available themes
func ListThemes() ([]string, error) {
	themeDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".config/gk/themes"),
		filepath.Join(os.Getenv("HOME"), ".gk/themes"),
		filepath.Join(".", "themes"),
	}

	var themes []string
	seen := make(map[string]bool)

	for _, dir := range themeDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
				name := entry.Name()[:len(entry.Name())-5] // Remove .json extension
				if !seen[name] {
					themes = append(themes, name)
					seen[name] = true
				}
			}
		}
	}

	return themes, nil
}
