//go:build !wasm
// +build !wasm

package unixid

import (
	"strconv"
	"sync"
	"time"
)

// timeServer implements time functions for server-side environments
type timeServer struct{}

// createUnixID implements the NewUnixID function for non-WebAssembly environments (server).
// It configures a UnixID for server use with a mutex for synchronization.
// In server environments, no user session handler is needed.
// If a sync.Mutex is provided as a parameter, that mutex will be used instead of creating a new one,
// which avoids potential deadlocks when integrating with other libraries using sync.
func createUnixID(params ...any) (*UnixID, error) {
	t := &timeServer{}

	c := &Config{
		Session:     &defaultEmptySession{}, // Use the default implementation that does nothing
		timeNano:    t,
		timeSeconds: t,
		syncMutex:   &sync.Mutex{}, // Default mutex
	}

	// Look for a mutex in the provided parameters
	for _, param := range params {
		switch mutex := param.(type) {
		case *sync.Mutex:
			// Use the provided mutex instead of the default
			c.syncMutex = mutex
		case sync.Mutex:
			// If a mutex is passed by value, convert it to a pointer
			// This is a copy of the original mutex, but it's better than nothing
			ptrMutex := &mutex
			c.syncMutex = ptrMutex
		case userSessionNumber:
			// If a user session handler is provided, use it
			c.Session = mutex
		}
	}

	return configCheck(c)
}

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

// UnixSecondsToTime converts a Unix timestamp in seconds to a formatted time string.
// Format: "15:04:05" (hour:minute:second)
// It accepts a parameter of type any and attempts to convert it to an int64 Unix timestamp.
// eg: 1624397134 -> "15:32:14"
// supported types: int64, int, float64, string
func (u UnixID) UnixSecondsToTime(input any) string {
	var unixSeconds int64
	switch v := input.(type) {
	case int64:
		unixSeconds = v
	case int:
		unixSeconds = int64(v)
	case float64:
		unixSeconds = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return ""
		}
		unixSeconds = parsed
	default:
		return ""
	}

	t := time.Unix(unixSeconds, 0)
	return t.Format("15:04:05")
}
