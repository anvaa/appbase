#!/bin/bash

echo "=== Quick Test of Both User Levels ==="

# Source test environment configuration
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$TEST_DIR/../test.env"

# Build server URL from test.env
BASE_URL="http://$SERVER_URL:$SERVER_PORT"

echo "Using configuration from: $TEST_DIR/../test.env"
echo "Server: $BASE_URL"

cd "$BASE_DIR"

# Start server
echo "Starting server on $SERVER_URL:$SERVER_PORT..."
go run cmd/main.go -tls false -port $SERVER_PORT -debug false &
SERVER_PID=$!
sleep $SERVER_STARTUP_INTERVAL

# Test Level 10 User
echo ""
echo "1. Testing $USER_2 (Level 10):"
echo "   Logging in..."
LOGIN_10=$(curl -s -c /tmp/test_user10.txt -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"username\":\"$USER_2\",\"password\":\"$PSW_2\"}" \
  $BASE_URL/login -o /dev/null)
echo "   Login response: $LOGIN_10"

if [ "$LOGIN_10" = "200" ]; then
    echo "   Testing /app access..."
    APP_10=$(curl -s -b /tmp/test_user10.txt -w "%{http_code}" $BASE_URL/app/ -o /dev/null)
    echo "   /app response: $APP_10"
    
    echo "   Testing /tools/titles access (should be blocked)..."
    TOOLS_10=$(curl -s -b /tmp/test_user10.txt -w "%{http_code}" $BASE_URL/tools/titles -o /dev/null)
    echo "   /tools/titles response: $TOOLS_10"
fi

# Test Level 40 User
echo ""
echo "2. Testing $USER_1 (Level 40):"
echo "   Logging in..."
LOGIN_40=$(curl -s -c /tmp/test_admin40.txt -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"username\":\"$USER_1\",\"password\":\"$PSW_1\"}" \
  $BASE_URL/login -o /dev/null)
echo "   Login response: $LOGIN_40"

if [ "$LOGIN_40" = "200" ]; then
    echo "   Testing /app access..."
    APP_40=$(curl -s -b /tmp/test_admin40.txt -w "%{http_code}" $BASE_URL/app/ -o /dev/null)
    echo "   /app response: $APP_40"
    
    echo "   Testing /tools/titles access (should work)..."
    TOOLS_40=$(curl -s -b /tmp/test_admin40.txt -w "%{http_code}" $BASE_URL/tools/titles -o /dev/null)
    echo "   /tools/titles response: $TOOLS_40"
fi

echo ""
echo "=== Test Summary ==="
if [ "$LOGIN_10" = "200" ] && [ "$APP_10" = "200" ] && [ "$TOOLS_10" = "403" -o "$TOOLS_10" = "401" ]; then
    echo "✅ Level 10 user working correctly (can access /app, blocked from tools)"
else
    echo "❌ Level 10 user has issues (login:$LOGIN_10, app:$APP_10, tools:$TOOLS_10)"
fi

if [ "$LOGIN_40" = "200" ] && [ "$APP_40" = "200" ] && [ "$TOOLS_40" = "200" ]; then
    echo "✅ Level 40 user working correctly (can access both /app and tools)"
else
    echo "❌ Level 40 user has issues (login:$LOGIN_40, app:$APP_40, tools:$TOOLS_40)"
fi

# Cleanup
echo ""
echo "Cleaning up..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null
rm -f /tmp/test_user10.txt /tmp/test_admin40.txt

echo "Test completed!"