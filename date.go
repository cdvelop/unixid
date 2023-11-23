package unixid

import (
	"time"
)

// ej: 1123466.42 = 2006-01-02 15:04
func UnixNanoToStringDate(unixNanoStr string) (date, err string) {

	unixNano, er := validateID(unixNanoStr)
	if er != "" {
		return "", er
	}

	// Convierte el Unixtime en segundos
	unixSeconds := unixNano / 1e9

	// Crea una instancia de time.Time a partir del Unixtime en segundos
	t := time.Unix(unixSeconds, 0)

	// Formatea la fecha en el formato deseado
	formattedTime := t.Format("2006-01-02 15:04")

	return formattedTime, ""
}

// https://chat.openai.com/c/4af98def-f8d9-4095-bf31-deaaad84c094
