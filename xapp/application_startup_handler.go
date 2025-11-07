package xapp

import (
	"context"
)

// ApplicationStartupHandler is the hook handler when application starting.
type ApplicationStartupHandler interface {
	// Name return the handler's name
	Name() string

	// BeforeLoadProps is invoked before any props（env、flag、configfile）is retrieved.
	BeforeLoadProps(ctx context.Context, app *airmidApplication, opt *option) error

	// AfterLoadProps is invoked after all props has retrieved.
	AfterLoadProps(ctx context.Context, app *airmidApplication, opt *option) error

	// BeforeStartRunner is invoked before AppRunner starting.
	BeforeStartRunner(ctx context.Context, app *airmidApplication, opt *option) error
}
