package proxy

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/vibium/clicker/internal/bidi"
	"github.com/vibium/clicker/internal/browser"
)

// Default timeout for actionability checks
const defaultTimeout = 30 * time.Second

// BrowserSession represents a browser session connected to a client.
type BrowserSession struct {
	LaunchResult *browser.LaunchResult
	BidiConn     *bidi.Connection
	BidiClient   *bidi.Client
	Client       *ClientConn
	mu           sync.Mutex
	closed       bool
	stopChan     chan struct{}

	// Internal command tracking for vibium: extension commands
	internalCmds   map[int]chan json.RawMessage // id -> response channel
	internalCmdsMu sync.Mutex
	nextInternalID int

	// WebSocket monitoring state
	wsPreloadScriptID string // "" if not installed
	wsSubscribed      bool   // whether script.message is subscribed

	// Download support
	downloadDir string // temp dir for downloads, cleaned up on close
}

// BiDi command structure for parsing incoming messages
type bidiCommand struct {
	ID     int                    `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// BiDi response structure for sending responses (follows WebDriver BiDi spec)
type bidiResponse struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"` // "success" or "error"
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Router manages browser sessions for connected clients.
type Router struct {
	sessions sync.Map // map[uint64]*BrowserSession (client ID -> session)
	headless bool
}

// NewRouter creates a new router.
func NewRouter(headless bool) *Router {
	return &Router{
		headless: headless,
	}
}

// OnClientConnect is called when a new client connects.
// It launches a browser and establishes a BiDi connection.
func (r *Router) OnClientConnect(client *ClientConn) {
	fmt.Printf("[router] Launching browser for client %d...\n", client.ID)

	// Launch browser
	launchResult, err := browser.Launch(browser.LaunchOptions{
		Headless: r.headless,
	})
	if err != nil {
		fmt.Printf("[router] Failed to launch browser for client %d: %v\n", client.ID, err)
		client.Send(fmt.Sprintf(`{"error":{"code":-32000,"message":"Failed to launch browser: %s"}}`, err.Error()))
		client.Close()
		return
	}

	fmt.Printf("[router] Browser launched for client %d, WebSocket: %s\n", client.ID, launchResult.WebSocketURL)

	// Connect to browser BiDi WebSocket
	bidiConn, err := bidi.Connect(launchResult.WebSocketURL)
	if err != nil {
		fmt.Printf("[router] Failed to connect to browser BiDi for client %d: %v\n", client.ID, err)
		launchResult.Close()
		client.Send(fmt.Sprintf(`{"error":{"code":-32000,"message":"Failed to connect to browser: %s"}}`, err.Error()))
		client.Close()
		return
	}

	fmt.Printf("[router] BiDi connection established for client %d\n", client.ID)

	// Create a BiDi client for handling custom commands
	bidiClient := bidi.NewClient(bidiConn)

	session := &BrowserSession{
		LaunchResult:   launchResult,
		BidiConn:       bidiConn,
		BidiClient:     bidiClient,
		Client:         client,
		stopChan:       make(chan struct{}),
		internalCmds:   make(map[int]chan json.RawMessage),
		nextInternalID: 1000000, // Start at high number to avoid collision with client IDs
	}

	r.sessions.Store(client.ID, session)

	// Start routing messages from browser to client
	go r.routeBrowserToClient(session)

	// Subscribe to events for onPage/onPopup, network interception, dialog handling, and downloads.
	// Events are forwarded to the JS client by routeBrowserToClient.
	go func() {
		r.sendInternalCommand(session, "session.subscribe", map[string]interface{}{
			"events": []string{
				"browsingContext.contextCreated",
				"network.beforeRequestSent",
				"network.responseCompleted",
				"browsingContext.userPromptOpened",
				"log.entryAdded",
				"browsingContext.downloadWillBegin",
				"browsingContext.downloadEnd",
			},
		})
		r.setupDownloads(session)
	}()
}

// OnClientMessage is called when a message is received from a client.
// It handles custom vibium: extension commands or forwards to the browser.
func (r *Router) OnClientMessage(client *ClientConn, msg string) {
	sessionVal, ok := r.sessions.Load(client.ID)
	if !ok {
		fmt.Printf("[router] No session for client %d\n", client.ID)
		return
	}

	session := sessionVal.(*BrowserSession)

	session.mu.Lock()
	if session.closed {
		session.mu.Unlock()
		return
	}
	session.mu.Unlock()

	// Parse the command to check for custom vibium: extension methods
	var cmd bidiCommand
	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		// Can't parse, forward as-is
		if err := session.BidiConn.Send(msg); err != nil {
			fmt.Printf("[router] Failed to send to browser for client %d: %v\n", client.ID, err)
		}
		return
	}

	// Handle vibium: extension commands (per WebDriver BiDi spec for extensions)
	switch cmd.Method {
	// Element interaction commands
	case "vibium:click":
		go r.handleVibiumClick(session, cmd)
		return
	case "vibium:dblclick":
		go r.handleVibiumDblclick(session, cmd)
		return
	case "vibium:fill":
		go r.handleVibiumFill(session, cmd)
		return
	case "vibium:type":
		go r.handleVibiumType(session, cmd)
		return
	case "vibium:press":
		go r.handleVibiumPress(session, cmd)
		return
	case "vibium:clear":
		go r.handleVibiumClear(session, cmd)
		return
	case "vibium:check":
		go r.handleVibiumCheck(session, cmd)
		return
	case "vibium:uncheck":
		go r.handleVibiumUncheck(session, cmd)
		return
	case "vibium:selectOption":
		go r.handleVibiumSelectOption(session, cmd)
		return
	case "vibium:hover":
		go r.handleVibiumHover(session, cmd)
		return
	case "vibium:focus":
		go r.handleVibiumFocus(session, cmd)
		return
	case "vibium:dragTo":
		go r.handleVibiumDragTo(session, cmd)
		return
	case "vibium:tap":
		go r.handleVibiumTap(session, cmd)
		return
	case "vibium:scrollIntoView":
		go r.handleVibiumScrollIntoView(session, cmd)
		return
	case "vibium:dispatchEvent":
		go r.handleVibiumDispatchEvent(session, cmd)
		return

	// Element finding commands
	case "vibium:find":
		go r.handleVibiumFind(session, cmd)
		return
	case "vibium:findAll":
		go r.handleVibiumFindAll(session, cmd)
		return

	// Element state commands
	case "vibium:el.text":
		go r.handleVibiumElText(session, cmd)
		return
	case "vibium:el.innerText":
		go r.handleVibiumElInnerText(session, cmd)
		return
	case "vibium:el.html":
		go r.handleVibiumElHTML(session, cmd)
		return
	case "vibium:el.value":
		go r.handleVibiumElValue(session, cmd)
		return
	case "vibium:el.attr":
		go r.handleVibiumElAttr(session, cmd)
		return
	case "vibium:el.bounds":
		go r.handleVibiumElBounds(session, cmd)
		return
	case "vibium:el.isVisible":
		go r.handleVibiumElIsVisible(session, cmd)
		return
	case "vibium:el.isHidden":
		go r.handleVibiumElIsHidden(session, cmd)
		return
	case "vibium:el.isEnabled":
		go r.handleVibiumElIsEnabled(session, cmd)
		return
	case "vibium:el.isChecked":
		go r.handleVibiumElIsChecked(session, cmd)
		return
	case "vibium:el.isEditable":
		go r.handleVibiumElIsEditable(session, cmd)
		return
	case "vibium:el.eval":
		go r.handleVibiumElEval(session, cmd)
		return
	case "vibium:el.screenshot":
		go r.handleVibiumElScreenshot(session, cmd)
		return
	case "vibium:el.waitFor":
		go r.handleVibiumElWaitFor(session, cmd)
		return

	// Page-level input commands
	case "vibium:keyboard.press":
		go r.handleKeyboardPress(session, cmd)
		return
	case "vibium:keyboard.down":
		go r.handleKeyboardDown(session, cmd)
		return
	case "vibium:keyboard.up":
		go r.handleKeyboardUp(session, cmd)
		return
	case "vibium:keyboard.type":
		go r.handleKeyboardType(session, cmd)
		return
	case "vibium:mouse.click":
		go r.handleMouseClick(session, cmd)
		return
	case "vibium:mouse.move":
		go r.handleMouseMove(session, cmd)
		return
	case "vibium:mouse.down":
		go r.handleMouseDown(session, cmd)
		return
	case "vibium:mouse.up":
		go r.handleMouseUp(session, cmd)
		return
	case "vibium:mouse.wheel":
		go r.handleMouseWheel(session, cmd)
		return
	case "vibium:touch.tap":
		go r.handleTouchTap(session, cmd)
		return

	// Page-level capture commands
	case "vibium:page.screenshot":
		go r.handlePageScreenshot(session, cmd)
		return
	case "vibium:page.pdf":
		go r.handlePagePDF(session, cmd)
		return

	// Page-level evaluation commands
	case "vibium:page.eval":
		go r.handlePageEval(session, cmd)
		return
	case "vibium:page.evalHandle":
		go r.handlePageEvalHandle(session, cmd)
		return
	case "vibium:page.addScript":
		go r.handlePageAddScript(session, cmd)
		return
	case "vibium:page.addStyle":
		go r.handlePageAddStyle(session, cmd)
		return
	case "vibium:page.expose":
		go r.handlePageExpose(session, cmd)
		return

	// Page-level waiting commands
	case "vibium:page.waitFor":
		go r.handlePageWaitFor(session, cmd)
		return
	case "vibium:page.wait":
		go r.handlePageWait(session, cmd)
		return
	case "vibium:page.waitForFunction":
		go r.handlePageWaitForFunction(session, cmd)
		return

	// Navigation commands
	case "vibium:page.navigate":
		go r.handlePageNavigate(session, cmd)
		return
	case "vibium:page.back":
		go r.handlePageBack(session, cmd)
		return
	case "vibium:page.forward":
		go r.handlePageForward(session, cmd)
		return
	case "vibium:page.reload":
		go r.handlePageReload(session, cmd)
		return
	case "vibium:page.url":
		go r.handlePageURL(session, cmd)
		return
	case "vibium:page.title":
		go r.handlePageTitle(session, cmd)
		return
	case "vibium:page.content":
		go r.handlePageContent(session, cmd)
		return
	case "vibium:page.waitForURL":
		go r.handlePageWaitForURL(session, cmd)
		return
	case "vibium:page.waitForLoad":
		go r.handlePageWaitForLoad(session, cmd)
		return

	// Page & context lifecycle commands
	case "vibium:browser.page":
		go r.handleBrowserPage(session, cmd)
		return
	case "vibium:browser.newPage":
		go r.handleBrowserNewPage(session, cmd)
		return
	case "vibium:browser.newContext":
		go r.handleBrowserNewContext(session, cmd)
		return
	case "vibium:context.newPage":
		go r.handleContextNewPage(session, cmd)
		return
	case "vibium:browser.pages":
		go r.handleBrowserPages(session, cmd)
		return
	case "vibium:context.close":
		go r.handleContextClose(session, cmd)
		return

	// Cookie & storage commands
	case "vibium:context.cookies":
		go r.handleContextCookies(session, cmd)
		return
	case "vibium:context.setCookies":
		go r.handleContextSetCookies(session, cmd)
		return
	case "vibium:context.clearCookies":
		go r.handleContextClearCookies(session, cmd)
		return
	case "vibium:context.storageState":
		go r.handleContextStorageState(session, cmd)
		return
	case "vibium:context.addInitScript":
		go r.handleContextAddInitScript(session, cmd)
		return

	// Frame commands
	case "vibium:page.frames":
		go r.handlePageFrames(session, cmd)
		return
	case "vibium:page.frame":
		go r.handlePageFrame(session, cmd)
		return

	// Emulation commands
	case "vibium:page.setViewport":
		go r.handlePageSetViewport(session, cmd)
		return
	case "vibium:page.viewport":
		go r.handlePageViewport(session, cmd)
		return
	case "vibium:page.emulateMedia":
		go r.handlePageEmulateMedia(session, cmd)
		return
	case "vibium:page.setContent":
		go r.handlePageSetContent(session, cmd)
		return
	case "vibium:page.setGeolocation":
		go r.handlePageSetGeolocation(session, cmd)
		return

	// Accessibility commands
	case "vibium:page.a11yTree":
		go r.handleVibiumPageA11yTree(session, cmd)
		return
	case "vibium:el.role":
		go r.handleVibiumElRole(session, cmd)
		return
	case "vibium:el.label":
		go r.handleVibiumElLabel(session, cmd)
		return

	case "vibium:browser.close":
		go r.handleBrowserClose(session, cmd)
		return
	case "vibium:page.activate":
		go r.handlePageActivate(session, cmd)
		return
	case "vibium:page.close":
		go r.handlePageClose(session, cmd)
		return

	// Network interception commands
	case "vibium:page.route":
		go r.handlePageRoute(session, cmd)
		return
	case "vibium:page.unroute":
		go r.handlePageUnroute(session, cmd)
		return
	case "vibium:network.continue":
		go r.handleNetworkContinue(session, cmd)
		return
	case "vibium:network.fulfill":
		go r.handleNetworkFulfill(session, cmd)
		return
	case "vibium:network.abort":
		go r.handleNetworkAbort(session, cmd)
		return
	case "vibium:page.setHeaders":
		go r.handlePageSetHeaders(session, cmd)
		return

	// Dialog commands
	case "vibium:dialog.accept":
		go r.handleDialogAccept(session, cmd)
		return
	case "vibium:dialog.dismiss":
		go r.handleDialogDismiss(session, cmd)
		return

	// WebSocket monitoring
	case "vibium:page.onWebSocket":
		go r.handlePageOnWebSocket(session, cmd)
		return

	// Download & file commands
	case "vibium:download.saveAs":
		go r.handleDownloadSaveAs(session, cmd)
		return
	case "vibium:el.setFiles":
		go r.handleVibiumElSetFiles(session, cmd)
		return
	}

	// Forward standard BiDi commands to browser
	if err := session.BidiConn.Send(msg); err != nil {
		fmt.Printf("[router] Failed to send to browser for client %d: %v\n", client.ID, err)
	}
}

// getContext retrieves the first browsing context.
func (r *Router) getContext(session *BrowserSession) (string, error) {
	resp, err := r.sendInternalCommand(session, "browsingContext.getTree", map[string]interface{}{})
	if err != nil {
		return "", err
	}

	var result struct {
		Result struct {
			Contexts []struct {
				Context string `json:"context"`
			} `json:"contexts"`
		} `json:"result"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse getTree response: %w", err)
	}
	if len(result.Result.Contexts) == 0 {
		return "", fmt.Errorf("no browsing contexts available")
	}
	return result.Result.Contexts[0].Context, nil
}

// sendSuccess sends a successful response to the client.
func (r *Router) sendSuccess(session *BrowserSession, id int, result interface{}) {
	resp := bidiResponse{ID: id, Type: "success", Result: result}
	data, _ := json.Marshal(resp)
	session.Client.Send(string(data))
}

// sendError sends an error response to the client (follows WebDriver BiDi spec).
func (r *Router) sendError(session *BrowserSession, id int, err error) {
	resp := bidiResponse{
		ID:      id,
		Type:    "error",
		Error:   "timeout",
		Message: err.Error(),
	}
	data, _ := json.Marshal(resp)
	session.Client.Send(string(data))
}

// OnClientDisconnect is called when a client disconnects.
// It closes the browser session.
func (r *Router) OnClientDisconnect(client *ClientConn) {
	sessionVal, ok := r.sessions.LoadAndDelete(client.ID)
	if !ok {
		return
	}

	session := sessionVal.(*BrowserSession)
	r.closeSession(session)
}

// routeBrowserToClient reads messages from the browser and forwards them to the client.
func (r *Router) routeBrowserToClient(session *BrowserSession) {
	for {
		select {
		case <-session.stopChan:
			return
		default:
		}

		msg, err := session.BidiConn.Receive()
		if err != nil {
			session.mu.Lock()
			closed := session.closed
			session.mu.Unlock()

			if !closed {
				fmt.Printf("[router] Browser connection closed for client %d: %v\n", session.Client.ID, err)
				// Browser died, close the client
				session.Client.Close()
			}
			return
		}

		// Check if this is a response to an internal command
		var resp struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal([]byte(msg), &resp); err == nil && resp.ID > 0 {
			session.internalCmdsMu.Lock()
			ch, isInternal := session.internalCmds[resp.ID]
			session.internalCmdsMu.Unlock()

			if isInternal {
				// Route to internal handler
				ch <- json.RawMessage(msg)
				continue
			}
		}

		// Check for WebSocket channel events (intercept, don't forward raw script.message)
		if r.isWsChannelEvent(session, msg) {
			continue
		}

		// Forward message to client
		if err := session.Client.Send(msg); err != nil {
			fmt.Printf("[router] Failed to send to client %d: %v\n", session.Client.ID, err)
			return
		}
	}
}

// sendInternalCommand sends a BiDi command and waits for the response.
func (r *Router) sendInternalCommand(session *BrowserSession, method string, params map[string]interface{}) (json.RawMessage, error) {
	session.internalCmdsMu.Lock()
	id := session.nextInternalID
	session.nextInternalID++
	ch := make(chan json.RawMessage, 1)
	session.internalCmds[id] = ch
	session.internalCmdsMu.Unlock()

	defer func() {
		session.internalCmdsMu.Lock()
		delete(session.internalCmds, id)
		session.internalCmdsMu.Unlock()
	}()

	// Send the command
	cmd := map[string]interface{}{
		"id":     id,
		"method": method,
		"params": params,
	}
	cmdBytes, _ := json.Marshal(cmd)
	if err := session.BidiConn.Send(string(cmdBytes)); err != nil {
		return nil, err
	}

	// Wait for response (with timeout)
	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response to %s", method)
	case <-session.stopChan:
		return nil, fmt.Errorf("session closed")
	}
}

// closeSession closes a browser session and cleans up resources.
func (r *Router) closeSession(session *BrowserSession) {
	session.mu.Lock()
	if session.closed {
		session.mu.Unlock()
		return
	}
	session.closed = true
	session.mu.Unlock()

	fmt.Printf("[router] Closing browser session for client %d\n", session.Client.ID)

	// Signal the routing goroutine to stop
	close(session.stopChan)

	// Close BiDi connection
	if session.BidiConn != nil {
		session.BidiConn.Close()
	}

	// Clean up download temp dir
	if session.downloadDir != "" {
		os.RemoveAll(session.downloadDir)
	}

	// Close browser
	if session.LaunchResult != nil {
		session.LaunchResult.Close()
	}

	fmt.Printf("[router] Browser session closed for client %d\n", session.Client.ID)
}

// CloseAll closes all browser sessions.
func (r *Router) CloseAll() {
	r.sessions.Range(func(key, value interface{}) bool {
		session := value.(*BrowserSession)
		r.closeSession(session)
		r.sessions.Delete(key)
		return true
	})
}
