package unixid

// always return ""
func (defaultAuthNumber) UserAuthNumber() (num, err string) {
	return "", ""
}
