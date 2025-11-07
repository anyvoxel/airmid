package beans

// BeanFactoryAware is to be implemented by beans that wish to
// be aware of their owning.
type BeanFactoryAware interface {
	// SetBeanFactory will set the owning factory to a bean instance.
	SetBeanFactory(beanFactory BeanFactory)
}
