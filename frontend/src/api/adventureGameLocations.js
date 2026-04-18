import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all locations for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocation[]>}
 */
export async function fetchAdventureGameLocations(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch locations');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new location for a game.
 * @param {string} gameId
 * @param {Partial<GameLocation>} data
 * @returns {Promise<GameLocation>}
 */
export async function createAdventureGameLocation(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create location');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location by ID.
 * @param {string} gameId
 * @param {string} locationId
 * @param {Partial<GameLocation>} data
 * @returns {Promise<GameLocation>}
 */
export async function updateAdventureGameLocation(gameId, locationId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations/${encodeURIComponent(locationId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update location');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location by ID.
 * @param {string} gameId
 * @param {string} locationId
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameLocation(gameId, locationId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations/${encodeURIComponent(locationId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location');
} 