import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function getMe() {
  const res = await apiFetch(`${baseUrl}/api/v1/me`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  await handleApiError(res, 'Failed to fetch account user');
  const data = await res.json();
  return data.data;
}

export async function deleteAccountUser(accountId, accountUserId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  await handleApiError(res, 'Failed to delete account');
  return true;
}

export async function getAccountContacts(accountId, accountUserId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}/contacts`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  await handleApiError(res, 'Failed to fetch account contacts');
  const data = await res.json();
  return data.data || [];
}

export async function getAccountContact(accountId, accountUserId, contactId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}/contacts/${contactId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  await handleApiError(res, 'Failed to fetch account contact');
  const data = await res.json();
  return data.data;
}

export async function createAccountContact(accountId, accountUserId, contactData) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}/contacts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(contactData)
  });
  await handleApiError(res, 'Failed to create account contact');
  const data = await res.json();
  return data.data;
}

export async function updateAccountContact(accountId, accountUserId, contactId, contactData) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}/contacts/${contactId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(contactData)
  });
  await handleApiError(res, 'Failed to update account contact');
  const data = await res.json();
  return data.data;
}

export async function deleteAccountContact(accountId, accountUserId, contactId) {
  const res = await apiFetch(`${baseUrl}/api/v1/accounts/${accountId}/users/${accountUserId}/contacts/${contactId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() }
  });
  await handleApiError(res, 'Failed to delete account contact');
  return true;
}
