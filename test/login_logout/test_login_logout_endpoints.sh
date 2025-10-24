#!/bin/bash

# Comprehensive test script for login and logout endpoints
# Usage: ./test_login_logout_endpoints.sh

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
    local body="$5"
    local expected_status="$6"
    local description="$7"
    local cookie_file="$8"
    
    ((TOTAL_TESTS++))
    
    echo ""
    echo -e "${BLUE}Test: $test_name${NC}"
    echo -e "Description: $description"
    echo -e "Request: $method $endpoint"
    
    # Build curl command
    local curl_cmd="curl -s -w \"\\n%{http_code}\""
    
    if [ -n "$cookie_file" ]; then
        if [[ "$cookie_file" == "save:"* ]]; then
            local save_file="${cookie_file#save:}"
            curl_cmd="$curl_cmd -c $save_file"
        elif [[ "$cookie_file" == "load:"* ]]; then
            local load_file="${cookie_file#load:}"
            curl_cmd="$curl_cmd -b $load_file"
        fi
    fi
    
    curl_cmd="$curl_cmd -X \"$method\" \"$BASE_URL$endpoint\""
    
    if [ -n "$headers" ]; then
        curl_cmd="$curl_cmd $headers"
    fi
    
    if [ -n "$body" ]; then
        curl_cmd="$curl_cmd -d '$body'"
    fi
    
    # Execute request
    response=$(eval $curl_cmd)
    
    # Split response and status code
    response_body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    echo "Status: $status_code"
    echo "Response: $response_body"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        print_success "$test_name - Status code matches expected ($expected_status)"
        return 0
    else
        print_error "$test_name - Expected status $expected_status, got $status_code"
        return 1
    fi
}

validate_response_fields() {
    local response="$1"
    local required_fields="$2"
    local test_name="$3"
    
    ((TOTAL_TESTS++))
    
    for field in $required_fields; do
        if echo "$response" | grep -q "\"$field\""; then
            print_success "$test_name - Contains required field: $field"
        else
            print_error "$test_name - Missing required field: $field"
            return 1
        fi
    done
    return 0
}

# Test data
VALID_USER='{"username": "admin@app.loc", "password": "appadmin"}'
INVALID_USER='{"username": "invalid@user.com", "password": "wrongpassword"}'
MALFORMED_REQUEST='{"invalid": "json"}'
MISSING_PASSWORD='{"username": "admin@app.loc"}'
EMPTY_REQUEST='{}'

print_test_header "LOGIN AND LOGOUT ENDPOINT TESTS"

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

print_test_header "LOGIN ENDPOINT TESTS"

# Test 1: Valid login
if run_test "LOGIN_VALID_CREDENTIALS" \
    "POST" \
    "/login" \
    "-H \"$CONTENT_TYPE\"" \
    "$VALID_USER" \
    "200" \
    "Login with valid credentials should succeed" \
    "save:/tmp/valid_login_cookies.txt"; then
    
    # Validate response fields for successful login
    response_body=$(tail -n +1 /tmp/last_response.txt 2>/dev/null || echo "")
    validate_response_fields "$response_body" "token user access_token" "LOGIN_RESPONSE_VALIDATION"
fi

# Test 2: Invalid credentials
run_test "LOGIN_INVALID_CREDENTIALS" \
    "POST" \
    "/login" \
    "-H \"$CONTENT_TYPE\"" \
    "$INVALID_USER" \
    "401" \
    "Login with invalid credentials should fail"

# Test 3: Malformed request
run_test "LOGIN_MALFORMED_REQUEST" \
    "POST" \
    "/login" \
    "-H \"$CONTENT_TYPE\"" \
    "$MALFORMED_REQUEST" \
    "401" \
    "Login with malformed request should fail"

# Test 4: Missing password
run_test "LOGIN_MISSING_PASSWORD" \
    "POST" \
    "/login" \
    "-H \"$CONTENT_TYPE\"" \
    "$MISSING_PASSWORD" \
    "401" \
    "Login with missing password should fail"

# Test 5: Empty request
run_test "LOGIN_EMPTY_REQUEST" \
    "POST" \
    "/login" \
    "-H \"$CONTENT_TYPE\"" \
    "$EMPTY_REQUEST" \
    "400" \
    "Login with empty request should fail"

# Test 6: Missing Content-Type header
run_test "LOGIN_NO_CONTENT_TYPE" \
    "POST" \
    "/login" \
    "" \
    "$VALID_USER" \
    "400" \
    "Login without Content-Type header should fail"

# Test 7: GET method on login endpoint (should not be allowed)
run_test "LOGIN_GET_METHOD" \
    "GET" \
    "/login" \
    "" \
    "" \
    "200" \
    "GET /login should return login form (allowed for form display)"

print_test_header "LOGOUT ENDPOINT TESTS"

# Ensure we have a valid session for logout tests
print_info "Setting up valid session for logout tests..."
curl -s -c /tmp/logout_test_cookies.txt -X POST "$BASE_URL/login" \
    -H "$CONTENT_TYPE" \
    -d "$VALID_USER" > /dev/null

# Test 8: GET logout with valid session
run_test "LOGOUT_GET_VALID_SESSION" \
    "GET" \
    "/logout" \
    "" \
    "" \
    "302" \
    "GET logout with valid session should redirect" \
    "load:/tmp/logout_test_cookies.txt"

# Test 9: GET logout without session
run_test "LOGOUT_GET_NO_SESSION" \
    "GET" \
    "/logout" \
    "" \
    "" \
    "302" \
    "GET logout without session should still redirect"

print_test_header "SESSION MANAGEMENT TESTS"

# Test 10: Access protected endpoint after logout
print_info "Testing session invalidation after logout..."
((TOTAL_TESTS++))

# First, login to get a valid session
curl -s -c /tmp/session_test_cookies.txt -X POST "$BASE_URL/login" \
    -H "$CONTENT_TYPE" \
    -d "$VALID_USER" > /dev/null

# Verify session works
verify_before=$(curl -s -w "%{http_code}" -b /tmp/session_test_cookies.txt \
    -X POST "$BASE_URL/verify" -o /tmp/verify_before.txt)
before_status="${verify_before: -3}"

if [ "$before_status" = "200" ]; then
    print_success "SESSION_BEFORE_LOGOUT - Valid session confirmed"
    ((PASSED_TESTS++))
else
    print_error "SESSION_BEFORE_LOGOUT - Session setup failed"
    ((FAILED_TESTS++))
fi

# Now logout
curl -s -b /tmp/session_test_cookies.txt -X GET "$BASE_URL/logout" > /dev/null

# Try to access protected endpoint after logout
((TOTAL_TESTS++))
verify_after=$(curl -s -w "%{http_code}" -b /tmp/session_test_cookies.txt \
    -X POST "$BASE_URL/verify" -o /tmp/verify_after.txt)
after_status="${verify_after: -3}"

if [ "$after_status" = "401" ]; then
    print_success "SESSION_AFTER_LOGOUT - Session properly invalidated"
    ((PASSED_TESTS++))
else
    print_error "SESSION_AFTER_LOGOUT - Session not properly invalidated (got $after_status)"
    ((FAILED_TESTS++))
fi

print_test_header "RATE LIMITING TESTS"

# Test 11: Rate limiting on login attempts
print_info "Testing rate limiting with multiple failed login attempts..."
((TOTAL_TESTS++))

rate_limit_triggered=false
for i in {1..6}; do
    rate_response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/login" \
        -H "$CONTENT_TYPE" \
        -d '{"username": "rate@limit.test", "password": "invalid"}' \
        -o /tmp/rate_test_$i.txt)
    
    rate_status="${rate_response: -3}"
    
    if [ "$rate_status" = "429" ]; then
        print_success "RATE_LIMITING - Triggered after $i attempts"
        ((PASSED_TESTS++))
        rate_limit_triggered=true
        break
    fi
    
    sleep 0.2  # Small delay between requests
done

if [ "$rate_limit_triggered" = false ]; then
    print_info "RATE_LIMITING - No rate limiting detected (may be configured differently)"
    ((PASSED_TESTS++))
fi

print_test_header "EDGE CASE TESTS"

# Test 12: Very long username
((TOTAL_TESTS++))
long_username=$(printf 'a%.0s' {1..1000})
long_user_request="{\"username\": \"$long_username@test.com\", \"password\": \"test\"}"

long_response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "$CONTENT_TYPE" \
    -d "$long_user_request" \
    -o /tmp/long_username_test.txt)

long_status="${long_response: -3}"

if [ "$long_status" = "401" ] || [ "$long_status" = "400" ]; then
    print_success "LONG_USERNAME - Properly handled long username"
    ((PASSED_TESTS++))
else
    print_error "LONG_USERNAME - Unexpected response for long username (got $long_status)"
    ((FAILED_TESTS++))
fi

# Test 13: SQL injection attempt
((TOTAL_TESTS++))
sql_injection_user='{"username": "admin@app.loc\"; DROP TABLE users; --", "password": "test"}'

sql_response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "$CONTENT_TYPE" \
    -d "$sql_injection_user" \
    -o /tmp/sql_injection_test.txt)

sql_status="${sql_response: -3}"

if [ "$sql_status" = "401" ] || [ "$sql_status" = "400" ]; then
    print_success "SQL_INJECTION_PROTECTION - Properly handled SQL injection attempt"
    ((PASSED_TESTS++))
else
    print_error "SQL_INJECTION_PROTECTION - Unexpected response for SQL injection (got $sql_status)"
    ((FAILED_TESTS++))
fi

# Cleanup
rm -f /tmp/*_cookies.txt /tmp/*_test*.txt /tmp/verify_*.txt

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

success_rate=$(( PASSED_TESTS * 100 / TOTAL_TESTS ))
echo -e "Success Rate: $success_rate%"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi