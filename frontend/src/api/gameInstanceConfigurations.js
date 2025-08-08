import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

// Get all configurations for a specific game instance
export async function listGameInstanceConfigurations(gameInstanceId, params = {}) {
  const queryParams = new URLSearchParams();
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instance configurations');
  return await res.json();
}

// Get a specific game instance configuration by ID
export async function getGameInstanceConfiguration(gameInstanceId, configurationId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations/${configurationId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instance configuration');
  return await res.json();
}

// Create a new game instance configuration
export async function createGameInstanceConfiguration(gameInstanceId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create game instance configuration');
  return await res.json();
}

// Update an existing game instance configuration
export async function updateGameInstanceConfiguration(gameInstanceId, configurationId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations/${configurationId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update game instance configuration');
  return await res.json();
}

// Delete a game instance configuration
export async function deleteGameInstanceConfiguration(gameInstanceId, configurationId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations/${configurationId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete game instance configuration');
  return await res.json();
}

// Bulk update game instance configurations
export async function bulkUpdateGameInstanceConfigurations(gameInstanceId, configurations) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-instances/${gameInstanceId}/configurations/bulk`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ configurations }),
  });
  if (!res.ok) throw new Error('Failed to update game instance configurations');
  return await res.json();
} 