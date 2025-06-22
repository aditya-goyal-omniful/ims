# 🏬 Inventory Management Service (IMS)

Inventory Management microservice built using **Go**, **Gin**, **GORM**, **PostgreSQL**, **Redis**, and **GoCommons**. This service supports **multi-tenant inventory tracking**, **Redis caching**, and **inter-service order validation**. It's part of a microservice-based architecture where it integrates with an Order Management System (OMS) and is cloud-native ready.

---

## 🚀 Features

* 🔑 **Multi-Tenant Support** via `X-Tenant-ID` header
* 🏢 **Tenant, Seller, Hub, SKU, Inventory** CRUD operations
* ⚡ **Redis Caching** for hub and SKU lookups
* 🔁 **Inventory Upsert** operation with atomic behavior
* 🔍 **Order Validation API** for inter-service inventory checks
* 🔐 **Middleware-based Authentication** for tenant isolation
* 📦 Modular & production-ready with cloud support (e.g., S3, Kafka, SQS)

---

## 📁 Project Structure

```
ims/
├── cmd/                    # Main entry point
├── pkg/
│   ├── configs/            # DB, Redis, and config loading logic
│   ├── controllers/        # Gin handler functions for all models
│   ├── models/             # GORM models and business logic
│   ├── routes/             # All route definitions using GoCommons HTTP
│   ├── middlewares/        # Auth middleware for multi-tenant
│   └── utils/              # Helper functions (e.g., UUID validation)
├── swagger.yaml            # API documentation in OpenAPI format
└── go.mod
```

---

## ⚙️ Configuration

This project uses **GoCommons Config**, allowing flexible configuration sources:

* `local.yaml` for development
* AWS AppConfig for production

### Sample `local.yaml`

```yaml
postgres:
  master:
    host: localhost
    port: 5432
    user: ims_user
    password: secret
    dbname: ims

redis:
  address: localhost:6379
  db: 0

app:
  env: local
  port: 8080
```

Set environment variable:

```bash
$env:CONFIG_SOURCE = "local"
```

---

## 🚦 Running Locally

```bash
go run cmd/main.go
```

---

## 🔒 Headers Required

Every protected route requires the following header:

```http
X-Tenant-ID: <UUID>
```

This ensures tenant-level data isolation in all operations.

---

## 📚 API Documentation

Swagger UI is available at:

👉 [`http://localhost:8080/swagger/index.html`](http://localhost:8080/swagger/index.html)

---

## ✅ Sample Endpoints

### GET /hubs

```http
GET /hubs
Headers: X-Tenant-ID: <uuid>
```

### POST /inventories/upsert

```http
POST /inventories/upsert
Content-Type: application/json
Headers: X-Tenant-ID: <uuid>

Body:
{
  "tenant_id": "<uuid>",
  "hub_id": "<uuid>",
  "sku_id": "<uuid>",
  "quantity": 10
}
```

## 🌐 Related Services

* 📦 **Order Management Service (OMS)** – Consumes inventory APIs, validates availability via `/validators/validate_order/:hub_id/:sku_id`
* 🗃️ **Kafka + SQS** ready for real-time order and inventory flow

---

## 🤝 Contributors

* Aditya Goyal

---