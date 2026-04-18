import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all location objects for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationObject[]>}
 */
export async function fetchAdventureGameLocationObjects(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch location objects');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new location object for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationObject>} data
 * @returns {Promise<GameLocationObject>}
 */
export async function createAdventureGameLocationObject(gameId, data) {
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
export async function updateAdventureGameLocationObject(gameId, locationObjectId, data) {
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
export async function deleteAdventureGameLocationObject(gameId, locationObjectId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location object');
}
