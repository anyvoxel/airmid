package xapp

import "github.com/anyvoxel/airmid/beans"

// ApplicationAwareProcessor implement BeanPostProcessor to automatically
// invoke aware interface.
type ApplicationAwareProcessor struct {
	app *airmidApplication
}

// PostProcessBeforeInitialization implement the BeanPostProcessor.PostProcessBeforeInitialization.
func (p *ApplicationAwareProcessor) PostProcessBeforeInitialization(obj any, _ string) (v any, err error) {
	if vobj := beans.IndirectTo[ApplicationEventPublisherAware](obj); vobj != nil {
		vobj.SetApplicationEventPublisher(p.app)
	}

	if vobj := beans.IndirectTo[ApplicationAware](obj); vobj != nil {
		vobj.SetApplication(p.app)
	}

	return obj, nil
}

// PostProcessAfterInitialization implement the BeanPostProcessor.PostProcessAfterInitialization.
func (*ApplicationAwareProcessor) PostProcessAfterInitialization(obj any, _ string) (v any, err error) {
	return obj, nil
}
