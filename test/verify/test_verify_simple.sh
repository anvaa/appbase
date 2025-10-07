#!/bin/bash

# Simple test script for POST /verify and GET /verify endpoints
# Usage: ./test_verify_simple.sh

# Use SERVER_URL from environment if available, otherwise default to HTTPS
BASE_URL="${SERVER_URL:-https://localhost:5443}"
ADMIN_USER="${USER_1:-admin@app.loc}"
ADMIN_PASS="${PSW_1:-appadmin}"
REGULAR_USER="${USER_2:-user@app.loc}"
REGULAR_PASS="${PSW_2:-password}"

echo "=== Testing /verify endpoints ==="
echo "Server: $BASE_URL"
echo "Environment SERVER_URL: ${SERVER_URL:-not set}"
echo "Environment MANAGED_BY_RUNNER: ${MANAGED_BY_RUNNER:-not set}"
echo ""

# Check if server is running, start if needed (only if not managed by parent runner)
echo "1. Checking server health..."
health_response=$(curl -s -k "$BASE_URL/health" 2>/dev/null || echo "ERROR")
if [[ "$health_response" == *"ok"* ]]; then
    echo "✓ Server is already running at $BASE_URL"
elif [[ -z "$MANAGED_BY_RUNNER" ]]; then
    echo "🚀 Starting server with make run..."
    cd ${BASE_DIR:-/home/anv/appbase}
    make run &
    SERVER_PID=$!
    
    echo "🎆 Server starting (PID: $SERVER_PID)"
    
    # Wait for server to be ready
    for i in {1..10}; do
        sleep 2
        health_response=$(curl -s -k "$BASE_URL/health" 2>/dev/null || echo "ERROR")
        if [[ "$health_response" == *"ok"* ]]; then
            echo "✓ Server ready"
            break
        fi
        echo "⏳ Waiting for server... ($i/10)"
    done
    
    if [[ "$health_response" != *"ok"* ]]; then
        echo "❌ Server failed to start"
        exit 1
    fi
else
    echo "❌ Server not running and managed by parent runner - this shouldn't happen"
    exit 1
fi

echo ""
echo "2. Logging in to get authentication token..."

# Login to get authentication cookie
login_response=$(curl -s -k -c /tmp/test_cookies.txt -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"password\": \"$ADMIN_PASS\"}")

if [[ "$login_response" == *"token"* ]]; then
    echo "✓ Login successful"
    # Extract token from response for display
    token=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "Token: ${token:0:50}..."
else
    echo "✗ Login failed"
    echo "Response: $login_response"
    exit 1
fi

echo ""
echo ""
echo "4. Testing POST /verify without authentication..."
verify_noauth=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/verify" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -o /tmp/verify_response.txt)
status_code="${verify_noauth: -3}"
response_content=$(cat /tmp/verify_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"
if [ "$status_code" = "401" ]; then
    echo "✓ Correctly rejected request without authentication"
else
    echo "✗ Expected 401, got $status_code"
fi

echo ""
echo ""
echo "5. Testing POST /verify with admin authentication..."
verify_auth=$(curl -s -k -b /tmp/test_cookies.txt -w "%{http_code}" -X POST "$BASE_URL/verify" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -o /tmp/verify_auth_response.txt)
status_code="${verify_auth: -3}"
response_content=$(cat /tmp/verify_auth_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"
if [ "$status_code" = "200" ]; then
    echo "✓ POST /verify works with authentication"
else
    echo "✗ Expected 200, got $status_code"
fi

echo ""
echo "6. Testing with regular user (user@app.loc)..."

# Login with regular user to test second user level
regular_login_response=$(curl -s -k -c /tmp/test_cookies_user.txt -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$REGULAR_USER\", \"password\": \"$REGULAR_PASS\"}")

if [[ "$user_login_response" == *"token"* ]]; then
    echo "✓ Regular user login successful"
else
    echo "⚠ Regular user login failed, testing with admin only"
    echo "Response: $user_login_response"
fi

echo ""
echo "7. Testing GET /verify with admin authentication..."
get_verify=$(curl -s -k -b /tmp/test_cookies.txt -w "%{http_code}" -X GET "$BASE_URL/verify" \
    -H "Accept: application/json" \
    -o /tmp/get_verify_response.txt)
status_code="${get_verify: -3}"
response_content=$(cat /tmp/get_verify_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"
if [ "$status_code" = "200" ]; then
    echo "✓ GET /verify works with authentication"
else
    echo "✗ GET /verify failed. Expected 200, got $status_code"
fi

echo ""
echo "6. Testing /user/verify endpoint..."
user_verify=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/user/verify" \
    -H "Content-Type: application/json" \
    -b /tmp/test_cookies.txt \
    -o /tmp/user_verify_response.txt)
status_code="${user_verify: -3}"
response_content=$(cat /tmp/user_verify_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"
if [ "$status_code" = "200" ]; then
    echo "✓ /user/verify works with authentication"
else
    echo "✗ /user/verify failed. Expected 200, got $status_code"
fi

echo ""
echo "8. Testing with invalid token..."
invalid_verify=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/verify" \
    -H "Cookie: AppBase_Auth=invalid.jwt.token" \
    -H "Accept: application/json" \
    -o /tmp/invalid_verify_response.txt)
status_code="${invalid_verify: -3}"
response_content=$(cat /tmp/invalid_verify_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"
if [ "$status_code" = "401" ]; then
    echo "✓ Correctly rejected invalid token"
else
    echo "✗ Expected 401 for invalid token, got $status_code"
fi

# Cleanup
rm -f /tmp/test_cookies_*.txt /tmp/verify_response.txt /tmp/verify_auth_response.txt 
rm -f /tmp/get_verify_response.txt /tmp/user_verify_response.txt /tmp/invalid_verify_response.txt

# Stop server if we started it (only if not managed by parent runner)
if [ ! -z "$SERVER_PID" ] && [[ -z "$MANAGED_BY_RUNNER" ]]; then
    echo ""
    echo "Stopping server (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
fi

echo ""
echo "=== Test completed ==="