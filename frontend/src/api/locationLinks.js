import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

/**
 * Fetch all location links for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationLink[]>}
 */
export async function fetchLocationLinks(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch location links');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new location link for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationLink>} data
 * @returns {Promise<GameLocationLink>}
 */
export async function createLocationLink(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to create location link');
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
export async function updateLocationLink(gameId, locationLinkId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links/${encodeURIComponent(locationLinkId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update location link');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location link by ID.
 * @param {string} gameId
 * @param {string} locationLinkId
 * @returns {Promise<void>}
 */
export async function deleteLocationLink(gameId, locationLinkId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-links/${encodeURIComponent(locationLinkId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete location link');
} 