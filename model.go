package unixid

import "github.com/cdvelop/model"

// ej: sync.Mutex{}
type lockHandler interface {
	Lock()
	Unlock()
}

// ej: time.Now()
type unixTimeHandler interface {
	UnixNano() int64
}

type UnixID struct {
	lastUnixIDatabase string
	lockHandler
	unixTimeHandler
	user model.UserAuthNumber
}

// always return ""
type defaultAuthNumber struct{}
