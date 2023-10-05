package unixid_test

import (
	"sync"
	"testing"

	"github.com/cdvelop/unixid"
)

func Test_GetNewID(t *testing.T) {
	idRequired := 1000
	wg := sync.WaitGroup{}
	wg.Add(idRequired)

	// PG := sqlite.NewConnection("test.PG", false)

	uid := unixid.NewHandler()

	idObtained := make(map[string]int)
	var esperar sync.Mutex

	for i := 0; i < idRequired; i++ {
		go func() {
			defer wg.Done()
			id := uid.GetNewID()
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
		t.Fail()
	}
}
