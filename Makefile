.PHONY: all build build-go build-js deps clean clean-bin clean-js clean-cache clean-all serve help

# Default target
all: build

# Build everything (Go + JS)
build: build-go build-js

# Build clicker binary
build-go: deps
	cd clicker && go build -o bin/clicker ./cmd/clicker

# Build JS client
build-js: deps
	cd clients/javascript && npm run build

# Install npm dependencies (skip if node_modules exists)
deps:
	@if [ ! -d "node_modules" ]; then npm install; fi

# Start the proxy server
serve: build-go
	./clicker/bin/clicker serve

# Clean clicker binaries
clean-bin:
	rm -rf clicker/bin

# Clean JS dist
clean-js:
	rm -rf clients/javascript/dist

# Clean cached Chrome for Testing
clean-cache:
	rm -rf ~/Library/Caches/vibium/chrome-for-testing
	rm -rf ~/.cache/vibium/chrome-for-testing

# Clean everything (binaries + JS dist + cache)
clean-all: clean-bin clean-js clean-cache

# Alias for clean-bin + clean-js
clean: clean-bin clean-js

# Show available targets
help:
	@echo "Available targets:"
	@echo "  make             - Build everything (default)"
	@echo "  make build-go    - Build clicker binary"
	@echo "  make build-js    - Build JS client"
	@echo "  make deps        - Install npm dependencies"
	@echo "  make serve       - Start proxy server on :9515"
	@echo "  make clean       - Clean binaries and JS dist"
	@echo "  make clean-cache - Clean cached Chrome for Testing"
	@echo "  make clean-all   - Clean everything"
	@echo "  make help        - Show this help"
