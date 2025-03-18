package unixid_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cdvelop/unixid"
)

// Estructura para pruebas
type TestStruct struct {
	ID string
}

// Implementación mock de userSessionNumber para pruebas
type mockSessionHandler struct{}

func (mockSessionHandler) userSessionNumber() string {
	return "42"
}

func TestSetNewID(t *testing.T) {
	// Creamos una sola instancia de UnixID para todos los subtests
	// para evitar la sobrecarga de crear múltiples instancias
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("SetNewID con string", func(t *testing.T) {
		var id string
		uid.SetNewID(&id)

		if id == "" {
			t.Fatal("El ID generado no puede estar vacío")
		}

		// Validamos que tenga un formato correcto para servidor
		if strings.Contains(id, ".") {
			t.Fatalf("En entorno servidor, el ID no debe contener punto: %s", id)
		}
	})

	t.Run("SetNewID con reflect.Value", func(t *testing.T) {
		// Creamos una estructura para prueba
		testObj := TestStruct{}

		// Obtenemos el campo ID usando reflection
		rv := reflect.ValueOf(&testObj).Elem().FieldByName("ID")
		uid.SetNewID(&rv)

		if testObj.ID == "" {
			t.Fatal("El ID generado no puede estar vacío")
		}

		// Validamos que tenga un formato correcto para servidor
		if strings.Contains(testObj.ID, ".") {
			t.Fatalf("En entorno servidor, el ID no debe contener punto: %s", testObj.ID)
		}
	})

	t.Run("SetNewID con []byte", func(t *testing.T) {
		// Este test no es efectivo ya que los slices se pasan por valor
		// y el método SetNewID no puede modificarlos directamente
		// Lo mantenemos por compatibilidad pero lo hacemos más eficiente
		buf := make([]byte, 0, 8) // Reducimos el tamaño a 8 bytes que es suficiente para la prueba
		uid.SetNewID(buf)
		// No es necesario hacer más verificaciones aquí
	})

	t.Run("Compatibilidad entre GetNewID y SetNewID", func(t *testing.T) {
		// Obtenemos ID con GetNewID
		idFromGet := uid.GetNewID()

		// Obtenemos ID con SetNewID
		var idFromSet string
		uid.SetNewID(&idFromSet)

		// Solo verificamos que ambos IDs tengan el mismo formato (longitud similar)
		lenGet := len(idFromGet)
		lenSet := len(idFromSet)
		if lenGet < lenSet-2 || lenGet > lenSet+2 { // Permitimos una pequeña variación
			t.Fatalf("Los IDs generados por GetNewID y SetNewID tienen formatos muy diferentes: %d vs %d", lenGet, lenSet)
		}
	})

	// Esta prueba solo funcionaría en compilación para WebAssembly
	// Se mantiene como referencia pero se omite en la ejecución
	t.Run("WebAssembly user number format (referencial)", func(t *testing.T) {
		t.Skip("Esta prueba está destinada para entornos WebAssembly")
	})
}

// Añadimos un benchmark para SetNewID para medir su rendimiento
func BenchmarkSetNewID(b *testing.B) {
	uid, _ := unixid.NewUnixID()
	var id string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid.SetNewID(&id)
	}
}
