# Chirpy Server

A small Twitter-like REST API built with Go.

This project was created as a backend practice project to implement authentication, database integration, and webhook handling using the Go standard library.

---

## 📌 About

Chirpy Server is a lightweight backend service that allows users to:

- Register and authenticate
- Create and manage short posts ("chirps")
- Use JWT authentication (access + refresh tokens)
- Handle webhook events securely
- Interact with a PostgreSQL database using sqlc
- Run database migrations with goose

This project focuses on backend fundamentals rather than production scalability.

---

## 🧠 Tech Stack

- Go (net/http)
- PostgreSQL
- sqlc
- goose (migrations)
- JWT authentication
- HMAC signature verification (webhooks)

---

## 🚀 Features

### Authentication
- User registration
- Login with password hashing
- Access token (short-lived)
- Refresh token (stored in DB)
- Token revocation

### Chirps
- Create chirp
- Get all chirps (sorted asc/desc)
- Get chirp by ID
- Delete chirp (authorized)

### Webhooks
- Event-driven callback endpoint
- HMAC signature validation
- Immediate 200 response
- Asynchronous processing
- Constant-time comparison using hmac.Equal

---

## 🗄️ Database Setup

Set environment variables:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
export JWT_SECRET="your_secret"
```

Run migrations:

```bash
goose postgres "$DATABASE_URL" up
```

---

## ▶ Running the Server

```bash
go run main.go
```

Default port: `8080`

Health check:

```
GET /healthz
```

---

## 🔐 Authentication Flow

1. User logs in → receives access + refresh token  
2. Access token is used for protected endpoints  
3. When expired → client calls `/refresh`  
4. Refresh token can be revoked via `/revoke`

---

## 📡 API Endpoints

```
POST   /api/users
POST   /api/login
POST   /api/refresh
POST   /api/revoke

POST   /api/chirps
GET    /api/chirps
GET    /api/chirps/{id}
DELETE /api/chirps/{id}
```

---

## 🎯 Learning Goals

This project was built to practice:

- Writing REST APIs with net/http
- Middleware patterns
- JWT authentication
- Secure webhook verification
- Database schema migrations
- Query generation with sqlc
- Proper HTTP status handling

---
