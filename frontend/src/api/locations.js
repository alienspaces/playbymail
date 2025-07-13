import { API_BASE_URL } from './baseUrl';

/**
 * Fetch all locations for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocation[]>}
 */
export async function fetchLocations(gameId) {
  const res = await fetch(`${API_BASE_URL}/game-locations?game_id=${encodeURIComponent(gameId)}`);
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
  const res = await fetch(`${API_BASE_URL}/game-locations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ...data, game_id: gameId })
  });
  if (!res.ok) throw new Error('Failed to create location');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location by ID.
 * @param {string} locationId
 * @param {Partial<GameLocation>} data
 * @returns {Promise<GameLocation>}
 */
export async function updateLocation(locationId, data) {
  const res = await fetch(`${API_BASE_URL}/game-locations/${encodeURIComponent(locationId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update location');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location by ID.
 * @param {string} locationId
 * @returns {Promise<void>}
 */
export async function deleteLocation(locationId) {
  const res = await fetch(`${API_BASE_URL}/game-locations/${encodeURIComponent(locationId)}`, {
    method: 'DELETE'
  });
  if (!res.ok) throw new Error('Failed to delete location');
} 