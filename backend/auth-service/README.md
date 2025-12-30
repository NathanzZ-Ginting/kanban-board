# Auth Service

Authentication and authorization microservice for Kanban Board.

## Features
- User registration and login
- JWT token generation and validation
- Password hashing (bcrypt)
- Token refresh
- Logout functionality

## API Endpoints

### POST /api/v1/auth/register
Register a new user
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}
```

### POST /api/v1/auth/login
Login user
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### POST /api/v1/auth/refresh
Refresh access token
```json
{
  "refreshToken": "token"
}
```

### POST /api/v1/auth/logout
Logout user (invalidate token)

## Tech Stack
- Go 1.21+
- MongoDB
- JWT
- Bcrypt

## Run Service
```bash
go run cmd/main.go
```
