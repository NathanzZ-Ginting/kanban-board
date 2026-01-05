#!/bin/bash

echo "🔄 Restarting all Kanban services..."

# Kill all running Go services
echo "🛑 Stopping running services..."
pkill -f "go run cmd/main.go"
sleep 2

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}✅ All services stopped${NC}"
echo ""
echo -e "${BLUE}🚀 To start services again, run:${NC}"
echo ""
echo "cd backend/auth-service && go run cmd/main.go"
echo "cd backend/user-service && go run cmd/main.go"
echo "cd backend/board-service && go run cmd/main.go"
echo "cd backend/task-service && go run cmd/main.go"
echo ""
echo -e "${BLUE}Or use start.sh to run all services${NC}"
