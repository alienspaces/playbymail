import { API_BASE_URL, getAuthHeaders } from './baseUrl';

/**
 * Fetch all items for a game.
 * @param {string} gameId
 * @returns {Promise<GameItem[]>}
 */
export async function fetchItems(gameId) {
  const res = await fetch(`${API_BASE_URL}/game-items?game_id=${encodeURIComponent(gameId)}`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch items');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new item for a game.
 * @param {string} gameId
 * @param {Partial<GameItem>} data
 * @returns {Promise<GameItem>}
 */
export async function createItem(gameId, data) {
  const res = await fetch(`${API_BASE_URL}/game-items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ ...data, game_id: gameId })
  });
  if (!res.ok) throw new Error('Failed to create item');
  const json = await res.json();
  return json.data;
}

/**
 * Update an item by ID.
 * @param {string} itemId
 * @param {Partial<GameItem>} data
 * @returns {Promise<GameItem>}
 */
export async function updateItem(itemId, data) {
  const res = await fetch(`${API_BASE_URL}/game-items/${encodeURIComponent(itemId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update item');
  const json = await res.json();
  return json.data;
}

/**
 * Delete an item by ID.
 * @param {string} itemId
 * @returns {Promise<void>}
 */
export async function deleteItem(itemId) {
  const res = await fetch(`${API_BASE_URL}/game-items/${encodeURIComponent(itemId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete item');
} 