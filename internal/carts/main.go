package cart

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	// connect to redis
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// Add to Cart
func AddToCart(product string) error {
	ctx := context.Background()
	_, err := Client.RPush(ctx, "cart", product).Result()
	return err
}
