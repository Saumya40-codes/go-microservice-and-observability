package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"mysite.com/products/models"
)

var productsTemplate = template.Must(template.ParseFiles("products.html"))
var cartServiceURL string

func init() {
	cartServiceURL = os.Getenv("CART_SERVICE_URL")
	if cartServiceURL == "" {
		cartServiceURL = "http://localhost:3002"
	}

	initTracer()
}

func initTracer() func(context.Context) error {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("product-service"),
			attribute.String("environment", "development"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Println("OpenTelemetry tracer initialized")

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
			return err
		}
		log.Println("Tracer provider shut down")
		return nil
	}
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("product-service").Start(r.Context(), "ProductsHandler")
	defer span.End()

	products, err := models.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	productsTemplate.Execute(w, products)
}

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("product-service").Start(r.Context(), "AddToCartHandler")
	defer span.End()

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

	cartRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, cartServiceURL+"/add-to-cart", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cartRequest.Header.Set("Content-Type", "application/json")

	// Inject the trace context into the outgoing request
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(cartRequest.Header))

	resp, err := http.DefaultClient.Do(cartRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

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
