package unixid_test

import (
	"sync"
	"testing"

	"github.com/cdvelop/timeserver"
	"github.com/cdvelop/unixid"
)

// always return "1"
type defaultAuthNumber struct{}

// always return ""
func (defaultAuthNumber) UserSessionNumber() (num, err string) {
	return "1", ""
}

func Test_GetNewID(t *testing.T) {
	idRequired := 1000
	wg := sync.WaitGroup{}
	wg.Add(idRequired)

	uid, err := unixid.NewHandler(timeserver.Add(), &sync.Mutex{}, defaultAuthNumber{})
	if err != "" {
		t.Fatal(err)
		return
	}

	idObtained := make(map[string]int)
	var esperar sync.Mutex

	for i := 0; i < idRequired; i++ {
		go func() {
			defer wg.Done()
			id, err := uid.GetNewID()
			if err != "" {
				t.Log(err)
				wg.Done()
				return
			}
			esperar.Lock()
			if cantId, exist := idObtained[id]; exist {
				idObtained[id] = cantId + 1
			} else {
				idObtained[id] = 1
			}
			esperar.Unlock()

		}()
	}
	wg.Wait()

	// fmt.Printf("total id requeridos: %v ob: %v\n", idRequired, len(idObtained))
	// fmt.Printf("%v", idObtained)
	if idRequired != len(idObtained) {
		t.Fatal("se esperaban:", idRequired, " ids pero se obtuvieron:", len(idObtained))
		t.Fail()
	}
}
