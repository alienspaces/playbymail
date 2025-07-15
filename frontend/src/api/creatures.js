import { baseUrl, getAuthHeaders } from './baseUrl';

/**
 * Fetch all creatures for a game.
 * @param {string} gameId
 * @returns {Promise<GameCreature[]>}
 */
export async function fetchCreatures(gameId) {
  const res = await fetch(`${baseUrl}/game-creatures?game_id=${encodeURIComponent(gameId)}`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch creatures');
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
  const res = await fetch(`${baseUrl}/game-creatures`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ ...data, game_id: gameId })
  });
  if (!res.ok) throw new Error('Failed to create creature');
  const json = await res.json();
  return json.data;
}

/**
 * Update a creature by ID.
 * @param {string} creatureId
 * @param {Partial<GameCreature>} data
 * @returns {Promise<GameCreature>}
 */
export async function updateCreature(creatureId, data) {
  const res = await fetch(`${baseUrl}/game-creatures/${encodeURIComponent(creatureId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update creature');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a creature by ID.
 * @param {string} creatureId
 * @returns {Promise<void>}
 */
export async function deleteCreature(creatureId) {
  const res = await fetch(`${baseUrl}/game-creatures/${encodeURIComponent(creatureId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete creature');
} 