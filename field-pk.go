package unixid

import . "github.com/cdvelop/tinystring"

// FieldType determines if a field is an ID field and/or a primary key field.
// This function delegates to tinystring.IsIDorPrimaryKey for unified logic.
//
// Parameters:
//   - tableName: The name of the table or entity that the field belongs to
//   - fieldName: The name of the field to analyze
//
// Returns:
//   - ID: true if the field is an ID field (starts with "id")
//   - PK: true if the field is a primary key based on naming conventions
//
// Examples:
//   - FieldType("user", "id") returns (true, true)
//   - FieldType("user", "iduser") returns (true, true)
//   - FieldType("user", "id_user") returns (true, true)
//   - FieldType("user", "idaddress") returns (true, false)
//   - FieldType("user", "userid") returns (true, true)
//   - FieldType("user", "user_id") returns (true, true)
func (u *UnixID) FieldType(tableName, fieldName string) (ID, PK bool) {
	return IsIDorPrimaryKey(tableName, fieldName)
}
