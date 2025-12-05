import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all creatures for a game.
 * @param {string} gameId
 * @returns {Promise<GameCreature[]>}
 */
export async function fetchCreatures(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creatures`, {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch creatures');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new creature for a game.
 * @param {string} gameId
 * @param {Partial<GameCreature>} data
 * @returns {Promise<GameCreature>}
 */
export async function createCreature(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creatures`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create creature');
  const json = await res.json();
  return json.data;
}

/**
 * Update a creature by ID.
 * @param {string} gameId
 * @param {string} creatureId
 * @param {Partial<GameCreature>} data
 * @returns {Promise<GameCreature>}
 */
export async function updateCreature(gameId, creatureId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creatures/${encodeURIComponent(creatureId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update creature');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a creature by ID.
 * @param {string} gameId
 * @param {string} creatureId
 * @returns {Promise<void>}
 */
export async function deleteCreature(gameId, creatureId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creatures/${encodeURIComponent(creatureId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete creature');
} 