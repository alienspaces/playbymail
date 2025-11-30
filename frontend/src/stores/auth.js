import { defineStore } from 'pinia';
import router from '../router';
import { refreshSession } from '../api/auth';

// Refresh session every 5 minutes (well before the 15-minute expiry)
const SESSION_REFRESH_INTERVAL_MS = 5 * 60 * 1000;

export const useAuthStore = defineStore('auth', {
  state: () => ({
    sessionToken: localStorage.getItem('session_token') || '',
    refreshIntervalId: null,
  }),
  actions: {
    setSessionToken(token) {
      this.sessionToken = token;
      if (token) {
        localStorage.setItem('session_token', token);
        console.log('[auth] Stored session token:', token);
        // Start session refresh polling when token is set
        this.startSessionRefresh();
      } else {
        localStorage.removeItem('session_token');
        console.log('[auth] Cleared session token');
        // Stop session refresh polling when token is cleared
        this.stopSessionRefresh();
      }
    },
    logout(code) {
      // Stop session refresh polling
      this.stopSessionRefresh();
      // Always clear the session token on logout
      this.sessionToken = '';
      localStorage.removeItem('session_token');
      console.log('[auth] Cleared session token');
      console.log('[auth] User logged out');
      // Redirect to login with optional code
      router.push({ path: '/login', query: code ? { code } : {} });
    },
    startSessionRefresh() {
      // Don't start if already running or no token
      if (this.refreshIntervalId || !this.sessionToken) {
        return;
      }
      console.log('[auth] Starting session refresh polling');
      this.refreshIntervalId = setInterval(() => {
        this.doRefreshSession();
      }, SESSION_REFRESH_INTERVAL_MS);
    },
    stopSessionRefresh() {
      if (this.refreshIntervalId) {
        console.log('[auth] Stopping session refresh polling');
        clearInterval(this.refreshIntervalId);
        this.refreshIntervalId = null;
      }
    },
    async doRefreshSession() {
      if (!this.sessionToken) {
        this.stopSessionRefresh();
        return;
      }
      try {
        const result = await refreshSession();
        console.log('[auth] Session refreshed, expires in', result.expires_in_seconds, 'seconds');
      } catch (error) {
        console.warn('[auth] Session refresh failed:', error.message);
        // Session is invalid or expired - logout and redirect to login
        this.logout('session_expired');
      }
    },
    // Initialize session refresh on app startup if token exists
    initializeSessionRefresh() {
      if (this.sessionToken) {
        console.log('[auth] Initializing session refresh for existing token');
        this.startSessionRefresh();
        // Also do an immediate refresh to validate the token on startup
        this.doRefreshSession();
      }
    },
  },
});
