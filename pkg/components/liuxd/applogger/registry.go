package applogger

import (
	"fmt"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/applog"
	"github.com/liuxd6825/dapr/pkg/components"
	"github.com/pkg/errors"
	"strings"
)

type (
	Logger struct {
		Name          string
		FactoryMethod func() applog.Logger
	}

	Registry struct {
		Logger       logger.Logger
		messageBuses map[string]func(log logger.Logger) applog.Logger
	}
)

// DefaultRegistry is the singleton with the registry.
var DefaultRegistry *Registry = NewRegistry()

func New(name string, factoryMethod func() applog.Logger) Logger {
	return Logger{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

func NewRegistry() *Registry {
	return &Registry{
		messageBuses: map[string]func(log logger.Logger) applog.Logger{},
	}
}

// Register registers one or more new message buses.
func (p *Registry) Register(factoryMethod func(log logger.Logger) applog.Logger, name string, version string) {
	key := createFullName(name, version)
	fmt.Println(key)
	p.messageBuses[key] = factoryMethod
}

// Create instantiates a applogger on `name`.
func (p *Registry) Create(typename, version string) (applog.Logger, error) {
	fmt.Println(typename)
	if method, ok := p.getAppLoggerSourcing(typename, version); ok {
		return method(p.Logger), nil
	}
	return nil, errors.Errorf("couldn't find %s/%s", typename, version)
}

func (p *Registry) getAppLoggerSourcing(typename, version string) (func(log logger.Logger) applog.Logger, bool) {
	nameLower := strings.ToLower(typename)
	versionLower := strings.ToLower(version)
	key := nameLower + "/" + versionLower
	fmt.Print(key)
	pubSubFn, ok := p.messageBuses[key]
	if ok {
		return pubSubFn, true
	}
	if components.IsInitialVersion(versionLower) {
		pubSubFn, ok = p.messageBuses[nameLower]
	}
	return pubSubFn, ok
}

func createFullName(name string, version string) string {
	return strings.ToLower(fmt.Sprintf("applogger.%v/%v", name, version))
}
