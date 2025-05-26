# Eventure: Secure Event Management API

Eventure is a secure backend API built with Go (Gin + GORM + PostgreSQL) that demonstrates authentication, role-based and permission-based access control (RBAC + ABAC), resource ownership enforcement, structured logging, rate limiting, and more.

## Features

- User registration & login with simulated token system
- Role-based access (`admin`, `user`)
- Fine-grained permission checks (e.g., `event:update`, `event:delete:any`)
- Ownership-based controls for event updates/deletion
- Secure headers, custom error masking, and rate limiting
- Structured access-denial logs
- Safe vs. unsafe SQL search demo
- Test coverage

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL

### Setup

1. Create a `.env` file in the root directory (replace with your actual DB credentials and port if needed):

```
DB_DSN=postgres://postgres:password@localhost:5432/eventure?sslmode=disable
PORT=8080
```

2. Make sure you have all required Go modules:

```
go mod tidy
```

### Running the App

1. Start PostgreSQL:

```
brew services start postgresql@14
psql -U postgres
```

2. Create the database and grant permissions (if needed):

```
CREATE DATABASE eventure;
GRANT ALL PRIVILEGES ON DATABASE eventure TO postgres;
```

3. Run the app to auto-migrate tables and seed roles + permissions:

```
go run cmd/api/main.go
```

## API Routes

### Auth

| Method | Path             | Description                 |
| ------ | ---------------- | --------------------------- |
| POST   | `/auth/register` | Register new user           |
| POST   | `/auth/login`    | Login user                  |
| POST   | `/auth/refresh`  | Refresh access token        |
| GET    | `/auth/profile`  | Get profile of current user |

### Admin (admin only)

| Method | Path           | Description               |
| ------ | -------------- | ------------------------- |
| GET    | `/admin/users` | List all registered users |

### Event

| Method | Path          | Description                                     |
| ------ | ------------- | ----------------------------------------------- |
| POST   | `/events`     | Create a new event                              |
| PUT    | `/events/:id` | Update an event (ownership/permission required) |
| DELETE | `/events/:id` | Delete an event (ownership/permission required) |

### Search

| Method | Path             | Description                                     |
| ------ | ---------------- | ----------------------------------------------- |
| GET    | `/events/search` | Search events by title (safe/unsafe comparison) |

### Security (test-only)

Run with: `go test ./internal/test -v`

| Method | Path      | Description         |
| ------ | --------- | ------------------- |
| GET    | `/health` | Simple health check |
