//go:build !wasm
// +build !wasm

package unixid

import "reflect"

// GetNewID retorna un id único para el ingreso a la base de datos tipo unix time
// formato 121212.2  (unix time + . + user number)
// 1 allocs
func (id *UnixID) GetNewID() (string, error) {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	return id.unixIdNano(), nil
}

func (id *UnixID) SetValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	*valueOut = id.unixIdNano()

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// agregamos el id al campo de la estructura origen
	rv.SetString(*valueOut)

	return nil
}
