#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing Kanban API Endpoints${NC}"
echo ""

# Test 1: Health checks
echo -e "${BLUE}1. Testing Health Endpoints...${NC}"
echo "Auth Service:"
curl -s http://localhost:8001/health | jq .
echo ""
echo "User Service:"
curl -s http://localhost:8002/health | jq .
echo ""
echo "Board Service:"
curl -s http://localhost:8003/health | jq .
echo ""
echo "Task Service:"
curl -s http://localhost:8004/health | jq .
echo ""

# Test 2: Login
echo -e "${BLUE}2. Testing Login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joe@gmail.com",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.accessToken // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo -e "${RED}❌ Login failed! No access token received${NC}"
  echo ""
  echo -e "${BLUE}3. Try to register first...${NC}"
  REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8001/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "email": "joe@gmail.com",
      "username": "joe",
      "password": "password123",
      "firstName": "Joe",
      "lastName": "Ginting"
    }')
  
  echo "$REGISTER_RESPONSE" | jq .
  
  # Try login again
  LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "email": "joe@gmail.com",
      "password": "password123"
    }')
  
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.accessToken // empty')
fi

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo -e "${RED}❌ Still no access token. Check auth service logs.${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Got access token!${NC}"
echo ""

# Test 3: Get Boards
echo -e "${BLUE}4. Testing Get Boards (Protected)...${NC}"
BOARDS_RESPONSE=$(curl -s http://localhost:8003/api/v1/boards \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$BOARDS_RESPONSE" | jq .
echo ""

# Test 5: Get specific board
echo -e "${BLUE}5. Testing Get Board Detail (Protected)...${NC}"
BOARD_ID=$(echo "$BOARDS_RESPONSE" | jq -r '.boards[0].id // empty')

if [ -n "$BOARD_ID" ] && [ "$BOARD_ID" != "null" ]; then
  echo "Getting board ID: $BOARD_ID"
  curl -s http://localhost:8003/api/v1/boards/$BOARD_ID \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
else
  echo -e "${RED}No boards found. Create one first.${NC}"
fi

echo ""
echo -e "${GREEN}✅ API Test Complete!${NC}"
