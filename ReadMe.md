# 🏬 Inventory Management Service (IMS)

Golang-based microservice to manage inventory storage, validation, and updates using PostgreSQL, Redis, and inter-service APIs. Supports multi-tenant environments with Redis-backed caching and i18n support.

---

## 🧩 Tech Stack

* **Language**: Go
* **Framework**: Gin
* **Database**: PostgreSQL (via GORM)
* **Cache**: Redis
* **Config**: go\_commons/config
* **HTTP Client**: go\_commons/httpclient
* **i18n**: go\_commons/i18n
* **Swagger**: Swagger UI via swaggo/gin-swagger

---

## 📂 Features

* Multi-tenant support via `X-Tenant-ID` header
* CRUD operations for Tenant, Seller, Hub, SKU, Inventory
* Redis caching for SKU and Hub validation
* Inventory Upsert endpoint for atomic updates
* Order validation API for inter-service communication with OMS
* Middleware-based tenant isolation
* i18n support for multilingual logs and errors
* Swagger docs hosted at `/swagger/index.html`

---

## 🧪 API Endpoints

| Method | Endpoint                         | Description                        |
| ------ | -------------------------------- | ---------------------------------- |
| GET    | `/hubs`                          | Get list of hubs (tenant isolated) |
| GET    | `/skus`                          | Get list of SKUs with filters      |
| POST   | `/inventories/upsert`            | Atomically upsert inventory        |
| GET    | `/validators/validate_order/...` | Validate order hub/sku for OMS     |

---

## ⚙️ How It Works

### 1. **Inventory Upsert**

* API: `POST /inventories/upsert`
* Accepts tenant\_id, hub\_id, sku\_id, and quantity
* Uses GORM for insert/update based on existence

### 2. **Order Validation (OMS Integration)**

* API: `GET /validators/validate_order/:hub_id/:sku_id`
* Called by OMS to verify inventory exists for a hub+sku combo
* Uses Redis caching for fast validation

### 3. **Redis Caching**

* Hubs and SKUs are cached using Redis keyed by tenant and entity ID
* Improves performance on frequent validations

---

## 🐳 Docker Setup

Not containerized by default, but supports Docker-ready components:

* PostgreSQL
* Redis

---

## 🛠 Run Locally

```powershell
$env:CONFIG_SOURCE = "local"
go run cmd/main.go
```

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

---

## 📎 Swagger UI

> View Swagger docs at:

👉 [`http://localhost:8087/swagger/index.html`](http://localhost:8087/swagger/index.html)

---

## 🔐 Authentication Header

All routes require tenant identification:

```http
X-Tenant-ID: <uuid>
```

---

## 📦 Directory Structure

```
ims/
├── cmd/                    # Entry point
├── pkg/
│   ├── configs/            # DB, Redis, config logic
│   ├── controllers/        # API handler logic
│   ├── models/             # GORM models
│   ├── routes/             # HTTP routes
│   ├── middlewares/        # Multi-tenant auth
│   └── utils/              # Helper methods
├── swagger.yaml            # Swagger API docs
└── go.mod
```

---

## 📈 Future Improvements

* Add Prometheus/Grafana integration
* Add unit/integration test coverage
* Add support for soft deletes

---

## 🧠 Developer Notes

* Redis is used for caching hubs and SKUs for faster validations
* All logs and errors are i18n-enabled for future multi-locale support
* Configuration can be toggled via local YAML or AWS AppConfig
* Swagger comments are generated using `swag init`

---

## 🔗 External Dependencies

* [go\_commons](https://github.com/omniful/go_commons)
* [GORM](https://gorm.io)
* [PostgreSQL Go Driver](https://github.com/lib/pq)
* [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)

---

## 🤝 Contributors

* Aditya Goyal

---