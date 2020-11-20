package workingon

import (
	"errors"
	"fmt"
)

type Source interface {
	Configure(config *Config) error

	GetName() string
	GetTask(key string) (*Task, error)
	GetTasks() ([]Task, error)
	GetProjects() ([]Project, error)
}

type registry struct {
	RegisteredSources []Source
}

var (
	Registry registry
)

func (r *registry) Register(source Source) error {
	r.RegisteredSources = append(r.RegisteredSources, source)
	return nil
}

func (r *registry) GetNames() []string {
	var names []string
	for _, source := range r.RegisteredSources {
		names = append(names, source.GetName())
	}
	return names
}

func (r *registry) GetTask(key string) (*Task, error) {
	for _, source := range r.RegisteredSources {
		t, err := source.GetTask(key)
		if err == nil {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Task with key %s not found in any source.", key))
}
