# Go CRUD API

A production-style RESTful CRUD API built with Go and PostgreSQL for managing users. The project follows a package-by-feature architecture with separate DTO, entity, handler, service, and repository layers, along with centralized configuration, database, logging, request validation, and HTTP response handling.

![Go CRUD API Demo](./docs/demo.gif)

## Features

* RESTful CRUD operations for users
* Create, read, update, and delete users
* PostgreSQL database integration
* Package-by-feature architecture
* DTO and entity separation
* Request body validation
* JSON request size limitation
* Unknown JSON field validation
* Trailing JSON data validation
* Structured logging with `log/slog`
* Database migrations
* Context-aware database queries
* Graceful server shutdown
* Centralized HTTP response handling
* Environment-based configuration
* Taskfile for development commands

## Tech Stack

* **Go 1.27+**
* **net/http**
* **PostgreSQL 18+**
* **database/sql**
* **pgx v5**
* **golang-migrate**
* **godotenv**
* **log/slog**
* **Task**

## Requirements

Before running the project, make sure you have:

* Go 1.27+
* PostgreSQL 18+
* Git
* Task
* Postman (optional, for API testing)

## How to Use

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-crud-api
```

### 2. Configure Environment Variables

Create a `.env` file using `.env.example` as a reference.

```bash
cp .env.example .env
```

Update the values according to your local PostgreSQL configuration.

> **Note:** Never commit your `.env` file. The `.env.example` file is provided as a configuration reference.

### 3. Install Dependencies

Using Task:

```bash
task tidy
```

Or using Go:

```bash
go mod download
```

### 4. Run Database Migrations

Apply the database migrations:

```bash
task migrate-up
```

### 5. Start the Server

```bash
task run
```

The API will be available at:

```text
http://localhost:8080
```

## Available Tasks

| Command                 | Description                    |
| ----------------------- | ------------------------------ |
| `task run`              | Start the application          |
| `task build`            | Build the application          |
| `task test`             | Run tests                      |
| `task fmt`              | Format Go code                 |
| `task vet`              | Run Go vet                     |
| `task tidy`             | Update and clean dependencies  |
| `task check`            | Run project checks             |
| `task migrate-up`       | Apply database migrations      |
| `task migrate-down`     | Roll back the latest migration |
| `task migrate-down-all` | Roll back all migrations       |

## API Endpoints

| Method   | Endpoint      | Description       |
| -------- | ------------- | ----------------- |
| `GET`    | `/health`     | Check API health  |
| `GET`    | `/users`      | Get all users     |
| `GET`    | `/users/{id}` | Get a user by ID  |
| `POST`   | `/users`      | Create a new user |
| `PUT`    | `/users/{id}` | Update a user     |
| `DELETE` | `/users/{id}` | Delete a user     |

## Example Request

### Create User

```http
POST /users
Content-Type: application/json
```

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 25
}
```

## Project Structure

```text
go-crud-api/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   │
│   ├── user/
│   │   ├── dto/
│   │   │   └── user_dto.go
│   │   │
│   │   ├── entity/
│   │   │   └── user.go
│   │   │
│   │   ├── handler/
│   │   │   └── user_handler.go
│   │   │
│   │   ├── service/
│   │   │   └── user_service.go
│   │   │
│   │   └── repository/
│   │       └── user_repository.go
│   │
│   ├── config/
│   │   └── ...
│   │
│   ├── database/
│   │   └── ...
│   │
│   ├── logger/
│   │   └── ...
│   │
│   ├── httprequest/
│   │   └── ...
│   │
│   └── httpresponse/
│       └── ...
│
├── migrations/
│
├── docs/
│   └── demo.gif
│
├── .env.example
├── .gitignore
├── Taskfile.yml
├── go.mod
├── go.sum
└── README.md
```

## Architecture

The project follows a **package-by-feature** architecture.

```text
                         HTTP Request
                              │
                              ▼
                       ┌─────────────┐
                       │   Handler   │
                       └──────┬──────┘
                              │
                              ▼
                       ┌─────────────┐
                       │   Service   │
                       └──────┬──────┘
                              │
                              ▼
                       ┌─────────────┐
                       │ Repository  │
                       └──────┬──────┘
                              │
                              ▼
                         PostgreSQL
```

### User Module

The `user` package contains everything related to user functionality.

* **DTO** — Defines API request and response data structures.
* **Entity** — Represents the core user data model.
* **Handler** — Handles HTTP requests, validation, and responses.
* **Service** — Contains application and business logic.
* **Repository** — Handles database operations.

### Shared Packages

* **config** — Loads and manages application configuration.
* **database** — Handles database connection and configuration.
* **logger** — Provides centralized structured logging.
* **httprequest** — Provides common HTTP request parsing and validation.
* **httpresponse** — Provides standardized HTTP success and error responses.

### Server

The `cmd/server` package contains the application entry point.

It is responsible for initializing configuration, database connections, dependencies, routes, and the HTTP server.

## Database

The application uses PostgreSQL for persistent data storage.

Database schema changes are managed using **golang-migrate**.

### Apply Migrations

```bash
task migrate-up
```

### Roll Back Latest Migration

```bash
task migrate-down
```

### Roll Back All Migrations

```bash
task migrate-down-all
```

## Environment Variables

The project provides an example environment configuration:

```text
.env.example
```

Create your local environment file:

```text
.env
```

Update the values according to your local environment.

> **Important:** Do not commit `.env` to the repository.

## Testing

Run the test suite:

```bash
task test
```

Run project checks:

```bash
task check
```

## Demo

A demo GIF showing the API in action:

![Go CRUD API Demo](./docs/demo.gif)

A YouTube demo can be added here later.

## License

This project is built for learning and educational purposes.
