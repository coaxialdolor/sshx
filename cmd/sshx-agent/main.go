package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/coaxialdolor/sshx/internal/config"
	"github.com/coaxialdolor/sshx/internal/protocol"
	"github.com/coaxialdolor/sshx/internal/util"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Listen on localhost:9234
	listener, err := net.Listen("tcp", protocol.AgentAddress)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", protocol.AgentAddress, err)
	}
	defer listener.Close()

	log.Printf("sshx-agent listening on %s", protocol.AgentAddress)

	// Handle SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		listener.Close()
		os.Exit(0)
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, cfg)
	}
}

func handleConnection(conn net.Conn, cfg *config.Config) {
	defer conn.Close()

	req, err := protocol.ParseDownloadRequest(conn)
	if err != nil {
		log.Printf("Failed to parse request: %v", err)
		protocol.SendResponse(conn, false, fmt.Sprintf("Parse error: %v", err))
		return
	}

	// Use defaults from config if not specified
	if req.RemoteUser == "" {
		req.RemoteUser = cfg.DefaultUser
	}
	if req.RemotePort == 0 {
		req.RemotePort = cfg.DefaultPort
	}

	// Resolve local target path
	localPath, err := util.ResolveLocalPath(req.LocalTarget, cfg)
	if err != nil {
		log.Printf("Failed to resolve local path: %v", err)
		protocol.SendResponse(conn, false, fmt.Sprintf("Path error: %v", err))
		return
	}

	// Create directory if it doesn't exist
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		log.Printf("Failed to create directory: %v", err)
		protocol.SendResponse(conn, false, fmt.Sprintf("Directory error: %v", err))
		return
	}

	// Build scp command
	// scp -P <port> <user>@<host>:<remote-path> <local-path>
	scpArgs := []string{
		"-P",
		fmt.Sprintf("%d", req.RemotePort),
		fmt.Sprintf("%s@%s:%s", req.RemoteUser, req.RemoteHost, req.RemotePath),
		localPath,
	}

	scpCmd := exec.Command("scp", scpArgs...)
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr

	log.Printf("Executing: scp %v", scpArgs)

	if err := scpCmd.Run(); err != nil {
		errorMsg := fmt.Sprintf("scp failed: %v", err)
		log.Printf(errorMsg)
		protocol.SendResponse(conn, false, errorMsg)
		return
	}

	successMsg := fmt.Sprintf("Downloaded %s to %s", req.RemotePath, localPath)
	log.Printf(successMsg)
	protocol.SendResponse(conn, true, successMsg)
}

// Simple response writer for stdout
type responseWriter struct {
	w io.Writer
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}

