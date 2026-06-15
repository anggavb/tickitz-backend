## Logo
<p align="center">
  <img src="https://raw.githubusercontent.com/anggavb/tickitz-frontend/refs/heads/main/src/assets/logo.png" width="200" alt="Tickitz Logo" />
</p>

# Tickitz Backend
[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Tickitz Backend adalah REST API untuk platform pemesanan tiket bioskop. Backend ini menangani autentikasi pengguna, manajemen film, jadwal tayang, kursi, pemesanan tiket, profile user, admin movie management, dan dokumentasi API melalui Swagger.

## Tech Stacks
- [![Go](https://img.shields.io/badge/Go-1.20-blue?logo=go&logoColor=white)](https://golang.org/)
- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=gin&logoColor=white)](https://gin-gonic.com/en/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.4-blue?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

## Description & Features
- REST API untuk mengelola film, showtime, pemesanan tiket, dan autentikasi pengguna.
- Fitur utama: register/login, manajemen film (CRUD), pemesanan, pembayaran, dan endpoint admin.


## How to Setup
1. Pastikan Go (v1.20+) dan PostgreSQL terpasang.
2. Buat file environment di root proyek `tickitz-backend/.env`:
```
DB_HOST={YOUR_DB_HOST}
DB_PORT={YOUR_DB_PORT}
DB_PASS={YOUR_DB_PASS}
DB_NAME={YOUR_DB_NAME}
DB_USER={YOUR_DB_USER}
```

## Quickstart
```bash
# Clone Project
git clone https://github.com/anggavb/tickitz-backend.git
# masuk ke folder backend
cd tickitz-backend

# install dependency
go mod download

# jalankan migrasi (jika memakai migrasi di folder db/migrations)
# contoh menggunakan tool migrasi yang Anda pakai

# jalankan server
go run cmd/main.go

# atau gunakan Makefile jika tersedia
make run
```

## Routes

### Public / Movies
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /movies/ | GET | List movies with filters (public) |
| /movies/:slug | GET | Get movie detail by slug |
| /movies/:slug/schedules | GET | Get schedules for a movie by slug |
| /movies/showtimes | GET | Get available showtimes |
| /movies/locations | GET | Get available locations |

### Auth
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /auth/signup | POST | Register new user |
| /auth/activate | POST | Activate user account (OTP/email) |
| /auth/otp | POST | Request new OTP |
| /auth/signin | POST | Login (returns JWT) |
| /auth/password | PATCH | Change user password |

### Profile
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /profile | GET | Get authenticated user's profile |

### Admin / Movies
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /admin/movies | GET | List movies (admin view) |
| /admin/movies/months | GET | List release months for filtering |
| /admin/movies/:id | GET | Get movie by id (admin) |
| /admin/movies | POST | Create movie (admin) |
| /admin/movies/:id | PATCH | Update movie (admin) |
| /admin/movies/:id | DELETE | Delete movie (admin) |
| /admin/categories | GET | List categories |
| /admin/casts | GET | List casts |

### Seats
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /movie-cinemas/:movie_cinema_id/seats | GET | Get seat map for a specific movie-cinema

### Other
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /swagger/*any | GET | Swagger UI / API documentation |
| /img/* | GET | Static images served from `public/img` |

### Documentation
Untuk dokumentasi lengkap buka `/swagger/index.html` pada server (jika Swagger diaktifkan).

## How to Contribute
- Fork repository, buat branch untuk fitur/bugfix, lalu buat pull request.
- Sertakan deskripsi perubahan dan langkah reproduksi jika perlu.

## License
This project is licensed under the MIT License

## Related Project
[Frontend](https://github.com/anggavb/tickitz-frontend.git)
