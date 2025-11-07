// Package pointer provide the helper function to convert pointer
package pointer

// StringPtr return the pointer to v.
func StringPtr(v string) *string {
	return &v
}

// IntPtr return the pointer to v.
func IntPtr(v int) *int {
	return &v
}
