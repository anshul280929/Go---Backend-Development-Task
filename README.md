# 🧪 User Management API (Go + Fiber + PostgreSQL)

A high-performance, containerized, and production-ready RESTful API for managing users with their **name** and **date of birth (DOB)**. The API dynamically calculates and returns a user's **age** when fetching their details, while ensuring strict data validation and structured logging.

---

## ✨ Features

- **CRUD Operations**: Manage user profiles with database persistence.
- **Dynamic Age Calculation**: Dynamically computes user age based on DOB, correctly handling leap years.
- **Pagination & Metadata**: Supports `GET /users` with cursor/page parameters (`page` and `page_size`) and returns metadata.
- **Strict Validation**: Request validation using `go-playground/validator` (valid dates, proper name bounds).
- **Structured Logging**: Uses `uber-go/zap` for fast, structured logging.
- **Request Tracing**: Middleware to inject and log a unique Request ID (`X-Request-ID`) for every API call.
- **Database Code-Gen**: Uses `sqlc` for compile-time safe SQL code generation.
- **Docker Ready**: Multi-stage `Dockerfile` and a configured `docker-compose.yml` for single-command deployment.

---

## 🛠️ Tech Stack

| Technology | Purpose |
| :--- | :--- |
| **[Go 1.22+](https://go.dev/)** | Core programming language |
| **[GoFiber v2](https://gofiber.io/)** | High-performance, low-overhead HTTP web framework |
| **[PostgreSQL 16](https://www.postgresql.org/)** | Relational database storage |
| **[SQLC](https://sqlc.dev/)** | Type-safe SQL compiler to generate Go code from SQL queries |
| **[pgx v5](https://github.com/jackc/pgx)** | High-performance PostgreSQL driver and connection pool |
| **[Uber Zap](https://github.com/uber-go/zap)** | Structured, fast logging |
| **[go-playground/validator](https://github.com/go-playground/validator)** | Input payload validation |
| **[Docker](https://www.docker.com/)** | Application containerization and service orchestration |

---

## 📁 Project Structure

```text
├── cmd/
│   └── server/
│       └── main.go             # Application entrypoint & dependency injection
├── config/
│   └── config.go           # Environment variables configuration loader
├── db/
│   ├── migrations/         # SQL migration scripts (Schema definition)
│   └── sqlc/               # SQLC config, raw SQL queries, and generated code
├── internal/
│   ├── handler/            # HTTP controller / endpoint handlers (Fiber)
│   ├── logger/             # Uber Zap logger wrapper initialization
│   ├── middleware/         # Custom Fiber middlewares (logger, request_id)
│   ├── models/             # DTOs, request/response models, pagination structure
│   ├── repository/         # DB repository layer (interacts with SQLC queries)
│   ├── routes/             # API routing setup
│   └── service/            # Business logic layer & validation rules
├── .env.example            # Environment variables template
├── Dockerfile              # Multi-stage production container build
├── docker-compose.yml      # Service docker definition (App + PostgreSQL)
├── go.mod                  # Go module definition
└── README.md               # Documentation
```

---

## 🚀 Getting Started

### Option 1: Using Docker Compose (Recommended)

The easiest way to run the entire stack (API server + PostgreSQL) is via Docker Compose:

1. Clone the repository and navigate to the directory:
   ```bash
   git clone https://github.com/anshul280929/Go---Backend-Development-Task.git
   cd Go---Backend-Development-Task
   ```
2. Build and start the services:
   ```bash
   docker-compose up --build
   ```
3. The API will start on **`http://localhost:3000`**, and PostgreSQL will run on **`localhost:5432`** (with migrations auto-applied on startup).

---

### Option 2: Local Development Setup

To run the API and PostgreSQL manually without Docker:

1. **Setup Database**:
   - Ensure a PostgreSQL server is running.
   - Create a database:
     ```sql
     CREATE DATABASE user_management;
     ```
   - Run the migration script to create the users table:
     ```bash
     psql -U postgres -d user_management -f db/migrations/000001_create_users.up.sql
     ```

2. **Configure Environment Variables**:
   - Copy `.env.example` to `.env`:
     ```bash
     cp .env.example .env
     ```
   - Adjust `.env` connection strings for your local database:
     ```env
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=postgres
     DB_PASSWORD=your_password
     DB_NAME=user_management
     SERVER_PORT=3000
     ```

3. **Install Dependencies & Start**:
   - Download the modules:
     ```bash
     go mod tidy
     ```
   - Run the application:
     ```bash
     go run cmd/server/main.go
     ```
   - The application will be listening on `http://localhost:3000`.

---

## 🚦 API Endpoints

### 1. Health Check
Checks if the application server is up.
- **URL**: `GET /health`
- **Response (`200 OK`)**:
  ```json
  {
    "status": "ok"
  }
  ```

---

### 2. Create User
Creates a new user profile.
- **URL**: `POST /users`
- **Request Body**:
  ```json
  {
    "name": "Bruce Wayne",
    "dob": "1939-05-27"
  }
  ```
- **Response (`201 Created`)**:
  ```json
  {
    "id": 1,
    "name": "Bruce Wayne",
    "dob": "1939-05-27"
  }
  ```

---

### 3. Get User by ID
Fetches a user profile. Calculates age dynamically in the response.
- **URL**: `GET /users/:id`
- **Response (`200 OK`)**:
  ```json
  {
    "id": 1,
    "name": "Bruce Wayne",
    "dob": "1939-05-27",
    "age": 87
  }
  ```

---

### 4. List Users (Paginated)
Retrieves a paginated list of all users.
- **URL**: `GET /users`
- **Query Parameters**:
  - `page`: Page number (Default: `1`)
  - `page_size`: Number of records per page (Default: `10`, Max: `100`)
- **Response (`200 OK`)**:
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "Bruce Wayne",
        "dob": "1939-05-27",
        "age": 87
      }
    ],
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  }
  ```

---

### 5. Update User
Updates an existing user's details.
- **URL**: `PUT /users/:id`
- **Request Body**:
  ```json
  {
    "name": "Batman",
    "dob": "1939-05-27"
  }
  ```
- **Response (`200 OK`)**:
  ```json
  {
    "id": 1,
    "name": "Batman",
    "dob": "1939-05-27"
  }
  ```

---

### 6. Delete User
Removes a user from the database.
- **URL**: `DELETE /users/:id`
- **Response (`204 No Content`)**: *(Empty Response Body)*

---

## 🧪 Testing

The project contains unit tests verifying edge cases in the dynamic age calculation (including leap-year support).

Run all tests:
```bash
go test ./... -v
```

---

## 🔄 Regenerating Database Code (SQLC)

If you modify the database schema in `db/migrations/` or write new SQL queries in `db/sqlc/query.sql`:

1. Ensure `sqlc` is installed:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```
2. Run generation from the root directory:
   ```bash
   sqlc generate
   ```
This updates the models and queries under `db/sqlc/` automatically.
