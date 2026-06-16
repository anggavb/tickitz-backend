# Tickitz Backend

<p align="center">
  <img src="https://raw.githubusercontent.com/anggavb/tickitz-frontend/refs/heads/main/src/assets/logo.png" width="200" alt="Tickitz Logo" />
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Backend service for Tickitz, an online movie ticket booking platform. This API provides authentication, movie management, scheduling, booking, payment, dashboard, and profile features.

## Tech Stack

- [![Go](https://img.shields.io/badge/Go-v1.26.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
- [![Gin Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-008ECF?logo=gin&logoColor=white)](https://gin-gonic.com/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.3-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-v8.4.3-DC382D?logo=redis&logoColor=white)](https://redis.io/)
- [![Golang JWT](https://img.shields.io/badge/Golang_JWT-v5.3.1-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
- [![Gin Swagger](https://img.shields.io/badge/Gin_Swagger-v1.6.1-85EA2D?logo=swagger&logoColor=black)](https://swagger.io/)
- [![Docker](https://img.shields.io/badge/Docker-latest-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Features

- User registration, account activation, sign in, logout, and password reset
- JWT authentication with Redis-backed token/session state
- Public movie listing, detail, schedule filters, locations, and showtimes
- Admin movie CRUD, cinema/category/cast lookup, and movie showtime management
- Seat map lookup, pending order creation, seat selection, payment, QR ticket, and order history
- Admin sales and ticket dashboard charts
- Swagger API documentation and static asset serving
- Database migrations and seeders

## Requirements

- Go 1.26.3
- PostgreSQL
- Redis
- `Make`
- `migrate` CLI for database migrations
- `psql` CLI for seeders
- Docker, optional

## Setup

### 1. Clone Repository

```bash
git clone https://github.com/anggavb/tickitz-backend.git
cd tickitz-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Create Environment File

```bash
cp example.env .env
```

Then update `.env` with your local values.

Example:

```env
# SMTP
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_SENDER=Tickitz <your_email@gmail.com>

# App
APP_HOST=localhost
APP_PORT=8081
CLIENT_URL=http://localhost:5173

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tickitz
DB_USER=postgres
DB_PASS=password
DB_URL=postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

# Redis
RDB_HOST=localhost
RDB_PORT=6379
RDB_USER=
RDB_PASS=
RDB_PREFIX=tickitz

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_ISSUER=tickitz-backend
```

### 4. Create Database

```sql
CREATE DATABASE tickitz;
```

### 5. Run Migrations

```bash
make migrate-up
```

Equivalent manual command:

```bash
migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/tickitz?sslmode=disable" up
```

### 6. Seed Database

```bash
make seed
```

To reset the database and apply migrations plus seeders:

```bash
make fresh
```

### 7. Run Application

```bash
go run cmd/main.go
```

The server runs at:

```text
http://localhost:8081
```

View available migration and seeder commands:

```bash
make help
```

## Docker

Build the application image:

```bash
docker build -t tickitz-backend .
```

Run the container with your environment file:

```bash
docker run --env-file .env -p 8081:8081 tickitz-backend
```

PostgreSQL and Redis must still be reachable from inside the container. Use hostnames that match your Docker/network setup.

## API Documentation

Swagger UI is available after the server starts:

```text
http://localhost:8081/swagger/index.html
```

Use the `Authorization: Bearer <token>` header for authenticated endpoints. Admin endpoints also require an admin user.

### Public / Movies

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/movies` | GET | Public | List movies with filters and pagination |
| `/movies/upcoming` | GET | Public | List upcoming movies |
| `/movies/:slug` | GET | Public | Get movie detail by slug |
| `/movies/:slug/schedule-options` | GET | Public | Get schedule filter options for a movie |
| `/movies/:slug/schedules` | GET | Public | Get schedules for a movie |
| `/movies/showtimes` | GET | Public | Get available showtimes |
| `/movies/locations` | GET | Public | Get available locations |

### Authentication

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/auth/signup` | POST | Public | Register a new user |
| `/auth/activate` | POST | Public | Activate user account with OTP |
| `/auth/otp` | POST | Public | Request a new activation OTP |
| `/auth/signin` | POST | Public | Sign in and get a JWT token |
| `/auth/logout` | DELETE | User | Logout authenticated user |
| `/auth/password` | PATCH | User | Change authenticated user password |
| `/auth/password/forgot` | POST | Public | Request password reset link |
| `/auth/password/reset` | POST | Public | Reset password with token |

### Profile

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/profile` | GET | User | Get authenticated user profile |
| `/profile/update` | PATCH | User | Update authenticated user profile |

### Admin / Movies

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/admin/movies` | GET | Admin | List movies for admin |
| `/admin/movies/months` | GET | Admin | List available release months |
| `/admin/movies/:id` | GET | Admin | Get movie detail by ID |
| `/admin/movies` | POST | Admin | Create a movie |
| `/admin/movies/:id` | PATCH | Admin | Update a movie |
| `/admin/movies/:id` | DELETE | Admin | Delete a movie |
| `/admin/movies/:id/showtimes` | GET | Admin | Get showtimes for a movie |
| `/admin/movies/:id/showtimes` | POST | Admin | Add or update showtime schedules for a movie |
| `/admin/cinemas` | GET | Admin | List cinemas |
| `/admin/categories` | GET | Admin | List categories |
| `/admin/casts` | GET | Admin | List casts |

### Admin / Dashboard

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/admin/dashboard/sales-chart` | GET | Admin | Get revenue chart data |
| `/admin/dashboard/ticket-sales` | GET | Admin | Get ticket sales chart data |

### Seats

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/movie-cinemas/:movie_cinema_id/seats` | GET | Public | Get seat map for a movie schedule |

### Orders

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/orders` | POST | User | Create or reuse a pending order |
| `/orders/history` | GET | User | Get authenticated user's order history |
| `/orders/:order_id` | GET | User | Get order detail |
| `/orders/:order_id/seats` | PATCH | User | Set selected seats for a pending order |
| `/orders/:order_id/payment-methods` | GET | User | Get available payment methods for an order |
| `/orders/:order_id/payment` | PATCH | User | Submit order payment |
| `/orders/:order_id/qr` | GET | User | Get ticket QR image for a paid order |

### Static Assets and Docs

| Endpoint | Method | Auth | Description |
| --- | --- | --- | --- |
| `/swagger/*any` | GET | Public | Swagger UI |
| `/img/*` | GET | Public | Uploaded movie and profile images |
| `/payment/*` | GET | Public | Payment method assets |
| `/cinema/*` | GET | Public | Cinema logo assets |
| `/poster/*` | GET | Public | Seeded movie poster assets |

## Project Structure

```text
.
├── cmd/                    # Application entrypoint
├── db/
│   ├── migrations/         # SQL migration files
│   └── seeds/              # SQL seed files
├── docs/                   # Generated Swagger files
├── internal/
│   ├── config/             # PostgreSQL and Redis connections
│   ├── controller/         # HTTP handlers
│   ├── dto/                # Request and response DTOs
│   ├── middleware/         # Auth and CORS middleware
│   ├── model/              # Domain models
│   ├── repository/         # Data access layer
│   ├── router/             # Route registration
│   └── service/            # Business logic
├── pkg/                    # Shared utility packages
├── public/                 # Static assets
├── Dockerfile
├── Makefile
├── example.env
└── go.mod
```

## Verification

Run package checks:

```bash
go test ./...
```

## How to Contribute

1. Fork this repository.
2. Clone your forked repository.

```bash
git clone https://github.com/your-username/tickitz-backend.git
```

3. Create a new branch.

```bash
git checkout -b feature/your-feature-name
```

4. Make your changes.
5. Commit your changes.

```bash
git commit -m "feat: add movie schedule endpoint"
```

6. Push your branch.

```bash
git push origin feature/your-feature-name
```

7. Create a Pull Request to the main repository.

## License

This project is licensed under the MIT License.

## Related Projects

[Frontend](https://github.com/anggavb/tickitz-frontend.git)
