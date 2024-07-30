package unixid

const PrefixNameID = "id_"

type unixTimeNano interface {
	UnixNano() int64
}

type unixTimeSeconds interface {
	UnixSecondsToDate(time_seconds int64) (date string)
}

// ej: sync.Mutex{}
type lockHandler interface {
	Lock()
	Unlock()
}

type UnixID struct {
	userNum string

	lastUnixNano int64

	correlativeNumber int64

	buf []byte

	*Config
}

type Config struct {
	Session     userSessionNumber // ej: userSessionNumber() string = 1,4,4000 .... if nil, always return 0.. id ej: 124663.0
	timeNano    unixTimeNano      // unixTimeNano ej: time.Now() = UnixNano() int64
	timeSeconds unixTimeSeconds   // ej: 15454454677767 to: 2006-01-02 15:04 UnixSecondsToDate(time_seconds int64) (date string)
	syncMutex   lockHandler       // lockHandler ej: sync.Mutex{} = Lock() Unlock()
}

type userSessionNumber interface {
	// ej: 1 or 2 or 34 or 400.....
	userSessionNumber() (number string, err error)
}
