package unixid

import (
	"strconv"
)

func (pk) ValidateField(data_in string, skip_validation bool, options ...string) (err string) {
	if !skip_validation {

		_, err := validateID(data_in)

		return err
	}
	return ""
}

func validateID(new_id_in string) (id int64, err string) {
	var id_out string
	const this = "validateID error "
	const msg = "id contiene caracteres no válidos"

	var point_count int
	var point_index = len(new_id_in)
	for i, char := range new_id_in {
		if char == '.' {
			point_count++
			point_index = i
			if point_count > 1 {
				return 0, this + "id contiene más de un punto"
			}
		} else if char < '0' || char > '9' {
			return 0, this + msg
		}
	}

	id_out = new_id_in[:point_index]

	// fmt.Println("ID SALIDA:", id_out)

	id, e := strconv.ParseInt(id_out, 10, 64)
	if e != nil {
		return 0, this + msg
	}

	return id, ""
}
