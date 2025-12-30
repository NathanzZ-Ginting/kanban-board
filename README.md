# Kanban Board - Monorepo

Modern Kanban Board application built with Next.js (Frontend) and Golang Microservices (Backend).

## Tech Stack

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **State Management**: Zustand
- **UI Framework**: Tailwind CSS + shadcn/ui
- **HTTP Client**: Axios
- **Architecture**: Component-based with Atomic Design

### Backend
- **Language**: Golang 1.21+
- **Architecture**: Clean Architecture + Microservices
- **Database**: MongoDB
- **Pattern**: Repository Pattern
- **API**: REST API with JSON

### Infrastructure
- **Container**: Docker + Docker Compose
- **Database**: MongoDB 7.0+

## Project Structure

```
kanban-monorepo/
├── frontend/                 # Next.js application
├── backend/                  # Golang microservices
│   ├── auth-service/
│   ├── user-service/
│   ├── board-service/
│   └── task-service/
├── docker/                   # Docker configurations
├── docs/                     # Documentation
└── scripts/                  # Utility scripts
```

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- MongoDB 7.0+

### Development

1. **Clone repository**
```bash
git clone <repo-url>
cd kanban-monorepo
```

2. **Setup environment variables**
```bash
cp .env.example .env
```

3. **Run with Docker Compose**
```bash
docker-compose up -d
```

4. **Run locally**

Frontend:
```bash
cd frontend
npm install
npm run dev
```

Backend (each service):
```bash
cd backend/auth-service
go mod download
go run cmd/main.go
```

## Environment Variables

See `.env.example` for required environment variables.

## Architecture

### Backend Services
- **Auth Service** (Port 8001): Authentication & authorization
- **User Service** (Port 8002): User management
- **Board Service** (Port 8003): Board/project management
- **Task Service** (Port 8004): Task management

### API Documentation
- Swagger UI available at: `http://localhost:8080/swagger`

## License

MIT
