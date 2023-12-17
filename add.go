package unixid

import "github.com/cdvelop/model"

// unixTimeHandler ej: time.Now() = UnixNano() int64
// lockHandler ej: sync.Mutex{} = Lock() Unlock()
// UserSessionNumber ej: UserSessionNumber() string = 1,4,4000 .... if nil, always return 0.. id ej: 124663.0
func NewHandler(t model.UnixTimeHandler, l lockHandler, u model.UserSessionNumber) (h *UnixID, err string) {
	const e = "en unixid debes de ingresar un Manejador valido de "
	if t == nil {
		return nil, e + "tiempo, que retorne el método UnixNano() int64"
	}

	if l == nil {
		return nil, e + "protección de escritura, con los métodos: Lock() Unlock() ej sync.Mutex{}"
	}

	if u == nil {
		return nil, e + " UserSessionNumber ej: UserSessionNumber() (number string, err string)"
	}

	idh := UnixID{
		lastUnixIDatabase: "",
		lockHandler:       l,
		UnixTimeHandler:   t,
		user:              u,
	}

	return &idh, ""
}
