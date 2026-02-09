package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/vibium/clicker/internal/log"
	"github.com/vibium/clicker/internal/mcp"
)

// StatusResult is returned by daemon/status.
type StatusResult struct {
	Version   string `json:"version"`
	PID       int    `json:"pid"`
	Uptime    string `json:"uptime"`
	Socket    string `json:"socket"`
	StartTime string `json:"startTime"`
}

// handleConnection processes a single client connection.
// Each connection sends one JSON-RPC request and receives one response.
func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	d.touchActivity()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	scanner := bufio.NewScanner(conn)
	// Increase scanner buffer for large requests
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	if !scanner.Scan() {
		return
	}
	line := scanner.Bytes()

	response := d.handleRequest(line)
	if response == nil {
		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Debug("marshal error", "error", err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	fmt.Fprintf(conn, "%s\n", data)
}

// handleRequest parses and routes a JSON-RPC request.
func (d *Daemon) handleRequest(data []byte) *mcp.Response {
	var req mcp.Request
	if err := json.Unmarshal(data, &req); err != nil {
		return &mcp.Response{
			JSONRPC: "2.0",
			Error: &mcp.Error{
				Code:    mcp.ParseError,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
	}

	if req.JSONRPC != "2.0" {
		return &mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &mcp.Error{
				Code:    mcp.InvalidRequest,
				Message: "Invalid Request",
				Data:    "jsonrpc must be '2.0'",
			},
		}
	}

	result, mcpErr := d.route(req)

	if req.ID == nil {
		return nil
	}

	if mcpErr != nil {
		return &mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   mcpErr,
		}
	}

	return &mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// route dispatches requests to the appropriate handler.
func (d *Daemon) route(req mcp.Request) (interface{}, *mcp.Error) {
	log.Debug("daemon request", "method", req.Method, "id", req.ID)

	switch req.Method {
	case "daemon/status":
		return d.handleStatus()
	case "daemon/shutdown":
		go d.Shutdown() // Shutdown asynchronously so we can send response
		return map[string]string{"status": "shutting down"}, nil
	case "tools/call":
		return d.handleToolsCall(req.Params)
	case "tools/list":
		return mcp.ToolsListResult{
			Tools: mcp.GetToolSchemas(),
		}, nil
	case "initialize":
		return d.handleInitialize()
	case "initialized", "notifications/initialized":
		return nil, nil
	default:
		return nil, &mcp.Error{
			Code:    mcp.MethodNotFound,
			Message: "Method not found",
			Data:    req.Method,
		}
	}
}

// handleStatus returns daemon status information.
func (d *Daemon) handleStatus() (interface{}, *mcp.Error) {
	return StatusResult{
		Version:   d.version,
		PID:       pidSelf(),
		Uptime:    time.Since(d.startTime).Truncate(time.Second).String(),
		Socket:    d.socketPath,
		StartTime: d.startTime.Format(time.RFC3339),
	}, nil
}

// handleInitialize handles the MCP initialize request.
func (d *Daemon) handleInitialize() (interface{}, *mcp.Error) {
	return mcp.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: mcp.ServerCapabilities{
			Tools: &mcp.ToolsCapability{},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    "vibium",
			Version: d.version,
		},
	}, nil
}

// handleToolsCall executes a tool and returns the result.
func (d *Daemon) handleToolsCall(params json.RawMessage) (interface{}, *mcp.Error) {
	var p mcp.ToolsCallParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &mcp.Error{
			Code:    mcp.InvalidParams,
			Message: "Invalid params",
			Data:    err.Error(),
		}
	}

	// Serialize handler access — handlers are not thread-safe
	d.mu.Lock()
	result, err := d.handlers.Call(p.Name, p.Arguments)
	d.mu.Unlock()

	if err != nil {
		return mcp.ToolsCallResult{
			Content: []mcp.Content{{Type: "text", Text: err.Error()}},
			IsError: true,
		}, nil
	}

	return result, nil
}

func pidSelf() int {
	return os.Getpid()
}
