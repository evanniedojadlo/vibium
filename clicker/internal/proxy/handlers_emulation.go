package proxy

import (
	"encoding/json"
	"fmt"
)

// handlePageSetViewport handles vibium:page.setViewport — sets the viewport size.
// Uses BiDi browsingContext.setViewport.
func (r *Router) handlePageSetViewport(session *BrowserSession, cmd bidiCommand) {
	context, err := r.resolveContext(session, cmd.Params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	width, _ := cmd.Params["width"].(float64)
	height, _ := cmd.Params["height"].(float64)
	if width == 0 || height == 0 {
		r.sendError(session, cmd.ID, fmt.Errorf("width and height are required"))
		return
	}

	params := map[string]interface{}{
		"context": context,
		"viewport": map[string]interface{}{
			"width":  int(width),
			"height": int(height),
		},
	}

	if dpr, ok := cmd.Params["devicePixelRatio"].(float64); ok && dpr > 0 {
		params["devicePixelRatio"] = dpr
	}

	resp, err := r.sendInternalCommand(session, "browsingContext.setViewport", params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}
	if bidiErr := checkBidiError(resp); bidiErr != nil {
		r.sendError(session, cmd.ID, bidiErr)
		return
	}

	r.sendSuccess(session, cmd.ID, map[string]interface{}{})
}

// handlePageViewport handles vibium:page.viewport — returns the current viewport size.
// Uses JS eval since BiDi has no viewport getter.
func (r *Router) handlePageViewport(session *BrowserSession, cmd bidiCommand) {
	context, err := r.resolveContext(session, cmd.Params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	val, err := r.evalSimpleScript(session, context,
		`() => JSON.stringify({ width: window.innerWidth, height: window.innerHeight })`)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	var size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	if err := json.Unmarshal([]byte(val), &size); err != nil {
		r.sendError(session, cmd.ID, fmt.Errorf("failed to parse viewport: %w", err))
		return
	}

	r.sendSuccess(session, cmd.ID, map[string]interface{}{
		"width":  size.Width,
		"height": size.Height,
	})
}

// handlePageEmulateMedia handles vibium:page.emulateMedia — overrides CSS media features.
// Uses JS matchMedia override since BiDi has no CSS media feature commands.
// Supports: media, colorScheme, reducedMotion, forcedColors, contrast.
func (r *Router) handlePageEmulateMedia(session *BrowserSession, cmd bidiCommand) {
	context, err := r.resolveContext(session, cmd.Params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	// Build the overrides object from params.
	// Each option can be a string value or null (to reset).
	overrides := map[string]interface{}{}

	for _, key := range []string{"media", "colorScheme", "reducedMotion", "forcedColors", "contrast"} {
		if val, exists := cmd.Params[key]; exists {
			if val == nil {
				overrides[key] = nil
			} else if s, ok := val.(string); ok {
				overrides[key] = s
			}
		}
	}

	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		r.sendError(session, cmd.ID, fmt.Errorf("failed to serialize overrides: %w", err))
		return
	}

	// JS function that installs/updates matchMedia overrides.
	// The override wraps native matchMedia once (idempotent) and intercepts
	// queries for configured CSS media features.
	// NOTE: In Go backtick strings, characters are literal — so \( in the JS
	// source must be written as \( here (no extra escaping).
	script := "(overridesJSON) => {\n" +
		"const overrides = JSON.parse(overridesJSON);\n" +
		"if (!window.__vibiumMediaOverrides) { window.__vibiumMediaOverrides = {}; }\n" +
		"const featureMap = {\n" +
		"  colorScheme: 'prefers-color-scheme',\n" +
		"  reducedMotion: 'prefers-reduced-motion',\n" +
		"  forcedColors: 'forced-colors',\n" +
		"  contrast: 'prefers-contrast'\n" +
		"};\n" +
		"for (const [key, value] of Object.entries(overrides)) {\n" +
		"  if (value === null) { delete window.__vibiumMediaOverrides[key]; }\n" +
		"  else { window.__vibiumMediaOverrides[key] = value; }\n" +
		"}\n" +
		"if (!window.__vibiumOriginalMatchMedia) {\n" +
		"  window.__vibiumOriginalMatchMedia = window.matchMedia.bind(window);\n" +
		"  window.matchMedia = function(query) {\n" +
		"    const original = window.__vibiumOriginalMatchMedia(query);\n" +
		"    const ov = window.__vibiumMediaOverrides || {};\n" +
		"    if (ov.media !== undefined) {\n" +
		"      const q = query.trim().toLowerCase();\n" +
		"      if (q === 'print' || q === '(print)') return makeResult(original, ov.media === 'print', query);\n" +
		"      if (q === 'screen' || q === '(screen)') return makeResult(original, ov.media === 'screen', query);\n" +
		"    }\n" +
		"    for (const [key, feature] of Object.entries(featureMap)) {\n" +
		"      if (ov[key] !== undefined) {\n" +
		"        const re = new RegExp('\\\\(' + feature + '\\\\s*:\\\\s*([^)]+)\\\\)');\n" +
		"        const m = query.match(re);\n" +
		"        if (m) { return makeResult(original, m[1].trim() === ov[key], query); }\n" +
		"      }\n" +
		"    }\n" +
		"    return original;\n" +
		"  };\n" +
		"}\n" +
		"function makeResult(original, matches, media) {\n" +
		"  return {\n" +
		"    matches: matches, media: media, onchange: original.onchange,\n" +
		"    addListener: original.addListener.bind(original),\n" +
		"    removeListener: original.removeListener.bind(original),\n" +
		"    addEventListener: original.addEventListener.bind(original),\n" +
		"    removeEventListener: original.removeEventListener.bind(original),\n" +
		"    dispatchEvent: original.dispatchEvent.bind(original)\n" +
		"  };\n" +
		"}\n" +
		"return 'ok';\n" +
		"}"

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": context},
		"arguments": []map[string]interface{}{
			{"type": "string", "value": string(overridesJSON)},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	resp, err := r.sendInternalCommand(session, "script.callFunction", params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}
	if bidiErr := checkBidiError(resp); bidiErr != nil {
		r.sendError(session, cmd.ID, bidiErr)
		return
	}

	r.sendSuccess(session, cmd.ID, map[string]interface{}{})
}

// handlePageSetContent handles vibium:page.setContent — replaces the page HTML.
// Uses document.open/write/close to fully replace the document.
func (r *Router) handlePageSetContent(session *BrowserSession, cmd bidiCommand) {
	context, err := r.resolveContext(session, cmd.Params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	html, _ := cmd.Params["html"].(string)

	script := `(html) => { document.open(); document.write(html); document.close(); }`

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": context},
		"arguments": []map[string]interface{}{
			{"type": "string", "value": html},
		},
		"awaitPromise":    true,
		"resultOwnership": "root",
	}

	resp, err := r.sendInternalCommand(session, "script.callFunction", params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}
	if bidiErr := checkBidiError(resp); bidiErr != nil {
		r.sendError(session, cmd.ID, bidiErr)
		return
	}

	r.sendSuccess(session, cmd.ID, map[string]interface{}{})
}

// handlePageSetGeolocation handles vibium:page.setGeolocation — overrides geolocation.
// Uses JS override of navigator.geolocation since BiDi emulation.setGeolocationOverride
// is not widely supported and requires granted permissions.
func (r *Router) handlePageSetGeolocation(session *BrowserSession, cmd bidiCommand) {
	context, err := r.resolveContext(session, cmd.Params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}

	lat, hasLat := cmd.Params["latitude"].(float64)
	lng, hasLng := cmd.Params["longitude"].(float64)

	if !hasLat || !hasLng {
		r.sendError(session, cmd.ID, fmt.Errorf("latitude and longitude are required"))
		return
	}

	accuracy := float64(1)
	if acc, ok := cmd.Params["accuracy"].(float64); ok {
		accuracy = acc
	}

	coordsJSON, _ := json.Marshal(map[string]float64{
		"latitude":  lat,
		"longitude": lng,
		"accuracy":  accuracy,
	})

	script := "(coordsJSON) => {\n" +
		"const coords = JSON.parse(coordsJSON);\n" +
		"const geo = navigator.geolocation;\n" +
		"geo.getCurrentPosition = function(success, error, options) {\n" +
		"  success({ coords: { latitude: coords.latitude, longitude: coords.longitude, accuracy: coords.accuracy,\n" +
		"    altitude: null, altitudeAccuracy: null, heading: null, speed: null }, timestamp: Date.now() });\n" +
		"};\n" +
		"geo.watchPosition = function(success, error, options) {\n" +
		"  success({ coords: { latitude: coords.latitude, longitude: coords.longitude, accuracy: coords.accuracy,\n" +
		"    altitude: null, altitudeAccuracy: null, heading: null, speed: null }, timestamp: Date.now() });\n" +
		"  return 0;\n" +
		"};\n" +
		"return 'ok';\n" +
		"}"

	params := map[string]interface{}{
		"functionDeclaration": script,
		"target":              map[string]interface{}{"context": context},
		"arguments": []map[string]interface{}{
			{"type": "string", "value": string(coordsJSON)},
		},
		"awaitPromise":    false,
		"resultOwnership": "root",
	}

	resp, err := r.sendInternalCommand(session, "script.callFunction", params)
	if err != nil {
		r.sendError(session, cmd.ID, err)
		return
	}
	if bidiErr := checkBidiError(resp); bidiErr != nil {
		r.sendError(session, cmd.ID, bidiErr)
		return
	}

	r.sendSuccess(session, cmd.ID, map[string]interface{}{})
}
