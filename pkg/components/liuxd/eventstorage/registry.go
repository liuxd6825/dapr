package eventstorage

import (
	"github.com/dapr/kit/logger"
	"strings"

	es "github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/dapr/pkg/components"
	"github.com/pkg/errors"
)

type (
	EventStorage struct {
		Name          string
		FactoryMethod func(logger logger.Logger) es.EventStorage
	}

	Registry struct {
		Logger       logger.Logger
		messageBuses map[string]func(logger logger.Logger) es.EventStorage
	}
)

// DefaultRegistry is the singleton with the registry.
var DefaultRegistry *Registry = NewRegistry()

func New(name string, factoryMethod func(logger logger.Logger) es.EventStorage) EventStorage {
	return EventStorage{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() *Registry {
	return &Registry{
		messageBuses: map[string]func(logger logger.Logger) es.EventStorage{},
	}
}

// Register registers one or more new message buses.
func (p *Registry) Register(components ...EventStorage) {
	for _, component := range components {
		p.messageBuses[createFullName(component.Name)] = component.FactoryMethod
	}
}

func (p *Registry) RegisterComponent(factoryMethod func(logger logger.Logger) es.EventStorage, name string) {
	p.messageBuses[createFullName(name)] = factoryMethod
}

// Create instantiates a pub/sub based on `name`.
func (p *Registry) Create(typeName string, version string) (es.EventStorage, error) {
	if method, ok := p.getEventStorage(typeName, version); ok {
		return method(p.Logger), nil
	}
	return nil, errors.Errorf("couldn't find message bus %s/%s", typeName, version)
}

func (p *Registry) getEventStorage(typeName, version string) (func(logger logger.Logger) es.EventStorage, bool) {
	nameLower := strings.ToLower(typeName)
	versionLower := strings.ToLower(version)
	pubSubFn, ok := p.messageBuses[nameLower+"/"+versionLower]
	if ok {
		return pubSubFn, true
	}
	if components.IsInitialVersion(versionLower) {
		pubSubFn, ok = p.messageBuses[nameLower]
	}
	return pubSubFn, ok
}

func createFullName(name string) string {
	return strings.ToLower("eventstorage." + name)
}
