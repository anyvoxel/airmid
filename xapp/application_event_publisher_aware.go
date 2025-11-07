package xapp

// ApplicationEventPublisherAware is to be implemented by beans that want to
// be aware of ApplicationEventPublisher.
type ApplicationEventPublisherAware interface {
	// SetApplicationEventPublisher will set the ApplicationEventPublisher this object runs in.
	SetApplicationEventPublisher(publisher ApplicationEventPublisher)
}
