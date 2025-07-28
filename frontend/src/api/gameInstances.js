import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

export async function listGameInstances(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instances');
  return await res.json();
}

export async function getGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch game instance');
  return await res.json();
}

export async function createGameInstance(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create game instance');
  return await res.json();
}

export async function updateGameInstance(gameId, instanceId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update game instance');
  return await res.json();
}

export async function deleteGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete game instance');
  return await res.json();
}

// Game instance runtime management
export async function startGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}/start`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to start game instance');
  return await res.json();
}

export async function pauseGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}/pause`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to pause game instance');
  return await res.json();
}

export async function resumeGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}/resume`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to resume game instance');
  return await res.json();
}

export async function cancelGameInstance(gameId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${instanceId}/cancel`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to cancel game instance');
  return await res.json();
} 