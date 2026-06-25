# 🚗 SpotSync – Smart Parking & EV Charging Reservation

SpotSync is a centralized, high-performance platform designed for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots with zero double-booking or concurrency race-conditions.

---

## 🚀 Features

- **Centralized Parking Management**: Supports multiple parking zone types (General, EV Charging, Covered).
- **Concurrency-Safe Reservations**: Implements row-level locking (`SELECT ... FOR UPDATE`) inside database transactions to completely prevent double-booking.
- **JWT & Role-Based Security**: Secured endpoints with customizable middlewares for `driver` and `admin` roles.
- **Centralized Error Handling**: A global HTTP error handler that maps custom sentinel errors to clean HTTP responses and hides internal database stack traces.
- **Dynamic Availability Tracking**: Automatically calculates real-time `available_spots` for each zone by subtracting active reservations.

---

## 🛠️ Technology Stack

| Component | Technology | Description / Note |
| --- | --- | --- |
| **Backend Language** | Go (Golang) | Version 1.22+ |
| **Web Framework** | Echo v4 | High-performance, minimalist framework |
| **ORM** | GORM | PostgreSQL database interaction |
| **Database** | PostgreSQL | Relational database (NeonDB, Supabase, or local) |
| **Validation** | Validator v10 | Struct validation integrated with Echo context |
| **Authentication** | JWT & Bcrypt | `golang-jwt/jwt/v5` and `bcrypt` with cost 10 |

---

## 🏛️ Clean Architecture Layer Flow

The project is structured under strict Separation of Concerns:

```mermaid
graph TD
    Client[HTTP Client] --> Routing[router/router.go]
    Routing --> Middleware[middleware/auth.go]
    Middleware --> Handlers[handler/handler.go]
    Handlers --> DTO[dto/dto.go]
    Handlers --> Services[service/service.go]
    Services --> Repository[repository/repository.go]
    Repository --> Models[models/models.go]
    Repository --> DB[(PostgreSQL)]
    Handlers -.-> Errors[utils/error_handler.go]
```

### Layer Responsibilities
1. **HTTP Layer (Handlers & Middlewares)**:
   - Middleware handles JWT parsing, token validation, and role assignment.
   - Handlers bind and validate incoming DTO requests, invoke the appropriate Service method, and return standardized JSON envelopes.
2. **Business Logic Layer (Services)**:
   - Orchestrates use-cases, applies validation rules, encrypts passwords, checks business integrity, and coordinates mappings.
3. **Data Access Layer (Repositories & Models)**:
   - Houses isolated SQL/GORM queries. It uses transactions and row locks to interact with GORM Models. No HTTP-related concepts exist in this layer.

---

## ⚙️ Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://username:password@host:port/database?sslmode=require
JWT_SECRET=your_jwt_secret_key
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

---

## 💻 Local Setup & Execution

### Prerequisites
- Go 1.22 or higher installed.
- PostgreSQL database instance running.

### Installation Steps
1. Clone the repository and navigate to the project directory.
2. Create and configure your `.env` file:
   ```bash
   cp .env.example .env
   ```
3. Fetch Go module dependencies:
   ```bash
   go mod tidy
   ```
4. Start the server (which automatically runs GORM database migrations):
   ```bash
   go run main.go
   ```
5. The API will start on the port specified in `.env` (default is `8080`). You can ping `GET http://localhost:8080/health` to verify.

---

## 🌐 API Endpoint Directory

All responses are wrapped in a standard JSON wrapper:
- **Success Response**: `{ "success": true, "message": "Success message", "data": ... }`
- **Error Response**: `{ "success": false, "message": "Reason", "errors": ... }`

| Method | Endpoint | Access | Description | Request Body Example |
| :--- | :--- | :--- | :--- | :--- |
| **POST** | `/api/v1/auth/register` | Public | Register driver/admin | `{ "name": "Alice", "email": "alice@gmail.com", "password": "secure123", "role": "driver" }` |
| **POST** | `/api/v1/auth/login` | Public | Verify credentials & obtain JWT | `{ "email": "alice@gmail.com", "password": "secure123" }` |
| **POST** | `/api/v1/zones` | Admin Only | Create a new parking zone | `{ "name": "Zone A", "type": "ev_charging", "total_capacity": 10, "price_per_hour": 150.00 }` |
| **GET** | `/api/v1/zones` | Public | Get all zones (w/ dynamic availability) | *None* |
| **GET** | `/api/v1/zones/:id` | Public | Details of a single zone | *None* |
| **POST** | `/api/v1/reservations` | Authenticated | Create a locked safe booking | `{ "zone_id": 1, "license_plate": "DHAKA-12345" }` |
| **GET** | `/api/v1/reservations/my-reservations` | Authenticated | Retrieve user's own reservations | *None* |
| **DELETE**| `/api/v1/reservations/:id` | Authenticated | Cancel a reservation (Owner/Admin only) | *None* |
| **GET** | `/api/v1/reservations` | Admin Only | Get all bookings in the system | *None* |

---

## ☁️ Cloud Deployment Steps (Render & Railway)

### 1. Database Provisioning
Provision a PostgreSQL database using cloud providers such as **NeonDB** or **Supabase**. Obtain the transaction connection string.

### 2. Deployment Configurations

#### Railway Deployment
1. Connect your GitHub repository to Railway.
2. In the service settings, define the following variables under **Variables**:
   - `PORT` = `8080` (or leave it to Railway to inject automatically)
   - `DB_URL` = *Your PostgreSQL Connection String* (e.g., `postgresql://...`)
   - `JWT_SECRET` = *Secure random string*
   - `ALLOWED_ORIGINS` = `*` (or your frontend domain URL)
3. Railway automatically detects Go files and deploys. If custom build is required:
   - **Build Command**: `go build -o bin/server main.go`
   - **Start Command**: `./bin/server`

#### Render Deployment
1. Create a new **Web Service** on Render and link your GitHub repository.
2. Select runtime environment as **Go**.
3. Configure the following build settings:
   - **Build Command**: `go build -o server main.go`
   - **Start Command**: `./server`
4. Under **Advanced**, add the environment variables:
   - `DB_URL`
   - `JWT_SECRET`
   - `PORT` (e.g., `10000`)
   - `ALLOWED_ORIGINS`
5. Click **Deploy Web Service**.
