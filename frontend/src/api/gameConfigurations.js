import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

// Get all game configurations (with optional filtering by game type)
export async function listGameConfigurations(params = {}) {
  const queryParams = new URLSearchParams();
  if (params.gameType) queryParams.append('game_type', params.gameType);
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/game-configurations${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game configurations');
  return await res.json();
}

// Get a specific game configuration by ID
export async function getGameConfiguration(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-configurations/${id}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game configuration');
  return await res.json();
}

// Create a new game configuration
export async function createGameConfiguration(data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-configurations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create game configuration');
  return await res.json();
}

// Update an existing game configuration
export async function updateGameConfiguration(id, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-configurations/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update game configuration');
  return await res.json();
}

// Delete a game configuration
export async function deleteGameConfiguration(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-configurations/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete game configuration');
  return await res.json();
} 