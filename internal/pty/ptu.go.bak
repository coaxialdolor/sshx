package pty

import (
	"io"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

// PTY represents a pseudo-terminal
type PTY struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

// NewPTY creates a new PTY and starts the command
func NewPTY(command string, args ...string) (*PTY, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, we'll use a simpler approach with exec.Command
		// Full PTY support on Windows is complex, so we'll use basic stdin/stdout
		cmd = exec.Command(command, args...)
	} else {
		cmd = exec.Command(command, args...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	// Set stdin to raw mode if not on Windows
	if runtime.GOOS != "windows" {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			stdin.Close()
			stdout.Close()
			stderr.Close()
			return nil, err
		}

		// Restore terminal state on exit
		defer func() {
			if err != nil {
				term.Restore(int(os.Stdin.Fd()), oldState)
			}
		}()
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, err
	}

	pty := &PTY{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	return pty, nil
}

// Read reads from the PTY's stdout
func (p *PTY) Read(data []byte) (int, error) {
	return p.stdout.Read(data)
}

// Write writes to the PTY's stdin
func (p *PTY) Write(data []byte) (int, error) {
	return p.stdin.Write(data)
}

// Close closes the PTY and cleans up
func (p *PTY) Close() error {
	p.stdin.Close()
	p.stdout.Close()
	p.stderr.Close()
	return p.cmd.Wait()
}

// Wait waits for the command to finish
func (p *PTY) Wait() error {
	return p.cmd.Wait()
}

// SetRawMode sets the terminal to raw mode
func SetRawMode() (*term.State, error) {
	if runtime.GOOS == "windows" {
		// Windows doesn't support raw mode the same way
		// We'll use a simpler approach
		return nil, nil
	}
	return term.MakeRaw(int(os.Stdin.Fd()))
}

// RestoreMode restores the terminal to its previous state
func RestoreMode(state *term.State) error {
	if runtime.GOOS == "windows" || state == nil {
		return nil
	}
	return term.Restore(int(os.Stdin.Fd()), state)
}

// GetSize returns the size of the terminal
func GetSize() (width, height int, err error) {
	if runtime.GOOS == "windows" {
		// Windows implementation
		return term.GetSize(int(os.Stdin.Fd()))
	}
	return term.GetSize(int(os.Stdin.Fd()))
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
