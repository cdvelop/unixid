package unixid

import "errors"

// ej: 1123466.42 = 2006-01-02 15:04
func (u *UnixID) UnixNanoToStringDate(unixNanoStr string) (string, error) {

	// unixNano = int64
	unixNano, err := validateID(unixNanoStr)
	if err != nil {
		return "", err
	}

	// Convierte unixNano int64 a segundos
	unixSeconds := unixNano / 1e9

	if u.timeSeconds == nil {
		return "", errors.New("adaptador unixTimeSeconds nil")
	}

	return u.timeSeconds.UnixSecondsToDate(unixSeconds), nil
}

// https://chat.openai.com/c/4af98def-f8d9-4095-bf31-deaaad84c094
