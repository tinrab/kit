package di

import (
	"reflect"
)

type Dependency interface {
	Open() error
	Close()
}

// Container contains dependencies by name.
type Container struct {
	dependencies map[string]interface{}
}

// New creates new Container instance.
func New() *Container {
	return &Container{
		dependencies: make(map[string]interface{}),
	}
}

// Provide registers a dependency.
func (c *Container) Provide(name string, dependency interface{}) {
	c.dependencies[name] = dependency
}

// GetByName returns a dependency by name.
func (c *Container) GetByName(name string) interface{} {
	return c.dependencies[name]
}

// GetByType returns a dependency by type.
func (c *Container) GetByType(typ interface{}) interface{} {
	t := reflect.TypeOf(typ)
	for _, d := range c.dependencies {
		dt := reflect.TypeOf(d)
		if t.AssignableTo(dt) {
			return d
		}
	}
	return nil
}

// Resolve decorates objects with dependencies and initializes them.
func (c *Container) Resolve() error {
	for _, d := range c.dependencies {
		c.inject(d)
	}
	var opened []Dependency
	for _, d := range c.dependencies {
		if dep, ok := d.(Dependency); ok {
			err := dep.Open()
			if err != nil {
				for _, od := range opened {
					od.Close()
				}
				return err
			}
			opened = append(opened, dep)
		}
	}
	return nil
}

// Close closes all dependencies.
func (c *Container) Close() {
	for _, d := range c.dependencies {
		if dep, ok := d.(Dependency); ok {
			dep.Close()
		}
	}
}

func (c *Container) inject(obj interface{}) {
	t := reflect.TypeOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		inject := field.Tag.Get("inject")
		if inject != "" {
			dependency := c.GetByName(inject)
			if dependency != nil {
				reflect.ValueOf(obj).Elem().Field(i).Set(reflect.ValueOf(dependency))
			}
		}

		for _, d := range c.dependencies {
			dt := reflect.TypeOf(d)
			if field.Type.AssignableTo(dt) {
				reflect.ValueOf(obj).Elem().Field(i).Set(reflect.ValueOf(d))
			}
		}
	}
}
