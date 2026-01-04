# 🎯 Kanban Board - Monorepo

Modern Kanban Board application built with Next.js (Frontend) and Golang Microservices (Backend).

> **� Cara run aplikasi: [RUN.md](./RUN.md)**

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
- **Database**: MySQL 8.0+
- **Cache**: Redis 7+

## 🚀 Quick Start

### Manual Setup
```bash
# Clone repository
git clone https://github.com/NathanzZ-Ginting/kanban-board.git
cd kanban-monorepo

# Setup services
chmod +x start.sh
./start.sh
```

See complete guide at [SETUP_GUIDE.md](./SETUP_GUIDE.md)

### Access Application
- **Frontend**: http://localhost:3000
- **Auth Service**: http://localhost:8001
- **User Service**: http://localhost:8002
- **Board Service**: http://localhost:8003
- **Task Service**: http://localhost:8004

## 📚 Documentation
- **Setup Guide**: [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Panduan lengkap setup & running
- **API Documentation**: Coming soon
- **Architecture**: Coming soon

## Project Structure

```
kanban-monorepo/
├── frontend/                 # Next.js application
├── backend/                  # Golang microservices
│   ├── auth-service/
│   ├── user-service/
│   ├── board-service/
│   └── task-service/
├── docs/                     # Documentation
└── scripts/                  # Utility scripts
```

## 🔧 Commands Cheatsheet

```bash
# Development
cd frontend && npm run dev        # Run frontend
cd backend/auth-service && go run cmd/main.go  # Run auth service

# Database
mysql -u kanban_user -p  # Access MySQL
```

Untuk command lengkap, lihat [SETUP_GUIDE.md](./SETUP_GUIDE.md)

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- MySQL 8.0+

### Quick Start

Ikuti langkah di [SETUP_GUIDE.md](./SETUP_GUIDE.md) untuk:
- ✅ Setup manual
- ✅ Environment variables
- ✅ Troubleshooting
- ✅ API endpoints

## Environment Variables

Buat file `.env` di root directory. Lihat detail lengkap di [SETUP_GUIDE.md](./SETUP_GUIDE.md#environment-variables)

```env
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
MYSQL_HOST=mysql
MYSQL_USER=kanban_user
MYSQL_PASSWORD=kanban_pass
MYSQL_DATABASE=kanban_db
```

## 🏗️ Architecture

### Backend Services
- **Auth Service** (Port 8001): Authentication & authorization
- **User Service** (Port 8002): User management
- **Board Service** (Port 8003): Board/project management
- **Task Service** (Port 8004): Task management

### API Documentation
- Swagger UI available at: `http://localhost:8080/swagger`
- Detail API endpoints di [SETUP_GUIDE.md](./SETUP_GUIDE.md#api-endpoints)

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines.

## 📄 License

MIT License - see LICENSE file for details.

---

Made with ❤️ by Nathan Ginting
