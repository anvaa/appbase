#!/bin/bash

# Master test runner for all endpoint tests
# Usage: ./run_tests.sh {tls|notls} [TEST_SUITE] [options]
# 
# The test server runs with: go run cmd/main.go -tls {true|false} -port SERVER_PORT -debug true
# Server port and URL are read from test.env
# 
# Examples:
#   ./run_tests.sh tls           # Run all tests with HTTPS on configured port
#   ./run_tests.sh notls         # Run all tests with HTTP on configured port
#   ./run_tests.sh tls verify    # Run verify tests with HTTPS
#   ./run_tests.sh notls verify  # Run verify tests with HTTP

set -e

# Initialize variables
START_TIME=$(date '+%Y-%m-%d_%H-%M-%S')
SERVER_ARGS="tls"
SERVER_PROTO="https"
TEST_SUITE="all"
VERBOSE=false
QUICK=false

# Kill any existing processes on test server port (preserves dev server on 5443)
kill_existing_server() {
    echo "🔍 Checking for existing processes on test port $SERVER_PORT..."
    
    # Find processes using test server port
    local pids=$(ss -tulpn 2>/dev/null | grep :$SERVER_PORT | awk '{print $7}' | grep -o '[0-9]*' | head -1)
    
    if [ ! -z "$pids" ]; then
        echo "🔪 Found process using port $SERVER_PORT: $pids"
        kill -9 $pids 2>/dev/null || true
        sleep 1
    fi
    
    # Also kill any go run processes
    pkill -f "go run cmd/main.go" 2>/dev/null || true
    
    echo "✅ Port $SERVER_PORT cleared for fresh start"
}

# Generate test report
generate_test_report() {
    local test_name="$1"
    local status="$2"
    local details="$3"
    local end_time=$(date '+%Y-%m-%d_%H-%M-%S')
    
    local report_file="${TEST_DIR}/${START_TIME}_${test_name}_report.txt"
    
    cat > "$report_file" << EOF
==========================================
APPBASE TEST REPORT
==========================================
Test Suite: $test_name
Start Time: $(date -d "${START_TIME//_/ }" -d "${START_TIME//-/:}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo $START_TIME)
End Time: $(date -d "${end_time//_/ }" -d "${end_time//-/:}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo $end_time)
Server Mode: $SERVER_ARGS ($SERVER_PROTO)
Server URL: $FULL_SERVER_URL
Status: $status
Working Directory: $BASE_DIR
Admin User: $USER_1
Regular User: $USER_2

==========================================
TEST DETAILS
==========================================
$details

==========================================
ENVIRONMENT INFO
==========================================
BASE_DIR: $BASE_DIR
TEST_DIR: $TEST_DIR
SERVER_URL: $FULL_SERVER_URL
SERVER_PORT: $SERVER_PORT
SERVER_ARGS: $SERVER_ARGS
MANAGED_BY_RUNNER: $MANAGED_BY_RUNNER

Generated on: $(date)
EOF

    echo "📋 Test report saved to: $report_file"
}

# Load test configuration
load_config() {
    local config_file="$(dirname "$0")/test.env"
    if [ -f "$config_file" ]; then
        echo "📋 Loading configuration from $config_file"
        # Source the config file
        source "$config_file"
        
        # Set defaults if not provided in config
        BASE_DIR="${BASE_DIR:-/home/anv/appbase}"
        TEST_DIR="${TEST_DIR:-$BASE_DIR/test}"
        SERVER_URL="${SERVER_URL:-localhost}"
        SERVER_PORT="${SERVER_PORT:-8888}"
        USER_1="${USER_1:-admin@app.loc}"
        PSW_1="${PSW_1:-appadmin}"
        USER_2="${USER_2:-user@app.loc}"
        PSW_2="${PSW_2:-password}"
        
        echo "✓ Configuration loaded successfully"
        echo "  Base Directory: $BASE_DIR"
        echo "  Test Directory: $TEST_DIR" 
        echo "  Admin User: $USER_1"
        echo "  Regular User: $USER_2"
    else
        echo "⚠️ Configuration file not found: $config_file"
        echo "   Using default values..."
        
        # Set default values
        BASE_DIR="/home/anv/appbase"
        TEST_DIR="$BASE_DIR/test"
        SERVER_URL="localhost"
        SERVER_PORT="8888"
        USER_1="admin@app.loc"
        PSW_1="appadmin"
        USER_2="user@app.loc"
        PSW_2="password"
    fi
    
    # SERVER_URL will be set after parsing arguments
}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}============================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_suite() {
    echo -e "${PURPLE}🧪 $1${NC}"
}

usage() {
    echo "Usage: $0 {tls|notls} [TEST_SUITE] [OPTIONS]"
    echo ""
    echo "Required Arguments:"
    echo "  tls              Use HTTPS (SSL/TLS enabled)"
    echo "  notls            Use HTTP (no SSL/TLS)"
    echo ""
    echo "Optional Arguments:"
    echo "  TEST_SUITE       Test suite to run (default: all)"
    echo "    all            Run all test suites"
    echo "    verify         Run verify endpoint tests"
    echo "    login          Run login/logout tests"
    echo "    signup         Run signup tests"
    echo ""
    echo "Options:"
    echo "  -h, --help       Show this help message"
    echo "  -v, --verbose    Enable verbose output"
    echo "  -q, --quick      Run only quick tests"
    echo ""
    echo "Examples:"
    echo "  $0 tls                # Run all tests with HTTPS"
    echo "  $0 notls              # Run all tests with HTTP"
    echo "  $0 tls verify         # Run verify tests with HTTPS"
    echo "  $0 notls verify       # Run verify tests with HTTP"
    echo "  $0 tls login          # Run login tests with HTTPS"
    echo ""
}

check_server() {
    echo "🚀 Checking server status..."
    
    local server_health_url="${FULL_SERVER_URL}/health"
    
    if curl -s -k "$server_health_url" > /dev/null 2>&1; then
        echo "✅ Server is already running at $FULL_SERVER_URL"
        return 0
    fi
    
    echo "🏁 Starting test server on port $SERVER_PORT with protocol $SERVER_PROTO..."
    cd "$BASE_DIR"
    
    # Determine TLS setting based on SERVER_ARGS
    local tls_setting="false"
    if [ "$SERVER_ARGS" = "tls" ]; then
        tls_setting="true"
    fi
    
    # Start server with command line arguments
    echo "🚀 Running: go run cmd/main.go -tls $tls_setting -port $SERVER_PORT -debug true"
    go run cmd/main.go -tls $tls_setting -port $SERVER_PORT -debug true &
    SERVER_PID=$!
    
    echo "🎆 Test server starting (PID: $SERVER_PID) - TLS: $tls_setting, Port: $SERVER_PORT, Debug: true"
    
    # Wait for server to be ready
    local timeout=${SERVER_STARTUP_TIMEOUT:-15}
    local interval=${SERVER_STARTUP_INTERVAL:-2}
    
    for i in $(seq 1 $timeout); do
        sleep $interval
        if curl -s -k "$server_health_url" > /dev/null 2>&1; then
            echo "✅ Test server ready at $FULL_SERVER_URL"
            return 0
        fi
        echo "⏳ Waiting for server... ($i/$timeout)"
    done
    
    echo "❌ Server failed to start after $timeout attempts"
    return 1
}

run_verify_tests() {
    print_suite "Running Verify Endpoint Tests"
    cd "$TEST_DIR/verify"
    
    local details=""
    if [ -f "test_verify_simple.sh" ]; then
        echo "Running simple verify tests..."
        if ./test_verify_simple.sh; then
            print_success "Verify tests passed"
            details="Verify tests completed successfully."
            generate_test_report "verify" "PASSED" "$details"
            return 0
        else
            print_error "Verify tests failed"
            details="Verify tests failed with exit code $?"
            generate_test_report "verify" "FAILED" "$details"
            return 1
        fi
    else
        print_error "Verify test script not found"
        details="Test script test_verify_simple.sh not found in $TEST_DIR/verify"
        generate_test_report "verify" "ERROR" "$details"
        return 1
    fi
}

run_login_tests() {
    print_suite "Running Login/Logout Tests"
    cd "$TEST_DIR/login_logout"
    
    local details=""
    if [ -f "test_login_logout_simple.sh" ]; then
        echo "Running login/logout tests..."
        if ./test_login_logout_simple.sh; then
            print_success "Login/logout tests passed"
            details="Login/logout tests completed successfully."
            generate_test_report "login" "PASSED" "$details"
            return 0
        else
            print_error "Login/logout tests failed"
            details="Login/logout tests failed with exit code $?"
            generate_test_report "login" "FAILED" "$details"
            return 1
        fi
    else
        print_error "Login/logout test script not found"
        details="Test script test_login_logout_simple.sh not found in $TEST_DIR/login_logout"
        generate_test_report "login" "ERROR" "$details"
        return 1
    fi
}

run_signup_tests() {
    print_suite "Running Signup Tests"
    
    local details=""
    local test_script="$TEST_DIR/signup/test_signup_simple.sh"
    
    if [ -f "$test_script" ]; then
        print_info "Running simple signup tests..."
        cd "$TEST_DIR/signup"
        
        # Run the signup test script and capture output
        local test_output
        if test_output=$(bash test_signup_simple.sh 2>&1); then
            print_success "Signup tests passed"
            details="Signup endpoint tests completed successfully:\n\n$test_output"
            generate_test_report "signup" "PASSED" "$details"
            return 0
        else
            print_error "Signup tests failed"
            details="Signup endpoint tests failed:\n\n$test_output"
            generate_test_report "signup" "FAILED" "$details"
            return 1
        fi
    else
        print_info "Signup test script not found: $test_script"
        details="Signup test script not found.\nExpected: $test_script"
        generate_test_report "signup" "NOT_FOUND" "$details"
        return 1
    fi
}

run_all_tests() {
    print_header "RUNNING ALL TEST SUITES"
    
    local failed_suites=()
    local total_suites=0
    local details=""
    
    # Run verify tests
    ((total_suites++))
    if ! run_verify_tests; then
        failed_suites+=("verify")
    fi
    
    echo ""
    
    # Run login tests
    ((total_suites++))
    if ! run_login_tests; then
        failed_suites+=("login")
    fi
    
    echo ""
    
    # Run signup tests
    ((total_suites++))
    if ! run_signup_tests; then
        failed_suites+=("signup")
    fi
    
    # Summary
    echo ""
    print_header "TEST SUITE SUMMARY"
    echo "Total Suites: $total_suites"
    echo "Passed: $((total_suites - ${#failed_suites[@]}))"
    echo "Failed: ${#failed_suites[@]}"
    
    # Stop server if we started it
    if [ ! -z "$SERVER_PID" ]; then
        print_info "Stopping server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
    fi
    
    if [ ${#failed_suites[@]} -eq 0 ]; then
        print_success "All test suites passed! 🎉"
        details="All test suites completed successfully:\n- Verify: PASSED\n- Login: PASSED\n- Signup: SKIPPED"
        generate_test_report "all" "PASSED" "$details"
        return 0
    else
        print_error "Failed test suites: ${failed_suites[*]}"
        details="Some test suites failed:\nFailed suites: ${failed_suites[*]}\nTotal suites: $total_suites\nPassed: $((total_suites - ${#failed_suites[@]}))\nFailed: ${#failed_suites[@]}"
        generate_test_report "all" "FAILED" "$details"
        return 1
    fi
}

# Check for help first
if [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    usage
    exit 0
fi

# First argument must be tls or notls
if [[ $# -lt 1 ]]; then
    echo "❌ Error: Missing required argument"
    echo ""
    echo "You must specify either 'tls' or 'notls' as the first argument."
    echo ""
    usage
    exit 1
fi

case "$1" in
    tls)
        SERVER_ARGS="tls"
        SERVER_PROTO="https"
        shift
        ;;
    notls)
        SERVER_ARGS="notls"
        SERVER_PROTO="http"
        shift
        ;;
    *)
        echo "❌ Error: Invalid first argument '$1'"
        echo ""
        echo "First argument must be either 'tls' or 'notls'."
        echo ""
        usage
        exit 1
        ;;
esac

# Parse remaining arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quick)
            QUICK=true
            shift
            ;;
        all|verify|login|signup)
            TEST_SUITE="$1"
            shift
            ;;
        *)
            echo "❌ Error: Unknown option '$1'"
            echo ""
            usage
            exit 1
            ;;
    esac
done

echo "🔧 Server mode: $SERVER_ARGS ($SERVER_PROTO)"

# Load configuration
load_config

# Kill existing test server for fresh start (after config is loaded)
kill_existing_server

# Set complete SERVER_URL based on final server args using config SERVER_URL and SERVER_PORT
if [ "$SERVER_ARGS" = "notls" ]; then
    FULL_SERVER_URL="http://$SERVER_URL:$SERVER_PORT"
else  
    FULL_SERVER_URL="https://$SERVER_URL:$SERVER_PORT"
fi

# Display configuration
echo ""
echo "📋 Test Configuration:"
echo "  Base Directory: $BASE_DIR"
echo "  Test Directory: $TEST_DIR"
echo "  Server Host: $SERVER_URL"
echo "  Server Port: $SERVER_PORT"
echo "  Full Server URL: $FULL_SERVER_URL"
echo "  Server Protocol: $SERVER_PROTO"
echo "  Server Arguments: $SERVER_ARGS"
echo "  Admin User: $USER_1"
echo "  Regular User: $USER_2"
echo "  Test Suite: $TEST_SUITE"
echo ""

# Export variables for child test scripts
export BASE_DIR
export TEST_DIR
export SERVER_URL="$FULL_SERVER_URL"
export SERVER_PORT
export USER_1
export PSW_1
export USER_2
export PSW_2
export MANAGED_BY_RUNNER=true

# Main execution
print_header "APPBASE API TEST RUNNER"
echo "Server: $FULL_SERVER_URL"
echo "Protocol: $SERVER_PROTO"
echo "Test Suite: $TEST_SUITE"
echo "Working Directory: $BASE_DIR"
echo "Admin Credentials: $USER_1 / ****"
echo "User Credentials: $USER_2 / ****"
echo ""

# Check if we're in the right directory
if [ ! -d "$TEST_DIR" ]; then
    print_error "Test directory not found: $TEST_DIR"
    print_info "Please run this script from the appbase root directory"
    exit 1
fi

# Check server availability
if ! check_server; then
    generate_test_report "$TEST_SUITE" "ERROR" "Failed to start server in $SERVER_ARGS mode"
    exit 1
fi

echo ""

# Run the requested test suite
cd "$BASE_DIR"
case $TEST_SUITE in
    all)
        run_all_tests
        ;;
    verify)
        run_verify_tests
        ;;
    login)
        run_login_tests
        ;;
    signup)
        run_signup_tests
        ;;
    *)
        print_error "Unknown test suite: $TEST_SUITE"
        generate_test_report "$TEST_SUITE" "ERROR" "Unknown test suite: $TEST_SUITE"
        usage
        exit 1
        ;;
esac