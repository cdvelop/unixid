//go:build wasm
// +build wasm

package unixid

import "reflect"

func (id *UnixID) GetNewID() (string, error) {

	outID := id.unixIdNano()

	if err := id.setUserNumber(); err != nil {
		return "", err
	}

	outID += "."
	outID += id.userNum

	return outID, nil
}

func (id *UnixID) SetValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {

	*valueOut = id.unixIdNano()

	if err := id.setUserNumber(); err != nil {
		return err
	}

	*valueOut += "."
	*valueOut += id.userNum

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// agregamos el id al campo de la estructura origen
	rv.SetString(*valueOut)

	return nil
}

func (id *UnixID) setUserNumber() (err error) {

	if id.userNum == "" {
		id.userNum, err = id.Session.userSessionNumber()
		if err != nil {
			return
		}

		if id.userNum == "" {
			err = erNumSes
		}
	}

	return
}
