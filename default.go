package unixid

// always return ""
func (defaultAuthNumber) UserSessionNumber() (num, err string) {
	return "", ""
}
