package unixid_test

import (
	"testing"
	"time"

	. "github.com/cdvelop/tinystring"
	"github.com/cdvelop/unixid"
)

func TestUnixNanoToTime(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	// Test con timestamp conocido en la zona horaria local
	// Usar time.Unix para obtener el tiempo esperado en la zona local
	testUnixSeconds := int64(1624397134) // 2021-06-22 15:32:14 UTC
	expectedTime := time.Unix(testUnixSeconds, 0)
	expected := expectedTime.Format("15:04:05")

	nanoTimestamp := testUnixSeconds * 1e9 // convertir a nanosegundos

	result := uid.UnixNanoToTime(nanoTimestamp)
	if result != expected {
		t.Errorf("UnixNanoToTime(%d) = %s; want %s", nanoTimestamp, result, expected)
	}

	// Test con string
	result = uid.UnixNanoToTime(Fmt("%d", nanoTimestamp))
	if result != expected {
		t.Errorf("UnixNanoToTime(string) = %s; want %s", result, expected)
	}

	// Test con timestamps secuenciales para verificar orden
	now := time.Now()
	baseNano := now.UnixNano()

	var results []string
	for i := 0; i < 3; i++ {
		nano := baseNano + int64(i)*int64(time.Second) // Incrementar 1 segundo
		timeStr := uid.UnixNanoToTime(nano)
		results = append(results, timeStr)
		t.Logf("Nano: %d -> Time: %s", nano, timeStr)
	}

	// Verificar que los tiempos están en orden
	for i := 1; i < len(results); i++ {
		if results[i] <= results[i-1] {
			t.Errorf("Los timestamps no están en orden: %s <= %s", results[i], results[i-1])
		}
	}
}

func TestUnixNanoToTimeWithDifferentTypes(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	now := time.Now()
	nanoTimestamp := now.UnixNano()

	// Test con diferentes tipos de entrada
	testCases := []struct {
		name  string
		input any
	}{
		{"int64", nanoTimestamp},
		{"int", int(nanoTimestamp)},
		{"float64", float64(nanoTimestamp)},
		{"string", Fmt("%d", nanoTimestamp)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := uid.UnixNanoToTime(tc.input)

			if result == "" {
				t.Errorf("UnixNanoToTime devolvió string vacío para tipo %s", tc.name)
			}

			t.Logf("Tipo %s: %v -> %s", tc.name, tc.input, result)
		})
	}

	// Test con tipo no soportado
	invalidResult := uid.UnixNanoToTime(make(chan int))
	if invalidResult != "" {
		t.Error("UnixNanoToTime debería devolver string vacío para tipos no soportados")
	}
}

// TestGetNewIDWithCorrectFormatting prueba el flujo completo como lo usa devtui
func TestGetNewIDWithCorrectFormatting(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	// Simular el flujo de devtui con el método correcto
	var timestamps []string

	for i := 0; i < 3; i++ {
		id := uid.GetNewID() // Devuelve nanosegundos como string

		// CORRECTO: Usar UnixNanoToTime para nanosegundos
		timeStr := uid.UnixNanoToTime(id)
		timestamps = append(timestamps, timeStr)

		t.Logf("Mensaje %d - ID: %s -> Time: %s", i+1, id, timeStr)

		// Pausa de 1 segundo para garantizar diferencias de tiempo visibles
		time.Sleep(1 * time.Second)
	}

	// Verificar orden cronológico (permitir timestamps iguales si están muy cerca)
	for i := 1; i < len(timestamps); i++ {
		if timestamps[i] < timestamps[i-1] {
			t.Errorf("Los timestamps NO están en orden cronológico: %s < %s",
				timestamps[i], timestamps[i-1])
		}
	}

	t.Log("✅ Timestamps están en orden cronológico correcto")
}
