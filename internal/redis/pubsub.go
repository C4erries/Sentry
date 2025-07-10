package redis

import (
	"context"
	"encoding/json"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	topic  string
}

func NewRedisPubSub(client *redis.Client, topic string) *RedisPubSub {
	return &RedisPubSub{
		client: client,
		topic:  topic,
	}
}

func (r *RedisPubSub) Publish(ctx context.Context, alert *model.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, r.topic, data).Err()
}

func (r *RedisPubSub) Subscribe(ctx context.Context, handler func(alert *model.Alert)) error {
	sub := r.client.Subscribe(ctx, r.topic)
	ch := sub.Channel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				sub.Close()
				return
			case msg := <-ch:
				var alert model.Alert
				if err := json.Unmarshal([]byte(msg.Payload), &alert); err == nil {
					handler(&alert)
				}
			}
		}
	}()

	return nil
}
