/**
 * Utilities for parsing backend logs to extract verification tokens
 */

/**
 * Extract verification token from backend logs
 * Looks for the pattern: "verification token >TOKEN< for account ID >ID<"
 */
export function extractVerificationToken(logs, email) {
  if (!logs || typeof logs !== 'string') {
    return null;
  }
  
  // Look for the specific log pattern your backend uses
  // Pattern: "verification token >TOKEN< for account ID >ID<"
  const pattern = new RegExp(`verification token >([^<]+)< for account ID >[^<]+<`, 'i');
  const match = logs.match(pattern);
  
  if (match && match[1]) {
    console.log(`✅ Found verification token: ${match[1]}`);
    return match[1];
  }
  
  // Alternative pattern: just look for "verification token >TOKEN<"
  const altPattern = /verification token >([^<]+)</i;
  const altMatch = logs.match(altPattern);
  
  if (altMatch && altMatch[1]) {
    console.log(`✅ Found verification token (alt): ${altMatch[1]}`);
    return altMatch[1];
  }
  
  console.log('❌ No verification token found in logs');
  console.log('Logs:', logs);
  return null;
}

/**
 * Wait for verification token to appear in backend logs
 * This function will be called from tests that have access to the webServer
 */
export async function waitForVerificationToken(webServer, email, timeout = 10000) {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    
    const checkLogs = async () => {
      try {
        // Get backend logs from the webServer
        const backendProcess = webServer[0]; // Backend server
        if (!backendProcess || !backendProcess.stdout) {
          reject(new Error('Backend process not available'));
          return;
        }
        
        // Read available logs
        const logs = backendProcess.stdout.read() || '';
        const token = extractVerificationToken(logs, email);
        
        if (token) {
          resolve(token);
          return;
        }
        
        if (Date.now() - startTime > timeout) {
          reject(new Error(`Timeout waiting for verification token for ${email}`));
          return;
        }
        
        // Check again in 100ms
        setTimeout(checkLogs, 100);
      } catch (error) {
        reject(error);
      }
    };
    
    checkLogs();
  });
}

/**
 * Get all available logs from backend process
 */
export function getBackendLogs(webServer) {
  if (!webServer || !webServer[0] || !webServer[0].stdout) {
    return '';
  }
  
  try {
    return webServer[0].stdout.read() || '';
  } catch (error) {
    console.log('Error reading backend logs:', error.message);
    return '';
  }
}

/**
 * Clear backend logs (useful for test isolation)
 */
export function clearBackendLogs(webServer) {
  if (webServer && webServer[0] && webServer[0].stdout) {
    try {
      // This is a bit hacky but works for testing
      webServer[0].stdout._readableState.buffer.clear();
    } catch (error) {
      // Ignore errors when clearing logs
    }
  }
}

/**
 * Alternative approach: Monitor logs by making API calls and checking responses
 * This doesn't require direct access to webServer processes
 */
export async function waitForVerificationTokenViaAPI(page, email, timeout = 10000) {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    
    const checkForToken = async () => {
      try {
        // Make a request to see if we can detect the token being generated
        // This is a fallback approach when direct log access isn't available
        
        if (Date.now() - startTime > timeout) {
          reject(new Error(`Timeout waiting for verification token for ${email}`));
          return;
        }
        
        // For now, we'll use a simple timeout approach
        // In a real implementation, you might check the database or make API calls
        setTimeout(checkForToken, 100);
      } catch (error) {
        reject(error);
      }
    };
    
    checkForToken();
  });
}
