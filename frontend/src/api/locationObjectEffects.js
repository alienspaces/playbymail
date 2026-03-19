import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all location object effects for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationObjectEffect[]>}
 */
export async function fetchLocationObjectEffects(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-object-effects`, {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch location object effects');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new location object effect for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationObjectEffect>} data
 * @returns {Promise<GameLocationObjectEffect>}
 */
export async function createLocationObjectEffect(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-object-effects`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to create location object effect');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location object effect by ID.
 * @param {string} gameId
 * @param {string} locationObjectEffectId
 * @param {Partial<GameLocationObjectEffect>} data
 * @returns {Promise<GameLocationObjectEffect>}
 */
export async function updateLocationObjectEffect(gameId, locationObjectEffectId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-object-effects/${encodeURIComponent(locationObjectEffectId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to update location object effect');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location object effect by ID.
 * @param {string} gameId
 * @param {string} locationObjectEffectId
 * @returns {Promise<void>}
 */
export async function deleteLocationObjectEffect(gameId, locationObjectEffectId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-object-effects/${encodeURIComponent(locationObjectEffectId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location object effect');
}
