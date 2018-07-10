package kit

import (
	"reflect"
)

type Dependency interface {
	Open() error
	Close()
}

// DependencyInjection contains dependencies by name
type DependencyInjection struct {
	container map[string]Dependency
}

// NewDependencyInjection creates new instance of DependencyInjection
func NewDependencyInjection() *DependencyInjection {
	return &DependencyInjection{
		container: make(map[string]Dependency),
	}
}

// Provide registers a dependency
func (di *DependencyInjection) Provide(name string, dependency Dependency) {
	di.container[name] = dependency
}

// Get returns a dependency by name
func (di *DependencyInjection) Get(name string) Dependency {
	return di.container[name]
}

// Resolve decorates objects with dependencies and initializes them
func (di *DependencyInjection) Resolve() error {
	for _, dep := range di.container {
		di.inject(dep)
	}
	for _, dep := range di.container {
		err := dep.Open()
		if err != nil {
			return err
		}
	}
	return nil
}

// Close closes all dependencies
func (di *DependencyInjection) Close() {
	for _, dep := range di.container {
		dep.Close()
	}
}

func (di *DependencyInjection) inject(obj interface{}) {
	t := reflect.TypeOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		inject := field.Tag.Get("inject")
		if inject == "" {
			continue
		}
		dependency := di.Get(inject)
		reflect.ValueOf(obj).Elem().Field(i).Set(reflect.ValueOf(dependency))
	}
}
