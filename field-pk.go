package unixid

import . "github.com/cdvelop/tinystring"

// FieldType determines if a field is an ID field and/or a primary key field.
// This function analyzes field names to identify ID fields and primary keys based on naming conventions.
//
// Parameters:
//   - tableName: The name of the table or entity that the field belongs to
//   - fieldName: The name of the field to analyze
//
// Returns:
//   - ID: true if the field is an ID field (starts with "id")
//   - PK: true if the field is a primary key (is named "id" or matches the pattern "id{tableName}" or "id_{tableName}")
//
// Examples:
//   - FieldType("user", "id") returns (true, true)
//   - FieldType("user", "iduser") returns (true, true)
//   - FieldType("user", "id_user") returns (true, true)
//   - FieldType("user", "idaddress") returns (true, false)
func (u *UnixID) FieldType(tableName, fieldName string) (ID, PK bool) {
	if len(fieldName) >= 2 {
		key_name := Convert(fieldName).ToLower().String()

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
		if Contains(key_name, prefixNameID) {
			key_without_id = Convert(key_name).Replace(prefixNameID, "").String() //remove _
		} else {
			key_without_id = key_name[2:] //remove id
		}

		if key_without_id == tableName {
			PK = true
		}
	}
	return
}
