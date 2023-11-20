/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by mockery v1.0.0.

package testing

import (
	"context"

	"github.com/liuxd6825/components-contrib/pubsub"
	state "github.com/liuxd6825/components-contrib/state"
	"github.com/liuxd6825/dapr/pkg/apis/components/v1alpha1"
	"github.com/liuxd6825/dapr/pkg/outbox"
)

// MockPubSubAdapter is mock for PubSubAdapter
type MockPubSubAdapter struct {
	PublishFn     func(ctx context.Context, req *pubsub.PublishRequest) error
	BulkPublishFn func(ctx context.Context, req *pubsub.BulkPublishRequest) (pubsub.BulkPublishResponse, error)
}

// Publish is an adapter method for the runtime to pre-validate publish requests
// And then forward them to the Pub/Sub component.
// This method is used by the HTTP and gRPC APIs.
func (a *MockPubSubAdapter) Publish(ctx context.Context, req *pubsub.PublishRequest) error {
	return a.PublishFn(ctx, req)
}

// Publish is an adapter method for the runtime to pre-validate publish requests
// And then forward them to the Pub/Sub component.
// This method is used by the HTTP and gRPC APIs.
func (a *MockPubSubAdapter) BulkPublish(ctx context.Context, req *pubsub.BulkPublishRequest) (pubsub.BulkPublishResponse, error) {
	return a.BulkPublishFn(ctx, req)
}

func (a *MockPubSubAdapter) Outbox() outbox.Outbox {
	return &outboxMock{}
}

type outboxMock struct{}

func (o *outboxMock) AddOrUpdateOutbox(stateStore v1alpha1.Component) {}

func (o *outboxMock) Enabled(stateStore string) bool {
	return false
}

func (o *outboxMock) PublishInternal(ctx context.Context, stateStore string, states []state.TransactionalStateOperation, source string) ([]state.TransactionalStateOperation, error) {
	return nil, nil
}

func (o *outboxMock) SubscribeToInternalTopics(ctx context.Context, appID string) error {
	return nil
}
