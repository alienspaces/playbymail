import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

export async function getMyAccount() {
  const res = await apiFetch(`${baseUrl}/api/v1/my-account`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to fetch account');
  const data = await res.json();
  return data.data;
}

export async function updateMyAccount(accountData) {
  const res = await apiFetch(`${baseUrl}/api/v1/my-account`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(accountData)
  });
  if (!res.ok) throw new Error('Failed to update account');
  const data = await res.json();
  return data.data;
}

export async function deleteMyAccount() {
  const res = await apiFetch(`${baseUrl}/api/v1/my-account`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to delete account');
  return true;
}

export async function getAccountContacts(accountId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/contacts`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to fetch account contacts');
  const data = await res.json();
  return data.data || [];
}

export async function getAccountContact(accountId, contactId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/contacts/${contactId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to fetch account contact');
  const data = await res.json();
  return data.data;
}

export async function createAccountContact(accountId, contactData) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/contacts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(contactData)
  });
  if (!res.ok) throw new Error('Failed to create account contact');
  const data = await res.json();
  return data.data;
}

export async function updateAccountContact(accountId, contactId, contactData) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/contacts/${contactId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(contactData)
  });
  if (!res.ok) throw new Error('Failed to update account contact');
  const data = await res.json();
  return data.data;
}

export async function deleteAccountContact(accountId, contactId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/contacts/${contactId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  if (!res.ok) throw new Error('Failed to delete account contact');
  return true;
}
