package subscriber

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/pubsub"
)

func (s *Subscriber) Publish(ctx context.Context, request *pubsub.PublishRequest) error {
	c, ok := s.compStore.GetPubSubComponent(request.PubsubName)
	if !ok {
		return errors.New(fmt.Sprintf("PubSub not found %s", request.PubsubName))
	}
	return c.Publish(ctx, request)
}

func (s *Subscriber) BulkPublish(ctx context.Context, request *pubsub.BulkPublishRequest) (pubsub.BulkPublishResponse, error) {
	c, ok := s.compStore.GetPubSubComponent(request.PubsubName)
	if !ok {
		return pubsub.BulkPublishResponse{}, errors.New(fmt.Sprintf("PubSub not found %s", request.PubsubName))
	}
	bp, ok := c.(pubsub.BulkPublisher)
	if !ok {
		return pubsub.BulkPublishResponse{}, errors.New(fmt.Sprintf("PubSub %s is not  BulkPublisher  ", request.PubsubName))
	}
	resp, err := bp.BulkPublish(ctx, request)
	return resp, err
}
