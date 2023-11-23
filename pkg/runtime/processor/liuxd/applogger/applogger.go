package applogger

import (
	"context"
	"fmt"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/common"
	compapi "github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	comps "github.com/liuxd6825/dapr/pkg/components/liuxd/applogger"
	compstate "github.com/liuxd6825/dapr/pkg/components/state"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	"github.com/liuxd6825/dapr/pkg/encryption"
	"github.com/liuxd6825/dapr/pkg/runtime/compstore"
	rterrors "github.com/liuxd6825/dapr/pkg/runtime/errors"
	"github.com/liuxd6825/dapr/pkg/runtime/meta"
	"github.com/liuxd6825/dapr/pkg/runtime/pubsub"
	"github.com/liuxd6825/dapr/utils"
	"io"
	"strings"
	"sync"
)

const (
	propertyKeyActorStateStore = "actorstatestore"
)

var log = logger.NewLogger("dapr.runtime.processor.appLogger")

type Options struct {
	Registry         *comps.Registry
	ComponentStore   *compstore.ComponentStore
	Meta             *meta.Meta
	PlacementEnabled bool
	PubsubAdapter    pubsub.Adapter
}

type appLogger struct {
	registry            *comps.Registry
	compStore           *compstore.ComponentStore
	meta                *meta.Meta
	lock                sync.RWMutex
	pubSubAdapter       pubsub.Adapter
	actorStateStoreName *string
	placementEnabled    bool
}

func New(opts Options) *appLogger {
	return &appLogger{
		registry:         opts.Registry,
		compStore:        opts.ComponentStore,
		meta:             opts.Meta,
		placementEnabled: opts.PlacementEnabled,
		pubSubAdapter:    opts.PubsubAdapter,
	}
}

func (s *appLogger) Init(ctx context.Context, comp compapi.Component) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	fName := comp.LogName()

	store, err := s.registry.Create(comp.Spec.Type, comp.Spec.Version)
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

		props := meta.Properties
		err = store.Init(ctx, common.Metadata{Properties: props}, func() pubsub.Adapter {
			return s.pubSubAdapter
		})

		if err != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "init", comp.ObjectMeta.Name)
			return rterrors.NewInit(rterrors.InitComponentFailure, fName, err)
		}

		s.compStore.AddAppLogger(comp.ObjectMeta.Name, store)
		err = compstate.SaveStateConfiguration(comp.ObjectMeta.Name, props)
		if err != nil {
			diag.DefaultMonitoring.ComponentInitFailed(comp.Spec.Type, "init", comp.ObjectMeta.Name)
			wrapError := fmt.Errorf("failed to save lock keyprefix: %s", err.Error())
			return rterrors.NewInit(rterrors.InitComponentFailure, fName, wrapError)
		}

		// when placement address list is not empty, set specified actor store.
		if s.placementEnabled {
			// set specified actor store if "actorStateStore" is true in the spec.
			actorStoreSpecified := false
			for k, v := range props {
				//nolint:gocritic
				if strings.ToLower(k) == propertyKeyActorStateStore {
					actorStoreSpecified = utils.IsTruthy(v)
					break
				}
			}

			if actorStoreSpecified {
				if s.actorStateStoreName == nil {
					log.Info("Using '" + comp.ObjectMeta.Name + "' as actor state store")
					s.actorStateStoreName = &comp.ObjectMeta.Name
				} else if *s.actorStateStoreName != comp.ObjectMeta.Name {
					return fmt.Errorf("detected duplicate actor state store: %s and %s", *s.actorStateStoreName, comp.ObjectMeta.Name)
				}
			}
		}

		diag.DefaultMonitoring.ComponentInitialized(comp.Spec.Type)
	}

	return nil
}

func (s *appLogger) Close(comp compapi.Component) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ss, ok := s.compStore.GetEventStorage(comp.Name)
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

func (s *appLogger) ActorStateStoreName() (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.actorStateStoreName == nil {
		return "", false
	}
	return *s.actorStateStoreName, true
}
