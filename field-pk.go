package unixid

import (
	"strings"
)

func (u *UnixID) FieldType(tableName, fieldName string) (ID, PK bool) {
	if len(fieldName) >= 2 {

		key_name := strings.ToLower(fieldName)

		if key_name[:2] != "id" {
			return
		} else {
			ID = true
		}

		if key_name == "id" {
			PK = true
			return
		}

		var key_without_id string
		if strings.Contains(key_name, PrefixNameID) {

			key_without_id = strings.Replace(key_name, PrefixNameID, "", 1) //remover _
		} else {

			key_without_id = key_name[2:] //remover id
		}

		if key_without_id == tableName {
			PK = true
		}

	}
	return
}
