package beans

// InitializingBean is to be implemented by beans that need to react once all their properties
// have been set.
type InitializingBean interface {
	// AfterPropertiesSet is invoked after all properties have been set.
	AfterPropertiesSet() error
}

// SmartInitializingSingleton is to be implemented by beans that need to react once the single (non-lazy)
// has been pre instantiated.
type SmartInitializingSingleton interface {
	AfterSingletonsInstantiated() error
}
