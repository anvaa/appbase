#!/bin/bash

# Test script to verify app endpoint access for different user levels

# Source test environment configuration
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$TEST_DIR/../test.env"

# Build server URL from test.env
BASE_URL="http://$SERVER_URL:$SERVER_PORT"

echo "=== Testing /app endpoint access for different user levels ==="
echo "Server: $BASE_URL"
echo "Test Directory: $TEST_DIR"
echo "Using configuration from: $TEST_DIR/../test.env"

# Start server in background
echo "Starting test server..."
cd "$BASE_DIR"
go run cmd/main.go -tls false -port $SERVER_PORT -debug false &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Wait for server to start
sleep $SERVER_STARTUP_INTERVAL
echo

# Function to test user access
test_user_access() {
    local username="$1"
    local password="$2" 
    local level="$3"
    local expect_tools_access="$4"
    
    echo "=========================================="
    echo "Testing User: $username (Level $level)"
    echo "=========================================="
    
    # Test 1: Login
    echo "Test 1: Login as $username (level $level)"
    RESPONSE=$(curl -s -c "cookies_${level}.txt" -w "\nHTTP_CODE:%{http_code}" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$username\",\"password\":\"$password\"}" \
        "$BASE_URL/login")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
    RESPONSE_BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

    echo "HTTP Code: $HTTP_CODE"
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Login successful"
    else
        echo "❌ Login failed with code $HTTP_CODE"
        echo "Response: $RESPONSE_BODY"
        return 1
    fi
    echo

    # Test 2: Access /app endpoint 
    echo "Test 2: Access /app endpoint"
    RESPONSE=$(curl -s -b "cookies_${level}.txt" -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/app/")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
    RESPONSE_BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

    echo "HTTP Code: $HTTP_CODE"
    
    if echo "$RESPONSE_BODY" | grep -q "start.html\|<title>" && [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Successfully accessed /app endpoint - HTML page returned"
    else
        echo "❌ Failed to access /app endpoint"
        echo "Response preview: $(echo "$RESPONSE_BODY" | head -c 200)..."
        return 1
    fi
    echo

    # Test 3: Access /v/myaccount (should work for both levels)
    echo "Test 3: Access /v/myaccount (personal account - should work for both levels)"
    RESPONSE=$(curl -s -b "cookies_${level}.txt" -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/v/myaccount")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
    RESPONSE_BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

    echo "HTTP Code: $HTTP_CODE"
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Successfully accessed personal account page"
    else
        echo "❌ Failed to access personal account page"
        echo "Response preview: $(echo "$RESPONSE_BODY" | head -c 200)..."
        return 1
    fi
    echo

    # Test 4: Access restricted tools endpoint
    echo "Test 4: Access /tools/titles (level 30+ required)"
    RESPONSE=$(curl -s -b "cookies_${level}.txt" -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/tools/titles")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
    RESPONSE_BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

    echo "HTTP Code: $HTTP_CODE"

    if [ "$expect_tools_access" = "true" ]; then
        if [ "$HTTP_CODE" = "200" ]; then
            echo "✅ Correctly allowed access to tools endpoint (admin privilege)"
        else
            echo "❌ Unexpectedly blocked admin user from tools endpoint"
            return 1
        fi
    else
        if [ "$HTTP_CODE" = "403" ] || [ "$HTTP_CODE" = "401" ]; then
            echo "✅ Correctly blocked access to restricted tools endpoint"
        elif [ "$HTTP_CODE" = "200" ]; then
            echo "❌ Unexpectedly allowed low-level user access to restricted endpoint"
            return 1
        else
            echo "⚠️  Unexpected response code $HTTP_CODE for restricted endpoint"
        fi
    fi
    echo

    # Test 5: Access user management endpoint (level 30+ required)
    echo "Test 5: Access /v/users (user management - level 30+ required)"
    RESPONSE=$(curl -s -b "cookies_${level}.txt" -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/v/users")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)

    echo "HTTP Code: $HTTP_CODE"

    if [ "$expect_tools_access" = "true" ]; then
        if [ "$HTTP_CODE" = "200" ]; then
            echo "✅ Correctly allowed access to user management (admin privilege)"
        else
            echo "❌ Unexpectedly blocked admin user from user management"
            return 1
        fi
    else
        if [ "$HTTP_CODE" = "403" ] || [ "$HTTP_CODE" = "401" ]; then
            echo "✅ Correctly blocked access to user management endpoint"
        elif [ "$HTTP_CODE" = "200" ]; then
            echo "❌ Unexpectedly allowed low-level user access to user management"
            return 1
        else
            echo "⚠️  Unexpected response code $HTTP_CODE for user management endpoint"
        fi
    fi
    echo

    # Cleanup
    rm -f "cookies_${level}.txt"

    echo "✅ All tests passed for user level $level"
    echo
    
    return 0
}

# Test Level 10 User (Standard User) - using test.env credentials
if ! test_user_access "$USER_2" "$PSW_2" "10" "false"; then
    echo "❌ Tests failed for user level 10 ($USER_2)"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Test Level 40 User (Admin) - using test.env credentials  
if ! test_user_access "$USER_1" "$PSW_1" "40" "true"; then
    echo "❌ Tests failed for user level 40 ($USER_1)"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Cleanup server
echo "Stopping test server..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo "=========================================="
echo "🎉 ALL TESTS COMPLETED SUCCESSFULLY! 🎉"
echo "=========================================="
echo "✅ Level 10 users can access /app and personal features"
echo "✅ Level 10 users are correctly blocked from admin features"  
echo "✅ Level 40 users can access all endpoints including admin features"
echo "✅ Authentication and authorization working as expected"