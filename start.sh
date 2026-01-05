#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Kanban Monorepo Services...${NC}"
echo ""

# Function to start a service
start_service() {
    local service_name=$1
    local port=$2
    local path=$3
    
    echo -e "${YELLOW}Starting ${service_name} on port ${port}...${NC}"
    cd "${path}" && go run cmd/main.go &
    sleep 2
}

# Kill existing processes
echo -e "${YELLOW}🛑 Stopping existing services...${NC}"
pkill -f "go run cmd/main.go" 2>/dev/null
sleep 2

# Start all services
start_service "Auth Service" "8001" "backend/auth-service"
start_service "User Service" "8002" "backend/user-service"
start_service "Board Service" "8003" "backend/board-service"
start_service "Task Service" "8004" "backend/task-service"

echo ""
echo -e "${GREEN}✅ All services started!${NC}"
echo ""
echo -e "${BLUE}📝 Services running on:${NC}"
echo "  • Auth Service:  http://localhost:8001"
echo "  • User Service:  http://localhost:8002"
echo "  • Board Service: http://localhost:8003"
echo "  • Task Service:  http://localhost:8004"
echo ""
echo -e "${YELLOW}💡 To stop all services, run: pkill -f 'go run cmd/main.go'${NC}"
echo -e "${YELLOW}💡 To view logs, check terminal outputs${NC}"
echo ""
