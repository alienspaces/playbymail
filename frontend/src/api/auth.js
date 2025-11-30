import { baseUrl, getAuthHeaders } from './baseUrl';

/**
 * Get test bypass headers for development mode.
 * When APP_ENV=develop and bypass headers are configured, returns the bypass
 * header that allows using email as the verification code.
 * @returns {Object} Headers object with bypass header if in develop mode
 */
function getTestBypassHeaders() {
  const appEnv = import.meta.env.VITE_APP_ENV;
  const bypassHeaderName = import.meta.env.VITE_TEST_BYPASS_HEADER_NAME;
  const bypassHeaderValue = import.meta.env.VITE_TEST_BYPASS_HEADER_VALUE;

  if (appEnv === 'develop' && bypassHeaderName && bypassHeaderValue) {
    return { [bypassHeaderName]: bypassHeaderValue };
  }
  return {};
}

export async function requestAuth(email) {
  const res = await fetch(`${baseUrl}/api/v1/request-auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ email }),
  });
  return res.ok;
}

export async function verifyAuth(email, verification_token) {
  const res = await fetch(`${baseUrl}/api/v1/verify-auth`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
      ...getTestBypassHeaders(),
    },
    body: JSON.stringify({ email, verification_token }),
  });
  if (!res.ok) throw new Error('Verification failed');
  const data = await res.json();
  return data.session_token;
}

/**
 * Refresh the current session token.
 * Returns the session status and expiry time in seconds.
 * @returns {Promise<{status: string, expires_in_seconds: number}>}
 * @throws {Error} if the session is invalid or expired
 */
export async function refreshSession() {
  const res = await fetch(`${baseUrl}/api/v1/refresh-session`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) {
    throw new Error('Session refresh failed');
  }
  return await res.json();
} 