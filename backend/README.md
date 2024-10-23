# Linktree API

A RESTful API service built with Go that provides functionality similar to Linktree, allowing users to create and manage a collection of personal links with analytics tracking.

## 🚀 Features

- User authentication and management
- Link creation and management
- Click tracking and analytics
- JWT-based authentication
- Swagger documentation
- Docker deployment support

## 🛠️ Tech Stack

- **Go** - Programming language
- **Gin** - Web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Docker** - Containerization
- **Swagger** - API documentation

## 📁 Project Structure

```
.
├── cmd
│   └── server
│       └── main.go            # Application entry point
├── docker-compose.yaml        # Docker compose configuration
├── Dockerfile                 # Docker build instructions
├── docs                       # Swagger documentation
├── internal
│   ├── api
│   │   ├── handlers          # Request handlers
│   │   ├── middleware        # Custom middleware
│   │   └── routes.go         # Route definitions
│   ├── database              # Database configuration
│   ├── models                # Data models
│   ├── services              # Business logic
│   └── utils                 # Helper functions
└── tests                     # Test files
```

## 🔧 Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Docker and Docker Compose (optional)

## ⚙️ Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret
```

## 🚀 Getting Started

### Running with Docker

1. Clone the repository:

```bash
git clone https://github.com/codescalersinternships/Linktree-MohamedFadel
cd Linktree-MohamedFadel
```

2. Start the application using Docker Compose:

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8188`

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Start PostgreSQL server

3. Run the application:

```bash
go run cmd/server/main.go
```

## 📚 API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8188/swagger/index.html
```

### 🔑 Key Endpoints

#### Users

- `POST /api/v1/users/signup` - Create new user account
- `POST /api/v1/users/login` - User login
- `GET /api/v1/users/:username` - Get user profile
- `PUT /api/v1/users` - Update user profile
- `DELETE /api/v1/users` - Delete user account

#### Links

- `POST /api/v1/links` - Create new link
- `PUT /api/v1/links/:id` - Update existing link
- `DELETE /api/v1/links/:id` - Delete link

#### Analytics

- `POST /api/v1/analytics/:id/click` - Track link click

## 🔒 Authentication

The API uses JWT for authentication. To access protected endpoints:

1. Obtain a token through the login endpoint
2. Include the token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## 🧪 Running Tests

Execute the test suite:

```bash
go test ./tests/...
```

## 📦 Project Dependencies

Key dependencies include:

- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/swaggo/swag` - Swagger documentation
- `github.com/joho/godotenv` - Environment configuration
