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
// If a sync.Mutex is provided as a parameter, the function will assume that external
// synchronization is being handled by the caller. In this case, a defaultNoOpMutex
// will be used internally to prevent potential deadlocks.
func createUnixID(params ...any) (*UnixID, error) {
	t := &timeServer{}

	c := &Config{
		Session:     &defaultEmptySession{}, // Use the default implementation that does nothing
		timeNano:    t,
		timeSeconds: t,
		syncMutex:   &sync.Mutex{}, // Default mutex
	}

	externalMutexProvided := false

	// Look for a mutex in the provided parameters
	for _, param := range params {
		switch mutex := param.(type) {
		case *sync.Mutex:
			// If external mutex is provided, use a no-op mutex internally
			// to prevent deadlocks when GetNewID is called inside another lock
			externalMutexProvided = true
		case sync.Mutex:
			// If a mutex is passed by value, we also consider that an external mutex was provided
			externalMutexProvided = true
		case userSessionNumber:
			// If a user session handler is provided, use it
			c.Session = mutex
		}
	}

	// If an external mutex was provided, use a no-op mutex to avoid deadlocks
	if externalMutexProvided {
		c.syncMutex = &defaultNoOpMutex{}
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

// UnixNanoToTime converts a Unix timestamp in nanoseconds to a formatted time string.
// Format: "15:04:05" (hour:minute:second)
// It accepts a parameter of type any and attempts to convert it to an int64 Unix timestamp in nanoseconds.
// eg: 1624397134562544800 -> "15:32:14"
// supported types: int64, int, float64, string
func (u UnixID) UnixNanoToTime(input any) string {
	var unixNano int64
	switch v := input.(type) {
	case int64:
		unixNano = v
	case int:
		unixNano = int64(v)
	case float64:
		unixNano = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return ""
		}
		unixNano = parsed
	default:
		return ""
	}

	// Convert nanoseconds to seconds
	unixSeconds := unixNano / 1e9
	t := time.Unix(unixSeconds, 0)
	return t.Format("15:04:05")
}
