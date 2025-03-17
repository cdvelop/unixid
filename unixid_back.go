//go:build !wasm
// +build !wasm

package unixid

import (
	"reflect"
	"strconv"
	"sync"
	"time"
)

// createUnixID implementa la función NewUnixID para entornos no-WebAssembly (servidor).
// Configura un UnixID para su uso en el servidor con un mutex para sincronización.
// En entornos de servidor, no se requiere ningún manejador de sesión de usuario.
func createUnixID(none ...any) (*UnixID, error) {
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
