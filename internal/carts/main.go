package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"mysite.com/carts/models"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

var cartTemplate = template.Must(template.ParseFiles("cart.html"))

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.AddToCart(product.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	cart, err := models.GetCart()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cartTemplate.Execute(w, cart)
}

func main() {
	http.HandleFunc("/carts", GetCartHandler)
	http.HandleFunc("/add-to-cart", AddToCartHandler)

	log.Println("Starting server on :3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
