package unixid

// always return ""
func (defaultAuthNumber) UserAuthNumber() (string, error) {
	return "", nil
}
