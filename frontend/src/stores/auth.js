import { defineStore } from 'pinia';
import router from '../router';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    sessionToken: localStorage.getItem('session_token') || '',
  }),
  actions: {
    setSessionToken(token) {
      this.sessionToken = token;
      if (token) {
        localStorage.setItem('session_token', token);
        console.log('[auth] Stored session token:', token);
      } else {
        localStorage.removeItem('session_token');
        console.log('[auth] Cleared session token');
      }
    },
    logout(code) {
      // Always clear the session token on logout
      this.sessionToken = '';
      localStorage.removeItem('session_token');
      console.log('[auth] Cleared session token');
      console.log('[auth] User logged out');
      // Redirect to login with optional code
      router.push({ path: '/login', query: code ? { code } : {} });
    },
  },
}); 