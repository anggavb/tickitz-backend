# Tickitz Backend

<p align="center">
  <img src="./public/img/logo.png" alt="Tickitz Logo" width="250"/>
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Backend service for Tickitz, an online movie ticket booking platform. This API provides authentication, movie management, scheduling, booking, payment processing, and administrative functionalities.

## Tech Stacks

- [![Go](https://img.shields.io/badge/Go-v1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v17-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-v8-DC382D?logo=redis&logoColor=white)](https://redis.io/)
- [![JWT](https://img.shields.io/badge/JWT-Authentication-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
- [![Docker](https://img.shields.io/badge/Docker-latest-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
- [![REST API](https://img.shields.io/badge/REST-API-green)](#)

## Design Philosophy

The backend is designed with:

- Clean Architecture principles
- Separation of Concerns
- Scalable service structure
- Maintainable codebase
- Secure authentication and authorization
- Efficient database access

## Features

- User Authentication & Authorization
- JWT-based Security
- Movie Management (CRUD)
- Cinema Management
- Showtime Management
- Ticket Booking System
- Payment Processing
- User Profile Management
- Order Management
- Admin Features
- Database Migration Support
- Redis Caching
- RESTful API

## How to Setup

### 1. Clone Repository

```bash
git clone https://github.com/your-organization/tickitz-backend.git
cd tickitz-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Create Environment File

Copy the provided example environment file and update the values:

```bash
cp example.env .env
```

Then edit `.env` with your own settings.

Example values:

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

### 5. Run Database Migration

```bash
make migrate-up
```

Or:

```bash
migrate -path migrations -database postgres://user:password@localhost:5432/tickitz?sslmode=disable up
```

### 6. Seed Database (Optional)

```bash
make seed
```

### 7. Run Application

```bash
go run cmd/main.go
```

Server will run at:

```text
http://localhost:8081
```

If you want to see task help for migrations and seeders, run:

```bash
make help
```

## API Documentation

### Public / Movies

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /movies | GET | List movies with filters (public) |
| /movies/upcoming | GET | List upcoming movies |
| /movies/:slug | GET | Get movie detail by slug |
| /movies/:slug/schedule-options | GET | Get schedule options for a movie by slug |
| /movies/:slug/schedules | GET | Get schedules for a movie by slug |
| /movies/showtimes | GET | Get available showtimes |
| /movies/locations | GET | Get available locations |

### Authentication

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /auth/signup | POST | Register new user |
| /auth/activate | POST | Activate user account with OTP |
| /auth/otp | POST | Request new OTP for activation |
| /auth/signin | POST | Login and get JWT token |
| /auth/logout | DELETE | Logout authenticated user |
| /auth/password | PATCH | Change user password (authenticated) |
| /auth/password/forgot | POST | Request password reset link |
| /auth/password/reset | POST | Reset user password with token |

### Profile

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /profile | GET | Get authenticated user profile |
| /profile/update | PATCH | Update authenticated user profile |

### Admin / Movies

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /admin/movies | GET | List movies (admin view) |
| /admin/movies/months | GET | List release months |
| /admin/movies/:id | GET | Get movie detail by ID |
| /admin/movies | POST | Create movie |
| /admin/movies/:id | PATCH | Update movie |
| /admin/movies/:id | DELETE | Delete movie |
| /admin/movies/:id/showtimes | GET | Get showtimes for movie by ID |
| /admin/movies/:id/showtimes | POST | Add showtime schedule for movie |
| /admin/cinemas | GET | List cinemas |
| /admin/categories | GET | List categories |
| /admin/casts | GET | List casts |

### Seats

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /movie-cinemas/:movie_cinema_id/seats | GET | Get seat map for a movie schedule |

### Other

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /swagger/*any | GET | Swagger UI documentation |
| /img/* | GET | Static image assets |

## Project Structure

```text
├── cmd/
├── controller/
├── service/
├── repository/
├── model/
├── dto/
├── router/
├── middleware/
├── migrations/
├── seeds/
└── utils/
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