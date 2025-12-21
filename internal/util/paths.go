package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/coaxialdolor/sshx/internal/config"
	"github.com/coaxialdolor/sshx/internal/shortcuts"
)

// ResolveLocalPath resolves a local target path according to the rules:
// - If none provided → use default download dir
// - If absolute (`/path`) → use as-is
// - If starts with `~` → expand home
// - If shortcut → map to OS folder
// - Else → treat as relative to default download dir
func ResolveLocalPath(localTarget string, cfg *config.Config) (string, error) {
	// If no target provided, use default download dir
	if localTarget == "" {
		return cfg.DefaultDownloadDir, nil
	}

	// If absolute path (Unix-style or Windows-style), use as-is
	if filepath.IsAbs(localTarget) {
		return localTarget, nil
	}

	// If starts with ~, expand home directory
	if strings.HasPrefix(localTarget, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if localTarget == "~" {
			return homeDir, nil
		}
		if strings.HasPrefix(localTarget, "~/") {
			return filepath.Join(homeDir, localTarget[2:]), nil
		}
		// Handle ~user format (basic)
		return filepath.Join(homeDir, localTarget[1:]), nil
	}

	// Check if it's a shortcut
	if resolved, err := shortcuts.ResolveShortcut(localTarget); err == nil && resolved != "" {
		return resolved, nil
	}

	// Otherwise, treat as relative to default download dir
	return filepath.Join(cfg.DefaultDownloadDir, localTarget), nil
}

// ExpandPath expands a path, resolving shortcuts and home directory
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return homeDir, nil
		}
		if strings.HasPrefix(path, "~/") {
			return filepath.Join(homeDir, path[2:]), nil
		}
		return filepath.Join(homeDir, path[1:]), nil
	}

	if runtime.GOOS == "windows" && strings.HasPrefix(path, "%") && strings.HasSuffix(path, "%") {
		// Handle Windows environment variables like %USERPROFILE%
		envVar := path[1 : len(path)-1]
		if val := os.Getenv(envVar); val != "" {
			return val, nil
		}
	}

	return path, nil
}

// LineBuffer is a buffer that accumulates characters until a newline
type LineBuffer struct {
	buffer []byte
}

// NewLineBuffer creates a new line buffer
func NewLineBuffer() *LineBuffer {
	return &LineBuffer{
		buffer: make([]byte, 0, 256),
	}
}

// Add adds a byte to the buffer and returns the line if complete
func (lb *LineBuffer) Add(b byte) (string, bool) {
	if b == '\n' || b == '\r' {
		if len(lb.buffer) > 0 {
			line := string(lb.buffer)
			lb.buffer = lb.buffer[:0]
			return line, true
		}
		return "", false
	}

	// Handle backspace
	if b == 127 || b == 8 { // DEL or BS
		if len(lb.buffer) > 0 {
			lb.buffer = lb.buffer[:len(lb.buffer)-1]
		}
		return "", false
	}

	// Add printable characters
	if b >= 32 && b < 127 {
		lb.buffer = append(lb.buffer, b)
	}

	return "", false
}

// GetCurrent returns the current buffer contents
func (lb *LineBuffer) GetCurrent() string {
	return string(lb.buffer)
}

// Reset clears the buffer
func (lb *LineBuffer) Reset() {
	lb.buffer = lb.buffer[:0]
}

