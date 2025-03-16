//go:build !wasm
// +build !wasm

package unixid

import (
	"reflect"
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

// GetNewID generates a new unique ID based on Unix nanosecond timestamp.
// In server environments, this returns just the Unix nanosecond timestamp value.
// The method is thread-safe and handles concurrent access through a mutex lock.
// Returns a string representation of the unique ID.
func (id *UnixID) GetNewID() (string, error) {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	return id.unixIdNano(), nil
}

// SetValue sets a unique ID value to a struct field using reflection.
// This is used internally to populate struct fields with unique IDs.
// Parameters:
//   - rv: A reflect.Value pointer to the struct field that will receive the ID
//   - valueOut: A pointer to a string that will store the generated ID
//   - sizeOut: A byte slice that will track the size of the generated value
//
// Returns nil on success or an error if the operation fails.
func (id *UnixID) SetValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	*valueOut = id.unixIdNano()

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// Add the ID to the struct field of the source structure
	rv.SetString(*valueOut)

	return nil
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

// UnixSecondsToTime converts a Unix timestamp in seconds to a formatted time string.
// Format: "15:04:05" (hour:minute:second)
func (timeServer) UnixSecondsToTime(unixSeconds int64) string {
	t := time.Unix(unixSeconds, 0)
	return t.Format("15:04:05")
}
