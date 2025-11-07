package beans

import (
	"reflect"
)

// BeanPriority is the priority of bean, more higher the priority order is, the index of a bean in the slice/array
// smaller is.
type BeanPriority int32

// BeanPriorityOrder define the priority order interface.
type BeanPriorityOrder interface {
	// GetPriority return the priority of a bean, more higher the priority order is,
	// the index of a bean in the slice/array smaller is.
	GetPriority() BeanPriority
}

// CandidateBeans is a slice type of reflect.Value.
type CandidateBeans []reflect.Value

func (c CandidateBeans) Len() int {
	return len(c)
}

// Less return true in following conditions:
// 1. c[j] implements BeanPriorityOrder, but c[i] not,
// 2. c[i] and c[j] implements BeanPriorityOrder, c[j]'s priority > c[i]'s priority.
func (c CandidateBeans) Less(i, j int) bool {
	vi := IndirectTo[BeanPriorityOrder](c[i].Interface())
	vj := IndirectTo[BeanPriorityOrder](c[j].Interface())

	if vi != nil && vj == nil {
		return true
	}
	if vi != nil && vj != nil {
		return vi.GetPriority() > vj.GetPriority()
	}

	return false
}

func (c CandidateBeans) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
