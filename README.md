# AP2 Assignment 2 — gRPC Migration

Two Go microservices communicating over pure gRPC with Protocol Buffers.

## Services

| Service | HTTP | gRPC | Database |
|---|---|---|---|
| purchase-service | :8080 | :9090 | purchase_db |
| transaction-service | :8081 | :9091 | transaction_db |

## Architecture

```
Client (REST)
    │
    ▼
purchase-service :8080
    │  gRPC (protobuf + HTTP/2)
    ▼
transaction-service :9091
```

- purchase-service exposes REST endpoints for clients and a gRPC streaming endpoint for order status updates
- transaction-service exposes a gRPC server for payment processing
- Inter-service communication is pure gRPC with protobuf binary encoding

## Running

1. Start PostgreSQL on `localhost:5432`
2. Create databases: `purchase_db`, `transaction_db`
3. Run migrations from `purchase-service/migrations/` and `transaction-service/migrations/`
4. Start transaction-service first: `cd transaction-service && go run ./cmd/transaction-service`
5. Start purchase-service: `cd purchase-service && go run ./cmd/purchase-service`

See [docs/postman-demo.md](docs/postman-demo.md) for full demo instructions.

## Code Generation

Proto files are in `proto/`. Generated code is in `pkg/pb/` inside each service.

To regenerate:
```bash
cd proto
buf generate --template buf.gen.order.yaml --path order/order.proto
buf generate --template buf.gen.payment.yaml --path payment/payment.proto
```

Requires `buf` installed: https://buf.build/docs/installation

## Documentation

- [Project Explanation](docs/project-explanation.md) — architecture, flow, design decisions
- [Postman Demo](docs/postman-demo.md) — step-by-step demo instructions
