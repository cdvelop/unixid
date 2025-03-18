package unixid_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

// TestExternalMutexUsage verifies that an external mutex can be passed to NewUnixID
// and that it is correctly used to avoid deadlocks when other libraries also use sync.Mutex
func TestExternalMutexUsage(t *testing.T) {
	// We create an external mutex that simulates being shared with another library
	externalMutex := &sync.Mutex{}

	// We create an instance of UnixID passing the external mutex
	uid, err := unixid.NewUnixID(externalMutex)
	if err != nil {
		t.Fatalf("Error creating UnixID with external mutex: %v", err)
		return
	}

	// We simulate a scenario where another library locks the mutex
	// and then our code uses it
	externalMutex.Lock()

	// We use a channel and a goroutine to check that GetNewID
	// blocks waiting for the external mutex to be released
	resultChan := make(chan string)
	timeoutChan := make(chan bool)

	go func() {
		// If the mutex is used correctly, this call will block
		// until we release the external mutex
		id := uid.GetNewID()
		resultChan <- id
	}()

	// We set up a timeout to detect if GetNewID doesn't block when it should
	go func() {
		// We wait a bit to ensure the previous goroutine had enough time to start
		timer := time.NewTimer(time.Millisecond * 100)
		<-timer.C
		timeoutChan <- true
	}()

	// We check if GetNewID blocked correctly waiting for the mutex
	select {
	case id := <-resultChan:
		// If we reach here before the timeout, it means that GetNewID did not wait for the mutex
		t.Fatalf("GetNewID did not wait for the external mutex. ID obtained: %s", id)
	case <-timeoutChan:
		// This is the expected behavior: GetNewID is blocking
	}

	// We release the external mutex
	externalMutex.Unlock()

	// Now GetNewID should complete
	select {
	case id := <-resultChan:
		// Expected behavior: an ID is generated after releasing the mutex
		if id == "" {
			t.Fatal("Generated an empty ID after releasing the mutex")
		}
	case <-time.After(time.Second):
		// If we reach here, GetNewID is still blocked even after releasing the external mutex
		t.Fatal("GetNewID is still blocked after releasing the external mutex")
	}

	// Additional test: verify that we can generate several IDs without issues
	for i := 0; i < 10; i++ {
		id := uid.GetNewID()
		if id == "" {
			t.Fatalf("Generated an empty ID on iteration %d", i)
		}
	}
}
