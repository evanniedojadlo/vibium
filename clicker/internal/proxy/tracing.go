package proxy

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// TracingStartOptions configures how tracing behaves.
type TracingStartOptions struct {
	Name        string `json:"name"`
	Screenshots bool   `json:"screenshots"`
	Snapshots   bool   `json:"snapshots"`
	Sources     bool   `json:"sources"`
	Title       string `json:"title"`
	Bidi        bool   `json:"bidi"`
}

// traceEvent is a generic trace event stored as a JSON-friendly map.
type traceEvent = map[string]interface{}

// TraceRecorder manages trace recording state for a browser session.
// It collects events, screenshots, and DOM snapshots, then packages
// them into a Playwright-compatible trace zip.
type TraceRecorder struct {
	mu            sync.Mutex
	recording     bool
	options       TracingStartOptions
	events        []traceEvent      // current chunk's trace events
	network       []traceEvent      // current chunk's network events
	resources     map[string][]byte // sha1 hex -> binary data (PNG/HTML)
	groupStack    []string          // nested group names
	chunkIndex    int
	startTime     int64 // unix ms
	actionCounter int   // monotonic counter for action/bidi callIds

	// Screenshot goroutine control
	screenshotStop chan struct{}
	screenshotWg   sync.WaitGroup
}

// NewTraceRecorder creates a new trace recorder.
func NewTraceRecorder() *TraceRecorder {
	return &TraceRecorder{
		resources: make(map[string][]byte),
	}
}

// IsRecording returns whether tracing is currently active.
func (t *TraceRecorder) IsRecording() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.recording
}

// Start begins trace recording with the given options.
func (t *TraceRecorder) Start(opts TracingStartOptions) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.recording = true
	t.options = opts
	t.events = nil
	t.network = nil
	t.resources = make(map[string][]byte)
	t.groupStack = nil
	t.chunkIndex = 0
	t.startTime = time.Now().UnixMilli()

	title := opts.Title
	if title == "" {
		title = opts.Name
	}

	// First event must be context-options (required by Playwright trace viewer)
	t.events = append(t.events, traceEvent{
		"type":        "context-options",
		"browserName": "chromium",
		"platform":    runtime.GOOS,
		"wallTime":    float64(t.startTime),
		"title":       title,
		"options":     map[string]interface{}{},
		"sdkLanguage": "javascript",
	})
}

// Stop stops recording and returns the trace zip data.
func (t *TraceRecorder) Stop() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return nil, fmt.Errorf("tracing is not started")
	}

	t.recording = false
	return t.buildZipLocked()
}

// StartChunk starts a new chunk within the current trace.
func (t *TraceRecorder) StartChunk(name, title string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.events = nil
	t.network = nil
	t.chunkIndex++

	chunkTitle := title
	if chunkTitle == "" {
		chunkTitle = name
	}

	t.events = append(t.events, traceEvent{
		"type":        "context-options",
		"browserName": "chromium",
		"platform":    runtime.GOOS,
		"wallTime":    float64(time.Now().UnixMilli()),
		"title":       chunkTitle,
		"options":     map[string]interface{}{},
		"sdkLanguage": "javascript",
	})
}

// StopChunk packages the current chunk into a zip and returns it.
// Tracing remains active for additional chunks.
func (t *TraceRecorder) StopChunk() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return nil, fmt.Errorf("tracing is not started")
	}

	return t.buildZipLocked()
}

// StartGroup adds a group-start marker to the trace.
func (t *TraceRecorder) StartGroup(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.groupStack = append(t.groupStack, name)
	now := float64(time.Now().UnixMilli())
	t.events = append(t.events, traceEvent{
		"type":      "before",
		"callId":    fmt.Sprintf("group@%d", len(t.events)),
		"apiName":   name,
		"class":     "Tracing",
		"method":    "group",
		"params":    map[string]interface{}{"name": name},
		"wallTime":  now,
		"startTime": now,
	})
}

// StopGroup adds a group-end marker to the trace.
func (t *TraceRecorder) StopGroup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.groupStack) > 0 {
		t.groupStack = t.groupStack[:len(t.groupStack)-1]
	}

	t.events = append(t.events, traceEvent{
		"type":    "after",
		"callId":  fmt.Sprintf("group-end@%d", len(t.events)),
		"endTime": float64(time.Now().UnixMilli()),
	})
}

// Options returns the current tracing options.
func (t *TraceRecorder) Options() TracingStartOptions {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.options
}

// apiNameFromMethod maps a vibium: method to (class, apiName) for trace display.
func apiNameFromMethod(method string) (string, string) {
	// Strip the "vibium:" prefix
	if len(method) <= 7 || method[:7] != "vibium:" {
		return "Vibium", method
	}
	name := method[7:] // e.g. "click", "page.navigate", "el.text"

	switch {
	// Element interaction: click, dblclick, fill, type, press, clear, check, uncheck, selectOption, hover, focus, dragTo, tap, scrollIntoView, dispatchEvent
	case name == "click" || name == "dblclick" || name == "fill" || name == "type" ||
		name == "press" || name == "clear" || name == "check" || name == "uncheck" ||
		name == "selectOption" || name == "hover" || name == "focus" || name == "dragTo" ||
		name == "tap" || name == "scrollIntoView" || name == "dispatchEvent":
		return "Element", "Element." + name

	// Element finding: find, findAll
	case name == "find" || name == "findAll":
		return "Page", "Page." + name

	// Element state: el.*
	case len(name) > 3 && name[:3] == "el.":
		return "Element", "Element." + name[3:]

	// Page commands: page.*
	case len(name) > 5 && name[:5] == "page.":
		return "Page", "Page." + name[5:]

	// Browser commands: browser.*
	case len(name) > 8 && name[:8] == "browser.":
		return "Browser", "Browser." + name[8:]

	// Context commands: context.*
	case len(name) > 8 && name[:8] == "context.":
		return "BrowserContext", "BrowserContext." + name[8:]

	// Keyboard: keyboard.*
	case len(name) > 9 && name[:9] == "keyboard.":
		return "Page", "Page." + name

	// Mouse: mouse.*
	case len(name) > 6 && name[:6] == "mouse.":
		return "Page", "Page." + name

	// Touch: touch.*
	case len(name) > 6 && name[:6] == "touch.":
		return "Page", "Page." + name

	// Network: network.*
	case len(name) > 8 && name[:8] == "network.":
		return "Network", "Network." + name[8:]

	// Dialog: dialog.*
	case len(name) > 7 && name[:7] == "dialog.":
		return "Dialog", "Dialog." + name[7:]

	// Clock: clock.*
	case len(name) > 6 && name[:6] == "clock.":
		return "Clock", "Clock." + name[6:]

	// Download: download.*
	case len(name) > 9 && name[:9] == "download.":
		return "Download", "Download." + name[9:]

	default:
		return "Vibium", name
	}
}

// RecordAction records a vibium command as an action marker in the trace.
func (t *TraceRecorder) RecordAction(method string, params map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	class, apiName := apiNameFromMethod(method)
	t.actionCounter++
	now := float64(time.Now().UnixMilli())
	t.events = append(t.events, traceEvent{
		"type":      "before",
		"callId":    fmt.Sprintf("action@%d", t.actionCounter),
		"apiName":   apiName,
		"class":     class,
		"method":    method,
		"params":    params,
		"wallTime":  now,
		"startTime": now,
	})
}

// RecordActionEnd records the end of a vibium command action in the trace.
func (t *TraceRecorder) RecordActionEnd() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	t.events = append(t.events, traceEvent{
		"type":    "after",
		"callId":  fmt.Sprintf("action-end@%d", t.actionCounter),
		"endTime": float64(time.Now().UnixMilli()),
	})
}

// RecordBidiCommand records a raw BiDi command sent to the browser (opt-in via bidi: true).
func (t *TraceRecorder) RecordBidiCommand(method string, params map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	t.actionCounter++
	now := float64(time.Now().UnixMilli())
	t.events = append(t.events, traceEvent{
		"type":      "before",
		"callId":    fmt.Sprintf("bidi@%d", t.actionCounter),
		"apiName":   method,
		"class":     "BiDi",
		"method":    method,
		"params":    params,
		"wallTime":  now,
		"startTime": now,
	})
}

// RecordBidiCommandEnd records the end of a BiDi command in the trace.
func (t *TraceRecorder) RecordBidiCommandEnd() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	t.events = append(t.events, traceEvent{
		"type":    "after",
		"callId":  fmt.Sprintf("bidi-end@%d", t.actionCounter),
		"endTime": float64(time.Now().UnixMilli()),
	})
}

// AddScreenshot stores a screenshot PNG and adds a screencast-frame event.
func (t *TraceRecorder) AddScreenshot(pngData []byte, pageID string, width, height int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	hash := sha1Hex(pngData)
	t.resources[hash] = pngData
	t.events = append(t.events, traceEvent{
		"type":      "screencast-frame",
		"pageId":    pageID,
		"sha1":      hash,
		"width":     width,
		"height":    height,
		"timestamp": float64(time.Now().UnixMilli()),
	})
}

// AddSnapshot stores a DOM snapshot HTML and adds a frame-snapshot event.
func (t *TraceRecorder) AddSnapshot(htmlData []byte, pageID, frameURL string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	hash := sha1Hex(htmlData)
	t.resources[hash] = htmlData
	now := float64(time.Now().UnixMilli())
	t.events = append(t.events, traceEvent{
		"type":      "frame-snapshot",
		"pageId":    pageID,
		"sha1":      hash,
		"frameUrl":  frameURL,
		"title":     "",
		"timestamp": now,
		"wallTime":  now,
		"callId":    fmt.Sprintf("snapshot@%d", len(t.events)),
	})
}

// RecordBidiEvent records a raw BiDi event from the browser into the trace.
func (t *TraceRecorder) RecordBidiEvent(msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.recording {
		return
	}

	var bidiEvent struct {
		Method string                 `json:"method"`
		Params map[string]interface{} `json:"params"`
	}
	if err := json.Unmarshal([]byte(msg), &bidiEvent); err != nil {
		return
	}

	// Only record events (not responses)
	if bidiEvent.Method == "" {
		return
	}

	now := float64(time.Now().UnixMilli())

	switch bidiEvent.Method {
	case "network.beforeRequestSent", "network.responseCompleted", "network.fetchError":
		t.network = append(t.network, traceEvent{
			"type":      "resource-snapshot",
			"method":    bidiEvent.Method,
			"params":    bidiEvent.Params,
			"timestamp": now,
		})
	default:
		t.events = append(t.events, traceEvent{
			"type":   "event",
			"method": bidiEvent.Method,
			"params": bidiEvent.Params,
			"time":   now,
			"class":  "BrowserContext",
		})
	}
}

// StopScreenshots signals the screenshot goroutine to stop and waits for it.
func (t *TraceRecorder) StopScreenshots() {
	t.mu.Lock()
	ch := t.screenshotStop
	t.screenshotStop = nil
	t.mu.Unlock()

	if ch != nil {
		close(ch)
		t.screenshotWg.Wait()
	}
}

// StartScreenshotLoop starts a background goroutine that captures screenshots periodically.
// captureFunc should return (base64-encoded PNG, pageID, error).
func (t *TraceRecorder) StartScreenshotLoop(captureFunc func() (string, string, error)) {
	t.mu.Lock()
	t.screenshotStop = make(chan struct{})
	stopCh := t.screenshotStop
	t.mu.Unlock()

	t.screenshotWg.Add(1)
	go func() {
		defer t.screenshotWg.Done()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				b64Data, pageID, err := captureFunc()
				if err != nil || b64Data == "" {
					continue
				}

				pngData, err := decodeBase64(b64Data)
				if err != nil {
					continue
				}

				w, h := pngDimensions(pngData)
				t.AddScreenshot(pngData, pageID, w, h)
			}
		}
	}()
}

// buildZipLocked creates the Playwright-compatible trace zip.
// Must be called with t.mu held.
func (t *TraceRecorder) buildZipLocked() ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Write trace events: <chunkIndex>-trace.trace
	traceName := fmt.Sprintf("%d-trace.trace", t.chunkIndex)
	tw, err := zw.Create(traceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace entry: %w", err)
	}
	for _, event := range t.events {
		data, err := json.Marshal(event)
		if err != nil {
			continue
		}
		tw.Write(data)
		tw.Write([]byte("\n"))
	}

	// Write network events: <chunkIndex>-trace.network
	netName := fmt.Sprintf("%d-trace.network", t.chunkIndex)
	nw, err := zw.Create(netName)
	if err != nil {
		return nil, fmt.Errorf("failed to create network entry: %w", err)
	}
	for _, event := range t.network {
		data, err := json.Marshal(event)
		if err != nil {
			continue
		}
		nw.Write(data)
		nw.Write([]byte("\n"))
	}

	// Write resources: resources/<sha1>.<ext>
	for hash, data := range t.resources {
		ext := ".png"
		if len(data) > 0 && data[0] == '<' {
			ext = ".html"
		}
		rw, err := zw.Create("resources/" + hash + ext)
		if err != nil {
			continue
		}
		rw.Write(data)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}

	return buf.Bytes(), nil
}

// sha1Hex returns the lowercase hex-encoded SHA1 hash of data.
func sha1Hex(data []byte) string {
	h := sha1.Sum(data)
	return fmt.Sprintf("%x", h)
}

// decodeBase64 decodes a standard base64 string.
func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// pngDimensions reads width and height from a PNG file's IHDR chunk.
// Returns (0, 0) if the data is not a valid PNG.
func pngDimensions(data []byte) (int, int) {
	// PNG header: 8 bytes signature + 4 bytes chunk length + 4 bytes "IHDR" + 4 bytes width + 4 bytes height
	if len(data) < 24 {
		return 0, 0
	}
	w := int(binary.BigEndian.Uint32(data[16:20]))
	h := int(binary.BigEndian.Uint32(data[20:24]))
	return w, h
}

// WriteTraceToFile writes trace zip data to a file, creating directories as needed.
func WriteTraceToFile(data []byte, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create trace dir: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
