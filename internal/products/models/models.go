package models

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

func init() {
	// connect to redis
	Client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redispw",
		DB:       0,
		Username: "default",
	})
}

// GetProducts returns a list of products
func GetProducts() ([]Product, error) {
	ctx := context.Background()
	productKeys, err := Client.Keys(ctx, "product:*").Result()
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, key := range productKeys {
		productMap, err := Client.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		products = append(products, Product{
			ID:    productMap["ID"],
			Name:  productMap["Name"],
			Price: productMap["Price"],
		})
	}

	return products, nil
}

// GetProductByID returns a product by ID
func GetProductByID(id string) (Product, bool) {
	ctx := context.Background()
	productMap, err := Client.HGetAll(ctx, "product:"+id).Result()
	if err != nil || len(productMap) == 0 {
		return Product{}, false
	}

	return Product{
		ID:    productMap["ID"],
		Name:  productMap["Name"],
		Price: productMap["Price"],
	}, true
}
