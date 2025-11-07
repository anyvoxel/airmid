// Package props provide the access to user defined property value
package props

// Properties hold the configurable properties name & value.
type Properties interface {
	// Get return the value for key, the key is case-sensitive.
	Get(key string, opts ...GetOption) (any, error)

	// Set the value for key, it will overwrite the old value if key is already exists.
	// The val will been transform to string to store with the container, so when
	//   user try to get the val of key, the string returned.
	Set(key string, val any) error
}
