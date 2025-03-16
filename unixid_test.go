package unixid_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/cdvelop/unixid"
)

func Test_GetNewID(t *testing.T) {
	idRequired := 10000
	wg := sync.WaitGroup{}
	wg.Add(idRequired)

	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	idObtained := make(map[string]int)
	var esperar sync.Mutex

	for i := 0; i < idRequired; i++ {
		go func() {
			defer wg.Done()

			id, err := uid.GetNewID()
			if err != nil {
				t.Log(err)
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
	if idRequired != len(idObtained) {
		fmt.Printf("%v", idObtained)
		t.Fatal("se esperaban:", idRequired, " ids pero se obtuvieron:", len(idObtained))
		t.Fail()
	}

}

func BenchmarkGetNewID(b *testing.B) {

	uid, _ := unixid.NewUnixID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		// var buf = []
		// buf.Grow(20)
		uid.GetNewID()
	}

}
