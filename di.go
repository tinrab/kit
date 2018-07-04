package util

// DependencyInjection contains dependencies by name
type DependencyInjection struct {
	container map[string]interface{}
}

// NewDependencyInjection creates new instance of DependencyInjection
func NewDependencyInjection() *DependencyInjection {
	return &DependencyInjection{
		container: make(map[string]interface{}),
	}
}

// Provide registers a dependency
func (di *DependencyInjection) Provide(name string, dependency interface{}) {
	di.container[name] = dependency
}

// ProvideWith registers a dependency by calling a constructor function
func (di *DependencyInjection) ProvideWith(name string, constructor func(di *DependencyInjection) interface{}) {
	di.container[name] = constructor(di)
}

// Get returns a dependency by name
func (di *DependencyInjection) Get(name string) interface{} {
	return di.container[name]
}
