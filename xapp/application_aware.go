package xapp

// ApplicationAware is to implement by beans for be aware of the capabilities of airmid application.
type ApplicationAware interface {
	// SetApplication is used to set bean's Application field
	SetApplication(application Application)
}
