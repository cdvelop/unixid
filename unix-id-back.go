//go:build !wasm
// +build !wasm

package unixid

// GetNewID retorna un id único para el ingreso a la base de datos tipo unix time
// formato 121212.2  (unix time + . + user number)
// 1 allocs
func (id *UnixID) GetNewID() (string, error) {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	return id.unixIdNano(), nil
}
