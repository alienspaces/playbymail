import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all location objects for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationObject[]>}
 */
export async function fetchLocationObjects(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects`, {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch location objects');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new location object for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationObject>} data
 * @returns {Promise<GameLocationObject>}
 */
export async function createLocationObject(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to create location object');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location object by ID.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @param {Partial<GameLocationObject>} data
 * @returns {Promise<GameLocationObject>}
 */
export async function updateLocationObject(gameId, locationObjectId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to update location object');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location object by ID.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @returns {Promise<void>}
 */
export async function deleteLocationObject(gameId, locationObjectId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location object');
}
