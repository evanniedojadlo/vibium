package proxy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vibium/clicker/internal/bidi"
)

// Session abstracts BiDi communication so that both the proxy (WebSocket router)
// and the MCP server (direct bidi.Client) can share the same browser-automation
// logic. All shared standalone functions (Navigate, Click, etc.) accept a Session.
type Session interface {
	// SendBidiCommand sends a BiDi command and returns the full response JSON.
	// The response format matches the proxy's sendInternalCommand output:
	//   {"result": { ... }}   (success)
	//   {"type":"error", "error":"...", "message":"..."}  (error)
	SendBidiCommand(method string, params map[string]interface{}) (json.RawMessage, error)

	// SendBidiCommandWithTimeout is like SendBidiCommand but with a custom timeout.
	SendBidiCommandWithTimeout(method string, params map[string]interface{}, timeout time.Duration) (json.RawMessage, error)

	// GetContextID returns a browsing context ID. If the session tracks a
	// "current" context it returns that; otherwise it fetches the first context
	// from browsingContext.getTree.
	GetContextID() (string, error)
}

// ---------------------------------------------------------------------------
// ProxySession — adapts Router + BrowserSession to Session.
// ---------------------------------------------------------------------------

// ProxySession wraps a Router and BrowserSession pair so that shared
// standalone functions can call sendInternalCommand through the Session
// interface.
type ProxySession struct {
	Router  *Router
	Session *BrowserSession
	Context string // optional explicit context override
}

// NewProxySession creates a ProxySession.
func NewProxySession(r *Router, s *BrowserSession, context string) *ProxySession {
	return &ProxySession{Router: r, Session: s, Context: context}
}

func (p *ProxySession) SendBidiCommand(method string, params map[string]interface{}) (json.RawMessage, error) {
	return p.Router.sendInternalCommand(p.Session, method, params)
}

func (p *ProxySession) SendBidiCommandWithTimeout(method string, params map[string]interface{}, timeout time.Duration) (json.RawMessage, error) {
	return p.Router.sendInternalCommandWithTimeout(p.Session, method, params, timeout)
}

func (p *ProxySession) GetContextID() (string, error) {
	if p.Context != "" {
		return p.Context, nil
	}
	return p.Router.getContext(p.Session)
}

// ---------------------------------------------------------------------------
// MCPSession — adapts *bidi.Client to Session.
// ---------------------------------------------------------------------------

// MCPSession wraps a bidi.Client so that shared standalone functions can send
// BiDi commands through the Session interface. The bidi.Client already handles
// error responses as Go errors, so checkBidiError on wrapped responses is a
// safe no-op.
type MCPSession struct {
	Client *bidi.Client
}

// NewMCPSession creates an MCPSession.
func NewMCPSession(client *bidi.Client) *MCPSession {
	return &MCPSession{Client: client}
}

func (m *MCPSession) SendBidiCommand(method string, params map[string]interface{}) (json.RawMessage, error) {
	msg, err := m.Client.SendCommand(method, params)
	if err != nil {
		return nil, err
	}
	// Wrap msg.Result as {"result": <msg.Result>} to match the proxy response
	// format that parseScriptResult, checkBidiError, etc. expect.
	wrapped, err := json.Marshal(map[string]json.RawMessage{"result": msg.Result})
	if err != nil {
		return nil, fmt.Errorf("failed to wrap bidi result: %w", err)
	}
	return wrapped, nil
}

func (m *MCPSession) SendBidiCommandWithTimeout(method string, params map[string]interface{}, timeout time.Duration) (json.RawMessage, error) {
	// bidi.Client.SendCommand blocks until a response arrives; the timeout is
	// handled at the caller level (polling loops). Delegate to the regular path.
	return m.SendBidiCommand(method, params)
}

func (m *MCPSession) GetContextID() (string, error) {
	tree, err := m.Client.GetTree()
	if err != nil {
		return "", fmt.Errorf("failed to get browsing context: %w", err)
	}
	if len(tree.Contexts) == 0 {
		return "", fmt.Errorf("no browsing contexts available")
	}
	return tree.Contexts[0].Context, nil
}
