# 🚀 Cara Run Kanban Board

Guide singkat untuk menjalankan aplikasi.

---

## 📋 Prerequisites

- Node.js 18+
- Go 1.21+
- MySQL 8.0+
- Git

---

## ⚡ Quick Start

### 1. Clone & Setup
```bash
git clone https://github.com/NathanzZ-Ginting/kanban-board.git
cd kanban-monorepo
```

### 2. Buat File Environment
Buat file `.env` di root folder:
```env
JWT_SECRET=your-secret-key-change-this
JWT_EXPIRY=24h
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=kanban_db
MYSQL_USER=kanban_user
MYSQL_PASSWORD=kanban_pass
```

### 3. Setup Database
```bash
# Install MySQL dan jalankan
# Buat database
mysql -u root -p
CREATE DATABASE kanban_db;
CREATE USER 'kanban_user'@'localhost' IDENTIFIED BY 'kanban_pass';
GRANT ALL PRIVILEGES ON kanban_db.* TO 'kanban_user'@'localhost';
FLUSH PRIVILEGES;
```

### 4. Run Backend Services
```bash
# Auth Service
cd backend/auth-service
go mod download
go run cmd/main.go

# User Service (terminal baru)
cd backend/user-service
go mod download
go run cmd/main.go

# Board Service (terminal baru)
cd backend/board-service
go mod download
go run cmd/main.go

# Task Service (terminal baru)
cd backend/task-service
go mod download
go run cmd/main.go
```

### 5. Run Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## 🌐 Access Application

- **Frontend**: http://localhost:3000
- **Auth Service**: http://localhost:8001
- **User Service**: http://localhost:8002
- **Board Service**: http://localhost:8003
- **Task Service**: http://localhost:8004

---

## 🛑 Stop Services

Press `Ctrl+C` in each terminal to stop the services.

---

## 🐛 Troubleshooting

### Port sudah dipakai?
```bash
# Check port
lsof -i :3000

# Kill process
kill -9 <PID>
```

### MySQL error?
```bash
# Restart MySQL service
# macOS
brew services restart mysql

# Linux
sudo systemctl restart mysql
```

### Go dependencies error?
```bash
# Clean and reinstall
cd backend/auth-service
go clean -modcache
go mod download
```

---

**Selesai! Happy coding! 🎉**
