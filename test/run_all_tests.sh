#!/bin/bash

# Master test runner for all endpoint tests
# Usage: ./run_all_tests.sh [test_suite] [options]

set -e

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
        SERVER_URL="${SERVER_URL:-https://localhost:5443}"
        USER_1="${USER_1:-admin@app.loc}"
        PSW_1="${PSW_1:-appadmin}"
        USER_2="${USER_2:-user@app.loc}"
        PSW_2="${PSW_2:-password}"
        
        echo "✓ Configuration loaded successfully"
        echo "  Base Directory: $BASE_DIR"
        echo "  Test Directory: $TEST_DIR" 
        echo "  Server URL: $SERVER_URL"
        echo "  Admin User: $USER_1"
        echo "  Regular User: $USER_2"
    else
        echo "⚠️  Configuration file not found: $config_file"
        echo "   Stopping execution."
        exit 1
    fi
}

# Load configuration first
load_config

# Export variables for child test scripts
export BASE_DIR
export TEST_DIR
export SERVER_URL
export USER_1
export PSW_1
export USER_2
export PSW_2
export MANAGED_BY_RUNNER=true

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
    echo "Usage: $0 [OPTIONS] [TEST_SUITE]"
    echo ""
    echo "Configuration:"
    echo "  Settings are loaded from ./test.env file"
    echo "  Override with command line options when needed"
    echo ""
    echo "Test Suites:"
    echo "  all          Run all test suites (default)"
    echo "  verify       Run verify endpoint tests"
    echo "  login        Run login/logout tests"
    echo "  signup       Run signup tests"
    echo ""
    echo "Options:"
    echo "  -h, --help   Show this help message"
    echo "  -v, --verbose Enable verbose output"
    echo "  -q, --quick  Run only quick tests"
    echo "  --server URL Override server URL from config"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all tests with config from test.env"
    echo "  $0 verify             # Run only verify tests"
    echo "  $0 --server http://localhost:8080 verify  # Override server URL"
    echo ""
}

check_server() {
    echo "🚀 Checking server status..."
    
    # Extract host and port from SERVER_URL
    local server_health_url="${SERVER_URL}/health"
    
    if curl -s -k "$server_health_url" > /dev/null 2>&1; then
        echo "✅ Server is already running at $SERVER_URL"
        return 0
    fi
    
    echo "🏁 Starting server with make run..."
    cd "$BASE_DIR"
    make run &
    SERVER_PID=$!
    
    echo "🎆 Server starting (PID: $SERVER_PID)"
    
    # Wait for server to be ready
    local timeout=${SERVER_STARTUP_TIMEOUT:-15}
    local interval=${SERVER_STARTUP_INTERVAL:-2}
    
    for i in $(seq 1 $timeout); do
        sleep $interval
        if curl -s -k "$server_health_url" > /dev/null 2>&1; then
            echo "✅ Server ready at $SERVER_URL"
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
    
    if [ -f "test_verify_simple.sh" ]; then
        echo "Running simple verify tests..."
        if ./test_verify_simple.sh; then
            print_success "Verify tests passed"
            return 0
        else
            print_error "Verify tests failed"
            return 1
        fi
    else
        print_error "Verify test script not found"
        return 1
    fi
}

run_login_tests() {
    print_suite "Running Login/Logout Tests"
    cd "$TEST_DIR/login_logout"
    
    if [ -f "test_login_logout_simple.sh" ]; then
        echo "Running login/logout tests..."
        if ./test_login_logout_simple.sh; then
            print_success "Login/logout tests passed"
            return 0
        else
            print_error "Login/logout tests failed"
            return 1
        fi
    else
        print_error "Login/logout test script not found"
        return 1
    fi
}

run_signup_tests() {
    print_suite "Running Signup Tests"
    cd "$TEST_DIR/signup"
    
    # Check if signup tests exist  
    if [ -d "$TEST_DIR/signup" ] && [ "$(ls -A $TEST_DIR/signup 2>/dev/null)" ]; then
        print_info "Signup test directory found, but no runner script detected"
        print_info "Skipping signup tests for now"
        return 0
    else
        print_info "Signup tests not yet implemented"
        return 0
    fi
}

run_all_tests() {
    print_header "RUNNING ALL TEST SUITES"
    
    local failed_suites=()
    local total_suites=0
    
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
        return 0
    else
        print_error "Failed test suites: ${failed_suites[*]}"
        return 1
    fi
}

# Parse command line arguments
VERBOSE=false
QUICK=false
TEST_SUITE="all"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quick)
            QUICK=true
            shift
            ;;
        --server)
            SERVER_URL="$2"
            shift 2
            ;;
        all|verify|login|signup)
            TEST_SUITE="$1"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Main execution
print_header "APPBASE API TEST RUNNER"
echo "Server: $SERVER_URL"
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
        usage
        exit 1
        ;;
esac