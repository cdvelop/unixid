//go:build !wasm
// +build !wasm

package unixid

import (
	"sync"
	"time"
)

// NewHandler creates a new UnixID handler for server-side environments.
// This version doesn't require a session handler as user numbers aren't needed in server environments.
// Returns an initialized UnixID instance ready to generate unique IDs.
func NewHandler(none ...any) (*UnixID, error) {

	t := &timeServer{}

	c := &Config{
		Session:     nil,
		timeNano:    t,
		timeSeconds: t,
		syncMutex:   &sync.Mutex{},
	}

	return configCheck(c)
}

// timeServer implements time functions for server-side environments
type timeServer struct{}

// UnixNano returns the current Unix time in nanoseconds
func (timeServer) UnixNano() int64 {
	return time.Now().UnixNano()
}

// UnixSecondsToDate converts a Unix timestamp in seconds to a formatted date string
// Format: "2006-01-02 15:04" (year-month-day hour:minute)
func (timeServer) UnixSecondsToDate(unixSeconds int64) (date string) {
	// Create a time.Time instance from the Unix timestamp in seconds
	t := time.Unix(unixSeconds, 0)

	// Format the date in the desired format
	return t.Format("2006-01-02 15:04")
}
