//go:build wasm
// +build wasm

package unixid

import (
	"syscall/js"

	"github.com/cdvelop/tinystring"
)

type timeCLient struct{}

// createUnixID implementa la función NewUnixID para entornos WebAssembly.
// Recibe un manejador de sesión de usuario opcional y configura un UnixID para su uso en el cliente.
// En entornos WebAssembly, se requiere un manejador userSessionNumber para generar IDs únicos entre sesiones.
func createUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	t := timeCLient{}

	// Crear la configuración con valores por defecto
	c := &Config{
		Session:     &defaultEmptySession{}, // Valor por defecto
		timeNano:    t,
		timeSeconds: t,
		syncMutex:   &defaultNoOpMutex{}, // Mutex que no hace nada para entornos WASM
	}

	// Reemplazar el manejador de sesión si se proporciona uno
	for _, u := range handlerUserSessionNumber {
		if usNumber, ok := u.(userSessionNumber); ok {
			c.Session = usNumber
		}
	}

	return configCheck(c)
}

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
	date = date[0:10] + " " + date[11:16]

	return
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
		// Parse string to int64
		parsed := int64(0)
		multiplier := int64(1)
		for i := len(v) - 1; i >= 0; i-- {
			if v[i] >= '0' && v[i] <= '9' {
				parsed += int64(v[i]-'0') * multiplier
				multiplier *= 10
			} else {
				return "" // Invalid character
			}
		}
		unixNano = parsed
	default:
		return ""
	}

	// Convert nanoseconds to seconds
	unixSeconds := unixNano / 1e9

	// Use JavaScript Date to format time
	jsDate := js.Global().Get("Date").New(unixSeconds * 1000) // JS Date expects milliseconds

	// Get hours, minutes, seconds
	hours := jsDate.Call("getHours").Int()
	minutes := jsDate.Call("getMinutes").Int()
	seconds := jsDate.Call("getSeconds").Int()

	// Format as HH:MM:SS
	return tinystring.Format("%02d:%02d:%02d", hours, minutes, seconds).String()
}
