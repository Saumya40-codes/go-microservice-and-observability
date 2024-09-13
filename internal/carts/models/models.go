package models

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Cart struct {
	Items []CartItem
	Total int
}

type CartItem struct {
	Product  Product
	Quantity int
}

func init() {
	// connect to redis
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redispw",
		DB:       0,
		Username: "default",
	})
}

// Add to Cart
func AddToCart(product string) error {
	ctx := context.Background()
	_, err := Client.SAdd(ctx, "cart", product).Result()
	return err
}

// GetCart returns the cart
func GetCart() (Cart, error) {
	ctx := context.Background()
	productIDs, err := Client.SMembers(ctx, "cart").Result()
	if err != nil {
		return Cart{}, err
	}

	var cart Cart
	for _, id := range productIDs {
		productMap, err := Client.HGetAll(ctx, "product:"+id).Result()
		if err != nil {
			continue
		}

		cart.Items = append(cart.Items, CartItem{
			Product: Product{
				ID:    productMap["ID"],
				Name:  productMap["Name"],
				Price: productMap["Price"],
			},
			Quantity: 1,
		})
		val, err := strconv.Atoi(productMap["Price"])
		if err != nil {
			continue
		}
		cart.Total += val
	}

	return cart, nil
}
