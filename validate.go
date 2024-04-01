package unixid

import (
	"errors"
	"strconv"
)

func validateID(new_id_in string) (id int64, err error) {
	var id_out string
	const e = "validateID "
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
