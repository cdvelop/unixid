package unixid

import "errors"

const erConfHead = "in unixid you must supply a valid handler for "

var (
	errConf   = errors.New("configuration (&Config = nil)")
	errNano   = errors.New(erConfHead + "time, which returns the UnixNano() int64 method")
	errSecond = errors.New(erConfHead + "time, which returns the UnixSecondsToDate(time_seconds int64) (date string) method")
	erSes     = errors.New(erConfHead + "e.g.: userSessionNumber() string")
	erNumSes  = errors.New(erConfHead + "user number must not be empty")
)
