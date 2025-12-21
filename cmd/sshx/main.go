package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/coaxialdolor/sshx/internal/alias"
	"github.com/coaxialdolor/sshx/internal/parser"
	"github.com/coaxialdolor/sshx/internal/protocol"
)

func main() {
	// Handle special commands
	if len(os.Args) >= 2 && os.Args[1] == "uninstall-alias" {
		handleUninstallAlias()
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: sshx [ssh-args...]\n")
		fmt.Fprintf(os.Stderr, "       sshx uninstall-alias\n")
		os.Exit(1)
	}

	// Find the real ssh binary
	sshPath, err := findSSH()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find ssh: %v\n", err)
		os.Exit(1)
	}

	// Extract host, user, and port from arguments
	host, user, port := extractSSHParams(os.Args[1:])

	// Create the ssh command
	sshCmd := exec.Command(sshPath, os.Args[1:]...)
	sshCmd.Stderr = os.Stderr

	stdin, err := sshCmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create stdin pipe: %v\n", err)
		os.Exit(1)
	}

	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		fmt.Fprintf(os.Stderr, "Failed to create stdout pipe: %v\n", err)
		os.Exit(1)
	}

	// Start the ssh process
	if err := sshCmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start ssh: %v\n", err)
		os.Exit(1)
	}

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		sshCmd.Process.Signal(sig)
		os.Exit(0)
	}()

	// Set up terminal for raw input if not on Windows
	var oldState interface{}
	if runtime.GOOS != "windows" {
		// We'll handle raw mode differently - read line by line
		// For simplicity, we'll use a buffered reader
	}

	// Copy stdin to ssh stdin, intercepting download commands
	// Buffer characters until newline to check for commands
	go func() {
		reader := bufio.NewReader(os.Stdin)
		lineBuffer := make([]byte, 0, 256)

		for {
			char, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					// Send any remaining buffer
					if len(lineBuffer) > 0 {
						stdin.Write(lineBuffer)
					}
					stdin.Close()
					return
				}
				fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
				return
			}

			// Add to buffer
			lineBuffer = append(lineBuffer, char)

			// Check if we have a complete line
			if char == '\n' || char == '\r' {
				line := string(lineBuffer[:len(lineBuffer)-1]) // Remove newline
				if char == '\r' {
					// Might be followed by \n, but handle \r\n or just \r
					line = strings.TrimRight(line, "\r\n")
				} else {
					line = strings.TrimRight(line, "\n")
				}

				// Check if this is a download command
				cmd, isDownload := parser.Parse(line)
				if isDownload {
					// Send download request to agent
					req := protocol.DownloadRequest{
						RemotePath:  cmd.RemotePath,
						LocalTarget: cmd.LocalTarget,
						RemoteHost:  host,
						RemoteUser:  user,
						RemotePort:  port,
					}

					if err := protocol.SendDownloadRequest(req); err != nil {
						fmt.Fprintf(os.Stderr, "\r[sshx] Download request failed: %v\n", err)
					} else {
						fmt.Fprintf(os.Stderr, "\r[sshx] Download request sent\n")
					}
					// Don't send the command to the remote
					lineBuffer = lineBuffer[:0]
					continue
				}

				// Send the complete line to ssh
				stdin.Write(lineBuffer)
				lineBuffer = lineBuffer[:0]
			}
			// If not a newline, we continue buffering
			// This means interactive programs won't work perfectly, but commands will be caught
			// For better interactivity, we'd need a more sophisticated approach
		}
	}()

	// Copy ssh stdout to our stdout
	io.Copy(os.Stdout, stdout)

	// Wait for ssh to finish
	sshCmd.Wait()

	if oldState != nil {
		// Restore terminal state if needed
	}
}

// findSSH finds the ssh binary in the system PATH
func findSSH() (string, error) {
	path, err := exec.LookPath("ssh")
	if err != nil {
		return "", fmt.Errorf("ssh not found in PATH: %w", err)
	}
	return path, nil
}

// extractSSHParams extracts host, user, and port from ssh arguments
func extractSSHParams(args []string) (host, user string, port int) {
	user = ""
	host = ""
	port = 22

	for i, arg := range args {
		switch arg {
		case "-l":
			if i+1 < len(args) {
				user = args[i+1]
			}
		case "-p":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &port)
			}
		case "-P":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &port)
			}
		default:
			// Check if it's a user@host format
			if !strings.HasPrefix(arg, "-") && host == "" {
				if strings.Contains(arg, "@") {
					parts := strings.Split(arg, "@")
					if len(parts) == 2 {
						user = parts[0]
						host = parts[1]
					}
				} else {
					host = arg
				}
			}
		}
	}

	return host, user, port
}

// handleUninstallAlias handles the uninstall-alias command
func handleUninstallAlias() {
	shell := alias.DetectShell()
	if shell == alias.ShellUnknown {
		fmt.Fprintf(os.Stderr, "Error: Could not detect your shell\n")
		os.Exit(1)
	}

	profilePath, err := alias.GetProfilePath(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	hasAlias, err := alias.HasAlias(profilePath, shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking profile: %v\n", err)
		os.Exit(1)
	}

	if !hasAlias {
		fmt.Printf("No SSHX alias found in %s\n", profilePath)
		return
	}

	err = alias.RemoveAlias(profilePath, shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing alias: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Alias removed from %s\n", profilePath)
	fmt.Println("Restart your shell or run: source", profilePath)
}

