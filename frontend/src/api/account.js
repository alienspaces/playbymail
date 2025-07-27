import { baseUrl, getAuthHeaders } from './baseUrl';

export async function getMyAccount() {
  const res = await fetch(`${baseUrl}/v1/my-account`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to fetch account');
  const data = await res.json();
  return data.data;
}

export async function updateMyAccount(accountData) {
  const res = await fetch(`${baseUrl}/v1/my-account`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(accountData)
  });
  if (!res.ok) throw new Error('Failed to update account');
  const data = await res.json();
  return data.data;
}

export async function deleteMyAccount() {
  const res = await fetch(`${baseUrl}/v1/my-account`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to delete account');
  return true;
} 