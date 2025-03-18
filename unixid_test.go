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
		uid.GetNewID()
	}
}

// Prueba adicional para verificar que no haya duplicados al generar muchos IDs
func TestNoDuplicateIDs(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Generar una cantidad moderada de IDs y verificar que no haya duplicados
	numIDs := 1000
	ids := make(map[string]bool)

	for i := 0; i < numIDs; i++ {
		id := uid.GetNewID()
		if _, exists := ids[id]; exists {
			t.Fatalf("ID duplicado encontrado: %s", id)
		}
		ids[id] = true
	}
}

// Prueba para verificar que se generen IDs secuenciales cuando hay colisiones de timestamp
func TestSequentialIDs(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Generar varios IDs rápidamente, algunos tendrán el mismo timestamp base
	// pero deberían tener números secuenciales añadidos
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		ids[i] = uid.GetNewID()
	}

	// Verificar que tengamos al menos algunos IDs diferentes
	uniqueIDs := make(map[string]bool)
	for _, id := range ids {
		uniqueIDs[id] = true
	}

	if len(uniqueIDs) < len(ids) {
		t.Fatalf("Se esperaban %d IDs únicos, pero se obtuvieron %d", len(ids), len(uniqueIDs))
	}
}
