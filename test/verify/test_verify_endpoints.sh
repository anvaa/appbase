#!/bin/bash

# Test script for POST /verify and GET /verify endpoints
# This script tests the verification endpoints with various scenarios

set -e

# Configuration
BASE_URL="https://localhost:5443"
CONTENT_TYPE="Content-Type: application/json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Helper functions
print_test_header() {
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}===========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((PASSED_TESTS++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((FAILED_TESTS++))
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

run_test() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local headers="$4"
    local expected_status="$5"
    local description="$6"
    
    ((TOTAL_TESTS++))
    
    echo ""
    echo -e "${BLUE}Test: $test_name${NC}"
    echo -e "Description: $description"
    echo -e "Request: $method $endpoint"
    
    if [ -n "$headers" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" $headers)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    fi
    
    # Split response and status code
    response_body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    echo "Status: $status_code"
    echo "Response: $response_body"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        print_success "$test_name - Status code matches expected ($expected_status)"
    else
        print_error "$test_name - Expected status $expected_status, got $status_code"
    fi
}

# Function to get a valid authentication token
get_auth_token() {
    print_info "Getting authentication token..."
    
    login_response=$(curl -s -X POST "$BASE_URL/login" \
        -H "$CONTENT_TYPE" \
        -d '{"username": "admin@app.loc", "password": "appadmin"}' \
        -c /tmp/cookies.txt)
    
    if echo "$login_response" | grep -q "token"; then
        print_success "Login successful"
        echo "$login_response" | jq -r '.token' 2>/dev/null || echo "TOKEN_IN_COOKIE"
    else
        print_error "Login failed"
        echo "$login_response"
        return 1
    fi
}

# Function to extract cookie from login
extract_cookie() {
    if [ -f /tmp/cookies.txt ]; then
        cookie_value=$(grep "AppBase_Auth" /tmp/cookies.txt | awk '{print $7}' 2>/dev/null || echo "")
        if [ -n "$cookie_value" ]; then
            echo "AppBase_Auth=$cookie_value"
        else
            echo ""
        fi
    else
        echo ""
    fi
}

print_test_header "VERIFY ENDPOINT TESTS"

# Start server using make run
print_info "Starting server with 'make run'..."
if ! pgrep -f "go run cmd/main.go" > /dev/null; then
    print_info "Starting server in background..."
    cd /home/anv/appbase
    make run > /tmp/server.log 2>&1 &
    SERVER_PID=$!
    print_info "Server PID: $SERVER_PID"
    sleep 3  # Wait for server to start
fi

# Check if server is running
print_info "Checking if server is running on $BASE_URL..."
for i in {1..10}; do
    if curl -s --connect-timeout 5 "$BASE_URL/health" > /dev/null 2>&1; then
        print_success "Server is running"
        break
    fi
    print_info "Waiting for server... ($i/10)"
    sleep 2
done

if ! curl -s --connect-timeout 5 "$BASE_URL/health" > /dev/null 2>&1; then
    print_error "Server failed to start. Check /tmp/server.log"
    exit 1
fi

# Get authentication token
AUTH_TOKEN=$(get_auth_token)
if [ "$AUTH_TOKEN" = "TOKEN_IN_COOKIE" ] || [ -n "$AUTH_TOKEN" ]; then
    COOKIE=$(extract_cookie)
    print_info "Using cookie authentication: $COOKIE"
else
    print_error "Failed to get authentication token"
    exit 1
fi

print_test_header "POST /verify ENDPOINT TESTS"

# Test 1: POST /verify without authentication
run_test "POST_VERIFY_NO_AUTH" \
    "POST" \
    "/verify" \
    "" \
    "401" \
    "POST /verify without authentication should return 401"

# Test 2: POST /verify with valid cookie
if [ -n "$COOKIE" ]; then
    run_test "POST_VERIFY_VALID_AUTH" \
        "POST" \
        "/verify" \
        "-H \"Cookie: $COOKIE\"" \
        "200" \
        "POST /verify with valid authentication should return 200"
else
    print_error "No cookie available for authenticated test"
    ((TOTAL_TESTS++))
    ((FAILED_TESTS++))
fi

# Test 3: POST /verify with invalid cookie
run_test "POST_VERIFY_INVALID_AUTH" \
    "POST" \
    "/verify" \
    "-H \"Cookie: AppBase_Auth=invalid_token_here\"" \
    "401" \
    "POST /verify with invalid token should return 401"

# Test 4: POST /verify with malformed cookie
run_test "POST_VERIFY_MALFORMED_AUTH" \
    "POST" \
    "/verify" \
    "-H \"Cookie: AppBase_Auth=not.a.valid.jwt\"" \
    "401" \
    "POST /verify with malformed token should return 401"

print_test_header "GET /verify ENDPOINT TESTS"

# Test 5: GET /verify without authentication
run_test "GET_VERIFY_NO_AUTH" \
    "GET" \
    "/verify" \
    "" \
    "401" \
    "GET /verify without authentication should return 401"

# Test 6: GET /verify with valid cookie
if [ -n "$COOKIE" ]; then
    run_test "GET_VERIFY_VALID_AUTH" \
        "GET" \
        "/verify" \
        "-H \"Cookie: $COOKIE\"" \
        "200" \
        "GET /verify with valid authentication should return 200"
else
    print_error "No cookie available for authenticated test"
    ((TOTAL_TESTS++))
    ((FAILED_TESTS++))
fi

# Test 7: GET /verify with invalid cookie
run_test "GET_VERIFY_INVALID_AUTH" \
    "GET" \
    "/verify" \
    "-H \"Cookie: AppBase_Auth=invalid_token_here\"" \
    "401" \
    "GET /verify with invalid token should return 401"

# Test 8: GET /verify with empty cookie
run_test "GET_VERIFY_EMPTY_AUTH" \
    "GET" \
    "/verify" \
    "-H \"Cookie: AppBase_Auth=\"" \
    "401" \
    "GET /verify with empty token should return 401"

print_test_header "/user/verify ENDPOINT TESTS"

# Test 9: POST /user/verify with valid authentication
if [ -n "$COOKIE" ]; then
    run_test "POST_USER_VERIFY_VALID_AUTH" \
        "POST" \
        "/user/verify" \
        "-H \"Cookie: $COOKIE\"" \
        "200" \
        "POST /user/verify with valid authentication should return 200"
else
    print_error "No cookie available for user verify test"
    ((TOTAL_TESTS++))
    ((FAILED_TESTS++))
fi

# Test 10: GET /user/verify with valid authentication
if [ -n "$COOKIE" ]; then
    run_test "GET_USER_VERIFY_VALID_AUTH" \
        "GET" \
        "/user/verify" \
        "-H \"Cookie: $COOKIE\"" \
        "200" \
        "GET /user/verify with valid authentication should return 200"
else
    print_error "No cookie available for user verify test"
    ((TOTAL_TESTS++))
    ((FAILED_TESTS++))
fi

print_test_header "RESPONSE CONTENT VALIDATION"

# Test 11: Validate response structure for successful verification
if [ -n "$COOKIE" ]; then
    print_info "Testing response content structure..."
    
    response=$(curl -s -X POST "$BASE_URL/verify" -H "Cookie: $COOKIE")
    
    # Check if response contains expected fields
    if echo "$response" | jq -e '.username' > /dev/null 2>&1; then
        print_success "Response contains username field"
        ((PASSED_TESTS++))
    else
        print_error "Response missing username field"
        ((FAILED_TESTS++))
    fi
    ((TOTAL_TESTS++))
    
    if echo "$response" | jq -e '.token' > /dev/null 2>&1; then
        print_success "Response contains token field"
        ((PASSED_TESTS++))
    else
        print_error "Response missing token field"
        ((FAILED_TESTS++))
    fi
    ((TOTAL_TESTS++))
    
    if echo "$response" | jq -e '.profile' > /dev/null 2>&1; then
        print_success "Response contains profile field"
        ((PASSED_TESTS++))
    else
        print_error "Response missing profile field"
        ((FAILED_TESTS++))
    fi
    ((TOTAL_TESTS++))
    
    print_info "Sample successful response:"
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
else
    print_error "Cannot test response content without valid authentication"
    ((TOTAL_TESTS += 3))
    ((FAILED_TESTS += 3))
fi

# Cleanup
rm -f /tmp/cookies.txt

# Stop server if we started it
if [ ! -z "$SERVER_PID" ]; then
    print_info "Stopping server (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
fi

print_test_header "TEST SUMMARY"

echo -e "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi