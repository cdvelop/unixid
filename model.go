package unixid

import "github.com/cdvelop/model"

// ej: sync.Mutex{}
type lockHandler interface {
	Lock()
	Unlock()
}

type UnixID struct {
	lastUnixIDatabase string
	lockHandler
	model.UnixTimeHandler
	user model.UserSessionNumber
}
