import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function listGameInstanceParameters(gameId, gameInstanceId, params = {}) {
  const queryParams = new URLSearchParams();
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/games/${gameId}/instances/${gameInstanceId}/parameters${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game instance parameters');
  return await res.json();
}

export async function getGameInstanceParameter(gameId, gameInstanceId, parameterId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${gameInstanceId}/parameters/${parameterId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game instance parameter');
  return await res.json();
}

export async function createGameInstanceParameter(gameId, gameInstanceId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${gameInstanceId}/parameters`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to create game instance parameter');
  return await res.json();
}

export async function updateGameInstanceParameter(gameId, gameInstanceId, parameterId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${gameInstanceId}/parameters/${parameterId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to update game instance parameter');
  return await res.json();
}

export async function deleteGameInstanceParameter(gameId, gameInstanceId, parameterId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/instances/${gameInstanceId}/parameters/${parameterId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete game instance parameter');
  return await res.json();
}
