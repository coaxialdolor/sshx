package alias

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ShellType represents the type of shell
type ShellType string

const (
	ShellZsh   ShellType = "zsh"
	ShellBash  ShellType = "bash"
	ShellFish  ShellType = "fish"
	ShellPwsh  ShellType = "powershell"
	ShellUnknown ShellType = "unknown"
)

// DetectShell detects the user's current shell
func DetectShell() ShellType {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			return ShellPwsh
		}
		return ShellUnknown
	}

	shellName := filepath.Base(shell)
	switch shellName {
	case "zsh":
		return ShellZsh
	case "bash":
		return ShellBash
	case "fish":
		return ShellFish
	default:
		return ShellUnknown
	}
}

// GetProfilePath returns the path to the shell profile file
func GetProfilePath(shell ShellType) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch shell {
	case ShellZsh:
		return filepath.Join(homeDir, ".zshrc"), nil
	case ShellBash:
		// Check for .bash_profile first (macOS), then .bashrc
		bashProfile := filepath.Join(homeDir, ".bash_profile")
		if _, err := os.Stat(bashProfile); err == nil {
			return bashProfile, nil
		}
		return filepath.Join(homeDir, ".bashrc"), nil
	case ShellFish:
		configDir := filepath.Join(homeDir, ".config", "fish")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", err
		}
		return filepath.Join(configDir, "config.fish"), nil
	case ShellPwsh:
		// PowerShell profile path
		if runtime.GOOS == "windows" {
			// Use PowerShell to get the profile path
			cmd := exec.Command("powershell", "-Command", "$PROFILE")
			output, err := cmd.Output()
			if err == nil {
				profilePath := strings.TrimSpace(string(output))
				if profilePath != "" {
					return profilePath, nil
				}
			}
		}
		// Fallback
		configDir := filepath.Join(homeDir, "Documents", "PowerShell")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", err
		}
		return filepath.Join(configDir, "Microsoft.PowerShell_profile.ps1"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// GetAliasLine returns the alias line for the given shell
func GetAliasLine(shell ShellType) string {
	switch shell {
	case ShellZsh, ShellBash:
		return `alias ssh="sshx"`
	case ShellFish:
		return `alias ssh="sshx"`
	case ShellPwsh:
		return `Set-Alias ssh sshx`
	default:
		return ""
	}
}

// HasAlias checks if the profile file already contains the alias
func HasAlias(profilePath string, shell ShellType) (bool, error) {
	data, err := os.ReadFile(profilePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	content := string(data)
	aliasLine := GetAliasLine(shell)

	// Check for various forms of the alias
	patterns := []string{
		aliasLine,
		`alias ssh="sshx"`,
		`alias ssh='sshx'`,
		`Set-Alias ssh sshx`,
	}

	for _, pattern := range patterns {
		if strings.Contains(content, pattern) {
			return true, nil
		}
	}

	return false, nil
}

// AddAlias adds the alias to the profile file
func AddAlias(profilePath string, shell ShellType) error {
	// Check if alias already exists
	hasAlias, err := HasAlias(profilePath, shell)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if hasAlias {
		return fmt.Errorf("alias already exists in %s", profilePath)
	}

	// Read existing content
	var content string
	if data, err := os.ReadFile(profilePath); err == nil {
		content = string(data)
	}

	// Add alias
	aliasLine := GetAliasLine(shell)
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += "# SSHX alias\n"
	content += aliasLine + "\n"

	// Write back
	return os.WriteFile(profilePath, []byte(content), 0644)
}

// RemoveAlias removes the alias from the profile file
func RemoveAlias(profilePath string, shell ShellType) error {
	data, err := os.ReadFile(profilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("profile file not found: %s", profilePath)
	}
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	skipNext := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this is the comment line
		if trimmed == "# SSHX alias" {
			skipNext = true
			continue
		}

		// Check if this is an alias line and we should skip it
		if skipNext && (strings.Contains(trimmed, "alias ssh") || strings.Contains(trimmed, "Set-Alias ssh")) {
			skipNext = false
			continue
		}

		// Also check if alias line appears without comment
		if strings.Contains(trimmed, "alias ssh") || strings.Contains(trimmed, "Set-Alias ssh") {
			// Make sure it's our alias, not a user's custom one
			if strings.Contains(trimmed, "sshx") {
				continue
			}
		}

		skipNext = false
		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	// Preserve original newline ending behavior
	originalContent := string(data)
	if strings.HasSuffix(originalContent, "\n") {
		if !strings.HasSuffix(newContent, "\n") {
			newContent += "\n"
		}
	} else if strings.HasSuffix(newContent, "\n") && !strings.HasSuffix(originalContent, "\n") {
		// Remove trailing newline if original didn't have one
		newContent = strings.TrimSuffix(newContent, "\n")
	}

	return os.WriteFile(profilePath, []byte(newContent), 0644)
}

