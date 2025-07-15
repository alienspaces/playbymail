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