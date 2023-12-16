package unixid

import (
	"strconv"
)

// GetNewID retorna un id único para el ingreso a la base de datos tipo unix time
func (id *UnixID) GetNewID() (new_id, err string) {
	id.lockHandler.Lock()
	idunix := strconv.FormatInt(id.UnixNano(), 10)

	// fmt.Println("ID OBTENIDO:", idunix)

	for {
		// si no esta el id lo agregamos
		if id.lastUnixIDatabase != idunix { //last id unix time
			break
		} else {
			//obtenemos nueva marca
			idunix = strconv.FormatInt(id.UnixNano(), 10)
			// log.Printf(">>>new time id slow %v ", idunix)
		}
	}
	// log.Printf("unix time maps %v", setting.lastUnixIDatabase)
	id.lastUnixIDatabase = idunix //actualizo id
	id.lockHandler.Unlock()

	user_num, err := id.user.UserSessionNumber()
	if user_num != "" && err == "" {
		idunix += "." + user_num
	}

	return idunix, err
}
