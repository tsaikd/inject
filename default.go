package inject

import (
	"sync"
)

var (
	defaultInjector Injector
	mut             sync.Mutex
)

func init() {
	mut.Lock()
	defer mut.Unlock()
	defaultInjector = New()
}

func Provide(provider interface{}) TypeMapper {
	mut.Lock()
	defer mut.Unlock()
	return defaultInjector.Provide(provider)
}

// Construct takes an object (structure or interface) given by a reference, takes its type and
// invokes appropriate constructor
// Uses default injector, see also (*Inject).Construct
func Construct(cr interface{}) error {
	mut.Lock()
	defer mut.Unlock()
	return defaultInjector.Construct(cr)
}

// Maps the concrete value of val to its dynamic type using reflect.TypeOf,
// It returns the TypeMapper registered in.
// Uses default injector, see also (*Inject).Map
func Map(val interface{}) TypeMapper {
	mut.Lock()
	defer mut.Unlock()
	return defaultInjector.Map(val)
}

func MapTo(val interface{}, ifacePtr interface{}) TypeMapper {
	mut.Lock()
	defer mut.Unlock()
	return defaultInjector.MapTo(val, ifacePtr)
}

func ConstructLater(creatable interface{}) error {
	mut.Lock()
	defer mut.Unlock()
	defaultInjector.ConstructLater(creatable)
	return nil
}

func FinishConstruct() error {
	return defaultInjector.FinishConstruct()
}
