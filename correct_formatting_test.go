package unixid_test

import (
	"testing"
	"time"

	"github.com/cdvelop/unixid"
)

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

	// t.Log("✅ Timestamps están en orden cronológico correcto")
}
