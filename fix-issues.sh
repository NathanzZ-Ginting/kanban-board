#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Kanban Board - Quick Fix Tool${NC}"
echo ""

# Check if backend services are running
echo -e "${YELLOW}1. Checking backend services...${NC}"
SERVICES=$(ps aux | grep "go run cmd/main.go" | grep -v grep | wc -l)

if [ "$SERVICES" -lt 4 ]; then
  echo -e "${RED}❌ Not all services are running (found: $SERVICES/4)${NC}"
  echo -e "${BLUE}   Start services with: ./start.sh${NC}"
else
  echo -e "${GREEN}✅ All backend services are running${NC}"
fi

echo ""

# Test health endpoints
echo -e "${YELLOW}2. Testing service health...${NC}"
for port in 8001 8002 8003 8004; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/health)
  if [ "$STATUS" = "200" ]; then
    echo -e "${GREEN}✅ Port $port is healthy${NC}"
  else
    echo -e "${RED}❌ Port $port is not responding${NC}"
  fi
done

echo ""

# Check if user can login
echo -e "${YELLOW}3. Testing authentication...${NC}"
echo "Enter your email (or press Enter to skip):"
read -r EMAIL

if [ -n "$EMAIL" ]; then
  echo "Enter your password:"
  read -rs PASSWORD
  
  echo ""
  echo "Attempting login..."
  
  LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
  
  TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.accessToken // .data.accessToken // empty')
  
  if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✅ Login successful!${NC}"
    echo ""
    echo "Access Token (copy this to browser localStorage):"
    echo "$TOKEN"
    echo ""
    echo -e "${BLUE}To use in browser:${NC}"
    echo "1. Open DevTools (F12)"
    echo "2. Go to Console tab"
    echo "3. Run: localStorage.setItem('accessToken', '$TOKEN')"
    echo "4. Refresh the page"
  else
    echo -e "${RED}❌ Login failed${NC}"
    echo "Response: $LOGIN_RESPONSE"
    echo ""
    echo -e "${YELLOW}💡 Try creating a new account:${NC}"
    echo ""
    echo "curl -X POST http://localhost:8001/api/v1/auth/register \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -d '{"
    echo "    \"email\": \"your@email.com\","
    echo "    \"username\": \"yourusername\","
    echo "    \"password\": \"password123\","
    echo "    \"firstName\": \"Your\","
    echo "    \"lastName\": \"Name\""
    echo "  }'"
  fi
fi

echo ""
echo -e "${BLUE}📝 Common Issues & Solutions:${NC}"
echo ""
echo "Issue: 'Failed to load board data'"
echo "  → Solution: Login again or check browser console (F12)"
echo ""
echo "Issue: 'prepared statement already exists'"
echo "  → Solution: Restart services with ./restart-services.sh"
echo ""
echo "Issue: 'Invalid or expired token'"
echo "  → Solution: Clear localStorage and login again"
echo ""
echo -e "${GREEN}✅ Quick Fix Complete!${NC}"
