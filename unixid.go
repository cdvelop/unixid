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
	userSessionNumber() (number string, err error)
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

// v := int64(42)
// b := unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v))
// fmt.Println(b, "id:", string(b))

// b := *(*[]byte)(unsafe.Pointer(&v)) // Cast directly to a []byte pointer
// fmt.Println(b, "id:", string(b))    // Output: [42 0 0 0 0 0 0 0] id: *

// buf := unsafe.Slice((*byte)(unsafe.Pointer(&currentUnixNano)), unsafe.Sizeof(currentUnixNano))

// fmt.Println("id:", string(buf))

// return unsafe.String(unsafe.SliceData(buf), len(buf))
// t := time.Now().UTC().UnixNano()
//     b := unsafe.Slice((*byte)(unsafe.Pointer(&t)), unsafe.Sizeof(t))
//     fmt.Println(b)

// https://stackoverflow.com/questions/76431857/what-is-the-fastest-way-to-convert-int64-to-byte-array
