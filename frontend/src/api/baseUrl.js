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

/**
 * Extracts error message from backend API error response
 * Backend returns errors as an array of { code, message } objects
 */
export async function handleApiError(res, defaultMessage) {
  if (res.status === 401) {
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

export async function apiFetch(url, options = {}) {
  const res = await fetch(url, options);
  if (res.status === 401) {
    useAuthStore().logout('session_expired');
    throw new Error('Session expired. Please log in again.');
  }
  return res;
} 