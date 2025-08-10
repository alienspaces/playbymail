import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

// Get all parameters for a specific game instance
export async function listGameInstanceParameters(gameInstanceId, params = {}) {
  const queryParams = new URLSearchParams();
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instance parameters');
  return await res.json();
}

// Get a specific game instance parameter by ID
export async function getGameInstanceParameter(gameInstanceId, parameterId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters/${parameterId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instance parameter');
  return await res.json();
}

// Create a new game instance parameter
export async function createGameInstanceParameter(gameInstanceId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create game instance parameter');
  return await res.json();
}

// Update an existing game instance parameter
export async function updateGameInstanceParameter(gameInstanceId, parameterId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters/${parameterId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update game instance parameter');
  return await res.json();
}

// Delete a game instance parameter
export async function deleteGameInstanceParameter(gameInstanceId, parameterId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters/${parameterId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete game instance parameter');
  return await res.json();
}

// Bulk update game instance parameters
export async function bulkUpdateGameInstanceParameters(gameInstanceId, parameters) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/parameters/bulk`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ parameters }),
  });
  if (!res.ok) throw new Error('Failed to update game instance parameters');
  return await res.json();
}
