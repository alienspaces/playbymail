import { useAuthStore } from '../stores/auth';

const isLocalhost =
  window.location.hostname === 'localhost' ||
  window.location.hostname === '127.0.0.1';

export const baseUrl = isLocalhost
  ? 'http://localhost:8080'
  : window.location.origin;

export function getAuthHeaders() {
  const authStore = useAuthStore();
  const token = authStore.sessionToken;

  return token ? { Authorization: `Bearer ${token}` } : {};
}

// Sentinel error class for unauthenticated responses. Callers that pass
// skipAutoLogout: true to apiFetch/handleApiError will receive this instead
// of being redirected to /login. Used by the player turn sheet flow, which
// re-exchanges the turn sheet token for a fresh session rather than logging in.
export class UnauthenticatedError extends Error {
  constructor() {
    super('unauthenticated');
    this.name = 'UnauthenticatedError';
  }
}

/**
 * Extracts error message from backend API error response.
 * Backend returns errors as an array of { code, message } objects.
 *
 * Options:
 *   skipAutoLogout {boolean} - When true, throws UnauthenticatedError on 401
 *     instead of calling logout() and redirecting. Use for flows where the
 *     caller handles session recovery (e.g. player turn sheet re-auth).
 */
export async function handleApiError(res, defaultMessage, { skipAutoLogout = false } = {}) {
  if (res.status === 401) {
    if (skipAutoLogout) {
      throw new UnauthenticatedError();
    }
    useAuthStore().logout('session_expired');
    throw new Error('Session expired. Please log in again.');
  }

  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    let errorMessage = defaultMessage || 'Request failed';

    // Backend returns array of errors with code and message
    if (Array.isArray(errorData) && errorData.length > 0) {
      errorMessage = errorData[0].message || errorMessage;
    } else if (errorData.error) {
      // Fallback for error object format
      errorMessage = errorData.error.message || errorData.error.detail || errorMessage;
    } else if (errorData.message) {
      // Fallback for direct message field
      errorMessage = errorData.message;
    }

    throw new Error(errorMessage);
  }

  return res;
}

/**
 * Options:
 *   skipAutoLogout {boolean} - See handleApiError above.
 */
export async function apiFetch(url, options = {}) {
  const { skipAutoLogout = false, ...fetchOptions } = options;
  const res = await fetch(url, fetchOptions);
  if (res.status === 401) {
    if (skipAutoLogout) {
      throw new UnauthenticatedError();
    }
    useAuthStore().logout('session_expired');
    throw new Error('Session expired. Please log in again.');
  }
  return res;
}