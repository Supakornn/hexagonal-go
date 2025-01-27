# 🛍️ E-Commerce REST APIs With GO

A modern e-commerce backend service built with Hexagonal Architecture in Go

## 📌 Features

- 🔐 Authentication & Authorization
- 👥 User Management
- 📦 Product Management
- 🛒 Order Processing
- 💳 Payment Integration
- 🏷️ Category Management

## 🚀 Technologies

- 🐹 Golang (Fiber) - Fast HTTP web framework
- 🐳 Docker - Containerization
- 🐘 PostgreSQL - Database

## 🏗️ Architecture

This project follows Hexagonal Architecture (Ports and Adapters) pattern:

- `core` - Business logic and domain models
- `handlers` - HTTP handlers (API endpoints)
- `repositories` - Database interactions
- `services` - Application services

## 📊 Database Schema

<a href="https://dbdiagram.io/d/hexagonal-go-66d444c6eef7e08f0e57c865"><img src="./pictures/db_diagram.png" alt="Database Schema"/></a>

## ⚙️ Getting Started

### Prerequisites

- Go 1.20 or higher
- Docker & Docker Compose
- PostgreSQL

### 🔧 Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/hexagonal-go.git
```

2. Navigate to project directory

```bash
cd hexagonal-go
```

3. Run with Docker

```bash
docker-compose up -d
```

## 🔧 Environment Setup

```bash
# Required environment variables
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=ecommerce
JWT_SECRET=your-secret-key
```

## 🛣️ API Endpoints

### Users API

```http
POST /api/v1/auth/register    # Register new user
POST /api/v1/auth/login       # Login user
GET /api/v1/users/profile     # Get user profile
```

### Products API

```http
GET /api/v1/products         # List all products
POST /api/v1/products        # Create product
GET /api/v1/products/:id     # Get product details
PUT /api/v1/products/:id     # Update product
DELETE /api/v1/products/:id  # Delete product
```

### Orders API

```http
POST /api/v1/orders         # Create order
GET /api/v1/orders          # List user orders
GET /api/v1/orders/:id      # Get order details
```

## 🏗️ Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── core/
│   ├── handlers/
│   ├── repositories/
│   └── services/
├── tests/
├── docker-compose.yml
└── README.md
```
