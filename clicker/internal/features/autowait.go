package features

import (
	"fmt"
	"time"

	"github.com/vibium/clicker/internal/bidi"
	errs "github.com/vibium/clicker/internal/errors"
)

// Default timeouts and intervals
const (
	DefaultTimeout  = 30 * time.Second
	DefaultInterval = 100 * time.Millisecond
)

// Check represents an actionability check type.
type Check int

const (
	CheckVisibleType Check = iota
	CheckStableType
	CheckReceivesEventsType
	CheckEnabledType
	CheckEditableType
)

// String returns the check name for error messages.
func (c Check) String() string {
	switch c {
	case CheckVisibleType:
		return "Visible"
	case CheckStableType:
		return "Stable"
	case CheckReceivesEventsType:
		return "ReceivesEvents"
	case CheckEnabledType:
		return "Enabled"
	case CheckEditableType:
		return "Editable"
	default:
		return "Unknown"
	}
}

// Predefined check sets for different actions
var (
	// ClickChecks are the checks required before clicking an element.
	ClickChecks = []Check{
		CheckVisibleType,
		CheckStableType,
		CheckReceivesEventsType,
		CheckEnabledType,
	}

	// TypeChecks are the checks required before typing into an element.
	TypeChecks = []Check{
		CheckVisibleType,
		CheckStableType,
		CheckReceivesEventsType,
		CheckEnabledType,
		CheckEditableType,
	}
)

// WaitOptions configures wait behavior.
type WaitOptions struct {
	Timeout  time.Duration
	Interval time.Duration
}

// DefaultWaitOptions returns the default wait configuration.
func DefaultWaitOptions() WaitOptions {
	return WaitOptions{
		Timeout:  DefaultTimeout,
		Interval: DefaultInterval,
	}
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

// WaitForActionable polls until all specified checks pass for the element.
func WaitForActionable(client *bidi.Client, context, selector string, checks []Check, opts WaitOptions) error {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}

	deadline := time.Now().Add(opts.Timeout)

	for {
		// Run all checks
		allPassed := true
		var failedCheck Check
		var checkErr error

		for _, check := range checks {
			passed, err := runCheck(client, context, selector, check)
			if err != nil {
				// Element not found or other error - keep waiting
				allPassed = false
				failedCheck = check
				checkErr = err
				break
			}
			if !passed {
				allPassed = false
				failedCheck = check
				break
			}
		}

		if allPassed {
			return nil // All checks passed
		}

		// Check if we've timed out
		if time.Now().After(deadline) {
			reason := fmt.Sprintf("check '%s' failed", failedCheck)
			if checkErr != nil {
				reason = fmt.Sprintf("check '%s' failed: %v", failedCheck, checkErr)
			}
			return &errs.TimeoutError{
				Selector: selector,
				Timeout:  opts.Timeout,
				Reason:   reason,
			}
		}

		// Wait before next poll
		time.Sleep(opts.Interval)
	}
}

// WaitForClick waits until an element is actionable for clicking.
func WaitForClick(client *bidi.Client, context, selector string, opts WaitOptions) error {
	// First wait for element to exist
	if err := WaitForSelector(client, context, selector, opts); err != nil {
		return err
	}
	// Then wait for click checks
	return WaitForActionable(client, context, selector, ClickChecks, opts)
}

// WaitForType waits until an element is actionable for typing.
func WaitForType(client *bidi.Client, context, selector string, opts WaitOptions) error {
	// First wait for element to exist
	if err := WaitForSelector(client, context, selector, opts); err != nil {
		return err
	}
	// Then wait for type checks
	return WaitForActionable(client, context, selector, TypeChecks, opts)
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

// runCheck executes a single actionability check.
func runCheck(client *bidi.Client, context, selector string, check Check) (bool, error) {
	switch check {
	case CheckVisibleType:
		return CheckVisible(client, context, selector)
	case CheckStableType:
		return CheckStable(client, context, selector)
	case CheckReceivesEventsType:
		return CheckReceivesEvents(client, context, selector)
	case CheckEnabledType:
		return CheckEnabled(client, context, selector)
	case CheckEditableType:
		return CheckEditable(client, context, selector)
	default:
		return false, fmt.Errorf("unknown check type: %d", check)
	}
}
