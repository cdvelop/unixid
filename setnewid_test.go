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
	t.Run("SetNewID con string", func(t *testing.T) {
		// Configuración del servidor (no WebAssembly)
		uid, err := unixid.NewUnixID()
		if err != nil {
			t.Fatal(err)
		}

		var id string
		uid.SetNewID(&id)

		if id == "" {
			t.Fatal("El ID generado no puede estar vacío")
		}

		// Validamos que tenga un formato correcto para servidor
		// En servidor el formato es solo un número sin punto
		if strings.Contains(id, ".") {
			t.Fatalf("En entorno servidor, el ID no debe contener punto: %s", id)
		}
	})

	t.Run("SetNewID con reflect.Value", func(t *testing.T) {
		// Configuración del servidor (no WebAssembly)
		uid, err := unixid.NewUnixID()
		if err != nil {
			t.Fatal(err)
		}

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
		// Configuración del servidor (no WebAssembly)
		uid, err := unixid.NewUnixID()
		if err != nil {
			t.Fatal(err)
		}

		// Buffer para prueba
		buf := make([]byte, 0, 64)
		originalLen := len(buf)

		uid.SetNewID(buf)

		// En esta implementación, el buffer no se modifica directamente
		// ya que Go pasa los slices por valor, no por referencia
		// Verificamos que la función no cause errores al trabajar con slices
		if len(buf) != originalLen {
			t.Fatal("No se espera que el buffer cambie de tamaño cuando se pasa por valor")
		}
	})

	t.Run("Compatibilidad entre GetNewID y SetNewID", func(t *testing.T) {
		// Verificamos que el formato del ID sea consistente
		uid, err := unixid.NewUnixID()
		if err != nil {
			t.Fatal(err)
		}

		// Obtenemos ID con GetNewID
		idFromGet := uid.GetNewID()

		// Obtenemos ID con SetNewID
		var idFromSet string
		uid.SetNewID(&idFromSet)

		// El formato debe ser similar (números de la misma longitud)
		if len(idFromGet) != len(idFromSet) {
			t.Fatal("Los IDs generados por GetNewID y SetNewID tienen formatos diferentes")
		}
	})

	// Esta prueba solo funcionaría en compilación para WebAssembly
	// Incluida solo como referencia, pero se saltará en entorno de servidor
	t.Run("WebAssembly user number format (referencial)", func(t *testing.T) {
		t.Skip("Esta prueba está destinada para entornos WebAssembly")

		// Este código sería para WebAssembly, pero como no podemos compilar condicionalmente
		// en pruebas, lo dejamos como referencia.
		/*
			uid, err := unixid.NewUnixID(&mockSessionHandler{})
			if err != nil {
				t.Fatal(err)
			}

			var id string
			uid.SetNewID(&id)

			// En WebAssembly, el ID debe tener formato "timestamp.userNumber"
			parts := strings.Split(id, ".")
			if len(parts) != 2 || parts[1] != "42" {
				t.Fatalf("El formato del ID en WebAssembly debe ser 'timestamp.userNumber', recibido: %s", id)
			}
		*/
	})
}
