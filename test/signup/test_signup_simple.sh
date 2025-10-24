#!/bin/bash

# Simple test script for POST /signup and GET /signup endpoints
# Usage: ./test_signup_simple.sh

# Use SERVER_URL from environment if available, otherwise default to HTTPS
BASE_URL="${SERVER_URL:-https://localhost:5443}"
ADMIN_USER="${USER_1:-admin@app.loc}"
ADMIN_PASS="${PSW_1:-appadmin}"

echo "=== Testing /signup endpoints ==="
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

# Generate unique test data
TIMESTAMP=$(date +%s)
TEST_USERNAME="testuser_${TIMESTAMP}@test.loc"
TEST_PASSWORD="TestPassword123!"
TEST_ORGNAME="Test_Org_${TIMESTAMP}"

echo ""
echo "2. Testing GET /signup (signup form page)..."
signup_form_response=$(curl -s -k "$BASE_URL/signup/0" -w "%{http_code}")
http_code="${signup_form_response: -3}"
signup_content="${signup_form_response%???}"

if [[ "$http_code" == "200" ]]; then
    if [[ "$signup_content" == *"signup"* ]]; then
        echo "✓ GET /signup returns signup form (200 OK)"
    else
        echo "⚠ GET /signup returned 200 but content might be incorrect"
    fi
else
    echo "❌ GET /signup failed with HTTP code: $http_code"
    echo "Response: $signup_content"
fi

echo ""
echo "3. Testing GET /signup/:count with admin count (121209)..."
signup_admin_response=$(curl -s -k "$BASE_URL/signup/121209" -w "%{http_code}")
admin_http_code="${signup_admin_response: -3}"
admin_content="${signup_admin_response%???}"

if [[ "$admin_http_code" == "200" ]]; then
    if [[ "$admin_content" == *"signup"* ]]; then
        echo "✓ GET /signup/121209 returns admin signup form (200 OK)"
    else
        echo "⚠ GET /signup/121209 returned 200 but content might be incorrect"
    fi
else
    echo "❌ GET /signup/121209 failed with HTTP code: $admin_http_code"
fi

echo ""
echo "4. Testing POST /signup with invalid data (mismatched passwords)..."
signup_invalid_response=$(curl -s -k -X POST "$BASE_URL/signup" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"password\": \"$TEST_PASSWORD\",
        \"password2\": \"DifferentPassword\",
        \"orgname\": \"$TEST_ORGNAME\",
        \"count\": 0
    }" \
    -w "%{http_code}")

invalid_http_code="${signup_invalid_response: -3}"
invalid_content="${signup_invalid_response%???}"

echo "Status: $invalid_http_code"
echo "Response: $invalid_content"

if [[ "$invalid_http_code" == "500" ]] && [[ "$invalid_content" == *"Passwords do not match"* ]]; then
    echo "✓ POST /signup correctly rejected mismatched passwords"
elif [[ "$invalid_content" == *"Passwords do not match"* ]]; then
    echo "✓ POST /signup correctly rejected mismatched passwords (status: $invalid_http_code)"
else
    echo "❌ POST /signup should reject mismatched passwords"
fi

echo ""
echo "5. Testing POST /signup with invalid username (missing @)..."
signup_invalid_user_response=$(curl -s -k -X POST "$BASE_URL/signup" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"invaliduser\",
        \"password\": \"$TEST_PASSWORD\",
        \"password2\": \"$TEST_PASSWORD\",
        \"orgname\": \"$TEST_ORGNAME\",
        \"count\": 0
    }" \
    -w "%{http_code}")

invalid_user_http_code="${signup_invalid_user_response: -3}"
invalid_user_content="${signup_invalid_user_response%???}"

echo "Status: $invalid_user_http_code"
echo "Response: $invalid_user_content"

if [[ "$invalid_user_http_code" == "500" ]] && [[ "$invalid_user_content" == *"message"* ]]; then
    echo "✓ POST /signup correctly rejected invalid username format"
else
    echo "⚠ POST /signup validation response for invalid username: $invalid_user_content"
fi

echo ""
echo "6. Testing POST /signup with valid data (new user registration)..."
signup_valid_response=$(curl -s -k -X POST "$BASE_URL/signup" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"password\": \"$TEST_PASSWORD\",
        \"password2\": \"$TEST_PASSWORD\",
        \"orgname\": \"$TEST_ORGNAME\",
        \"count\": 0
    }" \
    -w "%{http_code}")

valid_http_code="${signup_valid_response: -3}"
valid_content="${signup_valid_response%???}"

echo "Status: $valid_http_code"
echo "Response: $valid_content"

if [[ "$valid_http_code" == "200" ]] && [[ "$valid_content" == *"success"* ]] && [[ "$valid_content" == *"/login"* ]]; then
    echo "✓ POST /signup successfully created new user"
    USER_CREATED=true
else
    echo "❌ POST /signup failed to create user"
    echo "Expected: 200 status with success message and /login redirect"
    USER_CREATED=false
fi

echo ""
echo "7. Testing POST /signup with duplicate username..."
signup_duplicate_response=$(curl -s -k -X POST "$BASE_URL/signup" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"password\": \"$TEST_PASSWORD\",
        \"password2\": \"$TEST_PASSWORD\",
        \"orgname\": \"Another_Org\",
        \"count\": 0
    }" \
    -w "%{http_code}")

duplicate_http_code="${signup_duplicate_response: -3}"
duplicate_content="${signup_duplicate_response%???}"

echo "Status: $duplicate_http_code"
echo "Response: $duplicate_content"

if [[ "$duplicate_http_code" == "500" ]] && [[ "$duplicate_content" == *"InternalServerError"* ]]; then
    echo "✓ POST /signup correctly rejected duplicate username"
else
    echo "⚠ POST /signup duplicate handling: $duplicate_content"
fi

echo ""
echo "8. Testing POST /signup with admin count (121209)..."
admin_test_username="admintest_${TIMESTAMP}@test.loc"
signup_admin_post_response=$(curl -s -k -X POST "$BASE_URL/signup" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$admin_test_username\",
        \"password\": \"$TEST_PASSWORD\",
        \"password2\": \"$TEST_PASSWORD\",
        \"orgname\": \"Admin_Test_Org\",
        \"count\": 121209
    }" \
    -w "%{http_code}")

admin_post_http_code="${signup_admin_post_response: -3}"
admin_post_content="${signup_admin_post_response%???}"

echo "Status: $admin_post_http_code"
echo "Response: $admin_post_content"

if [[ "$admin_post_http_code" == "200" ]] && [[ "$admin_post_content" == *"success"* ]] && [[ "$admin_post_content" == *"/v/users"* ]]; then
    echo "✓ POST /signup with admin count redirects to /v/users"
else
    echo "⚠ POST /signup admin count handling: $admin_post_content"
fi

if [[ "$USER_CREATED" == true ]]; then
    echo ""
    echo "9. Verifying new user can login..."
    login_test_response=$(curl -s -k -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USERNAME\",
            \"password\": \"$TEST_PASSWORD\"
        }" \
        -w "%{http_code}")

    login_test_http_code="${login_test_response: -3}"
    login_test_content="${login_test_response%???}"

    echo "Login Status: $login_test_http_code"
    echo "Login Response: $login_test_content"

    if [[ "$login_test_http_code" == "401" ]]; then
        echo "ℹ New user login failed - this is expected as new users need admin approval (IsAuth=false)"
        echo "✓ Signup created user but requires admin authorization before login"
    elif [[ "$login_test_http_code" == "200" ]]; then
        echo "✓ New user can login immediately after signup"
    else
        echo "⚠ Unexpected login response for new user"
    fi
fi

echo ""
echo "=== Test completed ==="

# Summary
echo "📋 Signup Test Summary:"
echo "  - GET /signup form: $([ "$http_code" == "200" ] && echo "✓ PASS" || echo "❌ FAIL")"
echo "  - GET /signup admin form: $([ "$admin_http_code" == "200" ] && echo "✓ PASS" || echo "❌ FAIL")"
echo "  - POST invalid passwords: $([ "$invalid_content" == *"Passwords do not match"* ] && echo "✓ PASS" || echo "❌ FAIL")"
echo "  - POST invalid username: $([ "$invalid_user_http_code" == "500" ] && echo "✓ PASS" || echo "⚠ PARTIAL")"
echo "  - POST valid signup: $([ "$USER_CREATED" == true ] && echo "✓ PASS" || echo "❌ FAIL")"
echo "  - POST duplicate user: $([ "$duplicate_content" == *"InternalServerError"* ] && echo "✓ PASS" || echo "⚠ PARTIAL")"
echo "  - POST admin count redirect: $([ "$admin_post_content" == *"/v/users"* ] && echo "✓ PASS" || echo "⚠ PARTIAL")"

# Exit with appropriate code
if [[ "$http_code" == "200" ]] && [[ "$USER_CREATED" == true ]]; then
    echo "✓ Signup tests passed"
    exit 0
else
    echo "✗ Some signup tests failed"
    exit 1
fi