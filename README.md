# 🎯 Kanban Board

Kanban Board dengan Next.js & Golang Microservices + Supabase.

## 🚀 Quick Start

### 1. Setup Supabase
1. Buat project di [supabase.com](https://supabase.com)
2. Copy connection string dari Settings > Database

### 2. Environment Variables
Buat file `.env`:
```env
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
SUPABASE_DB_HOST=your-project.supabase.co
SUPABASE_DB_PORT=5432
SUPABASE_DB_USER=postgres
SUPABASE_DB_PASSWORD=your-password
SUPABASE_DB_NAME=postgres
```

### 3. Run Services

**Option 1: Manual (Recommended for Development)**
```bash
# Terminal 1 - Auth Service
cd backend/auth-service && go run cmd/main.go    # :8001

# Terminal 2 - User Service  
cd backend/user-service && go run cmd/main.go    # :8002

# Terminal 3 - Board Service
cd backend/board-service && go run cmd/main.go   # :8003

# Terminal 4 - Task Service
cd backend/task-service && go run cmd/main.go    # :8004

# Terminal 5 - Frontend
cd frontend && npm install && npm run dev         # :3000
```

**Option 2: Using start script (Background)**
```bash
./start.sh  # Start all backend services in background
cd frontend && npm run dev  # Start frontend
```

### 4. Access
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8001-8004

**First Time Setup:**
1. Register akun baru di http://localhost:3000/register
2. Login dengan akun yang baru dibuat
3. Create board dan mulai menggunakan aplikasi

**Jika ada masalah:**
```bash
./fix-issues.sh  # Tool untuk diagnosa dan fix masalah umum
```

## 🛠️ Tech Stack
- **Frontend**: Next.js 14, TypeScript, Tailwind CSS, Zustand
- **Backend**: Golang 1.21+, Clean Architecture
- **Database**: Supabase (PostgreSQL)

## 🔧 Troubleshooting

### Error: "prepared statement already exists (SQLSTATE 42P05)"
✅ **Fixed!** Prepared statements sudah dinonaktifkan di semua services.

Jika masih terjadi, restart services:
```bash
# Stop all services
pkill -f "go run cmd/main.go"

# Start again manually atau dengan ./start.sh
```

### Error: "relation already exists (SQLSTATE 42P07)"
✅ **Fixed!** Services sekarang akan skip migration error yang tidak critical.

### Error: "Failed to load board data" di Frontend

**Penyebab umum:**
1. **Token expired atau invalid** - Login ulang
2. **Backend service belum running** - Pastikan semua services running
3. **CORS issues** - Check browser console

**Cara debug:**
1. Buka browser DevTools (F12) → Console tab
2. Cari error messages:
   - "No access token found" → Login ulang
   - "401 Unauthorized" → Token expired, login ulang
   - "404 Not Found" → Board ID tidak valid

**Solusi:**
```bash
# 1. Logout dan login ulang di frontend
# 2. Check localStorage di browser DevTools:
#    - Application → Local Storage → http://localhost:3000
#    - Cari 'accessToken' dan 'refreshToken'

# 3. Test API manually:
./test-api.sh

# 4. Jika perlu, register user baru:
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "username": "testuser",
    "password": "test12345",
    "firstName": "Test",
    "lastName": "User"
  }'
```

### Restart Single Service
```bash
# Stop service (Ctrl+C di terminal yang menjalankan service)
# Lalu run ulang:
cd backend/board-service && go run cmd/main.go
```

### Clean Database (Reset semua tabel)
Jika ingin reset database, hapus semua tabel di Supabase Dashboard atau gunakan SQL:
```sql
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
```

### Check Services Status
```bash
# Check if services are running
ps aux | grep "go run cmd/main.go"

# Test health endpoints
curl http://localhost:8001/health  # Auth Service
curl http://localhost:8002/health  # User Service
curl http://localhost:8003/health  # Board Service
curl http://localhost:8004/health  # Task Service
```

---
Made with ❤️ by Nathan Ginting
