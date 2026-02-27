package proxy

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

// PipeClientConn implements ClientTransport over stdin/stdout pipes.
type PipeClientConn struct {
	writer *bufio.Writer
	mu     sync.Mutex
	closed bool
}

// NewPipeClientConn creates a PipeClientConn that writes protocol messages to w.
func NewPipeClientConn(w io.Writer) *PipeClientConn {
	return &PipeClientConn{
		writer: bufio.NewWriter(w),
	}
}

// ID returns a fixed client ID (pipe mode supports exactly one client).
func (c *PipeClientConn) ID() uint64 { return 1 }

// Send writes a JSON message followed by a newline to the pipe.
func (c *PipeClientConn) Send(msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("pipe closed")
	}

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}
	if err := c.writer.WriteByte('\n'); err != nil {
		return err
	}
	return c.writer.Flush()
}

// Close marks the pipe as closed.
func (c *PipeClientConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return nil
}
