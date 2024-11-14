# Linktree Clone

A full-stack application that allows users to create their own Linktree-style landing page with multiple links. Built with Vue.js, Gin (Go), and PostgreSQL.

## 🌟 Features

- User authentication (signup/login)
- Create and manage personal profile
- Add, edit, and delete links
- Track link click analytics
- Responsive design
- Swagger API documentation
- Kubernetes deployment support

## 🏗️ Technology Stack

### Frontend

- Vue.js 3
- Vite
- Tailwind CSS
- Axios for API calls

### Backend

- Go (Gin framework)
- GORM for database operations
- JWT for authentication
- Swagger for API documentation

### Database

- PostgreSQL 17

### DevOps

- Docker
- Kubernetes
- Helm

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- Node.js 20+
- PostgreSQL 17+
- Docker
- Kubernetes cluster
- Helm 3+

### Local Development Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/codescalersinternships/Linktree-MohamedFadel.git
   cd Linktree-MohamedFadel
   ```

2. **Backend Setup**

   ```bash
   cd backend

   # Create .env file

   # Install dependencies
   go mod download

   # Run the server
   go run cmd/server/main.go
   ```

3. **Frontend Setup**

   ```bash
   cd frontend

   # Install dependencies
   npm install

   # Run development server
   npm run dev
   ```

4. **Database Setup**
   ```bash
   # The application will automatically create the necessary tables
   # Just ensure PostgreSQL is running and the credentials in .env are correct
   ```

### Docker Deployment

1. **Build and run using Docker Compose**
   ```bash
   docker-compose up --build
   ```

### Kubernetes Deployment using Helm

1. **Build and push Docker images**

   ```bash
   # Build and push frontend
   docker build -t your-registry/frontend:latest ./frontend
   docker push your-registry/frontend:latest

   # Build and push backend
   docker build -t your-registry/backend:latest ./backend
   docker push your-registry/backend:latest
   ```

2. **Update values.yaml with your image repositories**

   ```yaml
   frontend:
     image:
       repository: your-registry/frontend
       tag: latest

   backend:
     image:
       repository: your-registry/backend
       tag: latest
   ```

3. **Deploy using Helm**
   ```bash
   helm install linktree-app ./helm
   ```

## 🔧 Configuration

### Backend Environment Variables

```env
DB_HOST=database
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=linktree
JWT_SECRET=your-secure-jwt-secret
```

### Frontend Environment Variables

```env
VITE_API_URL=http://your-backend-url/api
```

## 🔒 Security

- JWT-based authentication
- Password hashing
- CORS configuration
- Protected routes
- Environment variable configuration

## 🌐 Production Deployment

The application is designed to be deployed on Kubernetes using Helm charts. The deployment includes:

- Separate services for frontend, backend, and database
- LoadBalancer service types for external access
- Persistent volume for PostgreSQL data
- Environment variable configuration
