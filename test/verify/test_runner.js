#!/usr/bin/env node

/**
 * Test runner for verify endpoints
 * Usage: node test_runner.js [test-file.json]
 */

const https = require('https');
const http = require('http');
const fs = require('fs');

// Ignore self-signed SSL certificate warnings
process.env["NODE_TLS_REJECT_UNAUTHORIZED"] = 0;

class TestRunner {
  constructor(configPath) {
    let configData = fs.readFileSync(configPath, 'utf8');
    
    // Substitute environment variables in JSON config
    configData = configData.replace(/\${([^}]+)}/g, (match, varName) => {
      const [envVar, defaultValue] = varName.split(':-');
      return process.env[envVar] || defaultValue || match;
    });
    
    this.testConfig = JSON.parse(configData);
    this.results = [];
  }

  async runTests() {
    console.log(`🚀 Running ${this.testConfig.info.name}`);
    console.log(`📋 Description: ${this.testConfig.info.description}`);
    console.log(`🌐 Base URL: ${this.testConfig.baseUrl}\n`);

    // Start server if needed
    await this.startServer();

    for (const test of this.testConfig.tests) {
      await this.runTest(test);
    }

    // Stop server if we started it
    await this.stopServer();

    this.printSummary();
  }

  async runTest(test) {
    console.log(`🧪 ${test.name}`);
    console.log(`   ${test.description}`);

    try {
      const options = this.buildRequestOptions(test);
      const response = await this.makeRequest(options, test.body);
      
      const success = this.validateResponse(test, response);
      
      this.results.push({
        name: test.name,
        success,
        status: response.statusCode,
        expected: test.expectedStatus
      });

      if (test.saveToken && success && response.data) {
        this.extractToken(response);
      }

      console.log(`   ${success ? '✅' : '❌'} Status: ${response.statusCode} (expected: ${test.expectedStatus})\n`);
      
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

    if (test.useToken && this.cookieJar) {
      options.headers['Cookie'] = this.cookieJar;
    }

    if (test.body) {
      options.headers['Content-Type'] = 'application/json';
    }

    return options;
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
        req.write(JSON.stringify(body));
      }

      req.end();
    });
  }

  validateResponse(test, response) {
    if (response.statusCode !== test.expectedStatus) {
      return false;
    }

    if (test.validateResponse && response.data) {
      if (test.validateResponse.hasFields) {
        for (const field of test.validateResponse.hasFields) {
          if (!(field in response.data)) {
            console.log(`   ⚠️  Missing field: ${field}`);
            return false;
          }
        }
      }
    }

    return true;
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
            console.log(`   🍪 Token saved`);
            break;
          }
        }
      }
    }
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
    const failed = total - passed;

    console.log('📊 Test Summary:');
    console.log(`   Total: ${total}`);
    console.log(`   ✅ Passed: ${passed}`);
    console.log(`   ❌ Failed: ${failed}`);
    
    if (failed > 0) {
      console.log('\n❌ Failed Tests:');
      this.results
        .filter(r => !r.success)
        .forEach(r => {
          const reason = r.error || `Expected ${r.expected}, got ${r.status}`;
          console.log(`   - ${r.name}: ${reason}`);
        });
    }

    console.log(failed === 0 ? '\n🎉 All tests passed!' : '\n💥 Some tests failed!');
    process.exit(failed === 0 ? 0 : 1);
  }
}

// Run tests
const testFile = process.argv[2] || 'test_verify_endpoints.json';
const runner = new TestRunner(testFile);
runner.runTests().catch(console.error);