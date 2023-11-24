package eventstore

import (
	"context"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/common"
	contrib "github.com/liuxd6825/dapr-components-contrib/liuxd/eventstore"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstore/impl/gorm_impl"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstore/impl/mongo_impl"
	compapi "github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	components "github.com/liuxd6825/dapr/pkg/components/liuxd/eventstore"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	"github.com/liuxd6825/dapr/pkg/encryption"
	"github.com/liuxd6825/dapr/pkg/runtime/compstore"
	rterrors "github.com/liuxd6825/dapr/pkg/runtime/errors"
	"github.com/liuxd6825/dapr/pkg/runtime/meta"
	pubsub_adapter "github.com/liuxd6825/dapr/pkg/runtime/pubsub"

	"github.com/pkg/errors"
	"io"
	"sync"
)

const (
	propertyKeyActorStateStore = "actorstatestore"
)

var log = logger.NewLogger("dapr.runtime.processor.eventStorage")

type Options struct {
	Logger           logger.Logger
	Registry         *components.Registry
	ComponentStore   *compstore.ComponentStore
	Meta             *meta.Meta
	PubsubAdapter    pubsub_adapter.Adapter
	PlacementEnabled bool
}

type eventStorage struct {
	registry            *components.Registry
	compStore           *compstore.ComponentStore
	meta                *meta.Meta
	lock                sync.RWMutex
	pubsubAdapter       pubsub_adapter.Adapter
	actorStateStoreName *string
	placementEnabled    bool
}

func New(opts Options) *eventStorage {
	return &eventStorage{
		registry:         opts.Registry,
		compStore:        opts.ComponentStore,
		meta:             opts.Meta,
		placementEnabled: opts.PlacementEnabled,
		pubsubAdapter:    opts.PubsubAdapter,
	}
}

func (s *eventStorage) Init(ctx context.Context, comp compapi.Component) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	fName := comp.LogName()
	store, err := s.registry.Create(comp.Name, comp.Spec.Type, comp.Spec.Version)
	if err != nil {
		diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "creation", comp.ObjectMeta.Name)
		return rterrors.NewInit(rterrors.CreateComponentFailure, fName, err)
	}

	if store != nil {
		secretStoreName := s.meta.AuthSecretStoreOrDefault(&comp)
		secretStore, _ := s.compStore.GetSecretStore(secretStoreName)
		encKeys, encErr := encryption.ComponentEncryptionKey(comp, secretStore)
		if encErr != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "creation", comp.ObjectMeta.Name)
			return rterrors.NewInit(rterrors.CreateComponentFailure, fName, err)
		}

		if encKeys.Primary.Key != "" {
			ok := encryption.AddEncryptedStateStore(comp.ObjectMeta.Name, encKeys)
			if ok {
				log.Infof("automatic encryption enabled for state store %s", comp.ObjectMeta.Name)
			}
		}

		meta, err := s.meta.ToBaseMetadata(comp)
		if err != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "init", comp.ObjectMeta.Name)
			return rterrors.NewInit(rterrors.InitComponentFailure, fName, err)
		}

		getAdapter := func() pubsub_adapter.Adapter {
			return s.pubsubAdapter
		}

		props := meta.Properties
		var ops *contrib.Options
		switch comp.Spec.Type {
		case mongo_impl.ComponentSpecMongo:
			ops, err = mongo_impl.NewMongoOptions(comp.Name, s.registry.Logger, common.Metadata{Properties: props}, getAdapter)
		case gorm_impl.ComponentSpecMySql:
			ops, err = gorm_impl.NewMySqlOptions(comp.Name, s.registry.Logger, common.Metadata{Properties: props}, getAdapter)
		default:
			err = errors.Errorf("%v, %v, %v 不支持的配置类型", comp.Spec.Type, comp.Spec.Version, comp.Name)
			return err
		}
		if err != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "init", comp.ObjectMeta.Name)
			return rterrors.NewInit(rterrors.InitComponentFailure, fName, err)
		}

		err = store.Init(ctx, ops)
		if err != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "init", comp.ObjectMeta.Name)
			return rterrors.NewInit(rterrors.InitComponentFailure, fName, err)
		}

		s.compStore.AddEventStore(comp.ObjectMeta.Name, store)

		diag.DefaultMonitoring.ComponentInitialized(comp.Spec.Type)
	}

	return nil
}

func (s *eventStorage) Close(comp compapi.Component) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ss, ok := s.compStore.GetEventStore(comp.Name)
	if !ok {
		return nil
	}

	closer, ok := ss.(io.Closer)
	if ok && closer != nil {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	s.compStore.DeleteStateStore(comp.Name)

	return nil
}

func (s *eventStorage) ActorStateStoreName() (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.actorStateStoreName == nil {
		return "", false
	}
	return *s.actorStateStoreName, true
}
