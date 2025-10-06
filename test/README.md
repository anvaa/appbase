# API Test Suite

This directory contains comprehensive test suites for all API endpoints in the AppBase application.

## 📁 Directory Structure

```
test/
├── verify/                          # Verify endpoint tests
│   ├── test_verify_simple.sh        # Simple bash test script
│   ├── test_verify_endpoints.sh     # Comprehensive bash test suite
│   ├── test_verify_endpoints.json   # JSON test configuration
│   └── test_runner.js               # Node.js test runner
├── login_logout/                    # Login/logout endpoint tests
├── signup/                          # Signup endpoint tests
└── README.md                        # This file
```

## 🚀 Quick Start

### Run All Tests
```bash
# From the project root directory
./run_all_tests.sh
```

### Run Specific Test Suite
```bash
./run_all_tests.sh verify     # Run only verify tests
./run_all_tests.sh login      # Run only login tests  
./run_all_tests.sh signup     # Run only signup tests
```

### Run Individual Test Scripts
```bash
# Verify endpoint tests
cd test/verify
./test_verify_simple.sh       # Quick verification test
./test_verify_endpoints.sh    # Comprehensive test suite
node test_runner.js           # Node.js test runner
```

## 📋 Test Coverage

### ✅ Verify Endpoints (`/test/verify/`)
- **POST /verify** - Token verification endpoint
- **GET /verify** - Alternative verification method  
- **POST /user/verify** - Protected user verification
- **GET /user/verify** - Protected user verification

**Test Scenarios:**
- ✅ Authentication with valid tokens
- ❌ Authentication with invalid tokens  
- ❌ Requests without authentication
- 📋 Response structure validation
- 🍪 Cookie/session handling

### 🔄 Login/Logout Endpoints (`/test/login_logout/`)
- **POST /login** - User authentication
- **GET /logout** - User logout
- **POST /logout** - User logout via POST

*Status: Directory exists, tests to be implemented*

### 📝 Signup Endpoints (`/test/signup/`)
- **POST /signup** - User registration
- **GET /signup** - Signup form view

*Status: Directory exists, tests to be implemented*

## 🛠️ Test Tools

### Bash Scripts
- **Simple Tests:** Quick verification with basic curl commands
- **Comprehensive Tests:** Full test suites with colored output and detailed reporting
- **Cross-platform:** Works on Linux, macOS, and WSL

### Node.js Test Runner
- **JSON Configuration:** Easy to maintain and extend
- **Professional Output:** Clean reporting with emojis and colors
- **Cookie Management:** Automatic session handling
- **Response Validation:** Structure and content verification

### JSON Test Definitions
- **Declarative:** Define tests in structured JSON format
- **Flexible:** Support for different HTTP methods, headers, and validation
- **Reusable:** Can be used with multiple test runners

## 📊 Example Test Output

```
🧪 Verify Endpoint Tests
📋 Description: Tests for POST /verify and GET /verify endpoints
🌐 Base URL: http://localhost:5443

🧪 Health Check
   Verify server is running
   ✅ Status: 200 (expected: 200)

🧪 Login  
   Login to get authentication token
   🍪 Token saved
   ✅ Status: 200 (expected: 200)

🧪 POST /verify - Valid Auth
   POST /verify with valid authentication
   ✅ Status: 200 (expected: 200)

📊 Test Summary:
   Total: 9
   ✅ Passed: 9
   ❌ Failed: 0

🎉 All tests passed!
```

## ⚙️ Configuration

### Server Configuration
Default server URL: `http://localhost:5443`

To test against a different server:
```bash
./run_all_tests.sh --server http://localhost:8080 verify
```

### Authentication
Default test credentials:
- Username: `admin@app.loc`
- Password: `appadmin`

### Test Data
Tests use the default seeded data from the application. Ensure your test database has:
- Default admin user
- Default organizations
- Default auth levels

## 🔧 Prerequisites

### For Bash Tests
- `curl` - HTTP client for API requests
- `jq` - JSON processor (optional, for response parsing)
- `bash` - Shell interpreter

### For Node.js Tests
- `node` - Node.js runtime (v14 or later recommended)

### Server Requirements
- AppBase server running on configured port (default: 5443)
- Test database with seeded data
- All required endpoints available

## 📝 Adding New Tests

### Adding Bash Tests
1. Create new test script in appropriate directory
2. Use existing scripts as templates
3. Make script executable: `chmod +x test_script.sh`
4. Update master test runner if needed

### Adding JSON Tests  
1. Add test definitions to `test_*.json` files
2. Follow existing JSON schema
3. Test with Node.js runner: `node test_runner.js`

### Adding New Test Suites
1. Create new directory under `/test/`
2. Add test scripts and configuration
3. Update `run_all_tests.sh` to include new suite

## 🐛 Troubleshooting

### Server Not Running
```
✗ Server is not available at http://localhost:5443
ℹ Please start the server with: go run cmd/main.go
```
**Solution:** Start the AppBase server before running tests.

### Authentication Failed
```
✗ Login failed
Response: {"message":"user or password invalid"}
```
**Solutions:**
- Check if test user exists in database
- Verify credentials in test configuration
- Ensure user is activated (`is_auth: true`)

### Permission Denied
```
bash: ./test_script.sh: Permission denied
```
**Solution:** Make script executable: `chmod +x test_script.sh`

### Missing Dependencies
```
curl: command not found
```
**Solution:** Install required tools:
```bash
# Ubuntu/Debian
sudo apt-get install curl jq

# macOS
brew install curl jq

# Windows (WSL)
sudo apt-get install curl jq
```

## 📚 References

- [API Documentation](../docs/api.md)
- [Authentication Guide](../docs/auth.md)
- [Development Setup](../README.md)

---

*Happy Testing! 🧪✨*