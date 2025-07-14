package unixid

import (
	. "github.com/cdvelop/tinystring"
)

// ValidateID validates and parses a Unix timestamp ID string.
// It handles IDs in both server format (just timestamp) and client format (timestamp.userNumber).
// This function extracts the timestamp portion from the ID and returns it as an int64.
//
// Parameters:
//   - new_id_in: The ID string to validate (e.g., "1624397134562544800" or "1624397134562544800.42")
//
// Returns:
//   - id: The timestamp portion of the ID as an int64 value
//   - err: An error if the ID format is invalid
//
// Validation rules:
//   - The ID must contain only digits and at most one decimal point
//   - The timestamp portion (before the decimal point) must be a valid int64
func ValidateID(new_id_in string) (id int64, err error) {
	var id_out string
	// Mensajes usando el diccionario multilenguaje
	msg_invalid := Err(D.Character, D.Invalid, D.Not, D.Supported)

	if len(new_id_in) == 0 {
		return 0, msg_invalid
	}

	// No debe comenzar ni terminar con punto
	if new_id_in[0] == '.' || new_id_in[len(new_id_in)-1] == '.' {
		return 0, msg_invalid
	}

	var point_count int
	var point_index = len(new_id_in)
	for i, char := range new_id_in {
		if char == '.' {
			point_count++
			point_index = i
			if point_count > 1 {
				return 0, Err(D.Invalid, D.Format, D.Found, D.More, D.Point)
			}
		} else if char < '0' || char > '9' {
			return 0, msg_invalid
		}
	}

	id_out = new_id_in[:point_index]

	id, er := Convert(id_out).Int64()
	if er != nil {
		return 0, msg_invalid
	}

	return id, nil
}
