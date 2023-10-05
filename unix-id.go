package unixid

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type UnixID struct {
	lastUnixIDatabase string
	controlProcess    sync.Mutex
}

func NewHandler() *UnixID {
	return &UnixID{
		lastUnixIDatabase: "",
		controlProcess:    sync.Mutex{},
	}
}

// GetNewID retorna un id único para el ingreso a la base de datos tipo unix time
func (id *UnixID) GetNewID() string {
	id.controlProcess.Lock()
	idunix := fmt.Sprint(time.Now().UnixNano())
	for {
		// si no esta el id lo agregamos
		if id.lastUnixIDatabase != idunix { //last id unix time
			break
		} else {
			//obtenemos nueva marca
			idunix = strconv.FormatInt(time.Now().UnixNano(), 10)
			// log.Printf(">>>new time id slow %v ", idunix)
		}
	}
	// log.Printf("unix time maps %v", setting.lastUnixIDatabase)
	id.lastUnixIDatabase = idunix //actualizo id
	id.controlProcess.Unlock()
	return idunix
}
