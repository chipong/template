package redisCache

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

var (
	//pub *redis.PubSub
	sub *redis.PubSub
)

func Subscribe(ctx context.Context, channels ...string) {
	sub = write.PSubscribe(ctx, channels...)
}

func Channel() <-chan *redis.Message {
	return sub.Channel()
}

func PubSub(ctx context.Context) ([]string, error) {
	return write.PubSubChannels(ctx, "*").Result()
}

// unsubscribe

func Publish(ctx context.Context, channel, message string) error {
	_, err := write.Publish(ctx, channel, message).Result()
	return err
}
