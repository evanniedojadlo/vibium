package proxy

import (
	"encoding/json"
	"fmt"
)

// resolveContext extracts the "context" param or returns the first context from getTree.
func (r *Router) resolveContext(session *BrowserSession, params map[string]interface{}) (string, error) {
	if ctx, ok := params["context"].(string); ok && ctx != "" {
		return ctx, nil
	}
	return r.getContext(session)
}

// evalSimpleScript runs a no-argument script.callFunction and returns the string result.
// The fn should be a JS function declaration that returns a string, e.g. "() => document.title".
func (r *Router) evalSimpleScript(session *BrowserSession, context, fn string) (string, error) {
	params := map[string]interface{}{
		"functionDeclaration": fn,
		"target":              map[string]interface{}{"context": context},
		"arguments":           []map[string]interface{}{},
		"awaitPromise":        false,
		"resultOwnership":     "root",
	}

	resp, err := r.sendInternalCommand(session, "script.callFunction", params)
	if err != nil {
		return "", err
	}

	return parseScriptResult(resp)
}

// checkBidiError checks if a BiDi response is an error and returns it.
// BiDi error responses have: { "type": "error", "error": "...", "message": "..." }
func checkBidiError(resp json.RawMessage) error {
	var errResp struct {
		Type    string `json:"type"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp, &errResp); err != nil {
		return nil // Can't parse, assume not an error
	}
	if errResp.Type == "error" {
		return fmt.Errorf("%s: %s", errResp.Error, errResp.Message)
	}
	return nil
}

// parseScriptResult parses a BiDi script.callFunction response and returns the string value.
// Expected structure: { "result": { "result": { "type": "string", "value": "..." } } }
func parseScriptResult(resp json.RawMessage) (string, error) {
	var result struct {
		Result struct {
			Result struct {
				Type  string `json:"type"`
				Value string `json:"value,omitempty"`
			} `json:"result"`
		} `json:"result"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse script result: %w", err)
	}

	if result.Result.Result.Type == "null" || result.Result.Result.Type == "undefined" {
		return "", fmt.Errorf("script returned %s", result.Result.Result.Type)
	}

	return result.Result.Result.Value, nil
}
