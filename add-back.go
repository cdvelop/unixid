//go:build !wasm
// +build !wasm

package unixid

import (
	"sync"
	"time"
)

// config server not required
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

type timeServer struct{}

func (timeServer) UnixNano() int64 {
	return time.Now().UnixNano()
}
func (timeServer) UnixSecondsToDate(unixSeconds int64) (date string) {

	// Crea una instancia de time.Time a partir del Unixtime en segundos
	t := time.Unix(unixSeconds, 0)

	// Formatea la fecha en el formato deseado
	return t.Format("2006-01-02 15:04")
}
