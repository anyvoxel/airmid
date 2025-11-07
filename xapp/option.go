package xapp

// Option applies a configuration option value to a application.
type Option interface {
	apply(*option)
}

// option contains configuration options for a application.
type option struct {
	attrs []Attribute
}

// optionFunc applies a set of options to a option.
type optionFunc func(*option)

// apply the function with a option.
func (f optionFunc) apply(o *option) {
	f(o)
}

// Attribute holds a key and value pair.
type Attribute struct {
	Key   string
	Value string
}

// WithAttributes sets the attrs.
func WithAttributes(attrs ...Attribute) Option {
	return optionFunc(func(o *option) {
		o.attrs = append(o.attrs, attrs...)
	})
}

func newOption(options []Option) *option {
	o := &option{}
	for _, opt := range options {
		opt.apply(o)
	}

	return o
}
