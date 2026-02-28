package proxy

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

// pipeWriteQueueSize is the buffer size for outgoing messages.
// This prevents the browser-to-client routing goroutine from blocking
// on slow pipe writes (e.g., when the Python client is busy processing
// events), which would also block delivery of internal command responses
// and cause navigation timeouts.
const pipeWriteQueueSize = 4096

// PipeClientConn implements ClientTransport over stdin/stdout pipes.
type PipeClientConn struct {
	writer *bufio.Writer
	mu     sync.Mutex
	closed bool
	msgCh  chan string
	done   chan struct{}
}

// NewPipeClientConn creates a PipeClientConn that writes protocol messages to w.
func NewPipeClientConn(w io.Writer) *PipeClientConn {
	c := &PipeClientConn{
		writer: bufio.NewWriter(w),
		msgCh:  make(chan string, pipeWriteQueueSize),
		done:   make(chan struct{}),
	}
	go c.writeLoop()
	return c
}

// writeLoop drains the message channel and writes to the pipe.
func (c *PipeClientConn) writeLoop() {
	defer close(c.done)
	for msg := range c.msgCh {
		c.mu.Lock()
		if c.closed {
			c.mu.Unlock()
			return
		}
		c.writer.WriteString(msg)
		c.writer.WriteByte('\n')
		c.writer.Flush()
		c.mu.Unlock()
	}
}

// ID returns a fixed client ID (pipe mode supports exactly one client).
func (c *PipeClientConn) ID() uint64 { return 1 }

// Send queues a JSON message for writing to the pipe.
// This is non-blocking as long as the write queue is not full.
func (c *PipeClientConn) Send(msg string) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return fmt.Errorf("pipe closed")
	}
	c.mu.Unlock()

	select {
	case c.msgCh <- msg:
		return nil
	default:
		// Queue is full — drop the message to avoid blocking.
		// This can happen under extreme event pressure but is better
		// than blocking the browser message routing goroutine.
		return nil
	}
}

// Close marks the pipe as closed and drains the write queue.
func (c *PipeClientConn) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()
	close(c.msgCh)
	<-c.done
	return nil
}
