package beans

// Scope is the strategy to hold bean instances in.
type Scope interface {
	// Get return the object with the given name from the underlying scope.
	// It will create it if not found in the underlying storage mechanism.
	Get(name string, objectFactory ObjectFactory) (any, error)

	// Remove the object from the underlying scope.
	Remove(name string) error
}

// ObjectFactory is a factory which can return an Object instance when invoked.
type ObjectFactory interface {
	// Return an instance (possibly shared or independent).
	GetObject() (any, error)
}
