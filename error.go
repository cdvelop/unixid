package unixid

import "errors"

const erConfHead = "en unixid debes de ingresar un Manejador valido de "

var (
	errConf = errors.New("configuración (&Config = nil)")

	errNano = errors.New(erConfHead + "tiempo, que retorne el método UnixNano() int64")

	errSecond = errors.New(erConfHead + "tiempo, que retorne el método UnixSecondsToDate(time_seconds int64) (date string)")

	erSes = errors.New(erConfHead + "ej: userSessionNumber() (string, error)")

	erNumSes = errors.New(erConfHead + "no existe numero de usuario para generar id")
)
