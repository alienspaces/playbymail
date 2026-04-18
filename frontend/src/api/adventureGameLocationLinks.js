import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all location links for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationLink[]>}
 */
export async function fetchAdventureGameLocationLinks(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch location links');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new location link for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationLink>} data
 * @returns {Promise<GameLocationLink>}
 */
export async function createAdventureGameLocationLink(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create location link');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location link by ID.
 * @param {string} gameId
 * @param {string} locationLinkId
 * @param {Partial<GameLocationLink>} data
 * @returns {Promise<GameLocationLink>}
 */
export async function updateAdventureGameLocationLink(gameId, locationLinkId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links/${encodeURIComponent(locationLinkId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update location link');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location link by ID.
 * @param {string} gameId
 * @param {string} locationLinkId
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameLocationLink(gameId, locationLinkId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links/${encodeURIComponent(locationLinkId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location link');
} 