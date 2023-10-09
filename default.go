package unixid

// always return 0
func (defaultAuthNumber) UserAuthNumber() string {
	return "0"
}
