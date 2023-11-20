package runtime

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl"
	"github.com/liuxd6825/components-contrib/pubsub"
	components_v1alpha1 "github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	"github.com/liuxd6825/dapr/pkg/outbox"

	// es "github.com/liuxd6825/dapr/pkg/components/liuxd/eventstorage"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	pubsub_adapter "github.com/liuxd6825/dapr/pkg/runtime/pubsub"
	"github.com/pkg/errors"
	"strings"
)

/*
func WithEventStorage(eventsourdings ...es.EventStorage) Option {
	return func(o *runtimeOpts) {
		o.eventStorages = append(o.eventStorages, eventsourdings...)
	}
}*/

func (a *DaprRuntime) initEventStorage(c components_v1alpha1.Component) error {
	const name = "EventStorage"
	//eventStorage, err := a.eventStorageRegistry.Create(c.Spec.Type, c.Spec.Version)

	eventStorage, err := a.runtimeConfig.registry.EventStorage().Create(c.Spec.Type, c.Spec.Version)

	if err != nil {
		log.Warnf("error creating pub sub %s (%s/%s): %s", &c.ObjectMeta.Name, c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "creation", name)
		return err
	}

	properties := a.convertMetadataItemsToProperties(c.Spec.Metadata)
	consumerID := strings.TrimSpace(properties["consumerID"])
	if consumerID == "" {
		consumerID = a.runtimeConfig.id
	}
	properties["consumerID"] = consumerID

	getAdapter := func() pubsub_adapter.Adapter {
		return a.getPublishAdapter()
	}

	var opts *eventstorage.Options
	switch c.Spec.Type {
	case mongo_impl.ComponentSpecMongo:
		opts, err = mongo_impl.NewMongoOptions(c.Name, eventStorage.GetLogger(), common.Metadata{Properties: properties}, getAdapter)
	case gorm_impl.ComponentSpecMySql:
		opts, err = gorm_impl.NewMySqlOptions(c.Name, eventStorage.GetLogger(), common.Metadata{Properties: properties}, getAdapter)
	default:
		err = errors.Errorf("%v 不支持的配置类型", c.Spec.Type)
	}

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init", name)
		return err
	}
	ctx := context.Background()

	if err = eventStorage.Init(ctx, opts); err != nil {
		return err
	}

	if err != nil {
		log.Warnf("error initializing pub sub %s/%s: %s", c.Spec.Type, c.Spec.Version, err)
		diag.DefaultMonitoring.ComponentInitFailed(c.Spec.Type, "init", name)
		return err
	}
	a.eventStorage = eventStorage
	diag.DefaultMonitoring.ComponentInitialized(c.Spec.Type)
	return nil
}

type pubSubAdapter struct {
	runtimeConfig *internalConfig
	outbox        func() outbox.Outbox
}

func NewPubSubAdapter(runtimeConfig *internalConfig, outbox func() outbox.Outbox) pubsub_adapter.Adapter {
	return &pubSubAdapter{
		runtimeConfig: runtimeConfig,
		outbox:        outbox,
	}
}

// GetPubSub is an adapter method to find a pubsub by name.
func (a *pubSubAdapter) GetPubSub(pubsubName string) (pubsub.PubSub, bool) {
	fun, ok := a.runtimeConfig.registry.PubSubs().GetPubSub(pubsubName)
	if ok {
		return fun(), ok
	}
	return nil, false
}

// Publish is an adapter method for the runtime to pre-validate publish requests
// And then forward them to the Pub/Sub component.
// This method is used by the HTTP and gRPC APIs.
func (a *pubSubAdapter) Publish(ctx context.Context, req *pubsub.PublishRequest) error {
	thepubsub, _ := a.GetPubSub(req.PubsubName)
	if thepubsub == nil {
		return pubsub_adapter.NotFoundError{PubsubName: req.PubsubName}
	}
	return thepubsub.Publish(ctx, req)
}

func (a *pubSubAdapter) BulkPublish(ctx context.Context, request *pubsub.BulkPublishRequest) (pubsub.BulkPublishResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *pubSubAdapter) Outbox() outbox.Outbox {
	return a.outbox()
}
