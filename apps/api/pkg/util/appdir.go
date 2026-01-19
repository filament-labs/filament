package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// appDataDir returns the platform-appropriate directory for persistent app data (DB, logs).
// - Linux:   $XDG_DATA_HOME/<appName> or ~/.local/share/<appName>
// - macOS:   ~/Library/Application Support/<appName>
// - Windows: %LOCALAPPDATA%\<appName>
func AppDataDir(appName string) (string, error) {
	var dir string

	switch runtime.GOOS {
	case "linux":
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			dir = filepath.Join(xdg, appName)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			dir = filepath.Join(home, ".local", "share", appName)
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		dir = filepath.Join(home, "Library", "Application Support", appName)
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		dir = filepath.Join(localAppData, appName)
	default:
		base, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get cache directory: %w", err)
		}
		dir = filepath.Join(base, appName)
	}

	return dir, nil
}
