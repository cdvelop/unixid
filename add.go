package unixid

import "github.com/cdvelop/model"

// unixTimeHandler ej: time.Now() = UnixNano() int64
// lockHandler ej: sync.Mutex{} = Lock() Unlock()
// UserSessionNumber ej: UserSessionNumber() string = 1,4,4000 .... if nil, always return 0.. id ej: 124663.0
func NewHandler(t model.UnixTimeHandler, l lockHandler, u model.UserSessionNumber) (h *UnixID, err string) {

	if t == nil {
		return nil, "debes ingresar un tipo de Manejador de tiempo que retorne el método UnixNano() int64"
	}

	if l == nil {
		return nil, "debes ingresar un tipo de Manejador de protección de escritura con los métodos: Lock() Unlock() ej sync.Mutex{}"
	}

	if u == nil {
		return nil, "debes ingresar un  Manejador UserSessionNumber valido"
	}

	idh := UnixID{
		lastUnixIDatabase: "",
		lockHandler:       l,
		UnixTimeHandler:   t,
		user:              u,
	}

	return &idh, ""
}
