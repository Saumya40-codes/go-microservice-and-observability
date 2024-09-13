package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"mysite.com/products/models"
)

var productsTemplate = template.Must(template.ParseFiles("products.html"))
var cartServiceURL = "http://localhost:3002"

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	productsTemplate.Execute(w, products)
}

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := r.FormValue("product_id")
	product, ok := models.GetProductByID(productID)

	if !ok {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cartRequest, err := http.NewRequest(http.MethodPost, cartServiceURL+"/add-to-cart", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cartRequest.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(cartRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error adding to cart", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/products", ProductsHandler)
	http.HandleFunc("/add-to-cart", AddToCartHandler)

	log.Println("Starting server on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
