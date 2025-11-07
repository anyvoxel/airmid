package beans

// BeanNameAware is to be implemented by beans that want to be aware of their bean name
// in a bean factory.
type BeanNameAware interface {
	// SetBeanName willl set the name of the bean in bean factory that created this bean.
	// This will invoked after population of normal bean properties but before an init callback
	// such as InitializingBean.AfterPropertiesSet.
	SetBeanName(name string)
}
