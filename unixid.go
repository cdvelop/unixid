package unixid

import (
	"reflect"
	"strconv"
	"unsafe"
)

const prefixNameID = "id_"
const sizeBuf = int32(19)

// unixTimeNano represents a time provider that can return Unix time in nanoseconds
type unixTimeNano interface {
	UnixNano() int64
}

// unixTimeSeconds represents a time provider that can format Unix timestamps as human-readable dates
type unixTimeSeconds interface {
	UnixSecondsToDate(time_seconds int64) (date string)
}

// lockHandler represents a mutex-like locking mechanism for thread safety
// Typically a sync.Mutex or similar implementation is used
type lockHandler interface {
	Lock()
	Unlock()
}

// UnixID is the main struct for ID generation and handling
// It contains all configuration and state needed for ID generation
type UnixID struct {
	// userNum is the user session identifier (used in WebAssembly environments)
	userNum string

	// lastUnixNano stores the last generated timestamp to detect collisions
	lastUnixNano int64

	// correlativeNumber is incremented when two IDs would otherwise have the same timestamp
	correlativeNumber int64

	// buf is a pre-allocated buffer to minimize allocations during ID generation
	buf []byte

	// Config holds the external dependencies for the UnixID
	*Config
}

// Config holds the configuration and dependencies for a UnixID instance
type Config struct {
	// Session provides user session numbers in WebAssembly environments
	Session userSessionNumber // e.g., userSessionNumber() string = "1","4","4000" etc.

	// timeNano provides Unix timestamps at nanosecond precision
	timeNano unixTimeNano // e.g., time.Now().UnixNano()

	// timeSeconds formats Unix timestamps as human-readable dates
	timeSeconds unixTimeSeconds // e.g., converts 15454454677767 to "2006-01-02 15:04"

	// syncMutex provides thread safety for concurrent ID generation
	syncMutex lockHandler // e.g., sync.Mutex{}
}

// NewUnixID creates a new UnixID handler with appropriate configuration based on the runtime environment.
//
// For WebAssembly environments (client-side):
// - Requires a userSessionNumber handler to be passed as a parameter
// - Creates IDs with format: "[timestamp].[user_number]" (e.g., "1624397134562544800.42")
// - No mutex is used as JavaScript is single-threaded
//
// For non-WebAssembly environments (server-side):
// - Does not require any parameters
// - Creates IDs with format: "[timestamp]" (e.g., "1624397134562544800")
// - Uses a sync.Mutex for thread safety
//
// Parameters:
//   - handlerUserSessionNumber: Optional userSessionNumber implementation (required for WebAssembly)
//
// Returns:
//   - A configured *UnixID instance
//   - An error if the configuration is invalid
//
// Usage examples:
//
//	// Server-side usage:
//	idHandler, err := unixid.NewUnixID()
//
//	// WebAssembly usage:
//	type sessionHandler struct{}
//	func (sessionHandler) userSessionNumber() string { return "42" }
//	idHandler, err := unixid.NewUnixID(&sessionHandler{})
func NewUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	// The actual implementation is in the build-specific files
	// This function declaration allows for a unified API
	// Implementation details are in unixid_front.go and unixid_back.go
	// and are selected at build time based on the target platform
	return createUnixID(handlerUserSessionNumber...)
}

func configCheck(c *Config) (*UnixID, error) {
	if c == nil {
		return nil, errConf
	}

	if c.timeNano == nil {
		return nil, errNano
	}

	if c.timeSeconds == nil {
		return nil, errSecond
	}

	// Para entornos WebAssembly, verificamos si se requiere un Session
	if c.Session != nil {
		userNum := c.Session.userSessionNumber()
		if userNum == "" {
			return nil, erNumSes
		}
	}

	return &UnixID{
		userNum:           "",
		lastUnixNano:      0,
		correlativeNumber: 0,
		buf:               make([]byte, 0, sizeBuf),
		Config:            c,
	}, nil
}

// userSessionNumber is an interface to obtain the current user's session number
// This is primarily used in WebAssembly environments to uniquely identify client sessions
type userSessionNumber interface {
	// userSessionNumber returns a unique identifier for the current user session
	// e.g., "1" or "2" or "34" or "400" etc.
	userSessionNumber() string
}

// SetNewID sets a new unique ID value to various types of targets.
// It generates a new unique ID based on Unix nanosecond timestamp and assigns it to the provided target.
// This function can work with multiple target types including reflect.Value, string pointers, and byte slices.
//
// In WebAssembly environments, IDs include a user session number as a suffix (e.g., "1624397134562544800.42").
// In server environments, IDs are just the timestamp (e.g., "1624397134562544800").
//
// Parameters:
//   - target: The target to receive the new ID. Can be:
//   - *reflect.Value: For setting struct field values via reflection
//   - *string: For setting a string variable directly
//   - []byte: For appending the ID to a byte slice
//
// This function is thread-safe in server-side environments.
//
// Examples:
//
//	// Setting a struct field using reflection
//	rv := reflect.ValueOf(&myStruct).Elem().FieldByName("ID")
//	idHandler.SetNewID(&rv)
//
//	// Setting a string variable
//	var id string
//	idHandler.SetNewID(&id)
//
//	// Appending to a byte slice
//	buf := make([]byte, 0, 64)
//	idHandler.SetNewID(buf)
func (id *UnixID) SetNewID(target any) {
	// Apply locking if mutex is available (server-side environments)
	if id.syncMutex != nil {
		id.syncMutex.Lock()
		defer id.syncMutex.Unlock()
	}

	// Generate a new ID
	newID := id.unixIdNano()

	// In WebAssembly environments, append the user session number
	if id.Session != nil {
		// Get or update the user number
		if id.userNum == "" {
			id.userNum = id.Session.userSessionNumber()
		}

		// Only append if we have a valid user number
		if id.userNum != "" {
			newID += "."
			newID += id.userNum
		}
	}

	// Set the ID to the appropriate target type
	switch t := target.(type) {
	case *reflect.Value:
		// For struct fields via reflection
		t.SetString(newID)
	case *string:
		// For string variables
		*t = newID
	case []byte:
		// For byte slices, we append the ID
		// The caller is responsible for ensuring the slice has sufficient capacity
		_ = append(t, []byte(newID)...)
	}
}

func (id *UnixID) setValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {
	*valueOut = id.unixIdNano()

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// agregamos el id al campo de la estructura origen
	rv.SetString(*valueOut)

	return nil
}

func (id *UnixID) unixIdNano() string {
	currentUnixNano := id.timeNano.UnixNano()

	if currentUnixNano == id.lastUnixNano {
		//mientras sean iguales sumar numero correlativo
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	// actualizo la variable unix nano
	id.lastUnixNano = currentUnixNano

	currentUnixNano += id.correlativeNumber

	return strconv.FormatInt(currentUnixNano, 10)
}

func (id *UnixID) unixIdNanoLAB() string {
	// ...existing code...
	currentUnixNano := id.timeNano.UnixNano()

	if currentUnixNano == id.lastUnixNano {
		//mientras sean iguales sumar numero correlativo
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	// actualizo la variable unix nano
	id.lastUnixNano = currentUnixNano

	currentUnixNano += id.correlativeNumber

	id.buf = id.buf[:0]

	// fmt.Println("size buffer:", sizeBuf)

	id.buf = strconv.AppendInt(id.buf, currentUnixNano, 10)
	// id.buf = strconv.AppendUint(id.buf, currentUnixNano, 10)

	// fmt.Println("tmpBuf:", id.buf, "size id buffer:", len(id.buf))

	// return string(id.buf)
	return *(*string)(unsafe.Pointer(&id.buf))
	// return unsafe.String(unsafe.SliceData(id.buf), len(id.buf))
}
