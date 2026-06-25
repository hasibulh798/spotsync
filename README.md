# 🚗 SpotSync – Smart Parking & EV Charging Reservation

SpotSync is a centralized platform designed for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.

---

## 🛠️ Technology Stack

- **Backend Language:** Go (Golang)
- **Web Framework:** Echo (`github.com/labstack/echo/v4`)
- **ORM:** GORM (`gorm.io/gorm`)
- **Database:** PostgreSQL (with NeonDB / Supabase)
- **Validation:** Go-playground validator (`github.com/go-playground/validator/v10`)
- **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto/bcrypt`)

---

## 🏛️ Clean Architecture Layers

We strictly separate concerns into the following directories:

- **`dto/`**: Data Transfer Objects defining request/response structures.
- **`handler/`**: Handles HTTP request parsing, DTO validation, context claims retrieval, calling services, and writing JSON responses.
- **`service/`**: Holds business logic (e.g., computing capacity, hashing passwords, generating JWTs).
- **`repository/`**: Directly communicates with the database via GORM (CRUD operations, Transactions, and Row-level locking).
- **`models/`**: Defines the database schema GORM structs.
- **`middleware/`**: Custom Echo middlewares (e.g., JWT Auth, Role Guard).
- **`config/`**: Manages environment variables and application configurations.
- **`router/`**: Defines endpoints and binds them to their respective handlers.

---

## ⚙️ Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://username:password@host:port/database?sslmode=require
JWT_SECRET=your_jwt_secret_key
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL database

### Run Locally

1. Clone the repository.
2. Setup your `.env` file.
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

---

## 🌐 API Endpoints Specification

### Authentication
- `POST /api/v1/auth/register` (Public) - User registration
- `POST /api/v1/auth/login` (Public) - User login

### Parking Zones
- `POST /api/v1/zones` (Admin Only) - Create a new zone
- `GET /api/v1/zones` (Public) - Get all parking zones & availability
- `GET /api/v1/zones/:id` (Public) - Get details of a single parking zone

### Reservations
- `POST /api/v1/reservations` (Authenticated) - Reserve a parking spot (Concurrency-protected)
- `GET /api/v1/reservations/my-reservations` (Authenticated) - View own reservations
- `DELETE /api/v1/reservations/:id` (Authenticated) - Cancel own reservation
- `GET /api/v1/reservations` (Admin Only) - View all system reservations
