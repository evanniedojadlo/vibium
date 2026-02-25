package features

import (
	"time"

	"github.com/vibium/clicker/internal/bidi"
	errs "github.com/vibium/clicker/internal/errors"
)

// Default timeouts and intervals
const (
	DefaultTimeout  = 30 * time.Second
	DefaultInterval = 100 * time.Millisecond
)

// WaitOptions configures wait behavior.
type WaitOptions struct {
	Timeout  time.Duration
	Interval time.Duration
}

// WaitForSelector polls until an element matching the selector exists.
func WaitForSelector(client *bidi.Client, context, selector string, opts WaitOptions) error {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}

	deadline := time.Now().Add(opts.Timeout)

	for {
		// Check if element exists
		_, err := client.FindElement(context, selector)
		if err == nil {
			return nil // Element found
		}

		// Check if we've timed out
		if time.Now().After(deadline) {
			return &errs.TimeoutError{
				Selector: selector,
				Timeout:  opts.Timeout,
				Reason:   "element not found",
			}
		}

		// Wait before next poll
		time.Sleep(opts.Interval)
	}
}

// WaitForHidden polls until an element is either not found or not visible.
func WaitForHidden(client *bidi.Client, context, selector string, opts WaitOptions) error {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}

	deadline := time.Now().Add(opts.Timeout)

	for {
		// Check if element exists
		_, err := client.FindElement(context, selector)
		if err != nil {
			// Element not found — it's hidden
			return nil
		}

		// Element exists — check if it's visible
		visible, err := CheckVisible(client, context, selector)
		if err != nil || !visible {
			// Not visible or error checking — treat as hidden
			return nil
		}

		if time.Now().After(deadline) {
			return &errs.TimeoutError{
				Selector: selector,
				Timeout:  opts.Timeout,
				Reason:   "element still visible",
			}
		}

		time.Sleep(opts.Interval)
	}
}
