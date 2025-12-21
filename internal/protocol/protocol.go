package protocol

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const AgentPort = "9234"
const AgentAddress = "127.0.0.1:" + AgentPort

// DownloadRequest represents a download request
type DownloadRequest struct {
	RemotePath string
	LocalTarget string
	RemoteHost string
	RemoteUser string
	RemotePort int
}

// SendDownloadRequest sends a download request to the agent
func SendDownloadRequest(req DownloadRequest) error {
	conn, err := net.Dial("tcp", AgentAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	fmt.Fprintf(writer, "DOWNLOAD %s %s\n", req.RemotePath, req.LocalTarget)
	fmt.Fprintf(writer, "REMOTE_HOST %s\n", req.RemoteHost)
	fmt.Fprintf(writer, "REMOTE_USER %s\n", req.RemoteUser)
	fmt.Fprintf(writer, "REMOTE_PORT %d\n", req.RemotePort)
	fmt.Fprintf(writer, "END\n")

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	return nil
}

// ParseDownloadRequest parses a download request from the network
func ParseDownloadRequest(conn net.Conn) (*DownloadRequest, error) {
	req := &DownloadRequest{}
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "END" {
			break
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "DOWNLOAD":
			if len(parts) >= 2 {
				req.RemotePath = parts[1]
				if len(parts) >= 3 {
					req.LocalTarget = parts[2]
				}
			}
		case "REMOTE_HOST":
			if len(parts) >= 2 {
				req.RemoteHost = parts[1]
			}
		case "REMOTE_USER":
			if len(parts) >= 2 {
				req.RemoteUser = parts[1]
			}
		case "REMOTE_PORT":
			if len(parts) >= 2 {
				fmt.Sscanf(parts[1], "%d", &req.RemotePort)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return req, nil
}

// SendResponse sends a response back to the client
func SendResponse(conn net.Conn, success bool, message string) error {
	status := "FAIL"
	if success {
		status = "OK"
	}
	_, err := fmt.Fprintf(conn, "%s %s\n", status, message)
	return err
}

