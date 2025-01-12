package registry

import (
	"github.com/dapr/dapr/pkg/components/liuxd/applogger"
	"github.com/dapr/dapr/pkg/components/liuxd/eventstore"
)

// WithEventStore adds event storage to the runtime.
func (o *Options) WithEventStore(registry *eventstore.Registry) *Options {
	o.eventStore = registry
	return o
}

// WithAppLogger adds app logger to the runtime.
func (o *Options) WithAppLogger(registry *applogger.Registry) *Options {
	o.appLogger = registry
	return o
}
