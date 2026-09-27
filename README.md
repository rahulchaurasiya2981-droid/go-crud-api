<div align="center">

# 🚀 Go CRUD API

A production-ready RESTful CRUD API built with **Go** and **PostgreSQL** for managing users. Built following clean architecture principles, package-by-feature organization, and idiomatic Go design patterns.

[![Go Version](https://img.shields.io/badge/Go-1.27%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18%2B-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Package--by--Feature-orange?style=for-the-badge)](#-architecture)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

[Features](#-features) • [Tech Stack](#-tech-stack) • [Getting Started](#-getting-started) • [Available Tasks](#-available-tasks) • [API Endpoints](#-api-endpoints) • [Project Structure](#-project-structure) • [Architecture](#-architecture)

</div>

---

## 🎬 Demo

Check out the full walkthrough of the project, package structure breakdown, database migrations, and live API endpoints demonstration on YouTube:

> 📺 **Watch video on YouTube**: [Go CRUD API with PostgreSQL | Complete Backend Project Demo](https://www.youtube.com/watch?v=Lip1fP1kHDE)

## ✨ Features

- ⚡ **RESTful API**: Clean API endpoints for user CRUD operations (Create, Read, Update, Delete).
- 📦 **Package-by-Feature Architecture**: High modularity by organizing code around business domains.
- 🔀 **DTO & Entity Separation**: Strict boundary between API request/response contracts and internal models.
- 🛡️ **Robust HTTP Request Parsing & Validation**:
  - Content-Length & JSON payload size limits.
  - Rejection of unknown JSON fields.
  - Trailing JSON data detection.
  - DTO field-level validation.
- 📜 **Structured Logging**: Centralized logging powered by standard library `log/slog`.
- 🗄️ **Database Migrations**: Version-controlled SQL migrations using `golang-migrate`.
- ⚡ **Context-Aware Database Queries**: Timeout and cancellation propagation down to `pgx v5`.
- 🛑 **Graceful Server Shutdown**: Handles OS termination signals without dropping active HTTP requests.
- 🌐 **Standardized Responses**: Centralized HTTP success and error response wrappers.
- ⚙️ **Environment Configuration**: Flexible env loading with `godotenv`.
- 🛠️ **Automated Task Runner**: One-step migration and app start via `Taskfile`.

---

## 🛠️ Tech Stack

| Domain | Technology / Library | Description |
| :--- | :--- | :--- |
| **Language** | [Go 1.27+](https://golang.org/) | Core language |
| **HTTP Router** | Standard Library (`net/http`) | Lightweight HTTP server implementation |
| **Database** | [PostgreSQL 18+](https://www.postgresql.org/) | Relational database storage |
| **DB Driver / Pool** | [`database/sql`](https://pkg.go.dev/database/sql) & [`pgx v5`](https://github.com/jackc/pgx) | Native Go PostgreSQL driver & pooling |
| **Migrations** | [`golang-migrate`](https://github.com/golang-migrate/migrate) | Database schema versioning |
| **Logging** | [`log/slog`](https://pkg.go.dev/log/slog) | Native structured logging |
| **Environment** | [`godotenv`](https://github.com/joho/godotenv) | `.env` file loader |
| **Task Runner** | [`Task`](https://taskfile.dev/) | Task automation runner |

---

## 🚀 Getting Started

### 📋 Prerequisites

Ensure you have the following installed on your local environment:

- 🟢 **Go**: `v1.27` or higher
- 🐘 **PostgreSQL**: `v18` or higher
- 🐙 **Git**
- 🛠️ **Taskfile** *(optional, but recommended)*
- 🟧 **Postman** *(optional, for API testing)*

---

### 📥 1. Clone the Repository

```bash
git clone [https://github.com/your-username/go-crud-api.git](https://github.com/your-username/go-crud-api.git)
cd go-crud-api
