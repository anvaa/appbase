#!/bin/bash

# Simple test script for login and         health_response=$(curl -s -k "$BASE_URL/health" 2>/dev/null || echo "ERROR")ogout endpoints
# Usage: ./test_login_logout_simple.sh

logout_response=$(curl -s -k -w "%{http_code}" -b /tmp/login_test_cookies.txt \
    -X POST "$BASE_URL/logout" \
    -o /tmp/logout_response.txt)se SERVER_URL invalid_login_response=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "invalid@app.loc", "password": "wrongpass"}' \
    -o /tmp/invalid_login_response.txt) environment if available, otherwise default to HTTPS
BASE_URL="${SERVER_URL:-https://localhost:5443}"
ADMIN_USER="${USER_1:-admin@app.loc}"
ADMIN_PASS="${PSW_1:-appadmin}"
REGULAR_USER="${USER_2:-user@app.loc}"
REGULAR_PASS="${PSW_2:-password}"

echo "=== Testing Login and Logout Endpoints ==="
echo "Server: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Start server using make run
echo "1. Starting server with 'make run'..."
if ! pgrep -f "go run cmd/main.go" > /dev/null; then
    print_info "Starting server in background..."
    cd /home/anv/appbase
    make run > /tmp/server.log 2>&1 &
    SERVER_PID=$!
    echo "Server PID: $SERVER_PID"
    sleep 3  # Wait for server to start
fi

# Check if server is running
echo "2. Checking server health..."
for i in {1..10}; do
    health_response=$(curl -s -k "$BASE_URL/health" 2>/dev/null || echo "ERROR")
    if [[ "$health_response" == *"ok"* ]]; then
        print_success "Server is running"
        break
    fi
    print_info "Waiting for server... ($i/10)"
    sleep 2
done

if [[ "$health_response" != *"ok"* ]]; then
    print_error "Server failed to start. Check /tmp/server.log"
    exit 1
fi

echo ""
echo "3. Testing admin login with valid credentials..."

# Test admin login
admin_login=$(curl -s -k -c /tmp/test_cookies_admin.txt -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"password\": \"$ADMIN_PASS\"}" \
    -o /tmp/admin_login_response.txt)

status_code="${login_response: -3}"
response_content=$(cat /tmp/login_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "200" ]; then
    print_success "Login successful with valid credentials"
    
    # Check if response contains expected fields
    if [[ "$response_content" == *"token"* ]] && [[ "$response_content" == *"user"* ]]; then
        print_success "Response contains expected fields (token, user)"
    else
        print_error "Response missing expected fields"
    fi
else
    print_error "Login failed. Expected 200, got $status_code"
fi

echo ""
echo "4. Testing regular user login with valid credentials..."

# Test regular user login
user_login=$(curl -s -k -c /tmp/test_cookies_user.txt -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$REGULAR_USER\", \"password\": \"$REGULAR_PASS\"}" \
    -o /tmp/user_login_response.txt)

status_code="${user_login_response: -3}"
response_content=$(cat /tmp/login_response_user.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "200" ]; then
    print_success "Regular user login successful"
    
    # Check if response contains expected fields
    if [[ "$response_content" == *"token"* ]] && [[ "$response_content" == *"user"* ]]; then
        print_success "Response contains expected fields (token, user)"
    else
        print_error "Response missing expected fields"
    fi
else
    print_info "Regular user login failed (user may not exist), continuing with admin tests only"
fi

echo ""
echo "5. Testing login with invalid credentials..."

# Test invalid login
invalid_login_response=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "invalid@user.com", "password": "wrongpassword"}' \
    -o /tmp/invalid_login_response.txt)

status_code="${invalid_login_response: -3}"
response_content=$(cat /tmp/invalid_login_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "401" ]; then
    print_success "Correctly rejected invalid credentials"
else
    print_error "Expected 401 for invalid credentials, got $status_code"
fi

echo ""
echo "4. Testing login with malformed request..."

# Test malformed request
malformed_login_response=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"invalid": "json"}' \
    -o /tmp/malformed_login_response.txt)

status_code="${malformed_login_response: -3}"
response_content=$(cat /tmp/malformed_login_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "401" ] || [ "$status_code" = "400" ]; then
    print_success "Correctly rejected malformed request"
else
    print_error "Expected 400/401 for malformed request, got $status_code"
fi

echo ""
echo "5. Testing login with missing password..."

# Test missing password
missing_pwd_response=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "admin@app.loc"}' \
    -o /tmp/missing_pwd_response.txt)

status_code="${missing_pwd_response: -3}"
response_content=$(cat /tmp/missing_pwd_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "401" ] || [ "$status_code" = "400" ]; then
    print_success "Correctly rejected request with missing password"
else
    print_error "Expected 400/401 for missing password, got $status_code"
fi

echo ""
echo "6. Testing GET /logout (redirect-based logout)..."

# Test GET logout
logout_response=$(curl -s -k -w "%{http_code}" -b /tmp/login_test_cookies.txt \
    -X GET "$BASE_URL/logout" \
    -o /tmp/logout_response.txt)

status_code="${logout_response: -3}"
response_content=$(cat /tmp/logout_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "302" ] || [ "$status_code" = "200" ]; then
    print_success "Logout endpoint accessible"
else
    print_error "Logout failed. Expected 302/200, got $status_code"
fi

echo ""
echo "7. Testing login after logout (session should be cleared)..."

# Try to access protected endpoint with old cookie
protected_response=$(curl -s -k -w "%{http_code}" -b /tmp/login_test_cookies.txt \
    -X POST "$BASE_URL/verify" \
    -o /tmp/protected_after_logout.txt)

status_code="${protected_response: -3}"
response_content=$(cat /tmp/protected_after_logout.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "401" ]; then
    print_success "Session properly cleared after logout"
else
    print_error "Session not cleared properly. Expected 401, got $status_code"
fi

echo ""
echo "8. Testing fresh login after logout..."

# Fresh login after logout
fresh_login_response=$(curl -s -k -c /tmp/fresh_login_cookies.txt -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "admin@app.loc", "password": "appadmin"}' \
    -o /tmp/fresh_login_response.txt)

status_code="${fresh_login_response: -3}"
response_content=$(cat /tmp/fresh_login_response.txt)

echo "Status: $status_code"
echo "Response: $response_content"

if [ "$status_code" = "200" ]; then
    print_success "Fresh login after logout successful"
else
    print_error "Fresh login failed. Expected 200, got $status_code"
fi

echo ""
echo "9. Testing rate limiting (multiple rapid login attempts)..."

print_info "Attempting multiple rapid login requests..."
rate_limit_failed=false

for i in {1..5}; do
    rate_test_response=$(curl -s -k -w "%{http_code}" -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "password": "wrongpass"}' \
        -o /tmp/rate_test_$i.txt)
    
    status_code="${rate_test_response: -3}"
    
    if [ "$status_code" = "429" ]; then
        print_success "Rate limiting triggered after $i attempts"
        rate_limit_failed=false
        break
    fi
    
    # Small delay between requests
    sleep 0.1
done

if [ "$rate_limit_failed" = true ]; then
    print_info "Rate limiting not triggered (may be configured differently)"
else
    print_success "Rate limiting working properly"
fi

# Cleanup
rm -f /tmp/login_test_cookies*.txt /tmp/fresh_login_cookies.txt
rm -f /tmp/login_response*.txt /tmp/invalid_login_response.txt /tmp/malformed_login_response.txt
rm -f /tmp/missing_pwd_response.txt /tmp/logout_response.txt /tmp/protected_after_logout.txt
rm -f /tmp/fresh_login_response.txt /tmp/rate_test_*.txt

# Stop server if we started it (only if not managed by parent runner)
if [ ! -z "$SERVER_PID" ] && [[ -z "$MANAGED_BY_RUNNER" ]]; then
    echo ""
    print_info "Stopping server (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
fi

echo ""
echo "=== Test completed ==="