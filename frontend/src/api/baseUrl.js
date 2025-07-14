const isLocalhost =
  window.location.hostname === 'localhost' ||
  window.location.hostname === '127.0.0.1';

export const baseUrl = isLocalhost
  ? 'http://localhost:8080'
  : window.location.origin; 