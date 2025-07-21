import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

/**
 * Fetch all locations for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocation[]>}
 */
export async function fetchLocations(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch locations');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new location for a game.
 * @param {string} gameId
 * @param {Partial<GameLocation>} data
 * @returns {Promise<GameLocation>}
 */
export async function createLocation(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to create location');
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
export async function updateLocation(gameId, locationId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations/${encodeURIComponent(locationId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update location');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location by ID.
 * @param {string} gameId
 * @param {string} locationId
 * @returns {Promise<void>}
 */
export async function deleteLocation(gameId, locationId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/locations/${encodeURIComponent(locationId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete location');
} 