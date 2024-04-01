//go:build wasm
// +build wasm

package unixid

import (
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

type timeCLient struct{}

func (timeCLient) UnixNano() int64 {
	jsDate := js.Global().Get("Date").New()
	msTimestamp := jsDate.Call("getTime").Float()
	nanoTimestamp := int64(msTimestamp * 1e6)

	// js.Global().Get("console").Call("log", "tiempo unix ID:", nanoTimestamp)

	return nanoTimestamp
}

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
