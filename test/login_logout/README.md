# Login and Logout Endpoint Tests

This directory contains comprehensive test suites for authentication endpoints in the AppBase application.

## 📁 Test Files

```
test/login_logout/
├── test_login_logout_simple.sh          # Simple bash test script
├── test_login_logout_endpoints.sh       # Comprehensive bash test suite
├── test_login_logout_endpoints.json     # JSON test configuration
├── test_login_logout_runner.js          # Node.js test runner
└── README.md                            # This file
```

## 🎯 Test Coverage

### 🔐 Login Endpoint (`POST /login`)
- ✅ **Valid credentials** - Successful authentication
- ❌ **Invalid credentials** - Wrong username/password
- ❌ **Invalid email format** - Malformed email addresses
- ❌ **Missing fields** - Username or password missing
- ❌ **Empty request** - No data provided
- ❌ **Malformed JSON** - Invalid JSON syntax
- ❌ **Missing headers** - No Content-Type header
- 🔒 **SQL injection protection** - Security validation
- 🔒 **XSS protection** - Cross-site scripting prevention
- 📏 **Long input handling** - Very long usernames/passwords
- 🌍 **Unicode support** - International characters
- 📊 **Case sensitivity** - Username case handling

### 🚪 Logout Endpoint (`GET /logout`)
- ✅ **Valid session logout** - Logout with active session
- ✅ **No session logout** - Logout without active session
- 🔄 **Session invalidation** - Verify session is cleared
- 🔒 **Protected access after logout** - Verify access is blocked

### 🛡️ Security & Rate Limiting
- 🚦 **Rate limiting** - Multiple failed login attempts
- 🍪 **Session management** - Cookie handling and persistence
- 🔐 **Token validation** - JWT token verification
- 🛡️ **Security headers** - Proper HTTP security headers

### 📝 Login Form (`GET /login`)
- 🎨 **Form display** - Login form rendering
- 📱 **Content type** - Proper HTML response

## 🚀 Usage

### Quick Test
```bash
cd test/login_logout
./test_login_logout_simple.sh
```

### Comprehensive Test Suite
```bash
cd test/login_logout
./test_login_logout_endpoints.sh
```

### Node.js Test Runner
```bash
cd test/login_logout
node test_login_logout_runner.js
```

### From Project Root
```bash
# Run only login/logout tests
./run_all_tests.sh login

# Run all tests including login/logout
./run_all_tests.sh
```

## 📊 Example Test Output

### Simple Test Output
```
=== Testing Login and Logout Endpoints ===
Server: http://localhost:5443

1. Checking server health...
✓ Server is running

2. Testing login with valid credentials...
Status: 200
Response: {"token":"eyJhbGci...","user":{"username":"admin@app.loc"...}...}
✓ Login successful with valid credentials
✓ Response contains expected fields (token, user)

3. Testing login with invalid credentials...
Status: 401
Response: {"message":"user or password invalid"}
✓ Correctly rejected invalid credentials

...

=== Test completed ===
```

### Comprehensive Test Output
```
🔐 Running Login and Logout Endpoint Tests
📋 Description: Comprehensive tests for authentication endpoints
🌐 Base URL: http://localhost:5443

🧪 Test 1: Login - Valid Credentials
   Login with valid credentials should succeed
   ✅ Status: 200 (expected: 200)
   🍪 Authentication token saved
   ✅ Has required field: token
   ✅ Has required field: user
   ✅ User has required field: username

🧪 Test 2: Login - Invalid Credentials  
   Login with invalid credentials should fail
   ✅ Status: 401 (expected: 401)
   ✅ Message contains expected pattern: user or password invalid

...

📊 Test Summary:
   Total: 24
   ✅ Passed: 24
   ❌ Failed: 0
   📈 Success Rate: 100%

🚦 Rate Limiting Analysis:
   ✅ Rate limiting is working

🎉 All tests passed!
```

## 🔧 Configuration

### Test Credentials
Default credentials for testing:
```json
{
  "valid": {
    "username": "admin@app.loc",
    "password": "appadmin"
  },
  "invalid": {
    "username": "invalid@user.com",
    "password": "wrongpassword"
  }
}
```

### Server Configuration
- **Default URL:** `http://localhost:5443`
- **Rate Limit:** Configured per application settings
- **Session Timeout:** As configured in `usr.yaml`

### Expected Response Structure

#### Successful Login Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "session": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "url": "/",
  "user": {
    "username": "admin@app.loc",
    "email": "admin@app.loc",
    "roles": ["admin"],
    "profile": {
      "uuid": 1173481936,
      "auth_level": "admin",
      "orgs": [...],
      "note": "..."
    }
  }
}
```

#### Failed Login Response
```json
{
  "message": "user or password invalid"
}
```

## 🧪 Test Categories

### 🟢 Positive Tests
- Valid login with correct credentials
- Successful logout
- Session persistence
- Token generation and validation

### 🔴 Negative Tests  
- Invalid credentials
- Missing fields
- Malformed requests
- Unauthorized access after logout

### 🔒 Security Tests
- SQL injection attempts
- XSS protection
- Rate limiting
- Long input handling
- Unicode character handling

### ⚡ Performance Tests
- Rate limiting behavior
- Multiple concurrent requests
- Session cleanup efficiency

## 🐛 Troubleshooting

### Common Issues

#### Authentication Failed
```
✗ Login failed. Expected 200, got 401
Response: {"message":"user or password invalid"}
```
**Solutions:**
- Verify test credentials in database
- Check if user account is activated (`is_auth: true`)
- Confirm password matches database hash

#### Rate Limiting Not Working
```
⚠️ Rate limiting not observed
```
**Solutions:**
- Check rate limit configuration in `usr.yaml`
- Increase number of test attempts
- Verify rate limiting is enabled

#### Server Connection Failed
```
✗ Server is not running
```
**Solutions:**
- Start the server: `go run cmd/main.go`
- Check server port configuration
- Verify firewall settings

### Database Requirements

Ensure test database contains:
- Default admin user (`admin@app.loc`)
- Proper password hash for test credentials
- User with `is_auth: true`
- Required auth levels and organizations

### Dependencies

#### For Bash Tests
- `curl` - HTTP client
- `jq` - JSON parser (optional)
- `bash` - Shell interpreter

#### For Node.js Tests
- `node` - Node.js runtime (v14+)

## 📝 Adding New Tests

### Adding to JSON Configuration
```json
{
  "name": "New Test Case",
  "method": "POST",
  "url": "/login",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "username": "test@example.com",
    "password": "testpass"
  },
  "expectedStatus": 401,
  "description": "Description of what this test validates"
}
```

### Adding to Bash Scripts
```bash
# Test new scenario
run_test "TEST_NAME" \
    "POST" \
    "/login" \
    "-H \"Content-Type: application/json\"" \
    '{"username": "test", "password": "test"}' \
    "401" \
    "Test description"
```

## 🔗 Related Documentation

- [Main Test Suite README](../README.md)
- [API Documentation](../../docs/api.md)
- [Authentication Guide](../../docs/auth.md)
- [Verify Endpoint Tests](../verify/README.md)

---

*Secure Testing! 🔐✨*