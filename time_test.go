//go:build wasm
// +build wasm

package unixid_test

import (
	"reflect"
	"syscall/js"
	"testing"
)

func TestUnixSecondsToDate(t *testing.T) {

	testData := []struct {
		Comment      string
		DataIN       int64
		DataExpected any
	}{

		{"ok", 0, "1970-01-01 00:00"},
		{"ok", 1000, "1970-01-01 00:16"},        // 1000 segundos son 16 minutos y 40 segundos
		{"ok", 1613737200, "2021-02-19 16:00"},  // Ajuste a la zona horaria correcta
		{"ok", -1613737200, "1918-11-14 08:00"}, // Ajuste a la zona horaria correcta
	}

	message := func(comment string, expected, response any) {
		t.Fatalf("\n=> en %v la expectativa es:\n[%v]\n=> pero se obtuvo:\n[%v]\n", comment, expected, response)
	}

	compare := func(comment string, expected, response any) {
		if !reflect.DeepEqual(expected, response) {
			message(comment, expected, response)
		}
	}

	tc := timeCLient{}

	for _, test := range testData {
		t.Run(("\n" + test.Comment), func(t *testing.T) {
			Response := tc.UnixSecondsToDate(test.DataIN)
			compare(test.Comment, test.DataExpected, Response)
		})
	}
}

type timeCLient struct{}

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
