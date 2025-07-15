package unixid

import "github.com/cdvelop/tinyreflect"

// SetNewID sets a new unique ID value to various types of targets.
// It generates a new unique ID based on Unix nanosecond timestamp and assigns it to the provided target.
// This function can work with multiple target types including tinyreflect.Value, string pointers, and byte slices.
//
// In WebAssembly environments, IDs include a user session number as a suffix (e.g., "1624397134562544800.42").
// In server environments, IDs are just the timestamp (e.g., "1624397134562544800").
//
// Parameters:
//   - target: The target to receive the new ID. Can be:
//   - tinyreflect.Value or *tinyreflect.Value: For setting struct field values via tiny-reflection.
//   - *string: For setting a string variable directly.
//   - []byte: For appending the ID to a byte slice.
//
// This function is thread-safe in server-side environments.
//
// Examples:
//
//	// Setting a struct field using tiny-reflection
//	v := tinyreflect.ValueOf(&myStruct)
//	elem, _ := v.Elem()
//	field, _ := elem.Field(0) // Get field by index
//	idHandler.SetNewID(field)
//
//	// Setting a string variable
//	var id string
//	idHandler.SetNewID(&id)
//
//	// Appending to a byte slice
//	buf := make([]byte, 0, 64)
//	idHandler.SetNewID(buf)
func (id *UnixID) SetNewID(target any) {
	// Generate a new ID - this already has locking when necessary
	newID := id.GetNewID()

	// Set the ID to the appropriate target type
	switch t := target.(type) {
	case tinyreflect.Value:
		// For struct fields via reflection
		t.SetString(newID)
	case *tinyreflect.Value:
		// For struct fields via reflection (pointer)
		t.SetString(newID)
	case *string:
		// For string variables
		*t = newID
	case []byte:
		// For byte slices, we need to return the new slice
		// but since we can't modify the original slice reference,
		// we just copy the ID bytes into it if there's enough capacity
		if len(t) >= len(newID) {
			copy(t, newID)
		} else if cap(t) >= len(newID) {
			copy(t[:cap(t)], newID)
		}
	}
}
