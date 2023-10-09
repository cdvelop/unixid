package unixid

import (
	"strconv"

	"github.com/cdvelop/model"
)

func (pk) ValidateField(data_in string, skip_validation bool, options ...string) error {
	if !skip_validation {

		_, err := validateID(data_in)

		return err
	}
	return nil
}

func validateID(new_id_in string) (int64, error) {
	var id_out string

	var point_count int
	var point_index = len(new_id_in)
	for i, char := range new_id_in {
		if char == '.' {
			point_count++
			point_index = i
			if point_count > 1 {
				return 0, model.Error("error id contiene más de un punto")
			}
		} else if char < '0' || char > '9' {
			return 0, model.Error("error id contiene caracteres no válidos")
		}
	}

	id_out = new_id_in[:point_index]

	// fmt.Println("ID SALIDA:", id_out)

	id, err := strconv.ParseInt(id_out, 10, 64)
	if err != nil {
		return 0, model.Error("error id contiene caracteres no válidos")
	}

	return id, nil
}
