# Microservices Observability Project

This project demonstrates the implementation of observability in a microservices architecture using Go, OpenTelemetry, and Jaeger. It consists of two main services: a product service and a cart service, along with supporting infrastructure for tracing and data storage.

I made this simple project to understand build microservices in Go, observability using Open Telementry and Tracing using Jaeger 

## Project Overview

The project is designed to showcase:
- Microservices architecture with Go
- Distributed tracing using OpenTelemetry
- Trace visualization with Jaeger
- Docker containerization and orchestration with Docker Compose

### Components

1. **Product Service**: Handles product-related operations
2. **Cart Service**: Manages shopping cart functionality
3. **Redis**: Used for data storage
4. **Jaeger**: Collects and visualizes trace data
5. **Nginx**: Acts as a reverse proxy for the services

## Prerequisites

- Docker and Docker Compose
- Go (for local development)

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/Saumya40-codes/go-microservice-and-observability.git
   cd go-microservice-and-observability
   ```

2. Build and run the services:
   ```
   docker-compose up
   ```

3. Access the services:
   - Product Service: http://localhost/products
   - Cart Service: http://localhost/carts
   - Jaeger UI: http://localhost:16686

Both simple pages will look like:
![image](https://github.com/user-attachments/assets/bcbefb30-9f92-40ae-9154-2bdfcde4d70c)
![image](https://github.com/user-attachments/assets/7ac102bc-37af-4a6a-a617-bb26fa999d81)


## Observability Features

### OpenTelemetry Integration

Both the product and cart services are instrumented with OpenTelemetry. This allows for:
- Automatic tracing of HTTP requests
- Custom span creation for important operations
- Propagation of trace context between services

### Jaeger Tracing

Jaeger is used to collect and visualize the trace data:
- Access the Jaeger UI at http://localhost:16686
- View service dependencies, trace timelines, and span details
- Analyze request flows and performance bottlenecks

![image](https://github.com/user-attachments/assets/38d3b68e-ae57-478c-9ebb-d815c9493e43)

## Development

To modify the services:

1. Update the Go code in the respective service directories
2. Rebuild the Docker images:
   ```
   docker build -t saumyashah40/go-micro-and-observability-product:latest ./product-service
   docker build -t saumyashah40/go-micro-and-observability-cart:latest ./cart-service
   ```
3. Run `docker-compose up` to start the updated services

## Configuration

Environment variables in `docker-compose.yml` control various aspects of the services:
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Specifies the Jaeger collector endpoint
- `OTEL_SERVICE_NAME`: Sets the service name for tracing
- `CART_SERVICE_URL`: Configures the cart service URL for the product service
