#!/usr/bin/env node

/**
 * Test runner for login/logout endpoints
 * Usage: node test_login_logout_runner.js [test-file.json]
 */

const https = require('https');
const http = require('http');
const fs = require('fs');

// Ignore self-signed SSL certificate warnings
process.env["NODE_TLS_REJECT_UNAUTHORIZED"] = 0;

class LoginLogoutTestRunner {
  constructor(testFile = 'test_login_logout_endpoints.json') {
    let configData = fs.readFileSync(testFile, 'utf8');
    
    // Substitute environment variables in JSON config
    configData = configData.replace(/\${([^}]+)}/g, (match, varName) => {
      const [envVar, defaultValue] = varName.split(':-');
      return process.env[envVar] || defaultValue || match;
    });
    
    this.testConfig = JSON.parse(configData);
    this.token = null;
    this.cookieJar = '';
    this.results = [];
    this.testIndex = 0;
  }

  async runTests() {
    console.log(`🔐 Running ${this.testConfig.info.name}`);
    console.log(`📋 Description: ${this.testConfig.info.description}`);
    console.log(`🌐 Base URL: ${this.testConfig.baseUrl}\n`);

    // Start server if needed
    await this.startServer();

    for (const test of this.testConfig.tests) {
      await this.runTest(test);
      this.testIndex++;
      
      // Small delay between tests to avoid overwhelming the server
      await this.sleep(100);
    }

    // Stop server if we started it
    await this.stopServer();

    this.printSummary();
  }

  async runTest(test) {
    console.log(`🧪 Test ${this.testIndex + 1}: ${test.name}`);
    console.log(`   ${test.description}`);

    try {
      // Skip dependent tests if their dependency failed
      if (test.dependsOn && this.hasDependencyFailed(test.dependsOn)) {
        console.log(`   ⚠️  Skipped - dependency failed: ${test.dependsOn}\n`);
        this.results.push({
          name: test.name,
          success: false,
          skipped: true,
          reason: `Dependency failed: ${test.dependsOn}`
        });
        return;
      }

      const options = this.buildRequestOptions(test);
      const body = this.buildRequestBody(test);
      const response = await this.makeRequest(options, body);
      
      const success = this.validateResponse(test, response);
      
      this.results.push({
        name: test.name,
        success,
        status: response.statusCode,
        expected: test.expectedStatus,
        response: response.data
      });

      // Handle token saving
      if (test.saveToken && success && response.data) {
        this.extractToken(response);
      }

      // Handle multiple expected status codes
      const expectedStatuses = Array.isArray(test.expectedStatus) ? 
        test.expectedStatus : [test.expectedStatus];
      const expectedStr = expectedStatuses.join(' or ');

      console.log(`   ${success ? '✅' : '❌'} Status: ${response.statusCode} (expected: ${expectedStr})`);
      
      // Special handling for optional user login
      if (test.name.includes('Regular User') && response.statusCode === 401) {
        console.log('   ⚠️  Regular user not found in test database - this is expected');
      }
      
      // Additional validation output
      if (success && test.validateResponse) {
        await this.validateResponseContent(test, response);
      }

      console.log('');
      
    } catch (error) {
      console.log(`   ❌ Error: ${error.message}\n`);
      this.results.push({
        name: test.name,
        success: false,
        error: error.message
      });
    }
  }

  buildRequestOptions(test) {
    const url = new URL(test.url, this.testConfig.baseUrl);
    const options = {
      hostname: url.hostname,
      port: url.port || (url.protocol === 'https:' ? 443 : 80),
      path: url.pathname,
      method: test.method,
      headers: {
        ...test.headers
      }
    };

    // Add authentication if needed
    if (test.useToken && this.cookieJar) {
      options.headers['Cookie'] = this.cookieJar;
    }

    return options;
  }

  buildRequestBody(test) {
    if (test.bodyRaw) {
      return test.bodyRaw;
    }
    
    if (test.body) {
      if (!test.headers || !test.headers['Content-Type']) {
        // If no Content-Type header, don't set it (to test missing header scenario)
        return JSON.stringify(test.body);
      }
      return JSON.stringify(test.body);
    }
    
    return null;
  }

  makeRequest(options, body) {
    return new Promise((resolve, reject) => {
      const protocol = options.port === 443 ? https : http;
      
      const req = protocol.request(options, (res) => {
        let data = '';
        
        res.on('data', (chunk) => {
          data += chunk;
        });
        
        res.on('end', () => {
          try {
            const jsonData = data ? JSON.parse(data) : {};
            resolve({
              statusCode: res.statusCode,
              headers: res.headers,
              data: jsonData,
              rawData: data
            });
          } catch (e) {
            // If JSON parsing fails, return raw data
            resolve({
              statusCode: res.statusCode,
              headers: res.headers,
              data: null,
              rawData: data
            });
          }
        });
      });

      req.on('error', (error) => {
        reject(error);
      });

      if (body) {
        req.write(body);
      }

      req.end();
    });
  }

  validateResponse(test, response) {
    // Handle multiple expected status codes
    const expectedStatuses = Array.isArray(test.expectedStatus) ? 
      test.expectedStatus : [test.expectedStatus];
    
    return expectedStatuses.includes(response.statusCode);
  }

  async validateResponseContent(test, response) {
    if (!test.validateResponse || !response.data) return;

    const validation = test.validateResponse;

    // Check required fields
    if (validation.hasFields) {
      for (const field of validation.hasFields) {
        if (field in response.data) {
          console.log(`   ✅ Has required field: ${field}`);
        } else {
          console.log(`   ❌ Missing required field: ${field}`);
        }
      }
    }

    // Check user fields (nested validation)
    if (validation.userFields && response.data.user) {
      for (const field of validation.userFields) {
        if (field in response.data.user) {
          console.log(`   ✅ User has required field: ${field}`);
        } else {
          console.log(`   ❌ User missing required field: ${field}`);
        }
      }
    }

    // Check message pattern
    if (validation.messagePattern && response.data.message) {
      if (response.data.message.includes(validation.messagePattern)) {
        console.log(`   ✅ Message contains expected pattern: ${validation.messagePattern}`);
      } else {
        console.log(`   ❌ Message pattern not found. Expected: ${validation.messagePattern}, Got: ${response.data.message}`);
      }
    }
  }

  extractToken(response) {
    // Extract from Set-Cookie headers
    const setCookies = response.headers['set-cookie'];
    if (setCookies) {
      for (const cookie of setCookies) {
        if (cookie.includes('AppBase_Auth=')) {
          const match = cookie.match(/AppBase_Auth=([^;]+)/);
          if (match) {
            this.cookieJar = `AppBase_Auth=${match[1]}`;
            console.log(`   🍪 Authentication token saved`);
            return;
          }
        }
      }
    }

    // Fallback: extract from response body
    if (response.data && response.data.token) {
      this.token = response.data.token;
      console.log(`   🔑 Token saved from response body`);
    }
  }

  hasDependencyFailed(dependencyName) {
    const dependencyResult = this.results.find(r => r.name === dependencyName);
    return dependencyResult && !dependencyResult.success;
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
  
  async startServer() {
    const { spawn } = require('child_process');
    
    console.log('🚀 Starting server with make run...');
    
    // Check if server is already running
    try {
      const response = await this.makeRequest({
        hostname: 'localhost',
        port: 5443,
        path: '/health',
        method: 'GET'
      });
      
      if (response.statusCode === 200) {
        console.log('✅ Server already running');
        return;
      }
    } catch (e) {
      // Server not running, start it
    }
    
    // Start server
    this.serverProcess = spawn('make', ['run'], {
      cwd: '/home/anv/appbase',
      detached: true,
      stdio: ['ignore', 'pipe', 'pipe']
    });
    
    console.log(`🎆 Server starting (PID: ${this.serverProcess.pid})`);
    
    // Wait for server to be ready
    for (let i = 0; i < 10; i++) {
      await this.sleep(2000);
      try {
        const response = await this.makeRequest({
          hostname: 'localhost',
          port: 5443,
          path: '/health',
          method: 'GET'
        });
        
        if (response.statusCode === 200) {
          console.log('✅ Server ready');
          return;
        }
      } catch (e) {
        console.log(`⏳ Waiting for server... (${i + 1}/10)`);
      }
    }
    
    throw new Error('Server failed to start');
  }
  
  async stopServer() {
    if (this.serverProcess) {
      console.log('\n🛑 Stopping server...');
      process.kill(-this.serverProcess.pid, 'SIGTERM');
      this.serverProcess = null;
    }
  }

  printSummary() {
    const total = this.results.length;
    const passed = this.results.filter(r => r.success).length;
    const failed = this.results.filter(r => !r.success && !r.skipped).length;
    const skipped = this.results.filter(r => r.skipped).length;

    console.log('📊 Test Summary:');
    console.log(`   Total: ${total}`);
    console.log(`   ✅ Passed: ${passed}`);
    console.log(`   ❌ Failed: ${failed}`);
    if (skipped > 0) {
      console.log(`   ⚠️  Skipped: ${skipped}`);
    }

    const successRate = total > 0 ? Math.round((passed / total) * 100) : 0;
    console.log(`   📈 Success Rate: ${successRate}%`);
    
    if (failed > 0) {
      console.log('\n❌ Failed Tests:');
      this.results
        .filter(r => !r.success && !r.skipped)
        .forEach(r => {
          const reason = r.error || `Expected ${Array.isArray(r.expected) ? r.expected.join('/') : r.expected}, got ${r.status}`;
          console.log(`   - ${r.name}: ${reason}`);
        });
    }

    if (skipped > 0) {
      console.log('\n⚠️  Skipped Tests:');
      this.results
        .filter(r => r.skipped)
        .forEach(r => {
          console.log(`   - ${r.name}: ${r.reason}`);
        });
    }

    // Rate limiting analysis
    const rateLimitTests = this.results.filter(r => r.name.includes('Rate Limiting'));
    if (rateLimitTests.length > 0) {
      console.log('\n🚦 Rate Limiting Analysis:');
      const rateLimitTriggered = rateLimitTests.some(r => r.status === 429);
      if (rateLimitTriggered) {
        console.log('   ✅ Rate limiting is working');
      } else {
        console.log('   ⚠️  Rate limiting not observed (may need more requests or different configuration)');
      }
    }

    console.log(failed === 0 ? '\n🎉 All tests passed!' : '\n💥 Some tests failed!');
    process.exit(failed === 0 ? 0 : 1);
  }
}

// Run tests
const testFile = process.argv[2] || 'test_login_logout_endpoints.json';

if (!fs.existsSync(testFile)) {
  console.error(`❌ Test file not found: ${testFile}`);
  console.log('Available test files:');
  const files = fs.readdirSync('.').filter(f => f.endsWith('.json'));
  files.forEach(f => console.log(`   - ${f}`));
  process.exit(1);
}

const runner = new LoginLogoutTestRunner(testFile);
runner.runTests().catch(console.error);