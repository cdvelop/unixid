package unixid

// always return ""
type NoSessionNumber struct{}

// always return ""
func (NoSessionNumber) UserSessionNumber() (num, err string) {
	return "", ""
}
