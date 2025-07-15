//go:build wasm
// +build wasm

package unixid

import "github.com/cdvelop/tinytime"

// createUnixID para WASM ahora usa tinytime.TimeProvider
func createUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	t := tinytime.NewTimeProvider()

	c := &Config{
		Session:      &defaultEmptySession{},
		TimeProvider: t,
		syncMutex:    &defaultNoOpMutex{},
	}

	for _, u := range handlerUserSessionNumber {
		if usNumber, ok := u.(userSessionNumber); ok {
			c.Session = usNumber
		}
	}

	return configCheck(c)
}
