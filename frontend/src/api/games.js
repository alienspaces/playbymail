import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function listGames(options = {}) {
  const { subscriptionType, status } = options;
  const params = new URLSearchParams();
  if (subscriptionType) {
    params.append('subscription_type', subscriptionType);
  }
  if (status) {
    params.append('status', status);
  }
  const queryString = params.toString();
  const url = `${baseUrl}/api/v1/games${queryString ? `?${queryString}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch games');
  return await res.json();
}

export async function createGame({ name, game_type, turn_duration_hours, description }) {
  const res = await apiFetch(`${baseUrl}/api/v1/games`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ name, game_type, turn_duration_hours, description }),
  });
  await handleApiError(res, 'Failed to create game');
  return await res.json();
}

export async function updateGame(id, { name, game_type, turn_duration_hours, description }) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ name, game_type, turn_duration_hours, description }),
  });
  await handleApiError(res, 'Failed to update game');
  return await res.json();
}

export async function deleteGame(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete game');
  return await res.json();
}

export async function publishGame(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}/publish`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to publish game');
  return await res.json();
}

export async function validateGame(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}/validate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to validate game');
  return await res.json();
}