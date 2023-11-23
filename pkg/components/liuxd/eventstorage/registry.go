package eventstorage

import (
	"github.com/dapr/kit/logger"
	"strings"

	es "github.com/liuxd6825/dapr-components-contrib/liuxd/eventstorage"
	"github.com/pkg/errors"
)

type (
	EventStoreFactory struct {
		SpaceType     string
		Version       string
		FactoryMethod func(logger logger.Logger) es.EventStorage
	}

	Registry struct {
		Logger    logger.Logger
		factories map[string]*EventStoreFactory
	}
)

// DefaultRegistry is the singleton with the registry.
var DefaultRegistry *Registry = NewRegistry()

func New(spaceType string, factoryMethod func(logger logger.Logger) es.EventStorage) EventStoreFactory {
	return EventStoreFactory{
		SpaceType:     spaceType,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]*EventStoreFactory),
	}
}

// Register registers one or more new message buses.
func (p *Registry) Register(newFactory func(logger.Logger) es.EventStorage, spaceType string, version string) {
	p.factories[createFullName(spaceType, version)] = &EventStoreFactory{
		SpaceType:     spaceType,
		Version:       version,
		FactoryMethod: newFactory,
	}
}

// Create instantiates a EventStorage based on `name`.
func (p *Registry) Create(compName string, spaceType string, version string) (es.EventStorage, error) {
	if factory, ok := p.getFactory(spaceType, version); ok {
		return factory.FactoryMethod(p.Logger), nil
	}
	return nil, errors.Errorf("couldn't find eventStore %s/%s", spaceType, version)
}

func (p *Registry) getFactory(name, version string) (*EventStoreFactory, bool) {
	key := createFullName(name, version)
	factory, ok := p.factories[key]
	return factory, ok
}

func createFullName(name string, version string) string {
	return strings.ToLower(name) + "/" + strings.ToLower(version)
}
