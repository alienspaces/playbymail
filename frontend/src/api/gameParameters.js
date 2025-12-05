import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

// Get all game parameters (with optional filtering by game type)
export async function listGameParameters(params = {}) {
  const queryParams = new URLSearchParams();
  if (params.gameType) queryParams.append('game_type', params.gameType);
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/game-parameters${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game parameters');
  return await res.json();
}

// Get a specific game parameter by ID
export async function getGameParameter(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-parameters/${id}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game parameter');
  return await res.json();
}

// Create a new game parameter
export async function createGameParameter(data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-parameters`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to create game parameter');
  return await res.json();
}

// Update an existing game parameter
export async function updateGameParameter(id, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-parameters/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to update game parameter');
  return await res.json();
}

// Delete a game parameter
export async function deleteGameParameter(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-parameters/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete game parameter');
  return await res.json();
}
