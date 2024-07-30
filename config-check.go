package unixid

const sizeBuf = int32(19)

func configCheck(c *Config) (*UnixID, error) {

	if c == nil {
		return nil, errConf
	}

	if c.timeNano == nil {
		return nil, errNano
	}

	if c.timeSeconds == nil {
		return nil, errSecond
	}

	return &UnixID{
		userNum:           "",
		lastUnixNano:      0,
		correlativeNumber: 0,
		buf:               make([]byte, 0, sizeBuf),
		Config:            c,
	}, nil
}
