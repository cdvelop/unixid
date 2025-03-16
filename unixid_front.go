//go:build wasm
// +build wasm

package unixid

import (
	"reflect"
	"syscall/js"
)

// userSessionNumber() (number string, err error)
func NewHandler(handlerUserSessionNumber ...any) (*UnixID, error) {

	t := timeCLient{}

	c := &Config{
		Session:     nil,
		timeNano:    t,
		timeSeconds: t,
		syncMutex:   nil,
	}

	for _, u := range handlerUserSessionNumber {
		if usNumber, ok := u.(userSessionNumber); ok {
			c.Session = usNumber
		}
	}

	if c.Session == nil {
		return nil, erSes
	}

	return configCheck(c)
}

// GetNewID generates a new unique ID based on Unix nanosecond timestamp.
// In client-side WebAssembly environments, this appends a user session number
// to the timestamp, separated by a dot (e.g., "1624397134562544800.42").
// This helps ensure that IDs are unique even across different client sessions.
// Returns a string representation of the unique ID or an error if the user session number can't be obtained.
func (id *UnixID) GetNewID() (string, error) {

	outID := id.unixIdNano()

	if err := id.setUserNumber(); err != nil {
		return "", err
	}

	outID += "."
	outID += id.userNum

	return outID, nil
}

// SetValue sets a unique ID value to a struct field using reflection.
// For WASM environments, the ID includes the user session number as a suffix.
// Parameters:
//   - rv: A reflect.Value pointer to the struct field that will receive the ID
//   - valueOut: A pointer to a string that will store the generated ID
//   - sizeOut: A byte slice that will track the size of the generated value
//
// Returns nil on success or an error if the operation fails.
func (id *UnixID) SetValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {

	*valueOut = id.unixIdNano()

	if err := id.setUserNumber(); err != nil {
		return err
	}

	*valueOut += "."
	*valueOut += id.userNum

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// Add the ID to the struct field of the source structure
	rv.SetString(*valueOut)

	return nil
}

// setUserNumber obtains and caches the user session number.
// If the user number hasn't been set yet, it calls the userSessionNumber method
// provided in the Config to obtain it.
// Returns an error if the user session number can't be obtained.
func (id *UnixID) setUserNumber() (err error) {

	if id.userNum == "" {
		id.userNum, err = id.Session.userSessionNumber()
		if err != nil {
			return
		}

		if id.userNum == "" {
			err = erNumSes
		}
	}

	return
}

type timeCLient struct{}

// UnixNano retrieves the current Unix timestamp in nanoseconds.
// It creates a new JavaScript Date object, gets the timestamp in milliseconds,
// converts it to nanoseconds, and returns the result as an int64.
//
// Example usage:
//
//	// Assuming t is of type timeCLient (with method UnixNano)
//	var t timeCLient
//	timestamp := t.UnixNano()
//	fmt.Println("Current Unix timestamp in nanoseconds:", timestamp)
//
// Note: This method relies on JavaScript's Date API via js.Global(), so it
// works in a JavaScript-enabled environment (e.g., GopherJS or WebAssembly).
func (timeCLient) UnixNano() int64 {
	jsDate := js.Global().Get("Date").New()
	msTimestamp := jsDate.Call("getTime").Float()
	nanoTimestamp := int64(msTimestamp * 1e6)

	// js.Global().Get("console").Call("log", "tiempo unix ID:", nanoTimestamp)

	return nanoTimestamp
}

// UnixSecondsToDate converts a Unix timestamp (in seconds) to a formatted date string.
//
// This function creates a JavaScript Date object from the provided Unix seconds,
// then obtains an ISO 8601 string representation from it. It slices the resulting
// string to output a date in the "YYYY-MM-DD HH:MM:SS" format.
//
// Parameters:
//
//	unixSeconds int64 - the Unix timestamp in seconds (i.e., seconds elapsed since January 1, 1970 UTC)
//
// Returns:
//
//	string - a formatted date string in the "YYYY-MM-DD HH:MM:SS" format.
//
// Example:
//
//	// Example: Converting a Unix timestamp for January 1, 2021 at midnight UTC
//	ts := int64(1609459200) // January 1, 2021 00:00:00 UTC
//	formattedDate := UnixSecondsToDate(ts)
//	fmt.Println(formattedDate) // Output: "2021-01-01 00:00:00"
func (timeCLient) UnixSecondsToDate(unixSeconds int64) (date string) {
	// Crea una instancia de Date de JavaScript a partir de los segundos de Unix
	jsDate := js.Global().Get("Date").New(float64(unixSeconds) * 1000)

	// Llama al método toISOString para obtener la fecha formateada
	dateJSValue := jsDate.Call("toISOString")

	// Convierte el valor de JavaScript a una cadena de Go
	date = dateJSValue.String()

	// Formatea la cadena de fecha a "2006-01-02 15:04"
	date = date[0:10] + " " + date[11:19]

	return
}
