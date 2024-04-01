package unixid

import (
	"strings"
)

func IdPkTypeField(field_name, table_name string) (fieldTypeID, fieldTypePK bool) {
	if len(field_name) >= 2 {

		key_name := strings.ToLower(field_name)

		if key_name[:2] != "id" {
			return
		} else {
			fieldTypeID = true
		}

		if key_name == "id" {
			fieldTypePK = true
			return
		}

		var key_without_id string
		if strings.Contains(key_name, PrefixNameID) {

			key_without_id = strings.Replace(key_name, PrefixNameID, "", 1) //remover _
		} else {

			key_without_id = key_name[2:] //remover id
		}

		if key_without_id == table_name {
			fieldTypePK = true
		}

	}
	return
}
