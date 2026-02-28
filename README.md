# 🐦 Chirpy Server

A lightweight Twitter-like REST API built with Go — designed as a backend practice project covering authentication, database integration, and webhook handling.

---

## ✨ Features

- 👤 User registration & login with password hashing
- 🔐 JWT-based authentication (access token + refresh token)
- 🐦 Create, read, and delete short posts ("chirps")
- 🔄 Token refresh & revocation
- 🪝 Webhook endpoint with HMAC signature verification
- 🗄️ PostgreSQL database via **sqlc** (type-safe query generation)
- 📦 Database migrations with **goose**

---

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| Go (`net/http`) | HTTP server & routing |
| PostgreSQL | Database |
| sqlc | Type-safe SQL query generation |
| goose | Database migrations |
| JWT | Access & refresh token authentication |
| HMAC | Webhook signature verification |

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- [goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [sqlc](https://sqlc.dev/) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Setup

**1. Clone the repository**
```bash
git clone https://github.com/seminhnva/chirpy-server.git
cd chirpy-server
```

**2. Set environment variables**
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
export JWT_SECRET="your_jwt_secret"
export POLKA_KEY="your_webhook_secret"
```

**3. Run database migrations**
```bash
goose postgres "$DATABASE_URL" up
```

**4. Start the server**
```bash
go run .
```

The server runs on **port 8080** by default.

---

## 📡 API Endpoints

### Health
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/healthz` | Health check |

### Users
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/users` | Register a new user |
| `PUT` | `/api/users` | Update user info (auth required) |

### Authentication
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/login` | Login and receive tokens |
| `POST` | `/api/refresh` | Get a new access token |
| `POST` | `/api/revoke` | Revoke refresh token |

### Chirps
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/chirps` | Create a chirp (auth required) |
| `GET` | `/api/chirps` | Get all chirps (supports `?sort=asc\|desc`) |
| `GET` | `/api/chirps/{id}` | Get a chirp by ID |
| `DELETE` | `/api/chirps/{id}` | Delete a chirp (auth required, owner only) |

### Webhooks
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/polka/webhooks` | Receive webhook events (HMAC verified) |

---

## 🔐 Authentication Flow

1. Register via `POST /api/users`
2. Login via `POST /api/login` → receive **access token** (short-lived) + **refresh token**
3. Use the access token in the `Authorization: Bearer <token>` header for protected routes
4. When the access token expires, call `POST /api/refresh` with the refresh token to get a new one
5. Revoke a session via `POST /api/revoke`

---

## 🪝 Webhook Verification

Webhook requests are verified using **HMAC signatures**. The server compares the `Authorization` header against an expected signature using `hmac.Equal` (constant-time comparison) to prevent timing attacks. Invalid requests are rejected immediately with `401 Unauthorized`.

---

## 🎯 Learning Goals

This project was built to practice:

- Writing REST APIs with Go's standard `net/http` package
- Middleware patterns (logging, auth guards)
- JWT authentication & token lifecycle management
- Secure webhook event handling
- Database schema design & migrations with goose
- Type-safe SQL with sqlc

---

## 📄 License

MIT License — free to use and modify.
