package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"mysite.com/carts/models"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

var cartTemplate = template.Must(template.ParseFiles("cart.html"))

func init() {
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
			semconv.ServiceNameKey.String("cart-service"),
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

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	_, span := otel.Tracer("cart-service").Start(ctx, "AddToCartHandler")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
	w.Write([]byte("Product added to cart"))
}

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("cart-service").Start(r.Context(), "GetCartHandler")
	defer span.End()

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
