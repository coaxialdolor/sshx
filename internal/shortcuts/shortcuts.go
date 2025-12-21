package shortcuts

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ResolveShortcut resolves a shortcut name to an actual path
func ResolveShortcut(shortcut string) (string, error) {
	shortcut = strings.ToLower(strings.TrimSpace(shortcut))

	var homeDir string
	var err error

	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	switch shortcut {
	case "home":
		return homeDir, nil
	case "downloads":
		return filepath.Join(homeDir, "Downloads"), nil
	case "desktop":
		return filepath.Join(homeDir, "Desktop"), nil
	case "documents":
		return filepath.Join(homeDir, "Documents"), nil
	case "pictures":
		return filepath.Join(homeDir, "Pictures"), nil
	case "videos":
		if runtime.GOOS == "darwin" {
			// Check if Videos exists, otherwise use Movies
			videosPath := filepath.Join(homeDir, "Videos")
			if _, err := os.Stat(videosPath); err == nil {
				return videosPath, nil
			}
			return filepath.Join(homeDir, "Movies"), nil
		}
		return filepath.Join(homeDir, "Videos"), nil
	default:
		return "", nil // Not a shortcut
	}
}

// IsShortcut checks if a given string is a known shortcut
func IsShortcut(name string) bool {
	resolved, err := ResolveShortcut(name)
	return err == nil && resolved != ""
}

