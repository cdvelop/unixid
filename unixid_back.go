//go:build !wasm
// +build !wasm

package unixid

import (
	"sync"

	"github.com/cdvelop/tinytime"
)

// createUnixID para server usa tinytime.TimeProvider
func createUnixID(params ...any) (*UnixID, error) {
	t := tinytime.NewTimeProvider()

	c := &Config{
		Session:      &defaultEmptySession{},
		TimeProvider: t,
		syncMutex:    &sync.Mutex{},
	}

	externalMutexProvided := false

	for _, param := range params {
		switch mutex := param.(type) {
		case *sync.Mutex:
			externalMutexProvided = true
		case sync.Mutex:
			externalMutexProvided = true
		case userSessionNumber:
			c.Session = mutex
		}
	}

	if externalMutexProvided {
		c.syncMutex = &defaultNoOpMutex{}
	}

	return configCheck(c)
}
