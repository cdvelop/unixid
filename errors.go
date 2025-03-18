package unixid

import "errors"

// Error definitions
var (
	errConf   = errors.New("configuration options required")
	errNano   = errors.New("timeNano function required")
	errSecond = errors.New("timeSeconds function required")
	erNumSes  = errors.New("session number required (empty number detected)")
	erSes     = errors.New("session handler required")
	errMutex  = errors.New("sync mutex required")
)
