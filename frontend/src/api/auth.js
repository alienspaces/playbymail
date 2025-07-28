import { baseUrl, getAuthHeaders } from './baseUrl';

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
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ email, verification_token }),
  });
  if (!res.ok) throw new Error('Verification failed');
  const data = await res.json();
  return data.session_token;
} 