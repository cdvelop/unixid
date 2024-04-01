//go:build wasm
// +build wasm

package unixid

func (id *UnixID) GetNewID() (string, error) {
	var err error
	id.userNum, err = id.Session.userSessionNumber()
	if err != nil {
		return "", err
	}

	id.outID = id.unixIdNano()

	if id.userNum != "" {
		id.outID += "."
		id.outID += id.userNum
	}

	return id.outID, nil
}
