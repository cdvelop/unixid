package unixid

import (
	"errors"
	"strconv"
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
	const e = "ValidateID "
	const msg = "id contiene caracteres no válidos"

	var point_count int
	var point_index = len(new_id_in)
	for i, char := range new_id_in {
		if char == '.' {
			point_count++
			point_index = i
			if point_count > 1 {
				return 0, errors.New(e + "id contiene más de un punto")
			}
		} else if char < '0' || char > '9' {
			return 0, errors.New(e + msg)
		}
	}

	id_out = new_id_in[:point_index]

	// fmt.Println("ID SALIDA:", id_out)

	id, er := strconv.ParseInt(id_out, 10, 64)
	if er != nil {
		return 0, errors.New(e + msg)
	}

	return id, nil
}
